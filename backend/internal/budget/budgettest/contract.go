// Package budgettest provides reusable contract suites for the budget ports,
// so every storage adapter is verified against the same behaviour.
package budgettest

import (
	"testing"

	"github.com/jcyu0208/budget-app/backend/internal/budget"
)

// RunAccountRepositoryContract verifies an AccountRepository implementation.
// newRepo must return a fresh, empty repository each call. No assertions are
// made on specific ID values — only on behaviour.
func RunAccountRepositoryContract(t *testing.T, newRepo func() budget.AccountRepository) {
	t.Run("empty list", func(t *testing.T) {
		got, err := newRepo().List()
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("want empty, got %d", len(got))
		}
	})

	t.Run("save then list", func(t *testing.T) {
		repo := newRepo()
		a := budget.Account{ID: "a1", Name: "Checking", Kind: budget.KindCash, StartingBalance: 100000}
		if err := repo.Save(a); err != nil {
			t.Fatalf("Save: %v", err)
		}
		got, err := repo.List()
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if len(got) != 1 || got[0] != a {
			t.Errorf("want [%+v], got %+v", a, got)
		}
	})

	t.Run("insertion order", func(t *testing.T) {
		repo := newRepo()
		ids := []string{"a1", "a2", "a3"}
		for _, id := range ids {
			if err := repo.Save(budget.Account{ID: id, Name: "N", Kind: budget.KindCash}); err != nil {
				t.Fatalf("Save %s: %v", id, err)
			}
		}
		got, _ := repo.List()
		if len(got) != len(ids) {
			t.Fatalf("want %d, got %d", len(ids), len(got))
		}
		for i, id := range ids {
			if got[i].ID != id {
				t.Errorf("order[%d] = %s, want %s", i, got[i].ID, id)
			}
		}
	})

	t.Run("duplicate names allowed", func(t *testing.T) {
		repo := newRepo()
		_ = repo.Save(budget.Account{ID: "a1", Name: "Wallet", Kind: budget.KindCash})
		_ = repo.Save(budget.Account{ID: "a2", Name: "Wallet", Kind: budget.KindCash})
		got, _ := repo.List()
		if len(got) != 2 {
			t.Errorf("want 2 accounts with the same name, got %d", len(got))
		}
	})
}
