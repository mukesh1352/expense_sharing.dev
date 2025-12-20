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
	UserID     string  
	Amount     float64
	Percentage float64 
}

// ExpenseInput represents the input required to create an expense.
type ExpenseInput struct {
	ExpenseID    string   
	GroupID      string    
	PaidBy       string     
	TotalAmount  float64    
	SplitType    SplitType  
	Participants []string  
	Splits       []SplitInput 
	Description  string     
}
