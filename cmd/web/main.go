package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"snippetbox/pkg/models/postgresql"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *postgresql.SnippetModel
	templateCache map[string]*template.Template
}

type Config struct {
	Addr	  string
	StaticDir string
	DBHost    string
	DBUser    string
	DBPort    string
	DBName    string
	Pass      string
}

func main() {
	var cfg  = new(Config)
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Directory with static files")
	flag.Parse()
	
	infoLog := log.New(os.Stdout, "\u001b[36m[INFO]\u001b[0m\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "\u001b[31m[ERROR]\u001b[0m\t", log.Ldate|log.Ltime|log.Lshortfile)
	
	if err := initConfig(); err != nil {
		errorLog.Fatalf("couldn't initializing configs: %s", err.Error())
	}
		
	cfg.Addr = viper.GetString("host")+":"+viper.GetString("port")
	cfg.DBHost = viper.GetString("db.host")
	cfg.DBPort = viper.GetString("db.port")
	cfg.DBUser = viper.GetString("db.user")
	cfg.DBName = viper.GetString("db.dbname")
	cfg.Pass = viper.GetString("db.password")

	db, err := openDB(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.Pass, cfg.DBName))
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &postgresql.SnippetModel{DB: db},
		templateCache: templateCache,
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

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
