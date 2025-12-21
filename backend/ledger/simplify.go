package ledger

import "database/sql"

type balanceEdge struct {
	user   string
	amount float64
}

//Helps in the removal of the userid as the intermediate option
// X -> userID -> Y  ==>  X -> Y
func SimplifyUserBalances(tx *sql.Tx, userID string) error {

	// Incoming: X -> userID
	incomingBalances := []balanceEdge{}
	rowsIn, err := tx.Query(`
		SELECT from_user_id, amount
		FROM balances
		WHERE to_user_id = $1
		ORDER BY from_user_id
	`, userID)
	if err != nil {
		return err
	}
	defer rowsIn.Close()

	for rowsIn.Next() {
		var b balanceEdge
		if err := rowsIn.Scan(&b.user, &b.amount); err != nil {
			return err
		}
		incomingBalances = append(incomingBalances, b)
	}

	// Outgoing: userID -> Y
	outgoingBalances := []balanceEdge{}
	rowsOut, err := tx.Query(`
		SELECT to_user_id, amount
		FROM balances
		WHERE from_user_id = $1
		ORDER BY to_user_id
	`, userID)
	if err != nil {
		return err
	}
	defer rowsOut.Close()

	for rowsOut.Next() {
		var b balanceEdge
		if err := rowsOut.Scan(&b.user, &b.amount); err != nil {
			return err
		}
		outgoingBalances = append(outgoingBalances, b)
	}

	// Simplify X -> userID -> Y
	for i := 0; i < len(incomingBalances); i++ {
		for j := 0; j < len(outgoingBalances); j++ {

			transfer := min(incomingBalances[i].amount, outgoingBalances[j].amount)
			if transfer <= 0 {
				continue
			}

			// Reduce X -> userID
			_, err := tx.Exec(`
				UPDATE balances
				SET amount = amount - $1
				WHERE from_user_id = $2 AND to_user_id = $3
			`, transfer, incomingBalances[i].user, userID)
			if err != nil {
				return err
			}

			// Reduce userID -> Y
			_, err = tx.Exec(`
				UPDATE balances
				SET amount = amount - $1
				WHERE from_user_id = $2 AND to_user_id = $3
			`, transfer, userID, outgoingBalances[j].user)
			if err != nil {
				return err
			}

			// Remove zero balances
			_, err = tx.Exec(`
				DELETE FROM balances
				WHERE amount = 0
			`)
			if err != nil {
				return err
			}

			// Add X -> Y
			if err := applyBalanceDelta(
				tx,
				incomingBalances[i].user,
				outgoingBalances[j].user,
				transfer,
			); err != nil {
				return err
			}

			incomingBalances[i].amount -= transfer
			outgoingBalances[j].amount -= transfer
		}
	}

	return nil
}


func SimplifyBalances(tx *sql.Tx) error {

	rows, err := tx.Query(`
		SELECT DISTINCT user_id FROM (
			SELECT from_user_id AS user_id FROM balances
			UNION
			SELECT to_user_id AS user_id FROM balances
		) u
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var userID string
	for rows.Next() {
		if err := rows.Scan(&userID); err != nil {
			return err
		}

		if err := SimplifyUserBalances(tx, userID); err != nil {
			return err
		}
	}

	return nil
}
