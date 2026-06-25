package main

import (
	"log"

	"tg_bot_on_go/internal/config"
	"tg_bot_on_go/internal/scheduler"
	"tg_bot_on_go/internal/telegram"
)

func main() {
	log.Println("Запуск DevOps Duty Bot...")

	// 1. Загружаем конфиг через наш новый пакет
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Критическая ошибка конфигурации: %v", err)
	}

	// 2. Инициализируем бота, передавая токен из конфига
	runner, err := telegram.NewBotRunner(cfg.BotToken)
	if err != nil {
		log.Fatalf("Ошибка инициализации бота: %v", err)
	}

	// 3. Создаем планировщик, передавая ChatID из конфига
	sched := scheduler.NewScheduler(runner.TgBot, cfg.ChatID)

	// 4. Запускаем планировщик параллельно в горутине
	go sched.Start()

	// 5. Блокируем поток и слушаем входящие сообщения в ТГ
	runner.StartListening()
}