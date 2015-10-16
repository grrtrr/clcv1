/*
 * List the list of deployable networks mapped to an account in any Data Center.
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
	var acctAlias = flag.String("a", "",       "Account alias of the account in question")
	var simple    = flag.Bool("simple", false, "Use simple (debugging) output format")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]  [<Location>]\n", path.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "       Leave Location emtpy to mean home datacenter.\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	networks, err := client.GetDeployableNetworks(*acctAlias, flag.Arg(0))
	if err != nil {
		exit.Fatalf("Failed to list deployable networks: %s", err)
	}

	if len(networks) == 0 {
		println("Empty result.")
	} else if *simple {
		for _, l := range networks {
			utils.PrintStruct(l)
		}
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoFormatHeaders(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(false)

		table.SetHeader([]string{ "Gateway", "Name", "Description", "Account", "Location" })
		for _, l := range networks {
			table.Append([]string{ l.Gateway, l.Name, l.Description, l.AccountAlias, l.Location })
		}
		table.Render()
	}
}
