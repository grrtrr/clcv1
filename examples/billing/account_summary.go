/*
 * Print monthly and hourly charges and estimates for a given account or collection of accounts.
 */
package main

import (
	"github.com/grrtrr/clcv1/utils"
	"github.com/grrtrr/clcv1"
	"github.com/grrtrr/exit"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var acctAlias = flag.String("a", "", "Account alias to use")
	var simple = flag.Bool("simple", false, "Use simple (debugging) output format")

	flag.Parse()

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags|log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	acctSummary, err := client.GetAccountSummary(*acctAlias)
	if err != nil {
		exit.Fatalf("Failed to obtain account billing summary of %s: %s", *acctAlias, err)
	}

	if *simple {
		utils.PrintStruct(acctSummary)
	} else {
		fmt.Printf("Total charges during this hour:     $%.2f\n", acctSummary.CurrentHour)
		fmt.Printf("Total charges during previous hour: $%.2f\n", acctSummary.PreviousHour)
		fmt.Printf("Total hourly charges this month:    $%.2f\n", acctSummary.MonthToDate)
		fmt.Printf("Total one-time charges this month:  $%.2f\n", acctSummary.OneTimeCharges)
		fmt.Printf("Total overall charges this month:   $%.2f\n", acctSummary.MonthToDateTotal)
		fmt.Printf("Corresponding monthly estimate:     $%.2f\n", acctSummary.MonthlyEstimate)

	}
}
