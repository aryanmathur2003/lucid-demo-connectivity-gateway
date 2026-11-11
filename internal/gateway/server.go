package gateway

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var vinPattern = regexp.MustCompile(`^[A-HJ-NPR-Z0-9]{17}$`)

// Server holds the HTTP handlers for the gateway.
type Server struct {
	store Store
}

func NewServer(store Store) *Server { return &Server{store: store} }

// Routes returns the HTTP handler for the service.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/v1/vehicles/", s.handleConnectivity)
	return mux
}

type connectivityResponse struct {
	VIN        string    `json:"vin"`
	Link       LinkState `json:"link"`
	SignalBars int       `json:"signal_bars"`
	LastSeen   string    `json:"last_seen"`
	Firmware   string    `json:"firmware"`
	Locale     string    `json:"locale"`
	Message    string    `json:"message"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleConnectivity serves GET /v1/vehicles/{vin}/connectivity.
func (s *Server) handleConnectivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
		return
	}
	rest := strings.TrimPrefix(r.URL.Path, "/v1/vehicles/")
	parts := strings.Split(rest, "/")
	if len(parts) != 2 || parts[1] != "connectivity" {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "not found"})
		return
	}
	vin := strings.ToUpper(parts[0])
	if !vinPattern.MatchString(vin) {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid vin"})
		return
	}
	vehicle, ok := s.store.Get(vin)
	if !ok {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "vehicle not found"})
		return
	}
	locale := negotiateLocale(r.Header.Get("Accept-Language"))
	writeJSON(w, http.StatusOK, connectivityResponse{
		VIN:        vehicle.VIN,
		Link:       vehicle.Link,
		SignalBars: vehicle.SignalBars,
		LastSeen:   vehicle.LastSeen.UTC().Format(time.RFC3339),
		Firmware:   vehicle.FirmwareTrim,
		Locale:     locale.String(),
		Message:    localizedMessage(locale, vehicle.Link),
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
