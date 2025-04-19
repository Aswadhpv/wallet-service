# Wallet Service

## Overview

Simple Go service to deposit/withdraw and get wallet balance via REST.

### Prerequisites

- Go ≥1.19
- Docker & Docker Compose
- Visual Studio (or VS Code) with Go extension

## Setup

1. Clone:
   ```bash
   git clone https://github.com/Aswadhpv/wallet-service.git
   cd wallet-service

2. Configure config.env.

3. In Visual Studio:

  - Open folder.

  - Install Go plugin.

  - Run “Go: Install/Update Tools” from the Command Palette.

4. To run with Docker:
  ```bash
  docker-compose up --build
  ```
5. Endpoints:
  - POST /api/v1/wallet
  - GET /api/v1/wallets/{id}

## Tests
  - To run tests:
  ```bash
  go test ./...
  ```
## 🧪 Testing the API via Swagger

### 1. Start the stack

```bash
docker-compose up -d
```
### 2. Open Swagger UI

```bash
http://localhost:8080/swagger/index.html
```

```bash
- POST /api/v1/wallet
- GET /api/v1/wallets/{id}
```
### 3. Test Deposit

```bash
- Click POST /api/v1/wallet

- Click Try it out

- Paste this JSON payload: (here i put the test data you can put any different data of your choice for walletId and amount. for operationType choose DEPOSIT or WITHDRAW. Here is Test for DEPOSIT and i have given test for WITHDRAW)

{
  "walletId": "11111111-1111-1111-1111-111111111111",
  "operationType": "DEPOSIT",
  "amount": 1000
}
- Click Execute

You should see 204 No Content and no errors.
```
### 4. Test Withdraw

```bash
- Still in POST /api/v1/wallet

- Change payload to: (Test for WITHDRAW)
{
  "walletId": "11111111-1111-1111-1111-111111111111",
  "operationType": "WITHDRAW",
  "amount": 200
}
- Click Execute

- Expect 204 No Content (balance now 800).

If you withdraw more than the current balance, you’ll get 409 Conflict with insufficient funds.
```

### 5. Test Get Balance

```bash
- Click GET /api/v1/wallets/{id}

- Click Try it out

- Enter the Path parameter: (Put the id which you created before)
id = 11111111-1111-1111-1111-111111111111
- Click Execute

- You should get HTTP 200 and a JSON response: (This is based on how much money left after DEPOSIT or WITHDRAW)
{
  "balance": 800
}
```
### 6. When you’re done, tear down the stack with:
```bash
docker-compose down -v
```