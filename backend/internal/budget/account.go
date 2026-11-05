// Package budget holds the core budgeting rules. It imports no HTTP and no SQL;
// storage and transport are injected through the ports defined here.
package budget

// Money is an amount in whole cents. All budgeting math is integer math.
type Money int64

// Kind is the behavioural class of an account. Only cash is supported in the
// MVP; credit and loan accounts are on the roadmap (they move money differently).
type Kind string

const (
	// KindCash covers checking, savings, cash, and debit — accounts whose
	// spending immediately reduces real money.
	KindCash Kind = "cash"
)

// Account is the identity of a place money sits. It stores the starting balance
// (an input set at creation); the current balance is always derived, never stored.
type Account struct {
	ID              string
	Name            string
	Kind            Kind
	StartingBalance Money
}
