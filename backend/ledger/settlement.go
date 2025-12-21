package ledger

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

// SettleBalance records a real-world payment and updates the ledger.
func (l *Ledger) SettleBalance(
	ctx context.Context,
	fromUserID string,
	toUserID string,
	amount float64,
) error {

	return l.withTx(func(tx *sql.Tx) error {

		// 1️⃣ Validate input
		if fromUserID == "" || toUserID == "" {
			return errors.New("user IDs must be provided")
		}
		if fromUserID == toUserID {
			return errors.New("cannot settle balance with self")
		}
		if amount <= 0 {
			return errors.New("settlement amount must be positive")
		}

		// 2️⃣ Fetch existing balance
		var existing float64
		err := tx.QueryRow(`
			SELECT amount
			FROM balances
			WHERE from_user_id = $1 AND to_user_id = $2
		`, fromUserID, toUserID).Scan(&existing)

		if err == sql.ErrNoRows {
			return errors.New("no outstanding balance to settle")
		}
		if err != nil {
			return err
		}

		if amount > existing {
			return errors.New("settlement amount exceeds outstanding balance")
		}

		// 3️⃣ Reduce or remove balance
		if amount == existing {
			_, err = tx.Exec(`
				DELETE FROM balances
				WHERE from_user_id = $1 AND to_user_id = $2
			`, fromUserID, toUserID)
			if err != nil {
				return err
			}
		} else {
			_, err = tx.Exec(`
				UPDATE balances
				SET amount = amount - $1
				WHERE from_user_id = $2 AND to_user_id = $3
			`, amount, fromUserID, toUserID)
			if err != nil {
				return err
			}
		}

		// 4️⃣ Insert settlement record (immutable history)
		_, err = tx.Exec(`
			INSERT INTO settlements (id, from_user_id, to_user_id, amount)
			VALUES ($1, $2, $3, $4)
		`,
			uuid.NewString(),
			fromUserID,
			toUserID,
			amount,
		)
		if err != nil {
			return err
		}

		return nil
	})
}
