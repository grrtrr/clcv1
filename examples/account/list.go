/*
 * Lists all accounts.
 */
package main

import (
	"github.com/olekukonko/tablewriter"
	"github.com/grrtrr/clcv1"
	"github.com/grrtrr/exit"
	"flag"
	"log"
	"os"
)

func main() {
	flag.Parse()

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	accts, err := client.GetAccounts()
	if err != nil {
		exit.Fatalf("Failed to obtain account list: %s", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(true)

	table.SetHeader([]string{ "Account", "Parent", "Location", "Business Name" })
	for _, a := range accts {
		var acct = a.AccountAlias

		if !a.IsActive {
			acct += " (INACTIVE)"
		}
		table.Append([]string{ acct, a.ParentAlias, a.Location, a.BusinessName })
	}
	table.Render()
}
