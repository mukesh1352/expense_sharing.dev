package ledger

import (
	"context"
	"database/sql"
	"errors"
)

func (l *Ledger) CreateExpense(ctx context.Context, input ExpenseInput) error {
	return l.withTx(func(tx *sql.Tx) error {

		if input.TotalAmount <= 0 {
			return errors.New("total amount must be greater than 0")
		}
		if input.PaidBy == "" {
			return errors.New("paidBy must be provided")
		}
		if len(input.Participants) == 0 {
			return errors.New("at least one participant is required")
		}

		// insert expense
		_, err := tx.Exec(
			`INSERT INTO expenses (id, group_id, paid_by, amount, split_type, description)
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			input.ExpenseID,
			input.GroupID,
			input.PaidBy,
			input.TotalAmount,
			input.SplitType,
			input.Description,
		)
		if err != nil {
			return err
		}

		// calculate shares
		shares, err := calculateShares(input)
		if err != nil {
			return err
		}

		// insert expense_splits
		for userID, amount := range shares {
			_, err := tx.Exec(
				`INSERT INTO expense_splits (expense_id, user_id, amount)
				 VALUES ($1, $2, $3)`,
				input.ExpenseID,
				userID,
				amount,
			)
			if err != nil {
				return err
			}
		}

		// update balances
		for userID, amount := range shares {
			if userID == input.PaidBy {
				continue
			}
			err := applyBalanceDelta(tx, userID, input.PaidBy, amount)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
