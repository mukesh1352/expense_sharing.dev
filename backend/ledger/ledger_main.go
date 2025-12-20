package ledger

import (
	"context"
	"database/sql"
)

type Ledger struct{
	db *sql.DB
}

//creating a ledger instance
func New(db *sql.DB) *Ledger {
	return &Ledger{db: db}
}

// the below function helps in the commiting or reverting of the transaction state

func (l *Ledger) withTx(fn func(tx *sql.Tx) error) error{
	tx,err := l.db.BeginTx(
		context.Background(),
		&sql.TxOptions{
			Isolation: sql.LevelSerializable,
		},
	)
	if err!=nil{
		return err
	}
	// Ensuring rollback in case of any error or issue uprising 
	defer tx.Rollback()

if err := fn(tx); err != nil {
		return err
	}
	// if passing all the cases then commit 
	return tx.Commit()
}