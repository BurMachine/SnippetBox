package main

import (
	"database/sql"
	"flag"
	_ "github.com/lib/pq"
	"golangify.com/SnippetBox/pkg/models/postgresql"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type application1 struct {
	errorlog *log.Logger
	infolog  *log.Logger
	snippets *postsql.SnippetModel
}

/*
	http.NewServeMux - инициализация нового роутера
	mux.HandleFunc("/", home) регистрация home как обработчика url шаблона "/"
	http.ListenAndServe  - получает в качестве параметров TCP-адрес сети для прослушивания (localhost:4000)
	и созданный роутер
*/
func main() {
	addr := flag.String("addr", "127.0.0.1:4000", "Сетевоой адресс HTTP") // флаг командной строки

	// очень жестко(pg_hba.conf + нужно правильно прописать адрес
	dsn := flag.String("dsn", "postgresql://web:123@127.0.0.1:5433/userdb?sslmode=disable", "Название postSQL источника данных")

	flag.Parse()                                                                  // извлечение флага из командной строки(меняет по адресу addr)
	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)                  // создание логгера INFO в stdout
	errorlog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) // логгер ошибок ERROR

	db, err := openDB(*dsn) // инициализируем пул подключений к базе
	if err != nil {
		errorlog.Fatal(err)
	}
	defer db.Close()

	app := &application1{ // инициализация новой структуры, чтобы подтянуть методы
		errorlog: errorlog,
		infolog:  infolog,
		snippets: &postsql.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorlog,
		Handler:  app.routes(), // создает маршрутизатор и тп для декомпозиции
	}
	infolog.Printf("Запуск веб-сервера на http://%s", *addr)
	err = srv.ListenAndServe() // Запуск нового веб-сервера
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

// Функция openDB() обертывает sql.Open() и возвращает пул соединений sql.DB
// для заданной строки подключения (DSN).
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil { // проверка того что все настроено правильно
		return nil, err
	}
	return db, nil
}
