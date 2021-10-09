package main

import (
	"flag"
	"html/template"
	"log"
	"snippetbox/pkg/postgre"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/jackc/pgx"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets      *postgre.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	config := pgx.ConnConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "snippetbox_db",
		User:     "postgres",
		Password: "5779",
	}

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: 10,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}

	conn, err := pgx.NewConnPool(poolConfig)

	if err != nil {
		errorLog.Fatal(err)
	}
	defer conn.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &postgre.SnippetModel{Conn: conn},
		templateCache: templateCache,
	}

	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}