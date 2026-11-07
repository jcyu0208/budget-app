// Command server starts the budget-app backend HTTP server.
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jcyu0208/budget-app/backend/internal/budget"
	"github.com/jcyu0208/budget-app/backend/internal/budget/idgen"
	"github.com/jcyu0208/budget-app/backend/internal/budget/memory"
	"github.com/jcyu0208/budget-app/backend/internal/server"
)

func main() {
	svc := budget.NewService(memory.NewAccountRepository(), idgen.New())

	addr := ":" + port()
	srv := &http.Server{
		Addr:              addr,
		Handler:           server.NewRouter(svc),
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("budget-app backend listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

// port returns the port to listen on, from $PORT or a default.
func port() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}
