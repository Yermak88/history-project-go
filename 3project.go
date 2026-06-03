package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Event struct {
	Year        int    `json:"year"`
	Description string `json:"description"`
	Country     string `json:"country"`
}

func loadEvents(filename string) ([]Event, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var events []Event
	err = json.Unmarshal(file, &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func filterByCountry(events []Event, country string) []Event {
	result := []Event{}
	for _, event := range events {
		if event.Country == country {
			result = append(result, event)
		}
	}
	return result
}

func main() {
	events, err := loadEvents("events.json")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		country := r.URL.Query().Get("country")

		var result []Event
		if country == "" {
			result = events
		} else {
			result = filterByCountry(events, country)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	http.HandleFunc("/events/add", func(w http.ResponseWriter, r *http.Request) {
		var newEvent Event
		err := json.NewDecoder(r.Body).Decode(&newEvent)
		if err != nil {
			http.Error(w, "Неверный формат", 400)
			return
		}

		for _, event := range events {
			if event.Description == newEvent.Description {
				http.Error(w, "Это событие уже существует", 400)
				return
			}
		}

		events = append(events, newEvent)
		data, err := json.MarshalIndent(events, "", "   ")
		if err != nil {
			http.Error(w, "Ошибка сохранения", 500)
			return
		}

		err = os.WriteFile("events.json", data, 0644)
		if err != nil {
			http.Error(w, "Ошибка сохранения", 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newEvent)
	})

	fmt.Println("Сервер запущен: https://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
