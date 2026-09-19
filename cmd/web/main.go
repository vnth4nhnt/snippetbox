package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"snippetbox.example.com/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog	*log.Logger
	infoLog		*log.Logger
	snippets 	*models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	// command-line flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	// define a new command-line flag for the MySQL DSN string
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// to keep the main() function tidy, put the code for creating a connection pool into the separate openDB() function below. 
	// pass openDB() the DSN from the command-line flag
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// also  defer  a call to db.Close(), so that the connection pool is closed before the main() function exits
	defer db.Close()

	// Initialize a new template cache...
    templateCache, err := newTemplateCache()
    if err != nil {
        errorLog.Fatal(err)
    }

	app := &application{
		errorLog: 	errorLog,
		infoLog:  	infoLog,
		snippets:	&models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		// call the new app.routes() method to get the servemux containing our routes
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// the openDB() function wraps sql.Open() and returns a sql.DB connection pool fo a given DSN
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return  nil, err
	}
	return db, nil
}
