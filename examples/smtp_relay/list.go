/*
 * Lists all SMTP Aliases for the login account
 */
package main

import (
	"github.com/olekukonko/tablewriter"
	"github.com/grrtrr/clcv1/utils"
	"github.com/grrtrr/clcv1"
	"github.com/grrtrr/exit"
	"flag"
	"log"
	"fmt"
	"os"
)

func main() {
	var simple = flag.Bool("simple", false, "Use simple (debugging) output format")
	flag.Parse()


	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	aliases, err := client.ListAliases()
	if err != nil {
		exit.Fatalf("Failed to list queue requests: %s", err)
	}

	if len(aliases) == 0 {
		fmt.Println("Empty result.")
	} else if *simple {
		for _, l := range aliases {
			utils.PrintStruct(l)
		}
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoFormatHeaders(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(false)

		table.SetHeader([]string{ "Alias", "Password", "Status" })
		for _, a := range aliases {
			table.Append([]string{ a.Alias, a.Password, a.Status })
		}
		table.Render()
	}
}
