/*
 * Get blueprints according to specified criteria
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
	var visib = flag.Int("v", 1, "The visibility level of the Blueprint")
	flag.Parse()

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	blueprints, err := client.GetBlueprints(&clcv1.SearchBlueprintReq{Visibility: *visib})
	if err != nil {
		exit.Fatalf("Failed to list blueprints: %s", err)
	}

	for _, b := range blueprints {
		utils.PrintStruct(b)
	}
}
