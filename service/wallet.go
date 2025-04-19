package service

import (
	"context"

	"github.com/yourusername/wallet-service/repository"
)

type WalletService interface {
	Deposit(ctx context.Context, id string, amount int64) error
	Withdraw(ctx context.Context, id string, amount int64) error
	GetBalance(ctx context.Context, id string) (int64, error)
}

type walletService struct {
	repo repository.WalletRepo
}

func NewWalletService(r repository.WalletRepo) WalletService {
	return &walletService{repo: r}
}

func (s *walletService) Deposit(ctx context.Context, id string, amount int64) error {
	return s.repo.ApplyOperation(ctx, id, "DEPOSIT", amount)
}

func (s *walletService) Withdraw(ctx context.Context, id string, amount int64) error {
	return s.repo.ApplyOperation(ctx, id, "WITHDRAW", amount)
}

func (s *walletService) GetBalance(ctx context.Context, id string) (int64, error) {
	return s.repo.GetBalance(ctx, id)
}
