package main

import (
	"fmt"
	"math"
)

func main() {

	payment := 4000.0
	interest := 8.0
	years := 7.0

	interestRate := (interest / 100.0) / 12
	growthRate := 0.0
	periods := years * 12

	num := (interestRate - growthRate)
	denom := math.Pow(1.0+interestRate, periods) - math.Pow(1.0+growthRate, periods)

	futureValue := payment * denom / num

	other := 275000.0 * math.Pow(1+interestRate, periods)

	fmt.Printf(
		"Payment: %.2f\nInterest Rate: %.2f\nYears: %d\n\nPeriod Rate: %.4f\nPeriods: %d\nGrowth Rate: %.2f\n\nFuture: %.2f\nOther: %.2f\nNet: %.2f\n",
		payment,
		interest,
		int(years),
		interestRate,
		int(periods),
		growthRate,
		futureValue,
		other,
		futureValue+other,
	)
}
