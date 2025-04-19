package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/Aswadhpv/wallet-service/service"
)

// WalletHandler holds your service implementation
type WalletHandler struct {
	svc *service.WalletServiceImpl
}

// NewWalletHandler constructs the handler
func NewWalletHandler(s *service.WalletServiceImpl) *WalletHandler {
	return &WalletHandler{svc: s}
}

type operationRequest struct {
	WalletID      string `json:"walletId"`
	OperationType string `json:"operationType"`
	Amount        int64  `json:"amount"`
}

// CreateOperation godoc
// @Summary      Deposit or withdraw money
// @Description  Performs a deposit or withdrawal on a wallet
// @Tags         wallets
// @Accept       json
// @Produce      plain
// @Param        payload  body      operationRequest  true  "operation payload"
// @Success      204
// @Failure      400      {string}  string            "invalid request"
// @Failure      409      {string}  string            "insufficient funds"
// @Failure      500      {string}  string            "server error"
// @Router       /wallet [post]
func (h *WalletHandler) CreateOperation(w http.ResponseWriter, r *http.Request) {
	var req operationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var err error
	if req.OperationType == "DEPOSIT" {
		err = h.svc.Deposit(r.Context(), req.WalletID, req.Amount)
	} else {
		err = h.svc.Withdraw(r.Context(), req.WalletID, req.Amount)
	}

	if err != nil {
		if err == service.ErrInsufficientFunds {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetBalance godoc
// @Summary      Get wallet balance
// @Description  Returns the current balance for a given wallet ID
// @Tags         wallets
// @Produce      json
// @Param        id   path      string  true  "Wallet UUID"
// @Success      200  {object}  map[string]int64
// @Failure      500  {string}  string  "server error"
// @Router       /wallets/{id} [get]
func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	bal, err := h.svc.GetBalance(r.Context(), id)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]int64{"balance": bal})
}
