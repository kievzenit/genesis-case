package database

import "database/sql"

type TransactionManger struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) *TransactionManger {
	return &TransactionManger{db}
}

func (txManager *TransactionManger) ExecuteTx(fn func(*sql.Tx) error) error {
	tx, err := txManager.db.Begin()
    if err != nil {
        return err
    }
    
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        }
    }()
    
    if err := fn(tx); err != nil {
        tx.Rollback()
        return err
    }
    
    return tx.Commit()
}
