package service

import (
	"context"
	"errors"

	"github.com/Aswadhpv/wallet-service/repository"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

// WalletService defines the interface, still useful if you want to mock in tests
type WalletService interface {
	Deposit(ctx context.Context, walletID string, amount int64) error
	Withdraw(ctx context.Context, walletID string, amount int64) error
	GetBalance(ctx context.Context, walletID string) (int64, error)
}

// WalletServiceImpl is the concrete implementation
type WalletServiceImpl struct {
	repo repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) *WalletServiceImpl {
	return &WalletServiceImpl{repo: repo}
}

func (s *WalletServiceImpl) Deposit(ctx context.Context, walletID string, amount int64) error {
	return s.repo.ApplyOperation(ctx, walletID, "DEPOSIT", amount)
}

func (s *WalletServiceImpl) Withdraw(ctx context.Context, walletID string, amount int64) error {
	balance, err := s.repo.GetBalance(ctx, walletID)
	if err != nil {
		return err
	}
	if balance < amount {
		return ErrInsufficientFunds
	}
	return s.repo.ApplyOperation(ctx, walletID, "WITHDRAW", amount)
}

func (s *WalletServiceImpl) GetBalance(ctx context.Context, walletID string) (int64, error) {
	return s.repo.GetBalance(ctx, walletID)
}
