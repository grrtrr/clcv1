/*
 * Lists all users associated with a given account.
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
	var simple = flag.Bool("simple", false, "Use simple (debugging) output format")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s <Account-Alias>\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(0)
	}

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	users, err := client.GetUsers(flag.Arg(0))
	if err != nil {
		exit.Fatalf("Failed to list users of %s: %s", flag.Arg(0), err)
	}

	if len(users) == 0 {
		println("Empty result.")
	} else if *simple {
		for _, l := range users {
			utils.PrintStruct(l)
		}
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoFormatHeaders(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(true)

		table.SetHeader([]string{ "Username", "First", "Last", "Email", "Roles" })
		for _, u := range users {
			table.Append([]string{ u.UserName, u.FirstName, u.LastName,
				     u.EmailAddress, fmt.Sprint(u.Roles) })
		}
		table.Render()
	}
}
