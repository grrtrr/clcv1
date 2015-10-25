/*
 * Get the server-based hourly cost for a given time period.
 */
package main

import (
	"flag"
	"fmt"
	"github.com/grrtrr/clcv1"
	"github.com/grrtrr/exit"
	"github.com/olekukonko/tablewriter"
	"log"
	"os"
	"path"
)

func main() {
	var acctAlias = flag.String("a", "", "Account alias to use")
	var startDate = flag.String("start", "", "Start date of the query range")
	var endDate   = flag.String("end", "", "End date of the query range")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <Servername>\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags|log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	hr, err := client.GetServerHourlyCharges(flag.Arg(0), *startDate, *endDate, *acctAlias)
	if err != nil {
		exit.Fatalf("Failed to obtain server hourly charges for %s: %s", flag.Arg(0), err)
	}

	fmt.Printf("Hourly charges for %s (%s) from %s until %s:\n", hr.ServerName, hr.AccountAlias,
		hr.StartDate.Time.Format("2 Jan 2006 3:04PM"), hr.EndDate.Time.Format("2 Jan 2006 3:04PM"))

	fmt.Printf("\tCharges for the current hour:      $%.2f\n", hr.Summary.CurrentHour)
	fmt.Printf("\tCharges for the previous hour:     $%.2f\n", hr.Summary.PreviousHour)
	fmt.Printf("\tCurrent charges so far this month: $%.2f\n", hr.Summary.MonthToDate)
	fmt.Printf("\tEstimated monthly costs:           $%.2f\n", hr.Summary.MonthlyEstimate)

	if len(hr.HourlyCharges) > 0 {
		fmt.Println("\nHourly charges:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoFormatHeaders(false)
		table.SetAlignment(tablewriter.ALIGN_RIGHT)
		table.SetAutoWrapText(true)

		table.SetHeader([]string{"Hour", "CPU", "Memory", "Storage", "OS"})

		for _, ch := range hr.HourlyCharges {
			table.Append([]string{ch.Hour, ch.ProcessorCost,
				ch.MemoryCost, ch.StorageCost, ch.OSCost,
			})
		}
		table.Render()
	}
}
