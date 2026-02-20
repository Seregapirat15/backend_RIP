-- Установка пароля модератора
-- Сначала попробуйте: moderator / password
-- Если не подходит, сгенерируйте хеш: cd backend_RIP && go run ./scripts/gen_password.go
-- Скопируйте хеш для modpass123 и выполните:
--   UPDATE users SET password_hash = 'ПОДСТАВЬТЕ_ХЕШ' WHERE login = 'moderator';

-- Пароль "password" (стандартный для примеров):
UPDATE users SET password_hash = '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi' 
WHERE login IN ('moderator', 'admin', 'user1', 'user2');
