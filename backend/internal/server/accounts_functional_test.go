package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jcyu0208/budget-app/backend/internal/budget"
	"github.com/jcyu0208/budget-app/backend/internal/budget/memory"
	"github.com/jcyu0208/budget-app/backend/internal/server"
)

// testAPI drives the real HTTP entry point against the in-memory adapter.
type testAPI struct {
	t      *testing.T
	srv    *httptest.Server
	client *http.Client
}

func newTestAPI(t *testing.T) *testAPI {
	t.Helper()
	svc := budget.NewService(memory.NewAccountRepository(), memory.NewIDGenerator())
	srv := httptest.NewServer(server.NewRouter(svc))
	t.Cleanup(srv.Close)
	return &testAPI{t: t, srv: srv, client: srv.Client()}
}

type accountJSON struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Kind    string `json:"kind"`
	Balance int64  `json:"balance"`
}

type budgetJSON struct {
	ReadyToAssign int64         `json:"ready_to_assign"`
	Accounts      []accountJSON `json:"accounts"`
}

func (a *testAPI) createAccount(name, kind string, startingBalance int64) *http.Response {
	a.t.Helper()
	body, _ := json.Marshal(map[string]any{
		"name": name, "kind": kind, "starting_balance": startingBalance,
	})
	resp, err := a.client.Post(a.srv.URL+"/accounts", "application/json", bytes.NewReader(body))
	if err != nil {
		a.t.Fatalf("POST /accounts: %v", err)
	}
	return resp
}

func (a *testAPI) mustCreate(name, kind string, startingBalance int64) accountJSON {
	a.t.Helper()
	resp := a.createAccount(name, kind, startingBalance)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		a.t.Fatalf("create %q: status %d, want 201", name, resp.StatusCode)
	}
	var acc accountJSON
	if err := json.NewDecoder(resp.Body).Decode(&acc); err != nil {
		a.t.Fatalf("decode account: %v", err)
	}
	return acc
}

// budget reads the budget through the read side and checks the invariant.
func (a *testAPI) budget() budgetJSON {
	a.t.Helper()
	resp, err := a.client.Get(a.srv.URL + "/budget")
	if err != nil {
		a.t.Fatalf("GET /budget: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		a.t.Fatalf("GET /budget: status %d, want 200", resp.StatusCode)
	}
	var b budgetJSON
	if err := json.NewDecoder(resp.Body).Decode(&b); err != nil {
		a.t.Fatalf("decode budget: %v", err)
	}
	// Invariant: Σ(account balances) == Ready to Assign (category Available = 0 here).
	var sum int64
	for _, acc := range b.Accounts {
		sum += acc.Balance
	}
	if sum != b.ReadyToAssign {
		a.t.Fatalf("invariant broken: Σ balances=%d, ready_to_assign=%d", sum, b.ReadyToAssign)
	}
	return b
}

func TestEmptyBudget(t *testing.T) {
	api := newTestAPI(t)
	b := api.budget()
	if b.ReadyToAssign != 0 {
		t.Errorf("ready_to_assign = %d, want 0", b.ReadyToAssign)
	}
	if len(b.Accounts) != 0 {
		t.Errorf("accounts = %d, want 0", len(b.Accounts))
	}
}

func TestCreateAccountThenBudget(t *testing.T) {
	api := newTestAPI(t)

	acc := api.mustCreate("Checking", "cash", 100000)
	if acc.ID == "" {
		t.Error("expected an id")
	}
	if acc.Balance != 100000 {
		t.Errorf("balance = %d, want 100000", acc.Balance)
	}

	b := api.budget()
	if b.ReadyToAssign != 100000 {
		t.Errorf("ready_to_assign = %d, want 100000", b.ReadyToAssign)
	}
	if len(b.Accounts) != 1 {
		t.Fatalf("accounts = %d, want 1", len(b.Accounts))
	}

	api.mustCreate("Savings", "cash", 50000)
	if b := api.budget(); b.ReadyToAssign != 150000 {
		t.Errorf("ready_to_assign = %d, want 150000", b.ReadyToAssign)
	}
}

func TestZeroStartingBalance(t *testing.T) {
	api := newTestAPI(t)
	api.mustCreate("Empty", "cash", 0)
	b := api.budget()
	if b.ReadyToAssign != 0 {
		t.Errorf("ready_to_assign = %d, want 0", b.ReadyToAssign)
	}
	if len(b.Accounts) != 1 || b.Accounts[0].Balance != 0 {
		t.Errorf("want one account with balance 0, got %+v", b.Accounts)
	}
}

func TestNegativeStartingBalance(t *testing.T) {
	api := newTestAPI(t)
	api.mustCreate("Overdrawn", "cash", -20000)
	if b := api.budget(); b.ReadyToAssign != -20000 {
		t.Errorf("ready_to_assign = %d, want -20000", b.ReadyToAssign)
	}
}

func TestDuplicateNames(t *testing.T) {
	api := newTestAPI(t)
	a1 := api.mustCreate("Wallet", "cash", 1000)
	a2 := api.mustCreate("Wallet", "cash", 2000)
	if a1.ID == a2.ID {
		t.Error("duplicate-named accounts must have distinct ids")
	}
	if b := api.budget(); len(b.Accounts) != 2 {
		t.Errorf("accounts = %d, want 2", len(b.Accounts))
	}
}

func TestVeryLargeStartingBalance(t *testing.T) {
	api := newTestAPI(t)
	const big = int64(1_000_000_000_000) // $10 billion in cents
	api.mustCreate("Vault", "cash", big)
	if b := api.budget(); b.ReadyToAssign != big {
		t.Errorf("ready_to_assign = %d, want %d", b.ReadyToAssign, big)
	}
}

func TestInvalidKindRejected(t *testing.T) {
	api := newTestAPI(t)
	resp := api.createAccount("Card", "credit", 0)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", resp.StatusCode)
	}
}

func TestEmptyNameRejected(t *testing.T) {
	api := newTestAPI(t)
	resp := api.createAccount("   ", "cash", 0)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", resp.StatusCode)
	}
}
