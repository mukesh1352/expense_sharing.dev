package ledger


type BalanceView struct {
    FromUserID string  `json:"from_user_id"`
    ToUserID   string  `json:"to_user_id"`
    Amount     float64 `json:"amount"`
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

	balances := []BalanceView{}

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

	balances := []BalanceView{}
	for rows.Next() {
		var b BalanceView
		if err := rows.Scan(&b.FromUserID, &b.ToUserID, &b.Amount); err != nil {
			return nil, err
		}
		balances = append(balances, b)
	}

	return balances, nil
}

// UserView
type UserView struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}


type GroupView struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (l *Ledger) GetUsers() ([]UserView, error) {
	rows, err := l.db.Query(`
		SELECT id, name FROM users ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []UserView{}
	for rows.Next() {
		var u UserView
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (l *Ledger) GetGroups() ([]GroupView, error) {
	rows, err := l.db.Query(`
		SELECT id, name FROM groups ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := []GroupView{}
	for rows.Next() {
		var g GroupView
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

func (l *Ledger) GetGroupMembers(groupID string) ([]UserView, error) {
	rows, err := l.db.Query(`
		SELECT u.id, u.name
		FROM users u
		JOIN group_members gm ON gm.user_id = u.id
		WHERE gm.group_id = $1
		ORDER BY u.name
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []UserView{}
	for rows.Next() {
		var u UserView
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
