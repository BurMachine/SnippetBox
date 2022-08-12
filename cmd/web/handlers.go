package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

/*
	ОБРАБОТЧИКИ
	http.ResponseWriter предоствляет методы для объединения HTTP ответа и возвращение его пользователю
	*http.Request - это указател на структуру, которая содержит информацию о текущем запросе(GET, POST, DELETE, etc...)
	w - writer куда это все пишется, r - хранит все реквесты(запросы)
*/
func home(w http.ResponseWriter, r *http.Request) { // "/"
	if r.URL.Path != "/" { // Обработка неправильного URL
		http.NotFound(w, r)
		return
	}
	pageTemp, err := template.ParseFiles("../../ui/html/home.page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error1", 500)
		return
	}
	err = pageTemp.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error2", 500)
	}
	//w.Write([]byte("Привет здарова"))
}

// Отображает определенную заметку
func showSnippet(w http.ResponseWriter, r *http.Request) { // "/snippet"
	id, err := strconv.Atoi(r.URL.Query().Get("id")) // Считывание значения id. Затем проверка
	if err != nil || id < 0 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Отображение заметки с ID %d...", id)
}

// Создает новую заметку
func createSnippet(w http.ResponseWriter, r *http.Request) { // "/snippet/create"
	if r.Method != http.MethodPost {
		w.Header().Set("Allowed-method", http.MethodPost) // добавляет ключ:значение в карту HTTP
		http.Error(w, "Метод запрещен!", 405)             // отправляет в ResponceWriter строку и в карту HTTP еод ошибки (перед write обязательно)
		return
	}
	w.Write([]byte("Форма для создания новой заметки..."))
}
