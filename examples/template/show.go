/*
 * List the details of one template (a server really)
 * Note that the default location for templates is the 'Templates' system folder
 *
 */
package main

import (
	"github.com/olekukonko/tablewriter"
	"github.com/dustin/go-humanize"
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
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <Template Name>\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	template, err := client.GetServer(flag.Arg(0), *acctAlias)
	if err != nil {
		exit.Fatalf("Failed to list details of server %q: %s", flag.Arg(0), err)
	}

	if *simple {
		utils.PrintStruct(template)
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoFormatHeaders(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(false)

		table.SetHeader([]string{
			"Name", "Description", "OS",
			"CPU", "Disk",
			"Last Change",
		})
		table.Append([]string{
			template.Name, template.Description, fmt.Sprint(template.OperatingSystem),
			fmt.Sprint(template.Cpu),
			fmt.Sprintf("%d GB", template.TotalDiskSpaceGB),
			fmt.Sprintf("%s by %s", humanize.Time(template.DateModified.Time), template.ModifiedBy),
		})
		table.Render()
	}
}
