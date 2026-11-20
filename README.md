# Mortgage Calculator

A small full-stack mortgage calculator. A Go backend does the calculation and serves a plain HTML, CSS and JavaScript frontend. It estimates the **monthly**, **quarterly** and **annual** cost of a home loan, including property tax and insurance.

## Features

- Monthly, quarterly and annual payment estimates
- Breakdown into principal and interest, property tax and insurance
- Loan amount, down payment, total interest and total repayment
- Live updates as you type
- Light and dark mode, following your device's theme setting 
- Input validation with clear error messages
- Standard library only, with no third-party Go dependencies

## Project structure

```
.
├── main.go              HTTP server, routes and handlers
├── calculator.go        Calculation logic and input validation
├── calculator_test.go   Tests for the maths and validation
├── main_test.go         Tests for the API and frontend serving
├── go.mod
└── web/
    ├── index.html       Page structure
    ├── style.css        Styling
    └── app.js           Form handling and API calls
```

## Requirements

- [Go](https://go.dev/dl/) 1.22 or newer

## Run it

From the project folder (the one containing `go.mod`):

```
go run .
```

Then open http://localhost:8080.

To use another port, set the `PORT` environment variable:

```
PORT=3000 go run .
```

The server reads the `web/` folder from disk, so run it from the project root.

## API

### `POST /api/calculate`

Request body:

```json
{
  "homePrice": 400000,
  "downPercent": 20,
  "annualRate": 6.5,
  "termYears": 30,
  "taxPercent": 1.1,
  "annualIns": 1500
}
```

| Field         | Meaning                              |
| ------------- | ------------------------------------ |
| `homePrice`   | Home price                           |
| `downPercent` | Down payment, percent of price       |
| `annualRate`  | Yearly interest rate, percent        |
| `termYears`   | Loan length in years (1 to 50)       |
| `taxPercent`  | Yearly property tax, percent of price |
| `annualIns`   | Yearly home insurance cost           |

Try it:

```
curl -X POST localhost:8080/api/calculate \
  -d '{"homePrice":400000,"downPercent":20,"annualRate":6.5,"termYears":30,"taxPercent":1.1,"annualIns":1500}'
```

The response contains `monthly`, `quarterly` and `annual` objects, each with `principalAndInterest`, `propertyTax`, `insurance` and `total`, plus `downPayment`, `loanAmount`, `numPayments`, `totalInterest` and `totalRepayment`.

Errors return JSON like `{"error": "..."}` with status `400` for a malformed request or `422` for invalid values.

### `GET /api/health`

Returns `{"status":"ok"}`.

## How the calculation works

Monthly principal and interest use the standard amortization formula:

```
M = P * r * (1 + r)^n / ((1 + r)^n - 1)
```

- `P` is the loan amount (price minus down payment)
- `r` is the monthly interest rate (annual rate / 100 / 12)
- `n` is the number of monthly payments (years x 12)

Property tax and insurance are yearly costs divided by 12. The quarterly figure is the monthly total x 3 and the annual figure is the monthly total x 12.

## Tests

```
go test ./...
```

Add `-v` to see each test.

## Deploy

The app reads the `PORT` environment variable, which most hosting platforms set. Because the server serves `web/` from disk, deploy the `web/` folder alongside the program (for example, by copying both into a Docker image).

## Limitations

This is an estimate only. It does not include PMI, HOA fees, closing costs or changes to the interest rate over time.
