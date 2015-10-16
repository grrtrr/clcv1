/*
 * For a given Hardware Group name at a given Location, print the corresponding HW Group UUID.
 */
package main

import (
	"github.com/olekukonko/tablewriter"
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
	var location  = flag.String("l",    "",    "Data center location of @Group-Name")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <Group-Name>\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 || *location == "" {
		flag.Usage()
		os.Exit(1)
	}

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	groups, err := client.GetGroups(*location, *acctAlias)
	if err != nil {
		exit.Fatalf("Failed to obtain hardware groups: %s", err)
	}

	if len(groups) == 0 {
		exit.Errorf("Empty result.")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(true)
	table.SetHeader([]string{ "Name", "UUID", "Parent UUID", "System Group?"})
	for _, g := range groups {
		if g.Name == flag.Arg(0) {
			table.Append([]string{ g.Name, g.UUID, g.ParentUUID, fmt.Sprint(g.IsSystemGroup) })
		}
	}

	table.Render()
}
