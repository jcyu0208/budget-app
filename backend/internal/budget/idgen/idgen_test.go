package idgen_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/jcyu0208/budget-app/backend/internal/budget/idgen"
)

func TestNewID_ValidUniqueV7(t *testing.T) {
	g := idgen.New()
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := g.NewID()
		u, err := uuid.Parse(id)
		if err != nil {
			t.Fatalf("NewID returned a non-UUID %q: %v", id, err)
		}
		if u.Version() != 7 {
			t.Errorf("UUID version = %d, want 7", u.Version())
		}
		if seen[id] {
			t.Fatalf("duplicate id %q", id)
		}
		seen[id] = true
	}
}
