package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

type Config struct {
	Addr      string
	StaticDir string
	DBHost    string
	DBUser    string
	DBPort    string
	DBName    string
	Pass      string
}

func main() {
	cfg := new(Config)
	flag.StringVar(&cfg.Addr, "addr", ":80", "Server Address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Directory with static files")
	flag.StringVar(&cfg.DBHost, "host", "localhost", "The Hosy to connect to")
	flag.StringVar(&cfg.DBPort, "port", "5432", "The Port to bind to")
	flag.StringVar(&cfg.DBUser, "user", "web", "The user to sign in as")
	flag.StringVar(&cfg.Pass, "password", "1111", "The user's password")
	flag.StringVar(&cfg.DBName, "dbname", "SnippetBox", "The name of database to connect to")
	flag.Parse()

	infoLog := log.New(os.Stdout, "\u001b[36m[INFO]\u001b[0m\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "\u001b[31m[ERROR]\u001b[0m\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.Pass, cfg.DBName))
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     cfg.Addr,
		ErrorLog: errorLog,
		Handler:  app.routes(cfg),
	}

	infoLog.Printf("Server Start! at %s", cfg.Addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
