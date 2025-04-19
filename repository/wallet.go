package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

type WalletRepo interface {
	GetBalance(ctx context.Context, id string) (int64, error)
	ApplyOperation(ctx context.Context, id, opType string, amount int64) error
}

type walletRepo struct {
	db *sqlx.DB
}

func NewWalletRepo(db *sqlx.DB) WalletRepo {
	return &walletRepo{db: db}
}

func (r *walletRepo) GetBalance(ctx context.Context, id string) (int64, error) {
	var bal int64
	err := r.db.GetContext(ctx, &bal, "SELECT balance FROM wallets WHERE id=$1", id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return bal, err
}

func (r *walletRepo) ApplyOperation(ctx context.Context, id, opType string, amount int64) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Lock the row
	var bal int64
	err = tx.GetContext(ctx, &bal, "SELECT balance FROM wallets WHERE id=$1 FOR UPDATE", id)
	if err == sql.ErrNoRows {
		if opType == "WITHDRAW" {
			return ErrInsufficientFunds
		}
		bal = 0
		_, err = tx.ExecContext(ctx,
			"INSERT INTO wallets (id, balance) VALUES ($1, $2)", id, amount)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		newBal := bal
		if opType == "DEPOSIT" {
			newBal += amount
		} else {
			newBal -= amount
			if newBal < 0 {
				return ErrInsufficientFunds
			}
		}
		_, err = tx.ExecContext(ctx,
			"UPDATE wallets SET balance=$1 WHERE id=$2", newBal, id)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
