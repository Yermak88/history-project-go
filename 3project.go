package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Event struct {
	Year        int    `json:"year"`
	Description string `json:"description"`
	Country     string `json:"country"`
}

type UpdateRequest struct {
	OldDescription string `json:"oldDescription"`
	Year           int    `json:"year"`
	Description    string `json:"description"`
	Country        string `json:"country"`
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

func filterByYear(events []Event, year int) []Event {
	result := []Event{}
	for _, event := range events {
		if event.Year == year {
			result = append(result, event)
		}
	}
	return result
}

func loggerMiddleware(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, r.URL.Path)
		handler(w, r)
	}
}

func enableCORS(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		handler(w, r)
	}
}

func main() {
	events, err := loadEvents("events.json")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/events", loggerMiddleware(enableCORS(func(w http.ResponseWriter, r *http.Request) {
		yearStr := r.URL.Query().Get("year")

		var result []Event
		if yearStr == "" {
			result = events
		} else {
			year, err := strconv.Atoi(yearStr)
			if err != nil {
				http.Error(w, "Неверный формат", 400)
				return
			}
			result = filterByYear(events, year)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})))

	http.HandleFunc("/events/add", loggerMiddleware(enableCORS(func(w http.ResponseWriter, r *http.Request) {
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
	})))

	http.HandleFunc("/events/update", loggerMiddleware(enableCORS(func(w http.ResponseWriter, r *http.Request) {
		var req UpdateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Неверный формат", 400)
			return
		}

		found := false
		foundIndex := -1
		for i, event := range events {
			if event.Description == req.OldDescription {
				foundIndex = i
				found = true
				break
			}
		}

		if found {
			events[foundIndex] = Event{Year: req.Year, Description: req.Description, Country: req.Country}
		}

		if !found {
			http.Error(w, "Событие не найдено", 400)
			return
		}

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
		json.NewEncoder(w).Encode(events[foundIndex])

	})))

	http.HandleFunc("/events/delete", loggerMiddleware(enableCORS(func(w http.ResponseWriter, r *http.Request) {
		var eventToDelete Event
		err := json.NewDecoder(r.Body).Decode(&eventToDelete)
		if err != nil {
			http.Error(w, "Неверный формат", 400)
			return
		}

		found := false
		for i, event := range events {
			if event.Description == eventToDelete.Description {
				events = append(events[:i], events[i+1:]...)
				found = true
				break
			}
		}

		if !found {
			http.Error(w, "Событие не найдено", 404)
			return
		}

		data, err := json.MarshalIndent(events, "", "    ")
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
		json.NewEncoder(w).Encode("Событие удалено")

	})))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Сервер запускается на порту:", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Не удалось запустить сервер:", err)
	}
}
