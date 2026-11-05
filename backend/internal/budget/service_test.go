package budget_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jcyu0208/budget-app/backend/internal/budget"
)

// stubRepo is a minimal in-test AccountRepository (insertion-ordered).
type stubRepo struct{ saved []budget.Account }

func (r *stubRepo) Save(a budget.Account) error     { r.saved = append(r.saved, a); return nil }
func (r *stubRepo) List() ([]budget.Account, error) { return r.saved, nil }

// stubIDs is a deterministic IDGenerator.
type stubIDs struct{ n int }

func (g *stubIDs) NewID() string { g.n++; return fmt.Sprintf("acc-%d", g.n) }

func newService() *budget.Service {
	return budget.NewService(&stubRepo{}, &stubIDs{})
}

func TestCreateAccount_Valid(t *testing.T) {
	svc := newService()
	acc, err := svc.CreateAccount("  Checking  ", budget.KindCash, 100000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if acc.ID == "" {
		t.Error("expected an id to be assigned")
	}
	if acc.Name != "Checking" {
		t.Errorf("name should be trimmed, got %q", acc.Name)
	}
	if acc.Kind != budget.KindCash {
		t.Errorf("kind = %q, want cash", acc.Kind)
	}
	if acc.StartingBalance != 100000 {
		t.Errorf("starting balance = %d, want 100000", acc.StartingBalance)
	}
}

func TestCreateAccount_EmptyName(t *testing.T) {
	svc := newService()
	for _, name := range []string{"", "   "} {
		if _, err := svc.CreateAccount(name, budget.KindCash, 0); !errors.Is(err, budget.ErrEmptyName) {
			t.Errorf("name %q: expected ErrEmptyName, got %v", name, err)
		}
	}
}

func TestCreateAccount_InvalidKind(t *testing.T) {
	svc := newService()
	if _, err := svc.CreateAccount("Card", budget.Kind("credit"), 0); !errors.Is(err, budget.ErrInvalidKind) {
		t.Errorf("expected ErrInvalidKind, got %v", err)
	}
}

func TestCreateAccount_ZeroNegativeLarge(t *testing.T) {
	svc := newService()
	for _, bal := range []budget.Money{0, -20000, 9_000_000_000_000} {
		if _, err := svc.CreateAccount("A", budget.KindCash, bal); err != nil {
			t.Errorf("balance %d: unexpected error %v", bal, err)
		}
	}
}

func TestBudget_Empty(t *testing.T) {
	svc := newService()
	b, err := svc.Budget()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.ReadyToAssign != 0 {
		t.Errorf("ready to assign = %d, want 0", b.ReadyToAssign)
	}
	if len(b.Accounts) != 0 {
		t.Errorf("accounts = %d, want 0", len(b.Accounts))
	}
}

func TestBudget_SumsBalances(t *testing.T) {
	svc := newService()
	mustCreate(t, svc, "Checking", 100000)
	mustCreate(t, svc, "Savings", 50000)
	b, err := svc.Budget()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.ReadyToAssign != 150000 {
		t.Errorf("ready to assign = %d, want 150000", b.ReadyToAssign)
	}
	if len(b.Accounts) != 2 {
		t.Fatalf("accounts = %d, want 2", len(b.Accounts))
	}
	if b.Accounts[0].Balance != 100000 {
		t.Errorf("first balance = %d, want 100000", b.Accounts[0].Balance)
	}
}

func TestBudget_NegativeRTA(t *testing.T) {
	svc := newService()
	mustCreate(t, svc, "Overdrawn", -20000)
	b, _ := svc.Budget()
	if b.ReadyToAssign != -20000 {
		t.Errorf("ready to assign = %d, want -20000", b.ReadyToAssign)
	}
}

func mustCreate(t *testing.T, svc *budget.Service, name string, bal budget.Money) {
	t.Helper()
	if _, err := svc.CreateAccount(name, budget.KindCash, bal); err != nil {
		t.Fatalf("create %q: %v", name, err)
	}
}
