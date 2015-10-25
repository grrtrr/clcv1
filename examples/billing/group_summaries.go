/*
 * Get the charges for groups and servers within a given account and date range.
 */
package main

import (
	"github.com/grrtrr/clcv1"
	"github.com/grrtrr/exit"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var acctAlias = flag.String("a", "", "Account alias to use")
	var startDate = flag.String("start", "", "Start date of the query range")
	var endDate = flag.String("end", "", "End date of the query range")

	flag.Parse()

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags|log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	gSum, err := client.GetGroupSummaries(*startDate, *endDate, *acctAlias)
	if err != nil {
		exit.Fatalf("Failed to obtain hardware group summaries: %s", err)
	}

	fmt.Printf("Group summaries for %s from %s to %s:\n", gSum.AccountAlias, gSum.StartDate, gSum.EndDate)

	fmt.Printf("\tCharges for the current hour:      $%.2f\n", gSum.Summary.CurrentHour)
	fmt.Printf("\tCharges for the previous hour:     $%.2f\n", gSum.Summary.PreviousHour)
	fmt.Printf("\tCurrent charges so far this month: $%.2f\n", gSum.Summary.MonthToDate)
	fmt.Printf("\tEstimated monthly costs:           $%.2f\n", gSum.Summary.MonthlyEstimate)

	for _, gt := range gSum.GroupTotals {
		fmt.Printf("Totals for group %s (%d) at %s:\n", gt.GroupName, gt.GroupID, gt.LocationAlias)

		fmt.Printf("\tCharges for the current hour:      $%.2f\n", gt.CurrentHour)
		fmt.Printf("\tCharges for the previous hour:     $%.2f\n", gt.PreviousHour)
		fmt.Printf("\tCurrent charges so far this month: $%.2f\n", gt.MonthToDate)
		fmt.Printf("\tEstimated monthly costs:           $%.2f\n", gt.MonthlyEstimate)

		for _, srv := range gt.ServerTotals {
			fmt.Printf("\tServer totals for %s:\n", srv.ServerName)

			fmt.Printf("\t\tCharges for the current hour:      $%.2f\n", srv.CurrentHour)
			fmt.Printf("\t\tCharges for the previous hour:     $%.2f\n", srv.PreviousHour)
			fmt.Printf("\t\tCurrent charges so far this month: $%.2f\n", srv.MonthToDate)
			fmt.Printf("\t\tEstimated monthly costs:           $%.2f\n", srv.MonthlyEstimate)
		}
	}
}
