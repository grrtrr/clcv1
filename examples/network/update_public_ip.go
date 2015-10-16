/*
 * Configure firewall settings of a public IP Address.
 */
package main

import (
	"github.com/grrtrr/clcv1"
	"github.com/grrtrr/exit"
	"path"
	"flag"
	"log"
	"fmt"
	"os"
)

func main() {
	var acctAlias = flag.String("a",     "", "Account alias to use")
	var ipAddr    = flag.String("i",     "", "The public IP address to modify")

	var http      = flag.Bool("http",    false, "Allow HTTP requests (port 80) on the new IP")
	var http8080  = flag.Bool("httpAlt", false, "Allow HTTP requests (port 8080) on the new IP")
	var https     = flag.Bool("https",   false, "Allow HTTPS requests (port 443) on the new IP")
	var ftp       = flag.Bool("ftp",     false, "Allow FTP requests (port 21) on the new IP")
	var ftps      = flag.Bool("ftps",    false, "Allow FTPS requests (port 990) on the new IP")
	var ssh       = flag.Bool("ssh",     true,  "Allow SSH requests (port 22) on the new IP")
	var sftp      = flag.Bool("sftp",    true,  "Allow SFTP requests (port 22) on the new IP")
	var rdp       = flag.Bool("rdp",     false, "Allow RDP requests (port 3389) on the new IP")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <server-name>\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()
	/*
	 * Only server-name and public IP address are required arguments.
	 */
	if flag.NArg() != 1 || *ipAddr == "" {
		flag.Usage()
		os.Exit(1)
	}

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	req := clcv1.UpdatePublicIPAddressReq{
		// The name of the server.
		ServerName:          flag.Arg(0),

		// The public, mapped IP to manage
		PublicIPAddress:     *ipAddr,

		// The alias of the account that owns the server.
		AccountAlias:	     *acctAlias,

		// The public IP mapping will allow HTTP requests.
		AllowHTTP:           *http,

		// The public IP mapping will allow HTTP requests on port 8080.
		AllowHTTPonPort8080: *http8080,

		// The public IP mapping will allow HTTPS requests.
		AllowHTTPS:          *https,

		// The public IP mapping will allow FTP requests.
		AllowFTP:            *ftp,

		// The public IP mapping will allow FTPS requests.
		AllowFTPS:           *ftps,

		// The public IP mapping will allow SFTP requests.
		AllowSFTP:           *sftp,

		// The public IP mapping will allow SSH requests.
		AllowSSH:            *ssh,

		// The public IP mapping will allow RDP requests.
		AllowRDP: *          rdp,
	}

	reqId, err := client.UpdatePublicIPAddress(&req)
	if err != nil {
		exit.Fatalf("Failed to modify public IP %s on %q: %s", *ipAddr, flag.Arg(0), err)
	}

	fmt.Printf("Request ID for modifying public IP %s on %s: %d\n", *ipAddr, flag.Arg(0), reqId)
}
