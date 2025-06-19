package main

import (
	"log"
	"net/http"
)

// @title Система создания учебного плана
// @version 1.0
// @description Дает возможность создавать учебный план
// @host localhost:8085
// @BasePath /home
func main() {
	http.Handle("/", http.FileServer(http.Dir("templates")))
	//Обработка кнопок
	http.HandleFunc("/home/create_item", HandlerCreateItem)
	http.HandleFunc("/home/getPlan", getPlan)
	http.HandleFunc("/home/updatePlan", updatePlan)
	http.HandleFunc("/home/deletePlan", deletePlan)

	//Код запуска сервера
	log.Println("Запуск сервера")
	err := http.ListenAndServe(":8085", nil)
	if err != nil {
		log.Fatal(err)
	}
}
