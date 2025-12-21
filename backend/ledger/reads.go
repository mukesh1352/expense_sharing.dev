package ledger

// BalanceView represents a readable balance entry
type BalanceView struct {
	FromUserID string
	ToUserID   string
	Amount     float64
}


func (l *Ledger) GetUserBalances(userID string) ([]BalanceView, error) {

	rows, err := l.db.Query(`
		SELECT from_user_id, to_user_id, amount
		FROM balances
		WHERE from_user_id = $1 OR to_user_id = $1
		ORDER BY from_user_id, to_user_id
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []BalanceView
	for rows.Next() {
		var b BalanceView
		if err := rows.Scan(&b.FromUserID, &b.ToUserID, &b.Amount); err != nil {
			return nil, err
		}
		balances = append(balances, b)
	}

	return balances, nil
}

func (l *Ledger) GetGroupBalances(groupID string) ([]BalanceView, error) {

	rows, err := l.db.Query(`
		SELECT b.from_user_id, b.to_user_id, b.amount
		FROM balances b
		JOIN group_members gm1 ON gm1.user_id = b.from_user_id
		JOIN group_members gm2 ON gm2.user_id = b.to_user_id
		WHERE gm1.group_id = $1 AND gm2.group_id = $1
		ORDER BY b.from_user_id, b.to_user_id
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []BalanceView
	for rows.Next() {
		var b BalanceView
		if err := rows.Scan(&b.FromUserID, &b.ToUserID, &b.Amount); err != nil {
			return nil, err
		}
		balances = append(balances, b)
	}

	return balances, nil
}
