/*
 * Lists queue requests
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
	var status = flag.Int("s", 1, "Status type to look for: 1 - All, 2 - Pending, 3 - Complete, 4 - Error")
	flag.Parse()


	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	requests, err := client.ListQueueRequests(clcv1.ItemStatus(*status))
	if err != nil {
		exit.Fatalf("Failed to list queue requests: %s", err)
	}

	if *simple {
		for _, r := range requests {
			utils.PrintStruct(r)
		}
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoFormatHeaders(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(false)

		table.SetHeader([]string{ "ID", "Status", "%", "Step", "Title", "Date" })
		for _, r := range requests {
			table.Append([]string{ fmt.Sprint(r.RequestID), r.CurrentStatus,
				fmt.Sprint(r.PercentComplete), fmt.Sprint(r.StepNumber),
				r.RequestTitle,
				r.StatusDate.Format("Mon, _2 Jan 2006 15:04:05 MST"),
			})
		}

		table.Render()
	}
}
