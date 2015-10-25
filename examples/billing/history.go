/*
 * Print the entire billing history for a given account or collection of accounts.
 */
package main

import (
	"github.com/olekukonko/tablewriter"
	"github.com/grrtrr/clcv1"
	"github.com/grrtrr/exit"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var acctAlias = flag.String("a", "", "Account alias to use")

	flag.Parse()

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags|log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	bh, err := client.GetBillingHistory(*acctAlias)
	if err != nil {
		exit.Fatalf("Failed to obtain billing history: %s", err)
	}

	fmt.Printf("Billing history for %s:\n", bh.AccountAlias)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)
	table.SetAlignment(tablewriter.ALIGN_RIGHT)
	table.SetAutoWrapText(true)

	table.SetHeader([]string{ "Date", "Debit", "Credit", "ID", "Description" })

	for _, le := range bh.BillingHistory {
		table.Append([]string{
			le.Date.Time.Format("Jan 2006"),
			fmt.Sprintf("%.2f", le.Debit),
			fmt.Sprintf("%.2f", le.Credit),
			le.InvoiceID, le.Description,
		})
	}
	table.Render()
}
