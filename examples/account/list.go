/*
 * Lists all accounts.
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
	flag.Parse()

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	accts, err := client.GetAccounts()
	if err != nil {
		exit.Fatalf("Failed to obtain account list: %s", err)
	}

	for _, a := range accts {
		utils.PrintStruct(a)
	}
}
