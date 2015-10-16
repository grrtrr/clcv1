/*
 * List the AccountCustomFields associated with an account.
 */
package main

import (
	"github.com/grrtrr/clcv1/utils"
	"github.com/grrtrr/clcv1"
	"github.com/grrtrr/exit"
	"flag"
	"log"
	"os"
)

func main() {
	var acctAlias = flag.String("a", "", "Account alias to use")

	flag.Parse()

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	customFields, err := client.GetCustomFields(*acctAlias)
	if err != nil {
		exit.Fatalf("Failed to obtain Custom Fields of %s: %s", *acctAlias, err)
	}

	if len(customFields) == 0 {
		println("Empty result.")
	} else {
		for _, cf := range customFields {
			utils.PrintStruct(cf)
		}
	}
}
