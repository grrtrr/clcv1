/*
 * List the details of one server
 */
package main

import (
	"github.com/olekukonko/tablewriter"
	"github.com/dustin/go-humanize"
	"github.com/grrtrr/clcv1/utils"
	"github.com/grrtrr/clcv1"
	"github.com/grrtrr/exit"
	"encoding/hex"
	"strings"
	"path"
	"flag"
	"log"
	"fmt"
	"os"
)

func main() {
	var acctAlias = flag.String("a",    "",    "Account alias to use")
	var simple    = flag.Bool("simple", false, "Use simple (debugging) output format")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <Server Name>\n", path.Base(os.Args[0]))
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

	server, err := client.GetServer(flag.Arg(0), *acctAlias)
	if err != nil {
		exit.Fatalf("Failed to list details of server %q: %s", flag.Arg(0), err)
	}

	if *simple {
		utils.PrintStruct(server)
	} else {
		grp, err := client.GetGroupByUUID(server.HardwareGroupUUID, server.Location, *acctAlias)
		if err != nil {
			exit.Fatalf("Failed to resolve group UUID: %s", err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoFormatHeaders(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(true)

		table.SetHeader([]string{
			"Name", "Group", "Description", "OS",
			"CPU", "Disk",
			"IP", "Power", "Last Change",
		})

		/* Also display public IP address(es) if configured. */
		IPs := []string{server.IPAddress}
		for _, ip := range server.IPAddresses {
			if ip.IsPublic() {
				IPs = append([]string{ip.Address}, IPs...)
			}
		}

		modifiedStr := humanize.Time(server.DateModified.Time)
		/* The ModifiedBy field can be an email address, or an API Key (hex string) */
		if _, err := hex.DecodeString(server.ModifiedBy); err == nil {
			modifiedStr += " via API Key"
		} else if len(server.ModifiedBy) > 6 {
			modifiedStr += " by " + server.ModifiedBy[:6]
		} else {
			modifiedStr += " by " + server.ModifiedBy
		}

		table.Append([]string{
			server.Name, grp.Name, server.Description, fmt.Sprint(server.OperatingSystem),
			fmt.Sprint(server.Cpu),	fmt.Sprintf("%dGB", server.TotalDiskSpaceGB),
			strings.Join(IPs, " "),
			server.PowerState, modifiedStr,
		})
		table.Render()
	}
}
