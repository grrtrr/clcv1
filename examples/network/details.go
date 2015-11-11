/*
 * Get the details for a Network and its IP Addresses.
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
	var location  = flag.String("l", "",       "Data centre alias of the network (uses Home Data Centre by default)")
	var ips       = flag.Bool("ips",    false, "Also list IP addresses")
	var unclaimed = flag.Bool("free",   false, "Also list unclaimed IP addresses (implies -ips)")
	var simple    = flag.Bool("simple", false, "Use simple (debugging) output format")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <Network-Name>\n", path.Base(os.Args[0]))
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

	details, err := client.GetNetworkDetails(flag.Arg(0), *acctAlias, *location)
	if err != nil {
		exit.Fatalf("Failed to query network details of %s: %s", flag.Arg(0), err)
	}

	if *simple {
		utils.PrintStruct(details)
	} else {
		fmt.Printf("Details for %s (%s) at %s:\n", details.Name, details.Description, details.Location)
		fmt.Printf("Gateway: %s\n", details.Gateway)
		fmt.Printf("Netmask: %s\n", details.NetworkMask)

		if !*ips && !*unclaimed {
			return
		}
		fmt.Printf("IP addresses:\n")

		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoFormatHeaders(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(false)

		table.SetHeader([]string{ "Address", "Type", "Claimed", "Used by" })
		for _, ip := range details.IPAddresses {
			if ip.IsClaimed || *unclaimed {
				table.Append([]string{ ip.Address, ip.AddressType, fmt.Sprint(ip.IsClaimed), ip.ServerName })
			}
		}
		table.Render()
	}
}
