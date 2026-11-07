package budget

import (
	"errors"
	"strings"
)

// Validation errors returned by the service.
var (
	ErrEmptyName   = errors.New("account name is required")
	ErrInvalidKind = errors.New("account kind must be cash")
)

// Service holds the budgeting rules and depends only on ports.
type Service struct {
	accounts AccountRepository
	ids      IDGenerator
}

// NewService wires the core with its adapters.
func NewService(accounts AccountRepository, ids IDGenerator) *Service {
	return &Service{accounts: accounts, ids: ids}
}

// CreateAccount validates and stores a new account. The starting balance is an
// input and may be zero, negative, or very large. A zero starting balance is
// just an attribute of 0 — no transaction is created.
func (s *Service) CreateAccount(name string, kind Kind, startingBalance Money) (Account, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Account{}, ErrEmptyName
	}
	if kind != KindCash {
		return Account{}, ErrInvalidKind
	}
	acc := Account{
		ID:              s.ids.NewID(),
		Name:            name,
		Kind:            kind,
		StartingBalance: startingBalance,
	}
	if err := s.accounts.Save(acc); err != nil {
		return Account{}, err
	}
	return acc, nil
}

// AccountView is an account as seen on the read side, with its derived balance.
type AccountView struct {
	ID      string
	Name    string
	Kind    Kind
	Balance Money
}

// BudgetView is the read model for a budget: Ready to Assign plus the accounts.
type BudgetView struct {
	ReadyToAssign Money
	Accounts      []AccountView
}

// Budget derives the read model. Every number here is computed from records,
// never read from a stored total.
func (s *Service) Budget() (BudgetView, error) {
	accs, err := s.accounts.List()
	if err != nil {
		return BudgetView{}, err
	}
	views := make([]AccountView, 0, len(accs))
	var totalBalances Money
	for _, a := range accs {
		bal := currentBalance(a)
		views = append(views, AccountView{ID: a.ID, Name: a.Name, Kind: a.Kind, Balance: bal})
		totalBalances += bal
	}
	// Ready to Assign = total account money − money sitting in category Available.
	return BudgetView{
		ReadyToAssign: totalBalances - totalCategoryAvailable(),
		Accounts:      views,
	}, nil
}

// currentBalance derives an account's balance from its records. No transactions
// exist yet, so it is the starting balance; when transactions arrive this will
// add their sum.
func currentBalance(a Account) Money {
	return a.StartingBalance
}

// totalCategoryAvailable is the sum of all category Available balances. No
// categories exist yet — this is the extension point for the assigning ticket.
func totalCategoryAvailable() Money {
	return 0
}
