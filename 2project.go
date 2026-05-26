package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sort"
)

func Quiz(event string, correctYear int) bool {
	fmt.Println("B каком году произошло:", event)
	fmt.Println("Ответ:")

	var answer int
	fmt.Scan(&answer)

	if answer == correctYear {
		fmt.Println("Верно")
		return true
	} else {
		fmt.Println("Неверно. Правильный ответ:", correctYear)
		return false
	}
}

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

func filterByCountry(events []Event, country string) map[int]string {
	result := map[int]string{}
	for _, event := range events {
		if event.Country == country {
			result[event.Year] = event.Description
		}
	}
	return result
}

func main() {
	events, err := loadEvents("events.json")
	if err != nil {
		fmt.Println("Ошибка загрузки файла:", err)
		return
	}
	fmt.Println("Загружено событий:", len(events))

	{
		russia := filterByCountry(events, "Российская империя")

		france := filterByCountry(events, "Франция")

		england := filterByCountry(events, "Великобритания")

		for {
			var choice int
			fmt.Println("Здесь ты можешь ознакомиться c историческими событиями распределенными по датам. Для начала выбери страну:")
			fmt.Println("1 - Российская империя")
			fmt.Println("2 - Франция")
			fmt.Println("3 - Великобритания")
			fmt.Println("Выбор: ")
			fmt.Scan(&choice)

			var events map[int]string
			if choice == 1 {
				events = russia
			} else if choice == 2 {
				events = france
			} else if choice == 3 {
				events = england
			} else {
				fmt.Println("Неверный выбор")
				return
			}

			fmt.Println("Что сделать?")
			fmt.Println("1 - Найти событие по году")
			fmt.Println("2 - Показать все события по определенной стране")
			fmt.Println("3 - Случайное событие")
			fmt.Println("4 - Квиз - проверь свои знания")
			fmt.Println("Вариант: ")
			var action int
			fmt.Scan(&action)

			if action == 2 {
				list := make([]int, 0, len(events))
				for year := range events {
					list = append(list, year)
				}
				sort.Ints(list)
				for _, year := range list {
					fmt.Println(year, "-", events[year])
				}

			}

			if action == 1 {
				var year int
				fmt.Print("Год (1801-1810): ")
				fmt.Scan(&year)

				event, found := events[year]
				if found {
					fmt.Println("Событие: ", event)
				} else {
					fmt.Println("Событие не нашлось")
				}
			}

			if action == 3 {
				list := make([]int, 0, len(events))
				for year := range events {
					list = append(list, year)
				}
				randomIndex := rand.Intn(len(list))
				randomYear := list[randomIndex]
				fmt.Println(randomYear, "-", events[randomYear])
			}
			if action == 4 {
				list := make([]int, 0, len(events))
				for year := range events {
					list = append(list, year)
				}
				rand.Shuffle(len(list), func(x, y int) {
					list[x], list[y] = list[y], list[x]

				})

				score := 0

				for _, year := range list {
					correct := Quiz(events[year], year)
					if correct {
						score = score + 1
					}

					fmt.Println("Продолжить квиз? (1 - Да, 2 - Закончить квиз)")
					var next int
					fmt.Scan(&next)
					if next == 2 {
						break
					}
				}
				fmt.Println("Результат:", score, "из", len(list))
			}

			fmt.Println("Продолжить? (1 - Вернуться в главное меню, 2 - Закрыть программу)")
			var again int
			fmt.Scan(&again)
			if again == 2 {
				return
			}
		}
	}
}
