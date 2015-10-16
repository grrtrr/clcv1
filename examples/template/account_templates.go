/*
 * List server templates available to the account.
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
	var simple   = flag.Bool("simple",  false, "Use simple (debugging) output format")

	flag.Parse()

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	templates, err := client.GetServerTemplates()
	if err != nil {
		exit.Fatalf("Failed to list account server templates: %s", err)
	}

	if len(templates) == 0 {
		println("Empty result.")
	} else if *simple {
		for _, t := range templates {
			utils.PrintStruct(t)
		}
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoFormatHeaders(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(true)

		table.SetHeader([]string{ "Name", "Description", "OS", "#Disk", "Disk/GB", "Loc" })

		for _, t := range templates {
			table.Append([]string{
				t.Name,	t.Description, fmt.Sprint(t.OperatingSystem),
				fmt.Sprint(t.DiskCount), fmt.Sprint(t.TotalDiskSpaceGB),
				t.Location,
			})
		}
		table.Render()
	}
}
