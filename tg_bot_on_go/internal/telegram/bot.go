// package telegram

// import (
// 	"log"

// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// )

// // Структура нашего бота, которая будет хранить сам клиент ТГ
// type BotRunner struct {
// 	TgBot *tgbotapi.BotAPI
// }

// // NewBotRunner создает и авторизует нового бота
// func NewBotRunner(token string) (*BotRunner, error) {
// 	bot, err := tgbotapi.NewBotAPI(token)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Включаем дебаг, чтобы в консоли видеть все логи запросов к API Telegram
// 	bot.Debug = true

// 	log.Printf("Бот успешно авторизован под аккаунтом: %s", bot.Self.UserName)
// 	return &BotRunner{TgBot: bot}, nil
// }

// // StartListening запускает бесконечный цикл прослушивания сообщений
// func (br *BotRunner) StartListening() {
// 	u := tgbotapi.NewUpdate(0)
// 	u.Timeout = 60

// 	updates := br.TgBot.GetUpdatesChan(u)

// 	log.Println("Ждем сообщений от пользователя...")

// 	for update := range updates {
// 		// Нам интересны только текстовые сообщения
// 		if update.Message == nil || update.Message.Text == "" {
// 			continue
// 		}

// 		msgText := update.Message.Text
// 		chatID := update.Message.Chat.ID

// 		var replyText string

// 		// Если сообщение начинается с нашей команды обновления
// 		if len(msgText) >= 16 && msgText[:16] == "/update_schedule" {
// 			var err error
// 			// Вызываем парсер из файла handlers.go (он в этом же пакете, так что доступен напрямую)
// 			replyText, err = ParseAndSaveSchedule(msgText)
// 			if err != nil {
// 				replyText = "Произошла внутренняя ошибка при сохранении."
// 			}
// 		} else if msgText == "/start" {
// 			replyText = "Привет! Отправь мне команду /update_schedule и список дежурных списком."
// 		} else {
// 			replyText = "Я знаю только команду /update_schedule"
// 		}

// 		// Отправляем ответ пользователю в ТГ
// 		msg := tgbotapi.NewMessage(chatID, replyText)
// 		_, err := br.TgBot.Send(msg)
// 		if err != nil {
// 			log.Printf("Ошибка отправки сообщения: %v", err)
// 		}
// 	}
// }
package telegram

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Структура нашего бота, которая будет хранить сам клиент ТГ
type BotRunner struct {
	TgBot *tgbotapi.BotAPI
}

// NewBotRunner создает и авторизует нового бота
func NewBotRunner(token string) (*BotRunner, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	// Включаем дебаг, чтобы в консоли видеть все логи запросов к API Telegram
	bot.Debug = true

	log.Printf("Бот успешно авторизован под аккаунтом: %s", bot.Self.UserName)
	return &BotRunner{TgBot: bot}, nil
}

// StartListening запускает бесконечный цикл прослушивания сообщений
func (br *BotRunner) StartListening() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := br.TgBot.GetUpdatesChan(u)

	log.Println("Ждем сообщений от пользователя...")

	for update := range updates {
		// Нам интересны только текстовые сообщения
		if update.Message == nil || update.Message.Text == "" {
			continue
		}

		msgText := update.Message.Text
		chatID := update.Message.Chat.ID
		var replyText string

		// Проверяем, является ли сообщение встроенной телеграм-командой (начинается с /)
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			
			case "start":
				replyText = "Привет! Я бот-планировщик дежурств.\n\n" +
					"Доступные команды:\n" +
					"🔹 <b>/update_schedule</b> — обновить график\n" +
					"🔹 <b>/clear_schedule</b> — полностью очистить расписание"

			case "update_schedule":
				var err error
				// Передаем весь текст сообщения в парсер
				replyText, err = ParseAndSaveSchedule(msgText)
				if err != nil {
					log.Printf("Ошибка при обновлении расписания: %v", err)
					replyText = "❌ Произошла внутренняя ошибка при сохранении."
				}

			case "clear_schedule":
				var err error
				// Вызываем твой новый обработчик очистки
				replyText, err = ClearScheduleHandler()
				if err != nil {
					log.Printf("Ошибка при очистке JSON: %v", err)
					replyText = "❌ Произошла внутренняя ошибка при очистке файла."
				}

			default:
				replyText = "🤔 Неизвестная команда. Доступны: /update_schedule и /clear_schedule"
			}
		} else {
			// Если прислали обычный текст без слэша в начале
			// Если случайно прислали список дежурств, забыв написать /update_schedule в начале
			if strings.Contains(msgText, "@") {
				replyText = "💡 Чтобы обновить расписание, начни сообщение с команды `/update_schedule`, а затем с новой строки перечисли дежурных."
			} else {
				replyText = "🤖 Я реагирую только на команды. Введи `/start`, чтобы посмотреть, что я умею."
			}
		}

		// Отправляем ответ пользователю в ТГ
		msg := tgbotapi.NewMessage(chatID, replyText)
		msg.ParseMode = "HTML" // Добавим поддержку разметки, чтобы красиво подсвечивать команды
		
		_, err := br.TgBot.Send(msg)
		if err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
		}
	}
}