package app

import (
	"context"
	_ "fizzbuzz/app/docs"
	"fizzbuzz/fizzbuzz"
	"fizzbuzz/pkg/apiresponse"
	"fizzbuzz/pkg/migrations"
	"fizzbuzz/pkg/storage"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

const (
	EnvLstnAddrName    = "FIZZBUZZ_LSTN_ADDR"
	EnvLstnAddrDefault = ":8090"
	EnvDBFileName      = "FIZZBUZZ_DB_FILE"
	EnvDBFileDefault   = "fizzbuzz.db"
)

// getEnvValOrDefault gets the environment variable with the given name.
// If the variable is not set, it returns the default value.
func getEnvValOrDefault(name, dval string) string {
	if val, ok := os.LookupEnv(name); ok {
		return val
	}
	return dval
}

// App represents the application.
type App struct {
	listenAddr string
	dbFile     string
}

// New creates a new application instance.
func New() App {
	return App{
		listenAddr: getEnvValOrDefault(EnvLstnAddrName, EnvLstnAddrDefault),
		dbFile:     getEnvValOrDefault(EnvDBFileName, EnvDBFileDefault),
	}
}

// Run runs the application.
// @title FizzBuzz Swagger UI
// @version 1.0
// @description FizzBuzz
// @contact.name API Support
// @termsOfService demo.com
// @contact.url http://demo.com/support
// @contact.email support@swagger.io
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @BasePath /
// @Schemes http https
// @query.collection.format multi
// @in header
// @name Authorization
func (a App) Run() {
	ctx := context.Background()

	store, err := storage.NewSQLiteStore(a.dbFile)
	if err != nil {
		slog.Error("Error while connecting to database", "dbfile", a.dbFile, "error", err)
		return
	}
	defer store.Close()

	// run migration
	if _, err := store.Exec(ctx, migrations.InitSQL); err != nil {
		slog.Error("Error while running migration", "error", err)
		return
	}

	mux := http.NewServeMux()
	slog.Info("Server up and listening to", "addr", a.listenAddr)

	mux.Handle("GET /health", LoggerMiddleware(http.HandlerFunc(Health)))
	mux.Handle("GET /", LoggerMiddleware(RecoverMiddleware(fizzbuzz.Handler(ctx, store))))
	mux.Handle("GET /stat", LoggerMiddleware(RecoverMiddleware(fizzbuzz.StatHandler(ctx, store))))
	mux.Handle("GET /swagger/", http.StripPrefix("/swagger/", httpSwagger.Handler(httpSwagger.URL("doc.json"))))

	srv := http.Server{
		Addr:    a.listenAddr,
		Handler: mux,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Error while listening and serving", "addr", srv.Addr, "error", err)
		}
	}()

	// Wait for signal for a graceful shutdown
	<-quit
	slog.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	slog.Info("Server exiting")
}

func Health(w http.ResponseWriter, r *http.Request) {
	apiresponse.New(w, http.StatusOK, "OK", nil)
}
