package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"lab4/internal/database"
	"lab4/internal/middleware"
	"lab4/internal/models"
)

// GetServicesHandler - GET /api/services - список услуг с фильтрацией
func GetServicesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Получение параметров фильтрации
	filter := models.ServiceFilter{
		Type:     r.URL.Query().Get("type"),
		Status:   r.URL.Query().Get("status"),
		Search:   r.URL.Query().Get("search"),
		DateFrom: getQueryPtr(r, "date_from"),
		DateTo:   getQueryPtr(r, "date_to"),
	}

	if minAccuracy := getQueryFloat(r, "min_accuracy"); minAccuracy != nil {
		filter.MinAccuracy = minAccuracy
	}
	if maxAccuracy := getQueryFloat(r, "max_accuracy"); maxAccuracy != nil {
		filter.MaxAccuracy = maxAccuracy
	}

	services, err := database.GetServicesWithFilter(filter)
	if err != nil {
		log.Printf("Ошибка получения услуг: %v", err)
		http.Error(w, "Ошибка получения услуг", http.StatusInternalServerError)
		return
	}

	// Подставляем presigned URL для изображений из MinIO
	for i := range services {
		fillServiceImageURL(&services[i])
	}

	json.NewEncoder(w).Encode(services)
}

// GetServiceHandler - GET /api/services/{id} - одна услуга
func GetServiceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Извлечение ID из URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Неверный URL", http.StatusBadRequest)
		return
	}

	serviceID, err := strconv.Atoi(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "Неверный ID услуги", http.StatusBadRequest)
		return
	}

	service, err := database.GetServiceByID(serviceID)
	if err != nil {
		log.Printf("Ошибка получения услуги: %v", err)
		http.Error(w, "Услуга не найдена", http.StatusNotFound)
		return
	}

	fillServiceImageURL(service)
	json.NewEncoder(w).Encode(service)
}

// CreateServiceHandler - POST /api/services - добавление услуги
func CreateServiceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var service models.Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	// Устанавливаем системные поля
	service.CreatedAt = time.Now()
	service.IsDeleted = false
	service.Status = "Активен"

	serviceID, err := database.CreateService(service)
	if err != nil {
		log.Printf("Ошибка создания услуги: %v", err)
		http.Error(w, "Ошибка создания услуги", http.StatusInternalServerError)
		return
	}

	service.ID = serviceID
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(service)
}

// UpdateServiceHandler - PUT /api/services/{id} - изменение услуги
func UpdateServiceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Извлечение ID из URL
	pathParts := strings.Split(r.URL.Path, "/")
	serviceID, err := strconv.Atoi(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "Неверный ID услуги", http.StatusBadRequest)
		return
	}

	var service models.Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	service.ID = serviceID
	err = database.UpdateService(service)
	if err != nil {
		log.Printf("Ошибка обновления услуги: %v", err)
		http.Error(w, "Ошибка обновления услуги", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(service)
}

// DeleteServiceHandler - DELETE /api/services/{id} - удаление услуги
func DeleteServiceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Извлечение ID из URL
	pathParts := strings.Split(r.URL.Path, "/")
	serviceID, err := strconv.Atoi(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "Неверный ID услуги", http.StatusBadRequest)
		return
	}

	err = database.DeleteService(serviceID)
	if err != nil {
		log.Printf("Ошибка удаления услуги: %v", err)
		http.Error(w, "Ошибка удаления услуги", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddServiceToOrderHandler - POST /api/orders/services - добавление услуги в заявку
func AddServiceToOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Логируем запрос для отладки
	log.Printf("AddServiceToOrderHandler: Method=%s, Path=%s, AuthHeader=%s", r.Method, r.URL.Path, r.Header.Get("Authorization"))

	// Получаем ID пользователя из контекста
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		log.Printf("AddServiceToOrderHandler: Пользователь не авторизован")
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	log.Printf("AddServiceToOrderHandler: Пользователь авторизован, userID=%d", userID)

	var orderService models.OrderService
	if err := json.NewDecoder(r.Body).Decode(&orderService); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	// Получаем или создаем заявку-черновик
	order, err := database.GetOrCreateDraftOrder(userID)
	if err != nil {
		log.Printf("Ошибка получения/создания заявки: %v", err)
		http.Error(w, "Ошибка работы с заявкой", http.StatusInternalServerError)
		return
	}

	orderService.OrderID = order.ID
	err = database.AddServiceToOrder(orderService)
	if err != nil {
		log.Printf("Ошибка добавления услуги в заявку: %v", err)

		// Проверяем, является ли ошибка дубликатом
		if err.Error() == "услуга уже добавлена в заявку" {
			http.Error(w, "Эта услуга уже добавлена в заявку", http.StatusConflict)
			return
		}

		http.Error(w, "Ошибка добавления услуги в заявку", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"order_id":   order.ID,
		"service_id": orderService.ServiceID,
		"message":    "Услуга успешно добавлена в заявку",
	})
}

// UploadServiceImageHandler - POST /api/services/{id}/image - добавление изображения в MinIO
func UploadServiceImageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pathParts := strings.Split(r.URL.Path, "/")
	serviceID, err := strconv.Atoi(pathParts[len(pathParts)-2])
	if err != nil {
		http.Error(w, "Неверный ID услуги", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB
		http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Файл не найден (ожидается поле 'image')", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := fmt.Sprintf("service_%d_%d.jpg", serviceID, time.Now().Unix())

	objectName, err := database.UploadServiceImage(serviceID, fileName, file, handler)
	if err != nil {
		log.Printf("Ошибка загрузки изображения: %v", err)
		http.Error(w, "Ошибка загрузки изображения: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отдаём клиенту presigned URL для отображения
	imageURL, _ := database.GetTelescopeImagePresignedURL(objectName)
	if imageURL == "" {
		imageURL = "http://localhost:9000/telescope-images/" + objectName
	}

	json.NewEncoder(w).Encode(map[string]string{
		"image_url": imageURL,
	})
}

// fillServiceImageURL подставляет presigned URL для image_url из MinIO (если в БД хранится ключ объекта)
func fillServiceImageURL(s *models.Service) {
	if s == nil || s.ImageURL == nil || *s.ImageURL == "" {
		return
	}
	key := *s.ImageURL
	if strings.HasPrefix(key, "http") {
		return // уже полный URL
	}
	url, err := database.GetTelescopeImagePresignedURL(key)
	if err == nil && url != "" {
		s.ImageURL = &url
	}
}

// Вспомогательные функции для парсинга параметров

func getQueryPtr(r *http.Request, key string) *string {
	val := r.URL.Query().Get(key)
	if val == "" {
		return nil
	}
	return &val
}

func getQueryFloat(r *http.Request, key string) *float64 {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		return nil
	}
	if val, err := strconv.ParseFloat(valStr, 64); err == nil {
		return &val
	}
	return nil
}
