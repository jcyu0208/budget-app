package memory

import (
	"strconv"
	"sync/atomic"
)

// IDGenerator issues sequential IDs ("1", "2", …). Deterministic for tests;
// swap for UUIDs later without touching the core.
type IDGenerator struct{ n atomic.Int64 }

// NewIDGenerator returns a generator starting at 1.
func NewIDGenerator() *IDGenerator { return &IDGenerator{} }

// NewID returns the next sequential ID.
func (g *IDGenerator) NewID() string {
	return strconv.FormatInt(g.n.Add(1), 10)
}
