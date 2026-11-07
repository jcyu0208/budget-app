package budget

// AccountRepository persists and lists accounts. Implementations (in-memory,
// SQL, …) are adapters; the core depends only on this port. List must return
// accounts in a stable order (insertion order).
type AccountRepository interface {
	Save(Account) error
	List() ([]Account, error)
}

// IDGenerator produces unique account IDs. Injected so IDs are deterministic
// in tests and the policy (sequence, UUID, …) lives in an adapter.
type IDGenerator interface {
	NewID() string
}
