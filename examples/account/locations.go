/*
 * Lists all datacenter locations
 */
package main

import (
	"github.com/olekukonko/tablewriter"
	"github.com/grrtrr/clcv1/utils"
	"github.com/grrtrr/clcv1"
	"github.com/grrtrr/exit"
	"flag"
	"log"
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

	locations, err := client.GetLocations()
	if err != nil {
		exit.Fatalf("Failed to list queue requests: %s", err)
	}

	if len(locations) == 0 {
		println("Empty result.")
	} else if *simple {
		for _, l := range locations {
			utils.PrintStruct(l)
		}
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoFormatHeaders(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(false)

		table.SetHeader([]string{ "Alias", "Region" })
		for _, l := range locations {
			table.Append([]string{ l.Alias, l.Region })
		}
		table.Render()
	}
}
