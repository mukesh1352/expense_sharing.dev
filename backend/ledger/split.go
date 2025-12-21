package ledger

import (
	"errors"
	"math"
)

const epsilon = 0.01 

func calculateEqualSplit(
	input ExpenseInput,
	shares map[string]float64,
) (map[string]float64, error) {

	n := len(input.Participants)
	if n == 0 {
		return nil, errors.New("no participants provided")
	}

	share := input.TotalAmount / float64(n)
	for _, userID := range input.Participants {
		shares[userID] = share
	}

	return shares, nil
}

func calculateExactSplit(
	input ExpenseInput,
	shares map[string]float64,
) (map[string]float64, error) {

	if len(input.Splits) == 0 {
		return nil, errors.New("exact split requires split details")
	}

	var total float64
	participants := make(map[string]bool)

	for _, userID := range input.Participants {
		participants[userID] = true
	}

	for _, split := range input.Splits {
		if split.Amount <= 0 {
			return nil, errors.New("split amount must be positive")
		}
		if !participants[split.UserID] {
			return nil, errors.New("split user not in participants list")
		}

		shares[split.UserID] = split.Amount
		total += split.Amount
	}

	if math.Abs(total-input.TotalAmount) > epsilon {
		return nil, errors.New("sum of exact splits must equal total amount")
	}

	return shares, nil
}

func calculatePercentageSplit(
	input ExpenseInput,
	shares map[string]float64,
) (map[string]float64, error) {

	if len(input.Splits) == 0 {
		return nil, errors.New("percentage split requires split details")
	}

	var totalPercentage float64
	participants := make(map[string]bool)

	for _, userID := range input.Participants {
		participants[userID] = true
	}

	for _, split := range input.Splits {
		if split.Percentage <= 0 {
			return nil, errors.New("percentage must be positive")
		}
		if !participants[split.UserID] {
			return nil, errors.New("split user not in participants list")
		}

		totalPercentage += split.Percentage
	}

	if math.Abs(totalPercentage-100) > epsilon {
		return nil, errors.New("sum of percentages must be 100")
	}

	for _, split := range input.Splits {
		shares[split.UserID] =
			input.TotalAmount * (split.Percentage / 100)
	}

	return shares, nil
}

func calculateShares(input ExpenseInput) (map[string]float64, error) {
	shares := make(map[string]float64)

	switch input.SplitType {
	case SplitEqual:
		return calculateEqualSplit(input, shares)

	case SplitExact:
		return calculateExactSplit(input, shares)

	case SplitPercentage:
		return calculatePercentageSplit(input, shares)

	default:
		return nil, errors.New("invalid split type")
	}
}
