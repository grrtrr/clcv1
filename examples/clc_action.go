/*
 * Rewrite of the 'clc_action' bash script into go.
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
	"fmt"
	"log"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage:\n")
	fmt.Fprintf(os.Stderr, "\t%s [options]      <action>  <Server-Name|Group-UUID>\n", path.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "\t%s -l <Location>  <action>  <Group-Name>\n\n", path.Base(os.Args[0]))

	for _, r := range [][]string{
		{ "show",     "show current status of server/group (group requires -l to be set)" },
		{ "on",       "power on server/group (or resume from paused state)" },
		{ "off",      "power off server/group" },
		{ "shutdown", "OS-level shutdown followed by power-off for server/group" },
		{ "pause",    "pause server/group" },
		{ "reset",    "perform forced power-cycle on server/group" },
		{ "reboot",   "reboot server/group" },
		{ "snapshot", "snapshot server (not supported for groups)" },
		{ "archive",  "archive the server/group" },
		{ "delete",   "delete server/group (CAUTION)" },
		{ "help",     "print this help screen" },
	} {
		fmt.Fprintf(os.Stderr, "\t%-10s %s\n", r[0], r[1])
	}
	fmt.Fprintf(os.Stderr, "\n")

	flag.PrintDefaults()
	os.Exit(0)
}

func action_map(server bool, client *clcv1.Client) map[string]func(string, string) (int, error) {
	if server {
		/* Server Action */
		return map[string]func(string, string) (int, error){
			"on":       client.PowerOnServer,
			"off":      client.PowerOffServer,
			"pause":    client.PauseServer,
			"reset":    client.ResetServer,
			"reboot":   client.RebootServer,
			"shutdown": client.ShutdownServer,
			"archive":  client.ArchiveServer,
			"delete":   client.DeleteServer,
			"snapshot": client.SnapshotServer,
		}
	}
	/* Group Action */
	return map[string]func(string, string) (int, error){
		"on":       client.PowerOnHardwareGroup,
		"off":      client.PowerOffHardwareGroup,
		"pause":    client.PauseHardwareGroup,
		"reset":    client.ResetHardwareGroup,
		"reboot":   client.RebootHardwareGroup,
		"shutdown": client.ShutdownHardwareGroup,
		"archive":  client.ArchiveHardwareGroup,
		"delete":   client.DeleteHardwareGroup,
	}
}

func main() {
	var location  = flag.String("l", "", "Location to use for <Group-Name>")
	var acctAlias = flag.String("a", "", "Account alias to use (to override default)")
	var server_action bool
	var action, where string

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 2 {
		action, where = flag.Arg(0), flag.Arg(1)
	} else if flag.NArg() == 1 && flag.Arg(0) == "show" {
		if *location == "" {
			exit.Errorf("Showing group details requires location (-l) argument.")
		}
		action = flag.Arg(0)
	} else {
		usage()
	}

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags|log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	/* If the first argument decodes as a hex value, assume it is a Hardware Group UUID */
	if _, err := hex.DecodeString(where); err == nil {
		server_action = false
	} else if utils.LooksLikeServerName(where) {
		server_action = true
		if *location != "" {
			fmt.Fprintf(os.Stderr, "WARNING: location (%s) ignored for %s\n", *location, where)
		}
	} else if *location != "" && where != "" {
		if group, err := client.GetGroupByName(where, *location, *acctAlias); err != nil {
			exit.Errorf("Failed to resolve group name %q: %s", where, err)
		} else if group == nil {
			exit.Errorf("No group named %q was found on %s", where, *location)
		} else {
			where = group.UUID
		}
	} else {
		server_action = true
	}

	switch action {
	case "show":
		if server_action {
			showServer(client, where, *acctAlias)
		} else {
			showGroup(client, where, *acctAlias, *location)
		}
		os.Exit(0)
	case "help":
		usage()
	}

	/* Long-running commands that return a RequestID */
	handler, ok := action_map(server_action, client)[action]
	if !ok {
		exit.Fatalf("Unsupported action %s", action)
	}

	reqId, err := handler(where, *acctAlias)
	if err != nil {
		exit.Fatalf("Command %q failed: %s", action, err)
	}

	fmt.Printf("Request ID for %q action: %d\n", action, reqId)

	locationStr := *location
	if server_action {
		locationStr = utils.ExtractLocationFromServerName(where)
	}
	client.PollDeploymentStatus(reqId, locationStr, *acctAlias, 1)
}


// Show server details
// @client:    authenticated CLCv1 Client
// @servname:  server name
// @acctAlias: account alias to use (leave blank to use default)
func showServer(client *clcv1.Client, servname, acctAlias string) {
	server, err := client.GetServer(servname, acctAlias)
	if err != nil {
		exit.Fatalf("Failed to list details of server %q: %s", servname, err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(true)

	table.SetHeader([]string{
		"Name",	"Power", "Status",
		"OS", "Description",
		"CPU", "Disk",
		"IP", "Last Change",
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
		server.Name, server.PowerState, server.Status,
		fmt.Sprint(server.OperatingSystem),
		server.Description,
		fmt.Sprint(server.Cpu), fmt.Sprintf("%dGB", server.TotalDiskSpaceGB),
		strings.Join(IPs, " "),
		modifiedStr,
	})
	table.Render()
}

// Show group details
// @client:    authenticated CLCv1 Client
// @uuid:      hardware group UUID to use
// @acctAlias: account alias to use (leave blank to use default)
// @location:  data centre location (needed to resolve @uuid)
func showGroup(client *clcv1.Client, uuid, acctAlias, location string) {
	if location == "" {
		exit.Errorf("Location is required in order to show the group hierarchy starting at %s", uuid)
	}
	root, err := client.GetGroupHierarchy(location, acctAlias, true)
	if err != nil {
		exit.Fatalf("Failed to look up groups at %s: %s", location, err)
	}
	start := root
	if uuid != "" {
		start = clcv1.FindGroupNode(root, func(g *clcv1.GroupNode) bool {
			return g.UUID == uuid
		})
		if start == nil {
			exit.Fatalf("Failed to look up UUID %s at %s", uuid, location)
		}
	}
	clcv1.PrintGroupHierarchy(start, "")
}
