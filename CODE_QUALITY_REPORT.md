# Code Quality Analysis: goftx

## Overview

**goftx** — Go SDK for FTX cryptocurrency exchange API (REST + WebSocket).

- **Lines of code:** ~4300
- **Go version:** 1.14 (EOL, released 2020)
- **Files:** 35 `.go` files (18 source + 12 tests + 2 examples + 3 models)

---

## 1. Architecture

The library follows a clean modular pattern with a central `Client` struct embedding sub-clients for each API domain:

| Module | File | Description |
|--------|------|-------------|
| Account | `account.go` | Account info, leverage |
| Markets | `markets.go` | Market data, order books, candles |
| Orders | `orders.go` | Order placement, modification, cancellation |
| Fills | `fills.go` | Trade fill history |
| Futures | `futures.go` | Futures contracts, funding rates |
| SubAccounts | `subaccounts.go` | Sub-account management |
| SpotMargin | `spot_margin.go` | Spot margin trading |
| Wallet | `wallet.go` | Balances, withdrawals |
| Converts | `converts.go` | OTC quotes/conversions |
| Stream | `websocket.go` | Real-time WebSocket streams |

Data models are separated into `models/` package. Authentication uses HMAC-SHA256 signing.

**Strengths:** clear separation of concerns, consistent API patterns, good use of `decimal.Decimal` for financial data.

---

## 2. Critical Issues

### 2.1 Concurrency Bugs in WebSocket (`websocket.go`)

- **Race condition:** Two goroutines (reader at line ~120 and main select loop at ~173) access `conn` without mutex protection. Reconnection logic can cause data races.
- **Resource leak:** If `s.subscribe()` fails after `s.connect()`, the connection may not be closed properly.
- **Silent event loss:** Type assertions like `fill, ok := event.(*models.FillResponse)` silently discard events on failure with `return`.

### 2.2 Unused Context in `wallet.go`

```go
func (w *Wallet) Withdraw(ctx context.Context, ...) // ctx never used
```

### 2.3 Suppressed Errors

Multiple `// nolint:errcheck` directives suppress potentially important errors in `websocket.go` (HMAC `mac.Write()`, `conn.WriteJSON()`).

---

## 3. Error Handling Issues

| Issue | Location | Severity |
|-------|----------|----------|
| Nil pointer dereference possible | `account.go:40`, similar in other files | High |
| Missing input validation | All API methods | High |
| Raw API errors exposed to caller | `client.go:206` | Medium |
| `ioutil.ReadAll` deprecated since Go 1.16 | `client.go:194` | Low |
| Typo: `signture()` instead of `signature()` | `client.go:170` | Low |

---

## 4. Security Concerns

| # | Issue | Severity |
|---|-------|----------|
| 1 | API credentials loaded from `.env` in tests — risk of accidental commit | Critical |
| 2 | HMAC signature vulnerable to timing attacks (no `hmac.Equal()`) | High |
| 3 | No HTTP client timeout — `http.DefaultClient` used without timeout config | High |
| 4 | No TLS version enforcement or certificate pinning | High |
| 5 | Sub-account name added to HTTP headers without sanitization (header injection) | Medium |
| 6 | No rate limiting — easy to trigger API rate limits | Medium |
| 7 | `OrderType`, `Side` are plain strings — no enum validation | Low |

---

## 5. Test Quality

### `go vet` Results

```
websocket_test.go:18:  context cancel function discarded (context leak)
websocket_test.go:51:  context cancel function discarded (context leak)
websocket_test.go:83:  context cancel function discarded (context leak)
websocket_test.go:111: context cancel function discarded (context leak)
```

### Coverage Gaps

- **No unit tests** — all tests are integration tests requiring live API credentials
- **No mocked HTTP responses** — cannot test error paths
- **No error case testing** — only happy paths covered
- **WebSocket tests skipped** by `.golangci.yml` (`skip-files: websocket_test.go`)
- **Missing tests:** client initialization, auth signing, reconnection logic, context cancellation
- **Flaky design:** time-dependent tests (3–15s timeouts), vulnerable to network latency

---

## 6. Dependencies

| Dependency | Version | Status |
|------------|---------|--------|
| `go` | 1.14 | **EOL** — update to 1.21+ |
| `gorilla/websocket` | v1.4.2 | OK, but v1.5+ available |
| `joho/godotenv` | v1.3.0 | Stable, but only for dev |
| `pkg/errors` | v0.9.1 | Maintained |
| `shopspring/decimal` | v1.2.0 | OK, v1.3+ available |
| `stretchr/testify` | v1.6.1 | **Outdated** — update to v1.8+ |

---

## 7. Code Duplication

WebSocket subscribe methods (`SubscribeToTickers`, `SubscribeToOrderBooks`, etc.) have substantial duplication marked with `// nolint: dupl`. These could be refactored into a generic subscribe function with a channel type parameter.

---

## 8. Recommendations (Priority Order)

1. **Fix race conditions** in `websocket.go` — add mutex for `conn`, handle reconnection safely
2. **Add HTTP client timeout** and pass `context.Context` to HTTP requests
3. **Update Go version** to 1.21+ — get security patches, modern stdlib
4. **Add unit tests with mocked HTTP** — enable testing error paths without live API
5. **Fix `go vet` warnings** — context leak in `websocket_test.go`
6. **Replace deprecated `ioutil.ReadAll`** with `io.ReadAll`
7. **Add input validation** for API parameters
8. **Implement rate limiting** to prevent API throttling
9. **Update dependencies** (testify, gorilla/websocket)
10. **Refactor WebSocket subscriptions** to eliminate duplication

---

## Summary

| Category | Critical | High | Medium | Low |
|----------|----------|------|--------|-----|
| Concurrency | 2 | — | — | — |
| Security | 1 | 3 | 2 | 1 |
| Error Handling | — | 2 | 1 | 2 |
| Testing | — | 2 | 2 | — |
| Dependencies | — | 1 | 1 | — |
| Code Quality | — | — | 1 | 1 |
| **Total** | **3** | **8** | **7** | **4** |

The library is functional but has significant production readiness issues, primarily around **concurrency safety**, **missing timeouts**, and **lack of proper testing**. Addressing the critical and high issues should be the priority before production use.
