package main

import "net/http"

func (app *application1) routes() *http.ServeMux {
	mux := http.NewServeMux()     // новый роутер
	mux.HandleFunc("/", app.home) // регистрирует функцию home как обработчик для роутера mux ...
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// инициализация FileServer, который будет обрабатывать HTTP-запросы к статическим файлам из ./ui/static/
	//fileServer := http.FileServer(neuterdFileSystem{http.Dir("./static")})
	// регистрация всех запросов начинающихся со "/static"
	//mux.Handle("/static", http.NotFoundHandler())
	fileServer := http.FileServer(http.Dir("../../ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
