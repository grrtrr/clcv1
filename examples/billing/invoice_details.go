/*
 * Get the details for a given invoice within an account.
 */
package main

import (
	"github.com/olekukonko/tablewriter"
	"github.com/grrtrr/clcv1"
	"github.com/grrtrr/exit"
	"strings"
	"path"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var acctAlias = flag.String("a", "", "Account alias to use")
	var itemDetails = flag.Bool("details", false, "Print individual line item details also")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <Invoice-ID>\n", path.Base(os.Args[0]))
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

	id, err := client.GetInvoiceDetails(flag.Arg(0), *acctAlias)
	if err != nil {
		exit.Fatalf("Failed to obtain invoice details for %s: %s", flag.Arg(0), err)
	}

	fmt.Printf("Details of invoice %q for %s (%s):\n", id.Invoice.ID, id.Invoice.PricingAccountAlias, id.Invoice.CompanyName)
	fmt.Printf("Address:               %s, %s, %s, %s\n", id.Invoice.Address1, id.Invoice.City, id.Invoice.StateProvince, id.Invoice.PostalCode)
	fmt.Printf("Parent account alias:  %s\n", id.Invoice.ParentAccountAlias)
	fmt.Printf("Support level:         %s\n", strings.Title(id.SupportLevel))
	fmt.Printf("Billing contact email: %s\n", id.Invoice.BillingContactEmail)
	if id.Invoice.PONumber != "" {
		fmt.Printf("Purchase order:        %s\n", id.Invoice.PONumber)
	}
	fmt.Printf("Invoice date:          %s\n", id.Invoice.InvoiceDate.Time.Format("Mon, 2 Jan 2006 15:04:05 MST"))
	fmt.Printf("Terms:                 %s\n", id.Invoice.Terms)
	fmt.Printf("Total amount:          $%.2f\n", id.Invoice.TotalAmount)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_RIGHT)

	table.SetHeader([]string{
		"Description", "Location", "Quantity", "Unit Cost", "Total",
	})

	for _, li := range id.Invoice.LineItems {
		table.Append([]string{
			li.Description, li.ServiceLocation, fmt.Sprint(li.Quantity),
			fmt.Sprintf("%.2f", li.UnitCost), fmt.Sprintf("%.2f", li.ItemTotal),
		})
	}
	table.Render()

	if *itemDetails {
		fmt.Println("\n\nIndividual details:")
		for _, li := range id.Invoice.LineItems {
			if len(li.ItemDetails) != 0 {
				fmt.Printf("%s ($%.2f total):\n", li.Description, li.ItemTotal)
				for _, det := range li.ItemDetails {
					fmt.Printf("\t%-20.20s $%.2f\n", det.Description, det.Cost)
				}
			}
		}
	}
}
