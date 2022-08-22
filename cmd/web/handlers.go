package main

import (
	"errors"
	"fmt"
	"golangify.com/SnippetBox/pkg/models"
	"net/http"
	"strconv"
)

/*
	ОБРАБОТЧИКИ
	http.ResponseWriter предоствляет методы для объединения HTTP ответа и возвращение его пользователю
	*http.Request - это указател на структуру, которая содержит информацию о текущем запросе(GET, POST, DELETE, etc...)
	w - writer куда это все пишется, r - хранит все реквесты(запросы)
*/
func (app *application1) home(w http.ResponseWriter, r *http.Request) { // "/"
	if r.URL.Path != "/" { // Обработка неправильного URL
		app.notFound(w)
		return
	}
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

// Отображает определенную заметку
func (app *application1) showSnippet(w http.ResponseWriter, r *http.Request) { // "/snippet"
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	// Используем помощника render() для отображения шаблона.
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

// Создает новую заметку
func (app *application1) createSnippet(w http.ResponseWriter, r *http.Request) { // "/snippet/create"
	if r.Method != http.MethodPost {
		w.Header().Set("Allowed-method", http.MethodPost) // добавляет ключ:значение в карту HTTP
		app.clientError(w, http.StatusMethodNotAllowed)   // отправляет в ResponceWriter строку и в карту HTTP еод ошибки (перед write обязательно)
		return
	}

	// Создаем несколько переменных, содержащих тестовые данные. Мы удалим их позже.
	title := "История про улитку"
	content := "Улитка выползла из раковины,\nвытянула рожки,\nи опять подобрала их."
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
	//w.Write([]byte("Форма для создания новой заметки..."))
}
