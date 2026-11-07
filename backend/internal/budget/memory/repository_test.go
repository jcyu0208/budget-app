package memory_test

import (
	"testing"

	"github.com/jcyu0208/budget-app/backend/internal/budget"
	"github.com/jcyu0208/budget-app/backend/internal/budget/budgettest"
	"github.com/jcyu0208/budget-app/backend/internal/budget/memory"
)

func TestAccountRepository_Contract(t *testing.T) {
	budgettest.RunAccountRepositoryContract(t, func() budget.AccountRepository {
		return memory.NewAccountRepository()
	})
}

func TestIDGenerator_UniqueNonEmpty(t *testing.T) {
	g := memory.NewIDGenerator()
	seen := map[string]bool{}
	for i := 0; i < 100; i++ {
		id := g.NewID()
		if id == "" {
			t.Fatal("got empty id")
		}
		if seen[id] {
			t.Fatalf("duplicate id %q", id)
		}
		seen[id] = true
	}
}
