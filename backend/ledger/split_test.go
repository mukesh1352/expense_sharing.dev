package ledger

import "testing"

func TestCalculateExactSplit_Success(t *testing.T) {
	input := ExpenseInput{
		TotalAmount: 300,
		SplitType:   SplitExact,
		Participants: []string{"u1", "u2"},
		Splits: []SplitInput{
			{UserID: "u1", Amount: 100},
			{UserID: "u2", Amount: 200},
		},
	}

	shares, err := calculateShares(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if shares["u1"] != 100 || shares["u2"] != 200 {
		t.Errorf("incorrect split result: %v", shares)
	}
}

func TestCalculateExactSplit_InvalidSum(t *testing.T) {
	input := ExpenseInput{
		TotalAmount: 300,
		SplitType:   SplitExact,
		Participants: []string{"u1", "u2"},
		Splits: []SplitInput{
			{UserID: "u1", Amount: 100},
			{UserID: "u2", Amount: 150},
		},
	}

	_, err := calculateShares(input)
	if err == nil {
		t.Errorf("expected error for invalid exact split sum")
	}
}

func TestCalculatePercentageSplit_Success(t *testing.T) {
	input := ExpenseInput{
		TotalAmount: 200,
		SplitType:   SplitPercentage,
		Participants: []string{"u1", "u2"},
		Splits: []SplitInput{
			{UserID: "u1", Percentage: 50},
			{UserID: "u2", Percentage: 50},
		},
	}

	shares, err := calculateShares(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if shares["u1"] != 100 || shares["u2"] != 100 {
		t.Errorf("incorrect percentage split result: %v", shares)
	}
}

func TestCalculatePercentageSplit_InvalidTotal(t *testing.T) {
	input := ExpenseInput{
		TotalAmount: 200,
		SplitType:   SplitPercentage,
		Participants: []string{"u1", "u2"},
		Splits: []SplitInput{
			{UserID: "u1", Percentage: 60},
			{UserID: "u2", Percentage: 30},
		},
	}

	_, err := calculateShares(input)
	if err == nil {
		t.Errorf("expected error when percentage != 100")
	}
}
