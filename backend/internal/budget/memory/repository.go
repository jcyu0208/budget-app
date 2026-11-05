// Package memory provides in-memory adapters for the budget ports. It is the
// MVP storage; a SQL adapter can replace it later without touching the core.
package memory

import (
	"sync"

	"github.com/jcyu0208/budget-app/backend/internal/budget"
)

// AccountRepository is an in-memory budget.AccountRepository, safe for
// concurrent use. It preserves insertion order.
type AccountRepository struct {
	mu    sync.Mutex
	order []string
	byID  map[string]budget.Account
}

// NewAccountRepository returns an empty repository.
func NewAccountRepository() *AccountRepository {
	return &AccountRepository{byID: make(map[string]budget.Account)}
}

// Save stores an account, keeping first-seen order for stable listing.
func (r *AccountRepository) Save(a budget.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.byID[a.ID]; !exists {
		r.order = append(r.order, a.ID)
	}
	r.byID[a.ID] = a
	return nil
}

// List returns all accounts in insertion order.
func (r *AccountRepository) List() ([]budget.Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]budget.Account, 0, len(r.order))
	for _, id := range r.order {
		out = append(out, r.byID[id])
	}
	return out, nil
}
