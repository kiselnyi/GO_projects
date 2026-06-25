package telegram

import (
	"strings"
	"tg_bot_on_go/internal/storage"
)

// ParseAndSaveSchedule принимает весь текст сообщения, парсит его и сохраняет
func ParseAndSaveSchedule(text string) (string, error) {
	// Создаем мапу, куда будем складывать результат
	newDuties := make(map[string]string)

	// Разбиваем текст на отдельные строчки
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		// Очищаем пробелы по краям строки
		line = strings.TrimSpace(line)

		// Игнорируем саму команду /update_schedule или пустые строки
		if line == "" || strings.HasPrefix(line, "/") {
			continue
		}

		// Разбиваем строку на две части по разделителю " - "
		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			// Если строка какая-то кривая, просто пропускаем её (или можно выдать ошибку)
			continue
		}

		// Убираем лишние пробелы вокруг даты и никнейма
		date := strings.TrimSpace(parts[0])
		username := strings.TrimSpace(parts[1])

		// Переводим дату из формата "ДД.ММ.ГГГГ" в "ГГГГ-ММ-ДД", чтобы JSON красиво сортировался.
		// Для простоты сделаем это через разбиение строки, раз мы учимся.
		dateParts := strings.Split(date, ".")
		if len(dateParts) == 3 {
			// Пересобираем в ГГГГ-ММ-ДД
			date = dateParts[2] + "-" + dateParts[1] + "-" + dateParts[0]
		}

		// Сохраняем в мапу
		newDuties[date] = username
	}

	if len(newDuties) == 0 {
		return "Не удалось найти ни одной валидной строки с дежурством.", nil
	}

	// Читаем текущее расписание, чтобы не затереть прошлые даты, а ОБНОВИТЬ их
	currentDuties, err := storage.LoadSchedule()
	if err != nil {
		return "", err
	}

	// Добавляем новые данные в существующие
	for k, v := range newDuties {
		currentDuties[k] = v
	}

	// Сохраняем обновленную мапу в JSON
	err = storage.SaveSchedule(currentDuties)
	if err != nil {
		return "", err
	}

	return "Расписание успешно обновлено!", nil
}