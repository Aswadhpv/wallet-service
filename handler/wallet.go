package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/Aswadhpv/wallet-service/service"
)

// WalletHandler holds your service implementation
type WalletHandler struct {
	svc service.WalletService
}

// NewWalletHandler constructs the handler
func NewWalletHandler(s service.WalletService) *WalletHandler {
	return &WalletHandler{svc: s}
}

// operationRequest represents the parameters for deposit or withdrawal.
// swagger:model operationRequest
type operationRequest struct {
	// WalletID is the unique identifier of the wallet
	// example: 11111111-1111-1111-1111-111111111111
	WalletID string `json:"walletId" form:"walletId"`
	// OperationType defines whether the operation is DEPOSIT or WITHDRAW
	// enum: DEPOSIT,WITHDRAW
	// example: DEPOSIT
	OperationType string `json:"operationType" form:"operationType"`
	// Amount to deposit or withdraw (must be > 0)
	// example: 100
	Amount int64 `json:"amount" form:"amount"`
}

// CreateOperation godoc
// @Summary      Deposit or withdraw money
// @Description  Performs a deposit or withdrawal on a wallet
// @Tags         wallets
// @Accept       application/x-www-form-urlencoded,application/json
// @Produce      json
// @Param        walletId       formData  string  true  "Wallet ID"
// @Param        operationType  formData  string  true  "Operation Type"   Enums(DEPOSIT,WITHDRAW)
// @Param        amount         formData  int     true  "Amount (must be > 0)"
// @Success      204  "no content"
// @Failure      400  {string}  string  "invalid request or amount must be positive"
// @Failure      409  {string}  string  "insufficient funds"
// @Failure      500  {string}  string  "server error"
// @Router       /wallet [post]
func (h *WalletHandler) CreateOperation(w http.ResponseWriter, r *http.Request) {
	var req operationRequest

	ct := r.Header.Get("Content-Type")
	if ct == "" || strings.HasPrefix(ct, "application/json") {
		// JSON payload
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
	} else {
		// form-data (Swagger UI)
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		req.WalletID = r.FormValue("walletId")
		req.OperationType = r.FormValue("operationType")
		amtStr := r.FormValue("amount")
		amt, err := strconv.ParseInt(amtStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid amount", http.StatusBadRequest)
			return
		}
		req.Amount = amt
	}

	// Validate amount
	if req.Amount <= 0 {
		http.Error(w, "amount must be positive", http.StatusBadRequest)
		return
	}

	// Normalize operation type and execute
	var err error
	switch strings.ToUpper(req.OperationType) {
	case "DEPOSIT":
		err = h.svc.Deposit(r.Context(), req.WalletID, req.Amount)
	case "WITHDRAW":
		err = h.svc.Withdraw(r.Context(), req.WalletID, req.Amount)
	default:
		http.Error(w, "invalid operation type", http.StatusBadRequest)
		return
	}

	// Handle service errors
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
