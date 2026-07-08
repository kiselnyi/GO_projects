// package scheduler

// import (
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"time"

// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// 	"tg_bot_on_go/internal/storage"
// )

// type Scheduler struct {
// 	bot    *tgbotapi.BotAPI
// 	chatID int64
// }

// func NewScheduler(bot *tgbotapi.BotAPI, chatID int64) *Scheduler {
// 	return &Scheduler{
// 		bot:    bot,
// 		chatID: chatID,
// 	}
// }

// // Start запускает боевой цикл проверки времени (раз в минуту)
// func (s *Scheduler) Start() {
// 	log.Println("Боевой планировщик дежурств успешно запущен...")
	
// 	ticker := time.NewTicker(1 * time.Minute)
// 	defer ticker.Stop()

// 	// Флаг, чтобы бот не спамил каждую минуту в течение одного часа
// 	var lastSentHour = -1

// 	for range ticker.C {
// 		now := time.Now()
// 		currentHour := now.Hour()
// 		currentMinute := now.Minute()

// 		// Если наступил новый час — сбрасываем флаг отправки
// 		if currentMinute == 0 {
// 			lastSentHour = -1
// 		}

// 		// Загружаем актуальный список дежурных из JSON
// 		duties, err := storage.LoadSchedule()
// 		if err != nil {
// 			log.Printf("Ошибка планировщика при чтении JSON: %v", err)
// 			continue
// 		}

// 		// Форматируем текущую дату под ключ JSON (ГГГГ-ММ-ДД)
// 		todayStr := now.Format("2006-01-02")
// 		dutyUser, exists := duties[todayStr]

// 		// Если на сегодня дежурный назначен и в этот час мы ему еще ничего не слали
// 		if exists && currentHour != lastSentHour {
			
// 			// 1. Железный пуш в 5:00 утра
// 			if currentHour == 5 && currentMinute == 0 {
// 				s.sendAlert(dutyUser, "🚨 ПОДЪЁМ! Время 5:00 утра. Проверка бдительности! Твоя очередь дежурить.")
// 				lastSentHour = currentHour
// 			}

// 			// 2. Рандомные напоминания в течение рабочего дня (с 8 до 23 часов)
// 			if currentHour >= 8 && currentHour <= 23 {
// 				// Шанс ~0.35% каждую минуту. За 15 часов бодрствования это даст в среднем 3 случайных пуша.
// 				if rand.Float64() < 0.0035 {
// 					s.sendAlert(dutyUser, "🔔 Внимание! Не расслабляемся, ты всё ещё на дежурстве!")
// 					lastSentHour = currentHour
// 				}
// 			}
// 		}
// 	}
// }

// func (s *Scheduler) sendAlert(username, message string) {
// 	text := fmt.Sprintf("%s, %s", username, message)
// 	msg := tgbotapi.NewMessage(s.chatID, text)
	
// 	_, err := s.bot.Send(msg)
// 	if err != nil {
// 		log.Printf("Ошибка отправки пуша планировщиком: %v", err)
// 	} else {
// 		log.Printf("Планировщик успешно отправил напоминание для %s", username)
// 	}
// }
package scheduler

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg_bot_on_go/internal/storage"
)

type Scheduler struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func NewScheduler(bot *tgbotapi.BotAPI, chatID int64) *Scheduler {
	return &Scheduler{
		bot:    bot,
		chatID: chatID,
	}
}

// Start запускает боевой цикл проверки времени (раз в минуту)
func (s *Scheduler) Start() {
	log.Println("Боевой планировщик дежурств успешно запущен...")
	
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// Флаг, чтобы бот не спамил каждую минуту в течение одного часа
	var lastSentHour = -1

	for range ticker.C {
		now := time.Now()
		currentHour := now.Hour()
		currentMinute := now.Minute()

		// Если наступил новый час — сбрасываем флаг отправки
		if currentMinute == 0 {
			lastSentHour = -1
		}

		// Загружаем актуальный список дежурных из JSON
		duties, err := storage.LoadSchedule()
		if err != nil {
			log.Printf("Ошибка планировщика при чтении JSON: %v", err)
			continue
		}

		// Ищем дежурного с учетом циклического наследования прошлых недель
		dutyUser, exists := s.getDutyWithFallback(duties, now)

		// Если дежурный найден (в JSON или унаследован) и в этот час мы ему еще ничего не слали
		if exists && currentHour != lastSentHour {
			
			// 1. Железный пуш в 5:00 утра
			if currentHour == 5 && currentMinute == 0 {
				s.sendAlert(dutyUser, "🚨 ПОДЪЁМ! Время 5:00 утра. Проверка бдительности! Твоя очередь дежурить.")
				lastSentHour = currentHour
			}

			// 2. Рандомные напоминания в течение рабочего дня (с 8 до 23 часов)
			if currentHour >= 8 && currentHour <= 23 {
				// Шанс ~0.35% каждую минуту. За 15 часов бодрствования это даст в среднем 3 случайных пуша.
				if rand.Float64() < 0.0035 {
					s.sendAlert(dutyUser, "🔔 Внимание! Не расслабияемся, ты всё ещё на дежурстве!")
					lastSentHour = currentHour
				}
			}
		}
	}
}

// getDutyWithFallback ищет дежурного на указанную дату. 
// Если на эту дату никто не записан, метод отматывает время назад на 1, 2, 3 или 4 недели в поисках циклического графика.
func (s *Scheduler) getDutyWithFallback(duties map[string]string, targetTime time.Time) (string, bool) {
	// Делаем до 4 попыток найти дежурного на этот же день недели в прошлом
	for i := 0; i < 4; i++ {
		dateStr := targetTime.Format("2006-01-02")
		if user, found := duties[dateStr]; found && user != "" {
			return user, true
		}
		// Отматываем ровно на 7 дней назад
		targetTime = targetTime.AddDate(0, 0, -7)
	}
	return "", false
}

func (s *Scheduler) sendAlert(username, message string) {
	text := fmt.Sprintf("%s, %s", username, message)
	msg := tgbotapi.NewMessage(s.chatID, text)
	
	_, err := s.bot.Send(msg)
	if err != nil {
		log.Printf("Ошибка отправки пуша планировщиком: %v", err)
	} else {
		log.Printf("Планировщик успешно отправил напоминание для %s", username)
	}
}