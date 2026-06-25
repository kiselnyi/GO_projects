package config

import (
	"errors"
	"os"
	"strconv"
)

// Config хранит все конфигурационные данные нашего бота
type Config struct {
	BotToken string
	ChatID   int64
}

// LoadConfig собирает конфигурацию из переменных окружения,
// куда их безопасно прокидывает Vault при старте контейнера/сервиса.
func LoadConfig() (*Config, error) {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		return nil, errors.New("BOT_TOKEN не найден в переменных окружения")
	}

	chatIDStr := os.Getenv("CHAT_ID")
	if chatIDStr == "" {
		return nil, errors.New("CHAT_ID не найден в переменных окружения")
	}

	// Конвертируем строку из ENV в int64 для Telegram API
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return nil, errors.New("CHAT_ID должен быть валидным числом")
	}

	return &Config{
		BotToken: token,
		ChatID:   chatID,
	}, nil
}