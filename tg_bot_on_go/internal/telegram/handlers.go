package telegram

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
	"tg_bot_on_go/internal/storage"
)

// Регулярка: принимает дату ЛИБО буквы (пн, вт, mon...)
var dutyRegex = regexp.MustCompile(`^(\d{2}\.\d{2}\.\d{4}|[а-яА-Яa-zA-Z]{2,3})\s*[-—–]?\s*(@\w+)\s*$`)

// Мапа для конвертации дней недели
var daysMap = map[string]time.Weekday{
	"пн": time.Monday, "пон": time.Monday,
	"вт": time.Tuesday, "вто": time.Tuesday,
	"ср": time.Wednesday, "сре": time.Wednesday,
	"чт": time.Thursday, "чет": time.Thursday,
	"пт": time.Friday, "пят": time.Friday,
	"сб": time.Saturday, "суб": time.Saturday,
	"вс": time.Sunday, "вос": time.Sunday,
}

// Функция-помощник для перевода Weekday в ISO формат (Пн=1, Вт=2 ... Вс=7)
func isoWeekday(t time.Weekday) int {
	if t == time.Sunday {
		return 7
	}
	return int(t)
}

func ParseAndSaveSchedule(text string) (string, error) {
	newDuties := make(map[string]string)
	lines := strings.Split(text, "\n")
	now := time.Now()

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "/") {
			continue
		}

		matches := dutyRegex.FindStringSubmatch(line)
		if len(matches) != 3 {
			continue
		}

		dayOrDate := strings.ToLower(strings.TrimSpace(matches[1]))
		username := matches[2]
		var finalDate string

		// Если ввели день недели (пн, вт...)
		if weekday, isWeekday := daysMap[dayOrDate]; isWeekday {
			// Считаем разницу по ISO-стандарту (Пн-Вс)
			daysDiff := isoWeekday(weekday) - isoWeekday(now.Weekday())

			targetTime := now.AddDate(0, 0, daysDiff)
			finalDate = targetTime.Format("2006-01-02")
		} else {
			// Если ввели обычную дату (ДД.ММ.ГГГГ)
			dateParts := strings.Split(dayOrDate, ".")
			if len(dateParts) == 3 {
				finalDate = dateParts[2] + "-" + dateParts[1] + "-" + dateParts[0]
			} else {
				continue
			}
		}

		newDuties[finalDate] = username
	}

	if len(newDuties) == 0 {
		return "❌ Не удалось распознать расписание. Используй формат: `пн @username` или `ДД.ММ.ГГГГ @username`", nil
	}

	// Читаем текущее расписание
	currentDuties, err := storage.LoadSchedule()
	if err != nil {
		return "", err
	}
	if currentDuties == nil {
		currentDuties = make(map[string]string)
	}

	// Мержим новые дежурства
	for k, v := range newDuties {
		currentDuties[k] = v
	}

	// === 1. АВТООЧИСТКА СТАРЬЯ (Всё, что старше недели) ===
	cutoff := now.AddDate(0, 0, -7)
	for dateStr := range currentDuties {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err == nil && parsedDate.Before(cutoff) {
			delete(currentDuties, dateStr)
		}
	}

	// Сохраняем отфильтрованную мапу в JSON
	err = storage.SaveSchedule(currentDuties)
	if err != nil {
		return "", err
	}

	// === 2. СОРТИРОВКА ВЫВОДА ===
	// Собираем все ключи (даты) в слайс
	var keys []string
	for k := range currentDuties {
		keys = append(keys, k)
	}
	// Сортируем строки (так как формат ГГГГ-ММ-ДД, сорт отработает идеально хронологически)
	sort.Strings(keys)

	// Строим красивый упорядоченный ответ
	response := "✅ Расписание успешно обновлено!\n\n📋 <b>Актуальный список дежурств:</b>\n"
	for _, date := range keys {
		user := currentDuties[date]
		displayDate := date
		rParts := strings.Split(date, "-")
		if len(rParts) == 3 {
			displayDate = rParts[2] + "." + rParts[1] + "." + rParts[0]
		}
		response += fmt.Sprintf("📅 %s — %s\n", displayDate, user)
	}

	return response, nil
}
func ClearScheduleHandler() (string, error) {
	// Просто сохраняем пустую мапу в JSON
	emptySchedule := make(map[string]string)
	
	err := storage.SaveSchedule(emptySchedule)
	if err != nil {
		return "", fmt.Errorf("не удалось очистить файл: %w", err)
	}

	return "🗑 Расписание дежурств полностью очищено! Наследование прошлых недель также сброшено.", nil
}