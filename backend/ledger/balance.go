package ledger

import (
	"database/sql"
	"errors"
)

// applyBalanceDelta applies a single obligation:
// fromUserID owes toUserID amount.
//
// Ledger invariants enforced:
// - No self-debt
// - No negative balances
// - No bidirectional balances
// - Net obligations only
func applyBalanceDelta(
	tx *sql.Tx,
	fromUserID string,
	toUserID string,
	amount float64,
) error {

	// Guard conditions
	if fromUserID == toUserID {
		return nil
	}
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	var existing float64

	// 1. Check for reverse balance (to -> from)
	err := tx.QueryRow(`
		SELECT amount
		FROM balances
		WHERE from_user_id = $1 AND to_user_id = $2
	`, toUserID, fromUserID).Scan(&existing)

	if err == nil {
		// Reverse balance exists → net it
		switch {
		case existing > amount:
			// Reduce reverse balance
			_, err = tx.Exec(`
				UPDATE balances
				SET amount = amount - $1
				WHERE from_user_id = $2 AND to_user_id = $3
			`, amount, toUserID, fromUserID)
			return err

		case existing < amount:
			// Remove reverse balance
			_, err = tx.Exec(`
				DELETE FROM balances
				WHERE from_user_id = $1 AND to_user_id = $2
			`, toUserID, fromUserID)
			if err != nil {
				return err
			}

			// Insert remaining forward balance
			_, err = tx.Exec(`
				INSERT INTO balances (from_user_id, to_user_id, amount)
				VALUES ($1, $2, $3)
			`, fromUserID, toUserID, amount-existing)
			return err

		default:
			// existing == amount → cancel out
			_, err = tx.Exec(`
				DELETE FROM balances
				WHERE from_user_id = $1 AND to_user_id = $2
			`, toUserID, fromUserID)
			return err
		}
	}

	// If no reverse balance exists, add or increment forward balance
	if err != sql.ErrNoRows {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO balances (from_user_id, to_user_id, amount)
		VALUES ($1, $2, $3)
		ON CONFLICT (from_user_id, to_user_id)
		DO UPDATE SET amount = balances.amount + EXCLUDED.amount
	`, fromUserID, toUserID, amount)

	return err
}
