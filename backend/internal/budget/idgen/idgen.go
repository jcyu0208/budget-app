// Package idgen provides the production budget.IDGenerator. It generates
// UUIDv7 identifiers under the hood — the only place in the backend that
// depends on github.com/google/uuid. Callers depend on the role (ID
// generation), not the implementation.
package idgen

import (
	"github.com/google/uuid"

	"github.com/jcyu0208/budget-app/backend/internal/budget"
)

// Generator issues unique IDs. The current implementation is UUIDv7, which is
// time-ordered (good database index locality) and globally unique.
type Generator struct{}

var _ budget.IDGenerator = Generator{}

// New returns an ID generator.
func New() Generator { return Generator{} }

// NewID returns a new unique ID.
func (Generator) NewID() string {
	return uuid.Must(uuid.NewV7()).String()
}
