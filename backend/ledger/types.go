package ledger

// SplitType defines how an expense is divided among participants.
type SplitType string

const (
	SplitEqual      SplitType = "EQUAL"
	SplitExact      SplitType = "EXACT"
	SplitPercentage SplitType = "PERCENT"
)

// SplitInput represents how much a single participant owes.
type SplitInput struct {
	UserID     string  `json:"user_id"`
	Amount     float64 `json:"amount,omitempty"`
	Percentage float64 `json:"percentage,omitempty"`
}


// ExpenseInput represents the input required to create an expense.
type ExpenseInput struct {
	ExpenseID    string     `json:"expense_id"`
	GroupID      string     `json:"group_id"`
	PaidBy       string     `json:"paid_by"`
	TotalAmount  float64    `json:"total_amount"`
	SplitType    SplitType  `json:"split_type"`
	Participants []string   `json:"participants"`
	Splits       []SplitInput `json:"splits,omitempty"`
	Description  string     `json:"description,omitempty"`
}

