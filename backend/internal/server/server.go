// Package server wires up the budget-app HTTP API (an adapter over the core).
package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jcyu0208/budget-app/backend/internal/budget"
)

// NewRouter builds the HTTP handler for the backend API.
func NewRouter(svc *budget.Service) http.Handler {
	api := &api{svc: svc}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("POST /accounts", api.handleCreateAccount)
	mux.HandleFunc("GET /budget", api.handleBudget)
	return mux
}

// api holds the HTTP handlers' dependency on the core service.
type api struct {
	svc *budget.Service
}

// handleHealth reports that the service is up.
func handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type createAccountRequest struct {
	Name            string `json:"name"`
	Kind            string `json:"kind"`
	StartingBalance int64  `json:"starting_balance"`
}

type accountResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Kind    string `json:"kind"`
	Balance int64  `json:"balance"`
}

type budgetResponse struct {
	ReadyToAssign int64             `json:"ready_to_assign"`
	Accounts      []accountResponse `json:"accounts"`
}

func (a *api) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var req createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	acc, err := a.svc.CreateAccount(req.Name, budget.Kind(req.Kind), budget.Money(req.StartingBalance))
	if err != nil {
		if errors.Is(err, budget.ErrEmptyName) || errors.Is(err, budget.ErrInvalidKind) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "could not create account")
		return
	}
	// Balance is the derived current balance (= starting balance until there are transactions).
	writeJSON(w, http.StatusCreated, accountResponse{
		ID: acc.ID, Name: acc.Name, Kind: string(acc.Kind), Balance: int64(acc.StartingBalance),
	})
}

func (a *api) handleBudget(w http.ResponseWriter, _ *http.Request) {
	view, err := a.svc.Budget()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not read budget")
		return
	}
	resp := budgetResponse{
		ReadyToAssign: int64(view.ReadyToAssign),
		Accounts:      make([]accountResponse, 0, len(view.Accounts)),
	}
	for _, acc := range view.Accounts {
		resp.Accounts = append(resp.Accounts, accountResponse{
			ID: acc.ID, Name: acc.Name, Kind: string(acc.Kind), Balance: int64(acc.Balance),
		})
	}
	writeJSON(w, http.StatusOK, resp)
}

// writeError writes a JSON error body with the given status code.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// writeJSON writes data as a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
