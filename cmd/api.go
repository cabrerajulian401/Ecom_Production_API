package main

import (
	"log"
	"net/http"
	"time"

	repo "github.com/cabrerajulian401/ecom/internal/adapters/postgresql/sqlc"
	"github.com/cabrerajulian401/ecom/internal/products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

// mount the request onto the router
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID) // important for rate limiting. Key feature for limiting based off Request ID
	r.Use(middleware.RealIP)    // important for rate limiting and analytics and tracing
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // recover from crashes

	/* The above middleware auto inserts those data points into the context so they
	can be used in the later Clean Architecutre flow:
	user -> handler GET /products -> service getProducts -> repo SELECT * FROM Products ->throw error(RequestID)

	*/

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.

	// key funiton of Go context package
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good"))
	})

	productService := products.NewService(repo.New(app.db))
	productHandler := products.NewHandler(productService)
	r.Get("/products", productHandler.ListProducts)

	/* the handler register recives an instance of the interface*/
	return r
}

// run

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at addr %s", app.config.addr)

	return srv.ListenAndServe()
}

// this is a example of dependency injection. Such as what Yendo does with .Deps
type application struct {
	config config
	db     *pgx.Conn
	// logger
	// db driver
}

type config struct {
	addr string // the port address you are passing
	db   dbConfig
}

type dbConfig struct {
	dsn string // username for DB
}
