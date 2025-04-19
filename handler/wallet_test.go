package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/Aswadhpv/wallet-service/service"
)

type mockSvc struct{}

func (m *mockSvc) Deposit(_, _ string, _ int64) error   { return nil }
func (m *mockSvc) Withdraw(_, _ string, _ int64) error  { return nil }
func (m *mockSvc) GetBalance(_, _ string) (int64, error) { return 500, nil }

func TestCreateOperation(t *testing.T) {
	m := NewWalletHandler(&mockSvc{})
	body := map[string]interface{}{
		"walletId":      "abc",
		"operationType": "DEPOSIT",
		"amount":        100,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/wallet", bytes.NewReader(b))
	w := httptest.NewRecorder()

	m.CreateOperation(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestGetBalance(t *testing.T) {
	m := NewWalletHandler(&mockSvc{})
	req := httptest.NewRequest("GET", "/api/v1/wallets/abc", nil)
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/wallets/{id}", m.GetBalance).Methods("GET")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]int64
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, int64(500), resp["balance"])
}
