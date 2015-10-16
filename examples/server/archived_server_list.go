/*
 * Print list of archived servers.
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
	var simple    = flag.Bool("simple", false, "Use simple (debugging) output format")
	var location  = flag.String("l", "",       "Data center location")
	var acctAlias = flag.String("a", "",       "Account alias to use")

	flag.Parse()

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}


	//	servers, err := client.GetArchiveServers()
	servers, err := client.ListArchiveServers(*acctAlias, *location)
	if err != nil {
		exit.Fatalf("Failed to list archived servers: %s", err)
	}

	if *location != "" {
		fmt.Printf("Archived servers in %s:\n", *location)
	}
	if len(servers) == 0 {
		println("Empty result.")
	} else if *simple {
		for _, s := range servers {
			utils.PrintStruct(s)
		}
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoFormatHeaders(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(true)

		table.SetHeader([]string{ "Name", "Description" })
		for _, s := range servers {
			table.Append([]string{ s.Name, s.Description })
		}
		table.Render()
	}

}
