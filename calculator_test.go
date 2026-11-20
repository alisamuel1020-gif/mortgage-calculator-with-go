package main

import (
	"math"
	"testing"
)

func close(a, b float64) bool { return math.Abs(a-b) < 0.005 }

var defaultInput = Input{
	HomePrice: 400000, DownPercent: 20, AnnualRate: 6.5,
	TermYears: 30, TaxPercent: 1.1, AnnualIns: 1500,
}

func TestCalculateDefaultLoan(t *testing.T) {
	r := Calculate(defaultInput)

	checks := []struct {
		name      string
		got, want float64
	}{
		{"down payment", r.DownPayment, 80000},
		{"loan amount", r.LoanAmount, 320000},
		{"monthly P&I", r.Monthly.PrincipalAndInterest, 2022.62},
		{"monthly tax", r.Monthly.PropertyTax, 366.67},
		{"monthly insurance", r.Monthly.Insurance, 125},
		{"monthly total", r.Monthly.Total, 2514.28},
		{"quarterly total", r.Quarterly.Total, 7542.85},
		{"annual total", r.Annual.Total, 30171.41},
		{"total interest", r.TotalInterest, 408142.36},
	}
	for _, c := range checks {
		if !close(c.got, c.want) {
			t.Errorf("%s = %.2f, want %.2f", c.name, c.got, c.want)
		}
	}
	if r.NumPayments != 360 {
		t.Errorf("num payments = %d, want 360", r.NumPayments)
	}
}

func TestCalculateZeroRate(t *testing.T) {
	r := Calculate(Input{HomePrice: 120000, DownPercent: 0, AnnualRate: 0, TermYears: 10})
	if !close(r.Monthly.PrincipalAndInterest, 1000) {
		t.Errorf("monthly P&I = %.2f, want 1000", r.Monthly.PrincipalAndInterest)
	}
	if !close(r.TotalInterest, 0) {
		t.Errorf("total interest = %.2f, want 0", r.TotalInterest)
	}
}

func TestPeriodsAreConsistent(t *testing.T) {
	r := Calculate(defaultInput)
	if math.Abs(r.Quarterly.Total-r.Monthly.Total*3) > 3*0.005 {
		t.Errorf("quarterly %.2f is not about 3x monthly %.2f", r.Quarterly.Total, r.Monthly.Total)
	}
	if math.Abs(r.Annual.Total-r.Monthly.Total*12) > 12*0.005 {
		t.Errorf("annual %.2f is not about 12x monthly %.2f", r.Annual.Total, r.Monthly.Total)
	}
}

func TestValidate(t *testing.T) {
	bad := map[string]func(*Input){
		"zero price":        func(i *Input) { i.HomePrice = 0 },
		"negative price":    func(i *Input) { i.HomePrice = -1 },
		"100% down":         func(i *Input) { i.DownPercent = 100 },
		"negative down":     func(i *Input) { i.DownPercent = -5 },
		"negative rate":     func(i *Input) { i.AnnualRate = -1 },
		"rate over 100":     func(i *Input) { i.AnnualRate = 101 },
		"zero term":         func(i *Input) { i.TermYears = 0 },
		"term over 50":      func(i *Input) { i.TermYears = 51 },
		"negative tax":      func(i *Input) { i.TaxPercent = -1 },
		"negativeTaxInsure": func(i *Input) { i.AnnualIns = -1 },
	}
	for name, mutate := range bad {
		in := defaultInput
		mutate(&in)
		if err := in.Validate(); err == nil {
			t.Errorf("%s: expected an error, got nil", name)
		}
	}
	if err := defaultInput.Validate(); err != nil {
		t.Errorf("default input should be valid, got %v", err)
	}
}
