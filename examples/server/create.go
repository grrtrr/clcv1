/*
 * Create a new server
 */
package main

import (
	"github.com/grrtrr/clcv1"
	"github.com/grrtrr/exit"
	"encoding/hex"
	"path"
	"flag"
	"log"
	"fmt"
	"os"
)

func main() {
	var hwGroup    = flag.String("g", "",     "UUID or name (if unique) of the HW group to add this server to")
	var location   = flag.String("l", "",     "Data centre alias")
	var acctAlias  = flag.String("a", "",     "Account alias to use")
	var template   = flag.String("t", "",     "The name of the template to create the server from")
	var seed       = flag.String("s", "AUTO", "The seed for the server name (max 6 characters)")
	var desc       = flag.String("D", "",     "Description of the server")

	var net        = flag.String("net",  "",        "Name of the Network to use")
	var primDNS    = flag.String("dns1", "8.8.8.8", "Primary DNS to use")
	var secDNS     = flag.String("dns2", "8.8.4.4", "Secondary DNS to use")
	var password   = flag.String("pass", "",        "Desired password. Leave blank to auto-generate")

	var extraDrv   = flag.Int("drive",  0, "Extra drive (in GB) to add to server. Set to 0 to leave out")
	var numCpu     = flag.Int("cpu",    1, "Number of Cpus to use")
	var memGB      = flag.Int("memory", 4, "Amount of memory in GB")
	var serverType = flag.Int("type",   1, "The type of server to create (1: Standard, 2: Enterprise)")
	var servLevel  = flag.Int("level",  2, "Data storage service level (1: Premium, 2: Standard)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()
	if *hwGroup == "" || *location == "" || *template == "" || *seed == "" || *net == "" {
		flag.Usage()
		os.Exit(0)
	}

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	req := clcv1.CreateServerReq{
		// The alias of the account to own the server.
		// If not provided it will assume the account to which the API user is mapped.
		// Providing this value gives you the ability to create servers in your sub accounts. (optional)
		AccountAlias: *acctAlias,

		// The alias of the data center in which to create the server.
		// If not provided, will default to the API user's default data center.
		// NOTE: empty value can cause error if location is not the same as HWUUID's
		LocationAlias: *location,

		// The name of the template to create the server from (required)
		Template: *template,

		// The alias for the server. Limit 6 charcters (required)
		Alias: *seed,

		// An optional description for the server. If none is supplied the server name will be used. (optional)
		Description: *desc,

		// The unique identifier of the Hardware Group to add this server to (required)
		HardwareGroupUUID: *hwGroup,

		//The type of server to create (required)
		ServerType: *serverType,

		// The service level/performance for the underlying data store (required)
		ServiceLevel: *servLevel,

		// The number of processors to configure the server with (required)
		Cpu: *numCpu,

		// The number of GB of memory to configure the server with (required)
		MemoryGB: *memGB,

		// The size in GB of an additional drive to add to the server (required).
		// If no additional drive is needed, pass in a 0 value.
		ExtraDriveGB: *extraDrv,

		// The primary DNS to set on the server (optional)
		// If not supplied the default value set on the account will be used.
		PrimaryDns: *primDNS,

		// The secondary DNS to set on the server (optional)
		// If not supplied the default value set on the account will be used.
		SecondaryDns: *secDNS,

		// The name of the network to which to deploy the server.
		// If your account has not yet been assigned a network, leave this blank and one will be assigned automatically.
		// If one or more networks are available, the network name is required.
		Network: *net,

		// The desired Admin/Root password (optional).
		// Please note the password must meet the password strength policy.
		// Leave blank to have the system generate a password
		Password: *password,

		// A list of Custom Fields associated to this server
		// FIXME: not supported yet
		CustomFields: nil,
	}

	/* hwGroup may be hex uuid or group name */
	if _, err := hex.DecodeString(*hwGroup); err != nil {
		if group, err := client.GetGroupByName(*hwGroup, *location, *acctAlias); err != nil {
			exit.Errorf("Failed to resolve group name %q: %s", *hwGroup, err)
		} else if group == nil {
			exit.Errorf("No group named %q was found on %s", *hwGroup, *location)
		} else {
			req.HardwareGroupUUID = group.UUID
		}
	}

	reqId, err := client.CreateServer(&req)
	if err != nil {
		exit.Fatalf("Failed to create server: %s", err)
	}

	fmt.Println("Request ID for server creation:", reqId)
}
