/*
 * Configure an existing server
 */
package main

import (
	"github.com/grrtrr/clcv1"
	"github.com/grrtrr/exit"
	"path"
	"flag"
	"log"
	"fmt"
	"os"
)

func main() {
	var hwUUID    = flag.String("u", "", "Hardware UUID of the HW group this server belongs to")
	var acctAlias = flag.String("a", "", "Account alias to use")
	var memGB     = flag.Int("mem",   0, "Amount of memory in GB")
	var numCpu    = flag.Int("cpu",   0, "Number of Cpus to use")
	var extraDrv  = flag.Int("drv",   0, "Extra storage (in GB) to add to server")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <server-name>\n", path.Base(os.Args[0]))
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

	req := clcv1.ConfigureServerReq{ Name: flag.Arg(0) }
	if *acctAlias != "" {
		req.AccountAlias = *acctAlias
	}

	if *numCpu == 0 || *memGB == 0 || *hwUUID == "" {
		fmt.Printf("Fetching details of %s for re-configuration ...\n", flag.Arg(0))

		server, err := client.GetServer(flag.Arg(0), *acctAlias)
		if err != nil {
			exit.Fatalf("Failed to retrieve details of server %q: %s", flag.Arg(0), err)
		}
		req.HardwareGroupUUID = server.HardwareGroupUUID
		req.Cpu               = server.Cpu
		req.MemoryGB          = server.MemoryGB
	}

	if *hwUUID != "" {
		if req.HardwareGroupUUID != "" && *hwUUID != req.HardwareGroupUUID {
			fmt.Printf("Moving server from %s -> %s\n", req.HardwareGroupUUID, *hwUUID)
		}
		req.HardwareGroupUUID = *hwUUID
	}
	if *memGB != 0 {
		if req.MemoryGB != 0 && *memGB != req.MemoryGB {
			fmt.Printf("Changing memory from %d to %d GB\n", req.MemoryGB, *memGB)
		}
		req.MemoryGB = *memGB
	}
	if *numCpu != 0 {
		if req.Cpu != 0 && *numCpu != req.Cpu {
			fmt.Printf("Changing Cpu number from %d to %d\n", req.Cpu, *numCpu)
		}
		req.Cpu = *numCpu
	}
	if *extraDrv != 0 {
		fmt.Printf("Adding %d GB extra storage\n", *extraDrv)
		req.AdditionalStorageGB = *extraDrv
	}
	// FIXME: not addressing CustomFields in this revision

	reqId, err := client.ConfigureServer(&req)
	if err != nil {
		exit.Fatalf("Failed to configure server %s: %s", flag.Arg(0), err)
	}

	fmt.Println("Request ID for server configuration:", reqId)
}
