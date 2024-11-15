package app

import (
	"context"
	"fizzbuzz/fizzbuzz"
	"fizzbuzz/pkg/apiresponse"
	"fizzbuzz/pkg/storage"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
func (a App) Run() {
	ctx := context.Background()

	store, err := storage.NewSQLiteStore(a.dbFile)
	if err != nil {
		slog.Error("Error while connecting to database", "dbfile", a.dbFile, "error", err)
	}
	defer store.Close()

	mux := http.NewServeMux()
	slog.Info("Server up and listening to", "addr", a.listenAddr)

	mux.Handle("GET /health", LoggerMiddleware(http.HandlerFunc(Health)))
	mux.Handle("GET /", LoggerMiddleware(RecoverMiddleware(fizzbuzz.Handler(ctx, store))))
	mux.Handle("GET /stat", LoggerMiddleware(RecoverMiddleware(fizzbuzz.StatHandler(ctx, store))))

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
