/*
 * List server templates available at a location and/or for the account.
 */
package main

import (
	"github.com/olekukonko/tablewriter"
	"github.com/grrtrr/clcv1/utils"
	"github.com/grrtrr/clcv1"
	"github.com/grrtrr/exit"
	"path"
	"flag"
	"log"
	"fmt"
	"os"
)

func main() {
	var acctAlias = flag.String("a",        "", "Account alias to use")
	var simple    = flag.Bool("simple",  false, "Use simple (debugging) output format")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <location>\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	templates, err := client.ListAvailableServerTemplates(*acctAlias, flag.Arg(0))
	if err != nil {
		exit.Fatalf("Failed to list available templates: %s", err)
	}

	if len(templates) == 0 {
		println("Empty result.")
	} else if *simple {
		for _, t := range templates {
			utils.PrintStruct(t)
		}
	} else {
		fmt.Printf("Server templates available at %s:\n", templates[0].Location)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoFormatHeaders(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(true)

		table.SetHeader([]string{ "Name", "Description", "OS", "#Disk", "Disk/GB" })

		for _, t := range templates {
			table.Append([]string{
				t.Name,	t.Description, fmt.Sprint(t.OperatingSystem),
				fmt.Sprint(t.DiskCount), fmt.Sprint(t.TotalDiskSpaceGB),
			})
		}
		table.Render()
	}
}
