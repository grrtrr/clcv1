/*
 * Get estimated costs for a group of servers.
 */
package main

import (
	"github.com/grrtrr/clcv1/utils"
	"github.com/grrtrr/clcv1"
	"github.com/grrtrr/exit"
	"path"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var acctAlias = flag.String("a", "", "Account alias to use")
	var simple = flag.Bool("simple", false, "Use simple (debugging) output format")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <Hardware Group UUID>\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags|log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	grpEst, err := client.GetGroupEstimate(flag.Arg(0), *acctAlias)
	if err != nil {
		exit.Fatalf("Failed to obtain hardware group estimate for %s: %s", flag.Arg(0), err)
	}

	if *simple {
		utils.PrintStruct(grpEst)
	} else {
		fmt.Printf("Usage and corresponding monthly estimate for %s:\n", flag.Arg(0))
		fmt.Printf("Charges incurred this hour:   $%.2f\n", grpEst.CurrentHour)
		fmt.Printf("Charges during previous hour: $%.2f\n", grpEst.PreviousHour)
		fmt.Printf("Charges incurred up to today: $%.2f\n", grpEst.MonthToDate)
		fmt.Printf("Predicted monthly cost:       $%.2f\n", grpEst.MonthlyEstimate)
	}
}
