package main

import (
	"alex/pkg/models/mysql"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql" // new import
	"github.com/golangcollege/sessions"
)

// Add a snippets filed to the application struct.  This will allow us to
// make the SnippetModel object available to our handlers
type application struct {
	users         *mysql.UserModel
	session       *sessions.Session
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}
type contextKey string

var contextKeyUser = contextKey("user")

func main() {

	// Создаем новый флаг командной строки, значение по умолчанию: ":4000".
	// Добавляем небольшую справку, объясняющая, что содержит данный флаг.
	// Значение флага будет сохранено в переменной addr.
	addr := flag.String("addr", ":4000", "http adress")
	dsn := flag.String("dsn", "web:pass@/bookdb?parseTime=true", "Mysql database")
	// Define a new command-line flag for the session secret(a random key which
	//will be used to encry[t and authenticate session cookies).It should be
	// bytes long
	secret := flag.String("secret", "sGWh6Ndh+pPbnzHbs*+9Pk8qTzbpa@ge", "Secret")
	flag.Parse()
	// Мы вызываем функцию flag.Parse() для извлечения флага из командной строки.
	// Она считывает значение флага из командной строки и присваивает его содержимое
	// переменной. Вам нужно вызвать ее *до* использования переменной addr
	// иначе она всегда будет содержать значение по умолчанию ":4000".
	// Если есть ошибки во время извлечения данных - приложение будет остановлено.

	// Значение, возвращаемое функцией flag.String(), является указателем на значение
	// из флага, а не самим значением. Нам нужно убрать ссылку на указатель
	// то есть перед использованием добавьте к нему префикс *. Обратите внимание, что мы используем
	// функцию log.Printf() для записи логов в журнал работы нашего приложения.
	//
	// to keep the main() function tidy i've put the code for creating a connect
	// pool into the separate opendb()function below. We pass opendb() the dsn from the command line flag.
	db, err := opendb(*dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	/*Initialize a mysql.SnippetModel instance and add it to the application
	depande*/
	//Initialize a new template cache ...
	templatecahe, err := newTemplateCache("./ui/html/")
	if err != nil {
		log.Fatal(err)
		return
	}
	//Use the sessions.New() function to initialize a new session  manager
	//passing in the secret key as the parametr. Then we configure it so
	//sessions always expires after 12 second

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true //Set the Secure flag on our session cookies
	//Блокирует отправку сессионных Cookie вайлов браузером пользователя для всех
	//межсайтовых использований
	session.SameSite = http.SameSiteStrictMode
	app := &application{
		users:         &mysql.UserModel{DB: db},
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templatecahe,
	}
	//Initialize a tls.Config struct to hold the non-default TLS settings we
	//the server to use
	// tlsConfig := tls.Config{
	// 	PreferServerCipherSuites: true,
	// 	CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	// }
	//Set the server`s TLSConfig field to use the tlsConfig  variable we just
	//created
	s := http.Server{
		Addr:    *addr,
		Handler: app.Routes(),
		// TLSConfig:    &tlsConfig,
		IdleTimeout:  time.Minute, // что всё keep-alive-connections будут автоматические закрыты после 1 минута без действий
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	//Use the ListenAndServeTLS() method to start the HTTPS server. We
	//pass in the paths to the TLS certificate and corresponding private key a the
	//two parametrssss
	log.Printf("Запуск сервера на %s", *addr)
	s.ListenAndServe()
}
func opendb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn) // sql.Open() it does initialize the pool for future use
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil { // method to create aconnection and check for any errors
		return nil, err
	}
	return db, err
}
