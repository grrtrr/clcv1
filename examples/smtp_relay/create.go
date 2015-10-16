/*
 * Create a new SMTP Relay Alias
 */
package main

import (
	"github.com/grrtrr/clcv1"
	"github.com/grrtrr/exit"
	"flag"
	"log"
	"fmt"
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

	alias, password, err := client.CreateAlias()
	if err != nil {
		exit.Fatalf("Failed to create new SMTP relay alias: %s", err)
	}

	fmt.Printf("New Relay Alias:    %s\n", alias)
	fmt.Printf("Password for Alias: %s\n", password)
}
