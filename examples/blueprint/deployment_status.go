/*
 * Query the deployment status for a given RequestID
 */
package main

import (
	"strconv"
	"github.com/grrtrr/clcv1"
	"github.com/grrtrr/exit"
	"path"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var location  = flag.String("l", "", "The location of the deployment to retrieve status for (required)")
	var pollIntvl = flag.Int("i",     1, "Poll interval in seconds (to monitor progress")
	var acctAlias = flag.String("a", "", "Account alias to use")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <Request-ID>\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 1 || *location == "" {
		flag.Usage()
		os.Exit(1)
	}

	reqId, err := strconv.ParseUint(flag.Arg(0), 10, 32)
	if err != nil {
		exit.Errorf("Invalid Request ID %q: %s", flag.Arg(0), err)
	}

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	err = client.PollDeploymentStatus(int(reqId), *location, *acctAlias, *pollIntvl)
	if err != nil {
		exit.Fatalf("Failed to poll status of request ID %d: %s", reqId, err)
	}
}
