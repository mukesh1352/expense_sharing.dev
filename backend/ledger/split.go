package ledger

import (
	"errors"
)

func calculateEqualSplit(input ExpenseInput, shares map[string]float64)(map[string]float64,error){
	n := len(input.Participants)
	if n == 0 {
    return nil, errors.New("no participants provided")
}

	share := input.TotalAmount/float64(n)
	for _,userID := range input.Participants {
		shares[userID] = share
	}
	return shares,nil
}



func calculateExactSplit(input ExpenseInput, shares map[string]float64) (map[string]float64, error) {
	if len(input.Splits) == 0 {
		return nil, errors.New("exact split requires split details")
	}

	var total float64

	for _, split := range input.Splits {
		if split.Amount <= 0 {
			return nil, errors.New("split amount must be positive")
		}
		shares[split.UserID] = split.Amount
		total += split.Amount
	}

	if total != input.TotalAmount {
		return nil, errors.New("sum of exact splits must equal total amount")
	}

	return shares, nil
}




func calculatePercentageSplit(input ExpenseInput, shares map[string]float64) (map[string]float64, error) {
	if len(input.Splits) == 0 {
		return nil, errors.New("percentage split requires split details")
	}

	var totalPercentage float64

	for _, split := range input.Splits {
		if split.Percentage <= 0 {
			return nil, errors.New("percentage must be positive")
		}
		totalPercentage += split.Percentage
	}

	if totalPercentage != 100 {
		return nil, errors.New("sum of percentages must be 100")
	}

	for _, split := range input.Splits {
		shares[split.UserID] = input.TotalAmount * (split.Percentage / 100)
	}

	return shares, nil
}

func calculateShares(input ExpenseInput)(map[string]float64,error){
	shares:=make(map[string]float64)
	switch input.SplitType{
	case SplitEqual:
		return calculateEqualSplit(input,shares)
	case SplitExact:
		return calculateExactSplit(input,shares)
	case SplitPercentage:
		return calculatePercentageSplit(input,shares)
	default:
		return  nil,errors.New("invalid split type..")
	}
}