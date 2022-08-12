package main

import (
	"log"
	"net/http"
)

/*
	http.NewServeMux - инициализация нового роутера
	mux.HandleFunc("/", home) регистрация home как обработчика url шаблона "/"
	http.ListenAndServe  - получает в качестве параметров TCP-адрес сети для прослушивания (localhost:4000)
	и созданный роутер
*/
func main() {
	mux := http.NewServeMux() // новый роутер
	mux.HandleFunc("/", home) // регистрирует функцию home как обработчик для роутера mux ...
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	log.Println("Запуск веб-сервера на http://127.0.0.1:4000")
	err := http.ListenAndServe(":4000", mux) // Запуск нового веб-сервера
	log.Fatal(err)
}
