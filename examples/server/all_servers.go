/*
 * Deep list of all Servers (for a given HW group, for a given location)
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
	var acctAlias = flag.String("a", "", "Account alias of the account that owns the servers")
	var location  = flag.String("l", "", "The data center location")
	var hwGrpUUID = flag.String("u", "", "UUID of the Hardware Group")
	var simple    = flag.Bool("simple", false, "Use simple (debugging) output format")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]  [<location>]\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 0 && *location != "" {
		flag.Usage()
		os.Exit(1)
	} else if flag.NArg() == 1 {
		*location = flag.Arg(0)
	}

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	servers, err := client.GetAllServers(*acctAlias, *hwGrpUUID, *location)
	if err != nil {
		exit.Fatalf("Failed to list all servers: %s", err)
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

		table.SetHeader([]string{
			"Name", "Description",
			"#CPU", "#Disk", "Disk",
			"OS", "IP", "Power", "Who modified", "Modified date",
		})
		for _, s := range servers {
			table.Append([]string{
				s.Name, s.Description,
				fmt.Sprint(s.Cpu), fmt.Sprint(s.DiskCount), fmt.Sprint(s.TotalDiskSpaceGB),
				fmt.Sprintf("%25.25s", s.OperatingSystem), s.IPAddress,
				s.PowerState, fmt.Sprintf("%12.12s", s.ModifiedBy),
				s.DateModified.Format("Jan _2/06 15:04"),
			})
		}
		table.Render()
	}
}
