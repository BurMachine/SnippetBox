package main

import (
	"fmt"
	"html/template"
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
	files := []string{
		"../../ui/html/home.page.tmpl",
		"../../ui/html/base.layout.tmpl",
		"../../ui/html/footer.partial.tmpl",
	}
	pageTemp, err := template.ParseFiles(files...) // любой путь. Тут читается файл шаблона
	if err != nil {
		app.errorlog.Println(err.Error())
		app.serverError(w, err)
		return
	}
	err = pageTemp.Execute(w, nil) // Записываем содержимое шаблона в тело HTTP ответа, nil для отправки динамических данных в шаблон
	if err != nil {
		app.errorlog.Println(err.Error())
		app.serverError(w, err)
	}
}

// Отображает определенную заметку
func (app *application1) showSnippet(w http.ResponseWriter, r *http.Request) { // "/snippet"
	id, err := strconv.Atoi(r.URL.Query().Get("id")) // Считывание значения id. Затем проверка
	if err != nil || id < 0 {
		app.notFound(w)
		return
	}
	fmt.Fprintf(w, "Отображение заметки с ID %d...", id)
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
