package storage

import (
	"encoding/json"
	"os"
)

// Структура для JSON. Оставляем имя поля с большой буквы (Duties), 
// чтобы к нему можно было обращаться из других пакетов (например, из бота).
type Schedule struct {
	Duties map[string]string `json:"duties"`
}

// Файл будет лежать в корне проекта, откуда запускается бот
const filename = "schedule.json"

// SaveSchedule сохраняет мапу дежурных в JSON-файл
func SaveSchedule(duties map[string]string) error {
	s := Schedule{Duties: duties}

	// MarshalIndent делает красивый JSON с отступами
	fileData, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, fileData, 0644)
}

// LoadSchedule читает JSON-файл и возвращает мапу дежурных
func LoadSchedule() (map[string]string, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return make(map[string]string), nil
	}

	fileData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var s Schedule
	err = json.Unmarshal(fileData, &s)
	if err != nil {
		return nil, err
	}

	return s.Duties, nil
}