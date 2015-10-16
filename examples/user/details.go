/*
 * Prints details of a particular user
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
	var acctAlias = flag.String("a",    "",    "Account alias to use")
	var simple    = flag.Bool("simple", false, "Use simple (debugging) output format")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <Username>\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 || *acctAlias == "" {
		flag.Usage()
		os.Exit(1)
	}

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	user, err := client.GetUserDetails(flag.Arg(0), *acctAlias)
	if err != nil {
		exit.Fatalf("Failed to list user details: %s", err)
	}

	if *simple {
		utils.PrintStruct(user)
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoFormatHeaders(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(false)

		table.SetHeader([]string{ "User", "SAML", "First", "Last", "Email", "Roles" })
		table.Append([]string{
				user.UserName, user.SAMLUserName,
				user.FirstName, user.LastName,
				user.EmailAddress, fmt.Sprint(user.Roles),
		})
		table.Render()
	}
}
