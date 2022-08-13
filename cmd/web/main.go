package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

/*
	http.NewServeMux - инициализация нового роутера
	mux.HandleFunc("/", home) регистрация home как обработчика url шаблона "/"
	http.ListenAndServe  - получает в качестве параметров TCP-адрес сети для прослушивания (localhost:4000)
	и созданный роутер
*/
func main() {
	addr := flag.String("addr", "127.0.0.1:4000", "Сетевоой адресс HTTP")         // флаг командной строки
	flag.Parse()                                                                  // извлечение флага из командной строки(меняет по адресу addr)
	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)                  // создание логгера INFO в stdout
	errorlog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) // логгер ошибок ERROR
	mux := http.NewServeMux()                                                     // новый роутер
	mux.HandleFunc("/", home)                                                     // регистрирует функцию home как обработчик для роутера mux ...
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// инициализация FileServer, который будет обрабатывать HTTP-запросы к статическим файлам из ./ui/static/
	fileServer := http.FileServer(neuterdFileSystem{http.Dir("./static")})
	// регистрация всех запросов начинающихся со "/static"
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorlog,
		Handler:  mux,
	}
	infolog.Printf("Запуск веб-сервера на http://%s", *addr)
	err := srv.ListenAndServe() // Запуск нового веб-сервера
	errorlog.Fatal(err)
}

/*
 Создаем настраиваемый тип, который включает в себя http.FileSystem
*/
type neuterdFileSystem struct {
	fs http.FileSystem
}

/*
	Open():
	Создаем метод для neuteredFileSystem, который вызывается каждый раз, когда http.FileServer получает запрос
	Возвращает файл, усли он не существует вернет ошибку os.ErrNotExist (преобразуется в 404)
	Также закрывает файл
*/
func (new_fs neuterdFileSystem) Open(path string) (http.File, error) {
	f, err := new_fs.fs.Open(path) // открываем вызываемый путь
	if err != nil {
		return nil, err
	}
	s, err := f.Stat() // os.File предоставляет доступ к информации о файле/пути os.FileInfo
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := new_fs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return f, nil
}
