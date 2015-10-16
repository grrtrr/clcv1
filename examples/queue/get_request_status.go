/*
 * Prints the request status of a single queue request
 */
package main

import (
	"github.com/grrtrr/clcv1"
	"github.com/grrtrr/exit"
	"strconv"
	"path"
	"flag"
	"log"
	"fmt"
	"os"
)

func main() {
	var reqId int64

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <RequestID>\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	reqId, err := strconv.ParseInt(flag.Arg(0), 10, 64)
	if err != nil {
		exit.Fatalf("Invalid RequestId %q", flag.Arg(0))
	}

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	request, err := client.GetRequestStatus(int(reqId))
	if err != nil {
		exit.Fatalf("Failed to list queue requests: %s", err)
	}

	fmt.Printf("Request:  %d\n", request.RequestID)
	fmt.Printf("Title:    %s\n", request.RequestTitle)
	fmt.Printf("Status:   %s\n", request.CurrentStatus)
	fmt.Printf("Progress: %s\n", request.ProgressDesc)
	fmt.Printf("Complete: %d%%\n", request.PercentComplete)
	fmt.Printf("Step:     %d\n", request.StepNumber)
	fmt.Printf("When:     %s\n", request.StatusDate.Format("Mon, _2 Jan 2006 15:04:05 MST"))
}
