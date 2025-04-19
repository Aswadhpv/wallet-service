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