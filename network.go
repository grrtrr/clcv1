package clcv1

/*
 * Network Listing
 */
type Network struct {
	// The name of the Network.
	// This value should be used on other API methods that require a network reference.
	Name		string

	// The friendly name of the network if one has been configured.

	Description	string
	// The default gateway for the network.
	Gateway		string

	// The Network's home datacenter alias.
	Location	string

	// The alias of the account that owns the network.
	AccountAlias	string
}

// Get the list of deployable Networks mapped to an account, in any Data Center.
// @acctAlias: The alias of the account that owns the network (optional).
//             If not provided it will assume the Account the API user is mapped to.
//             Providing this value gives you the ability to get networks in your sub accounts.
// @location:  The Network's home datacenter alias.
//             If blank, the account home datacenter location is assumed.
func (c *Client) GetDeployableNetworks(acctAlias, location string) (nets []Network, err error) {
	req := struct{ AccountAlias, Location string } { acctAlias, location }
	err = c.getResponse("/Network/GetDeployableNetworks/JSON", &req, &struct {
		BaseResponse
		Networks	*[]Network
	} { Networks: &nets })
	return
}

// Get the list of Networks mapped to an account, in any Data Center.
// @acctAlias: The alias of the account that owns the network (optional).
// @location:  The Network's home datacenter alias (optional).
func (c *Client) GetAccountNetworks(acctAlias, location string) (nets []Network, err error) {
	req := struct{ AccountAlias, Location string } { acctAlias, location }
	err = c.getResponse("/Network/GetAccountNetworks/JSON", &req, &struct {
		BaseResponse
		Networks	*[]Network
	} { Networks: &nets })
	return
}

// Get the list of Networks mapped to the account in its Primary Data Center.
// NOTE: The API docs mention acctAlias/location arguments, but their use did not show
//       any influence on the result (which referred to account default and home data
//       center). This command seems to be subsumed by GetAccountNetworks()
func (c *Client) GetNetworks() (nets []Network, err error) {
	err = c.getResponse("/Network/GetNetworks/JSON", nil, &struct {
		BaseResponse
		Networks	*[]Network
	} { Networks: &nets })
	return
}

/*
 * Network Details
 */
type NetworkDetails struct {
	// The Network name.
	Name		string

	// The friendly name of the network if one has been configured.
	Description	string

	// The default gateway for the network.
	Gateway		string

	// The network mask.
	NetworkMask	string

	// The Network's home datacenter alias.
	Location	string

	// A list of IPAddress objects:
	IPAddresses	[]struct {
		// The IP Address.
		Address		string

		// The type of the IP Address, one of:
		//   RIP - Real IP (internal IP configured on the VLAN)
		//   MIP - Mapped IP (external IP configured on the Firewall)
		//   VIP - Virtual IP (external IP configured on the Load Balancer)
		AddressType	string

		// Indicates if the address is claimed or available.
		IsClaimed	bool

		// The name of the Server using this IP address, if applicable.
		ServerName	string
	}
}

// Get the details for a Network and its IP Addresses.
// @name:      The Network name (required)
// @acctAlias: The alias of the account that owns the network (optional).
// @location:  The Network's home datacenter alias (optional).
func (c *Client) GetNetworkDetails(name, acctAlias, location string) (details NetworkDetails, err error) {
	req := struct{ Name, AccountAlias, Location string } { name, acctAlias, location }
	err = c.getResponse("/Network/GetNetworkDetails/JSON", &req, &struct {
		BaseResponse
		NetworkDetails *NetworkDetails
	} { NetworkDetails: &details })
	return
}


/*
 * Management of Public IP Addresses
 */
type AddPublicIPAddressReq struct {
	// The name of the server to add the IP address to (only required argument).
	ServerName	string

	/*
	 *	Optional Arguments
	 */
	// The alias of the account that owns the server.
	AccountAlias		string

	// An existing internal IP Address on the server to use for the mapping.
	// Leaving this blank will assign a new internal IP Address.
	IPAddress		string

	// The existing password, for authentication.
	// Required only if existing internal IP was not provided and new IP must be assigned.
	ServerPassword		string

	// The public IP mapping will allow HTTP requests.
	AllowHTTP		bool

	// The public IP mapping will allow HTTP requests on port 8080.
	AllowHTTPonPort8080	bool

	// The public IP mapping will allow HTTPS requests.
	AllowHTTPS		bool

	// The public IP mapping will allow FTP requests.
	AllowFTP		bool

	// The public IP mapping will allow FTPS requests.
	AllowFTPS		bool

	// The public IP mapping will allow SFTP requests.
	AllowSFTP		bool

	// The public IP mapping will allow SSH requests.
	AllowSSH		bool

	// The public IP mapping will allow RDP requests.
	AllowRDP		bool
}

// Map a public IP Address to a Server
func (c *Client) AddPublicIPAddress(req *AddPublicIPAddressReq) (reqId int, err error) {
	err = c.getResponse("/Network/AddPublicIPAddress/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

type UpdatePublicIPAddressReq struct {
	// The name of the server.
	ServerName		string

	// The public, mapped IP to manage
	PublicIPAddress		string

	/*
	 * Optional Arguments
	 */
	// The alias of the account that owns the server.
	AccountAlias		string

	// The public IP mapping will allow HTTP requests.
	AllowHTTP		bool

	// The public IP mapping will allow HTTP requests on port 8080.
	AllowHTTPonPort8080	bool

	// The public IP mapping will allow HTTPS requests.
	AllowHTTPS		bool

	// The public IP mapping will allow FTP requests.
	AllowFTP		bool

	// The public IP mapping will allow FTPS requests.
	AllowFTPS		bool

	// The public IP mapping will allow SFTP requests.
	AllowSFTP		bool

	// The public IP mapping will allow SSH requests.
	AllowSSH		bool

	// The public IP mapping will allow RDP requests.
	AllowRDP		bool
}

// Configure firewall settings on a public IP Address.
func (c *Client) UpdatePublicIPAddress(req *UpdatePublicIPAddressReq) (reqId int, err error) {
	err = c.getResponse("/Network/UpdatePublicIPAddress/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}
