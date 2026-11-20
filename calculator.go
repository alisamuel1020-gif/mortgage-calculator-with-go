package main

import (
	"errors"
	"math"
)

type Input struct {
	HomePrice   float64 `json:"homePrice"`
	DownPercent float64 `json:"downPercent"`
	AnnualRate  float64 `json:"annualRate"`
	TermYears   int     `json:"termYears"`
	TaxPercent  float64 `json:"taxPercent"`
	AnnualIns   float64 `json:"annualIns"`
}

type Breakdown struct {
	PrincipalAndInterest float64 `json:"principalAndInterest"`
	PropertyTax          float64 `json:"propertyTax"`
	Insurance            float64 `json:"insurance"`
	Total                float64 `json:"total"`
}

type Result struct {
	DownPayment    float64   `json:"downPayment"`
	LoanAmount     float64   `json:"loanAmount"`
	NumPayments    int       `json:"numPayments"`
	TotalInterest  float64   `json:"totalInterest"`
	TotalRepayment float64   `json:"totalRepayment"`
	Monthly        Breakdown `json:"monthly"`
	Quarterly      Breakdown `json:"quarterly"`
	Annual         Breakdown `json:"annual"`
}

func (in Input) Validate() error {
	switch {
	case in.HomePrice <= 0 || in.HomePrice > 1e10:
		return errors.New("home price must be greater than 0 and at most 10,000,000,000")
	case in.DownPercent < 0 || in.DownPercent >= 100:
		return errors.New("down payment must be at least 0 and less than 100 percent")
	case in.AnnualRate < 0 || in.AnnualRate > 100:
		return errors.New("interest rate must be between 0 and 100 percent")
	case in.TermYears < 1 || in.TermYears > 50:
		return errors.New("term must be between 1 and 50 years")
	case in.TaxPercent < 0 || in.TaxPercent > 100:
		return errors.New("property tax must be between 0 and 100 percent")
	case in.AnnualIns < 0 || in.AnnualIns > 1e9:
		return errors.New("insurance must be between 0 and 1,000,000,000 per year")
	}
	return nil
}

func Calculate(in Input) Result {
	downPayment := in.HomePrice * in.DownPercent / 100
	loan := in.HomePrice - downPayment

	monthlyRate := in.AnnualRate / 100 / 12
	n := in.TermYears * 12

	var pi float64
	if monthlyRate == 0 {
		pi = loan / float64(n)
	} else {
		growth := math.Pow(1+monthlyRate, float64(n))
		pi = loan * monthlyRate * growth / (growth - 1)
	}

	tax := in.HomePrice * in.TaxPercent / 100 / 12
	ins := in.AnnualIns / 12

	monthly := Breakdown{PrincipalAndInterest: pi, PropertyTax: tax, Insurance: ins, Total: pi + tax + ins}
	totalInterest := pi*float64(n) - loan

	return Result{
		DownPayment:    round2(downPayment),
		LoanAmount:     round2(loan),
		NumPayments:    n,
		TotalInterest:  round2(totalInterest),
		TotalRepayment: round2(loan + totalInterest),
		Monthly:        monthly.scaled(1),
		Quarterly:      monthly.scaled(3),
		Annual:         monthly.scaled(12),
	}
}

func (b Breakdown) scaled(months float64) Breakdown {
	return Breakdown{
		PrincipalAndInterest: round2(b.PrincipalAndInterest * months),
		PropertyTax:          round2(b.PropertyTax * months),
		Insurance:            round2(b.Insurance * months),
		Total:                round2(b.Total * months),
	}
}

func round2(v float64) float64 { return math.Round(v*100) / 100 }
