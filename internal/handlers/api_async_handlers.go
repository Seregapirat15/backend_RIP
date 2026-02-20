package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"lab4/internal/database"
)

const serviceTokenEnv = "SERVICE_TOKEN"
const asyncServiceURLEnv = "ASYNC_SERVICE_URL"
const defaultServiceToken = "Lab8Token"
const defaultAsyncURL = "http://localhost:8001"

// SubmitMassResultHandler - POST /api/internal/submit-mass-result
// Принимает результат от асинхронного сервиса. Псевдо-авторизация через X-Service-Token.
func SubmitMassResultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("X-Service-Token")
	expected := os.Getenv(serviceTokenEnv)
	if expected == "" {
		expected = defaultServiceToken
	}
	if token != expected {
		http.Error(w, `{"error":"Неверный токен"}`, http.StatusUnauthorized)
		return
	}

	var req struct {
		OrderID       int     `json:"order_id"`
		InstrumentID  int     `json:"instrument_id"`
		CalculatedMass float64 `json:"calculated_mass"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Неверный JSON"}`, http.StatusBadRequest)
		return
	}

	err := database.UpdateMMCalculatedMass(req.OrderID, req.InstrumentID, req.CalculatedMass)
	if err != nil {
		log.Printf("SubmitMassResult: %v", err)
		http.Error(w, `{"error":"Ошибка обновления"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{"ok": true})
}

// callAsyncService вызывает асинхронный сервис для расчёта массы
func callAsyncService(orderID int) error {
	items, err := database.GetOrderMMForAsync(orderID)
	if err != nil || len(items) == 0 {
		return err
	}

	url := os.Getenv(asyncServiceURLEnv)
	if url == "" {
		url = defaultAsyncURL
	}
	url += "/calculate"

	mmItems := make([]map[string]any, len(items))
	for i, it := range items {
		mmItems[i] = map[string]any{
			"instrument_id":      it.InstrumentID,
			"star_mass":          it.StarMass,
			"orbital_period":     it.OrbitalPeriod,
			"velocity_amplitude": it.VelocityAmplitude,
			"inclination":        it.Inclination,
		}
	}
	body := map[string]any{"order_id": orderID, "items": mmItems}
	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 202 {
		log.Printf("Async service вернул %d", resp.StatusCode)
		return nil // не блокируем, async мог быть недоступен
	}
	return nil
}

// TriggerAsyncCalculationHandler - POST /api/admin/orders/{id}/trigger-calculation
// Ручной запуск асинхронного расчёта (кнопка в UI модератора)
func TriggerAsyncCalculationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"Неверный ID"}`, http.StatusBadRequest)
		return
	}

	if err := callAsyncService(orderID); err != nil {
		log.Printf("TriggerAsyncCalculation: %v", err)
		http.Error(w, `{"error":"Не удалось запустить расчёт"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Расчёт запущен",
	})
}
