// Command connectivity-gateway serves vehicle connectivity status to the
// mobile app and the owner portal.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aryanmathur2003/lucid-demo-connectivity-gateway/internal/gateway"
)

func main() {
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	srv := gateway.NewServer(gateway.NewMemoryStore(gateway.SeedVehicles()))
	log.Printf("connectivity-gateway listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, srv.Routes()))
}
