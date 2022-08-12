package main

import (
	"log"
	"net/http"
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
	w.Write([]byte("Привет здарова"))
}

// Отображает определенную заметку
func showSnippet(w http.ResponseWriter, r *http.Request) { // "/snippet"
	w.Write([]byte("Оображение заетки..."))
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
