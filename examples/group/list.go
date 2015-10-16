/*
 * Lists all hardware groups associated with a given location (and account alias).
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
	var acctAlias = flag.String("a",    "",    "Account alias of the account in question")
	var simple    = flag.Bool("simple", false, "Use simple (debugging) output format")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <Location>\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()

	/* The Location argument is always required */
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

	groups, err := client.GetGroups(flag.Arg(0), *acctAlias)
	if err != nil {
		exit.Fatalf("Failed to obtain hardware groups: %s", err)
	}

	if len(groups) == 0 {
		println("Empty result.")
	} else if *simple {
		for _, g := range groups {
			utils.PrintStruct(g)
		}
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoFormatHeaders(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(true)

		table.SetHeader([]string{ "Name", "UUID", "Parent UUID", "System Group?"})
		for _, g := range groups {
			table.Append([]string{ g.Name, g.UUID, g.ParentUUID, fmt.Sprint(g.IsSystemGroup) })
		}
		table.Render()
	}
}
