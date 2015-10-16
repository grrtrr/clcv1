package clcv1

import (
	"github.com/grrtrr/clcv1/microsoft"
)

// Server Object
type Server struct {
	// The full name of the Server.
	Name			string

	// The unique identifier of the containing Group.
	HardwareGroupUUID	string

	// The description of the Server as provided on creation.
	Description		string

	// Modification origin
	ModifiedBy		string

	// Modification date
	DateModified		microsoft.Timestamp

	// The DNS name of the Server.
	DnsName			string

	// The number of processors configured on the Server.
	Cpu			int

	// Total GB of RAM configured on the Server.
	MemoryGB		int

	// Total number of disks configured on the Server.
	DiskCount		int

	// Total space across all disk configured on the Server.
	TotalDiskSpaceGB	int

	// True if the Server is a template, else False.
	IsTemplate		bool

	// True if the Server is a Hyperscale instance, else False.
	IsHyperscale		bool

	// Active, Archived, Deleted, UnderConstruction, QueuedForArchive, QueuedForDelete, or QueuedForRestore
	Status			string

	// The type of server. Standard or Premium
	ServerType		int

	// The service level/performance for the underlying data store.  Standard or Premium
	ServiceLevel		int

	// Operating System of the server (see below).
	OperatingSystem		OperatingSystem

	// StringThe current power state of the Server (Stopped, Started, Paused).
	PowerState		string

	// Indicates if the Server is in Maintenance Mode.
	InMaintenanceMode	bool

	// Home datacenter of the Server.
	Location		string

	// The primary IP address of the Server.
	IPAddress		string

	// A list of all IP Addresses assigned to the server (see below)
	IPAddresses		[]IPAddress

	// A list of Custom Fields associated to this server (see below)
	CustomFields		[]CustomField

	// The ID of the Server. DEPRECATED - Value is -1.
	ID			int

	// The legacy ID of the containing Group.
	// DEPRECATED - Not available after May 6, 2015. Use UUID instead.
	HardwareGroupID		int
}


type IPAddress struct {
	// The IP Address
	Address		string

	// The type of the IP Address:
	//   RIP - Real IP    (internal IP configured on the VLAN)
	//   MIP - Mapped IP  (external IP configured on the Firewall)
	//   VIP - Virtual IP (external IP configured on the Load Balancer)
	AddressType	string
}

func (i IPAddress) String() string {
	return i.Address + "/" + i.AddressType
}

func (i *IPAddress) IsPublic() bool {
	return i.AddressType == "MIP" || i.AddressType == "VIP"
}

type CustomField struct {
	// Unique identifier that is associated with the Account Custom Field.
	// Call Account/GetCustomFields for a list of all custom fields set at the account level.
	ID		string

	// Name for the Custom Field.
	Name		string

	// Type of custom field: Text, Option or Checkbox.
	Type		string

	// For Text: Any value;
	// For Option values, call Account/GetCustomFields to see possible values to pass in.
	// Checkbox values should be "true" or "false".
	Value		string

	// Deprecated. Value is -1. Use CustomFieldType instead.
	CustomFieldID	int
}

/*
 * Server Lists
 */
// List all servers of a given Hardware Group.
// @hwGrpUUID: The unique identifier of the Hardware Group.
// @acctAlias:
func (c *Client) GetServers(hwGrpUUID, acctAlias string) (servers []Server, err error) {
	req := struct { AccountAlias, HardwareGroupUUID string } { acctAlias, hwGrpUUID }
	err = c.getResponse("/Server/GetServers/JSON", &req, &struct {
		BaseResponse
		Servers		*[]Server
	} { Servers: &servers })
	return
}

// Get a deep list of all Servers for a given Hardware Group and its sub groups,
// or all Servers for a given location.
// All of the following arguments are optional (empty string):
// @acctAlias: The alias of the account that owns the servers.
// @hwGrpUUID: The unique identifier of the containing Hardware Group.
// @location:  The data center location.
func (c *Client) GetAllServers(acctAlias, hwGrpUUID, location string) (servers []Server, err error) {
	req := struct {
		AccountAlias, HardwareGroupUUID, Location string
	} { acctAlias, hwGrpUUID, location }
	err = c.getResponse("/Server/GetAllServers/JSON", &req, &struct {
		BaseResponse
		Servers		*[]Server
	} { Servers: &servers })
	return
}

// Like GetAllServers, but with additional begin/end dates added to query.
// @acctAlias: see GetAllServers
// @hwGrpUUID: see GetAllServers
// @location:  see GetAllServers
// @beginDate: Beginning of date range for querying modified servers.
//             Can be a partial DateTime (e.g. 2013-05-10) or a full DateTime (e.g. 2013:05-10T14:30:12).
//             If date is missing, then the value equals today minus one day.
// @endDate:   End of date range for querying modified servers. Can be a partial DateTime (e.g. 2013-05-10)
//             or a full DateTime (e.g. 2013:05-10T14:30:12).
//             If date is missing, then the value is set to the current date time.
func (c *Client) GetAllServersByModifiedDates(acctAlias, hwGrpUUID, location string,
					      beginDate, endDate string) (servers []Server, err error) {
	req := struct {
		AccountAlias, HardwareGroupUUID, Location	string
		BeginDate, EndDate				string
	} { acctAlias, hwGrpUUID, location, beginDate, endDate }
	err = c.getResponse("/Server/GetAllServersByModifiedDates/JSON", &req, &struct {
		BaseResponse
		Servers		*[]Server
	} { Servers: &servers })
	return
}

type AccountServer struct {
	AccountAlias	string
	Servers		[]Server
}

// Get a deep list of all servers for a given account hierarchy within a given data center.
// Use this operation to get a full list of all servers contained within an account and all its subaccounts.
// @acctAlias: The alias of the account that contains sub-accounts with servers (optional).
// @location:  The data center location. Otherwise leave blank and let it default to account's primary data center.
func (c *Client) GetAllServersForAccountHierarchy(acctAlias, location string) (servers []AccountServer, err error) {
	req := struct {
		AccountAlias, Location	string
	} { acctAlias, location }
	err = c.getResponse("/Server/GetAllServersForAccountHierarchy/JSON", &req, &struct {
		BaseResponse
		AccountServers	*[]AccountServer
	} { AccountServers: &servers })
	return
}

/*
 * Server Templates
 */
type ServerTemplate struct {
	// The ID of the Server.
	// Deprecated. Value is -1.
	ID			int

	// The full name of the Server.
	Name			string

	// The description of the Serve.
	Description		string

	// The number of processors.
	Cpu			int

	// Total GB of RAM configured on the Server.
	MemoryGB		int

	// Total number of disks configured on the Server.
	DiskCount		int

	// Total space across all disk configured on the Server.
	TotalDiskSpaceGB	int

	// Operating System of the server (see below).
	OperatingSystem		OperatingSystem

	// Home datacenter of the template.
	Location		string
}

// Gets the list of Templates available to the account (default location)
func (c *Client) GetServerTemplates() (templates []ServerTemplate, err error) {
	err = c.getResponse("/Server/GetServerTemplates/JSON", nil, &struct {
		BaseResponse
		Templates	*[]ServerTemplate
	} { Templates: &templates })
	return
}


// Gets the list of Templates available to the account and location.
// @acctAlias: The alias of the account that owns the server.
//             If not provided it will assume the account to which the API user is mapped.
//             Providing this value gives you the ability to access servers in your sub accounts.
// @location:  The data center of the server templates.
func (c *Client) ListAvailableServerTemplates(acctAlias, location string) (templates []ServerTemplate, err error) {
	req := struct { AccountAlias, Location string } { acctAlias, location }
	err = c.getResponse("/Server/ListAvailableServerTemplates/JSON", &req, &struct {
		BaseResponse
		Templates	*[]ServerTemplate
	} { Templates: &templates })
	return
}

// Convert the server to a template.
// @name:       The name of the Server.
// @password:   The administrator/root password for the server to convert.
// @templAlias: The alias for the Template to create.
// @acctAlias:  The alias of the account that owns the server (optional).
func (c *Client) ConvertServerToTemplate(name, password, templAlias, acctAlias string) (reqId int, err error) {
	req := struct {
		Name, AccountAlias	string
		Password		string
		TemplateAlias		string
	} { name, acctAlias, password, templAlias }
	err = c.getResponse("/Server/ConvertServerToTemplate/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

// Convert the template to a server.
// @name:       The name of the Template.
// @password:   The new administrator/root password for the converted server.
// @hwGrpUUID:  The unique identifier of the hardware group to add the converted server to.
// @network:    The name of the network to add the converted server to.
// @acctAlias:  The alias of the account that owns the server (optional).
func (c *Client) ConvertTemplateToServer(name, password, hwGrpUUID, network, acctAlias string) (reqId int, err error) {
	req := struct {
		Name, AccountAlias	string
		Password		string
		HardwareGroupUUID	string
		Network			string
	} { name, acctAlias, password, hwGrpUUID, network }
	err = c.getResponse("/Server/ConvertTemplateToServer/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

// Delete the Template with the specified name.
// @name:      The name of the Template to delete.
// @acctAlias: The alias of the account that owns the template (optional).
func (c *Client) DeleteTemplate(name, acctAlias string) (reqId int, err error) {
	req := struct {	Name, AccountAlias string } { name, acctAlias }
	err = c.getResponse("/Server/DeleteTemplate/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

/*
 * Invidual Servers
 */

// Gets the detail for one server.
// @name:      The name of the server
// @acctAlias: The alias of the account that owns the servers (optional)
func (c *Client) GetServer(name, acctAlias string) (server Server, err error) {
	req := struct { AccountAlias, Name string } { acctAlias, name }
	err = c.getResponse("/Server/GetServer/JSON", &req, &struct {
		BaseResponse
		Server		*Server
	} { Server: &server })
	return
}

/*
 * Credentials
 */
type ServerCredentials struct {
	// The administrator or root user name for the server.
	Username	string

	// The password associated with the account.
	Password	string
}


// Get the credentials for the specified server.
// @name:      The name of the server
// @acctAlias: The alias of the account that owns the servers (optional)
func (c *Client) GetServerCredentials(name, acctAlias string) (creds ServerCredentials, err error) {
	req := struct { AccountAlias, Name string } { acctAlias, name }
	err = c.getResponse("/Server/GetServerCredentials/JSON", &req, &struct {
		BaseResponse
		*ServerCredentials
	} { ServerCredentials: &creds })
	return
}
// Update the Admin/Root password for a Server.
// @name:      The name of the server
// @acctAlias: The alias of the account that owns the servers (optional)
// @curPass:   The existing password, for authentication.
// @newPass:   The new password to apply.
func (c *Client) ServerChangePassword(name, acctAlias, curPass, newPass string) error {
	req := struct {
		AccountAlias, Name	string
		CurrentPassword		string
		NewPassword		string
	} { acctAlias, name, curPass, newPass }
	return c.getResponse("/Server/ChangePassword/JSON", &req, new(BaseResponse))
}

/*
 * Disk-related Functions
 */
type DiskInfo struct {
	// Name assigned to the disk.
	Name		string

	// The SCSI bus ID of the disk.
	ScsiBusID	string

	// The SCSI device ID of the disk.
	ScsiDeviceID	string

	// Size of the disk in GB
	SizeGB		int
}

// List the disks on a Server.
// @name:      The name of the server
// @acctAlias: The alias of the account that owns the servers (optional)
// @diskNames: Set to true in order to retreive disk mount points / drive letters
func (c *Client) ListDisks(name, acctAlias string, diskNames bool) (hasSnapshot bool, disks []DiskInfo, err error) {
	req := struct {
		AccountAlias, Name	string
		QueryGuestDiskNames	bool
	} { acctAlias, name, diskNames }
	err = c.getResponse("/Server/ListDisks/JSON", &req, &struct {
		BaseResponse
		Server		string
		HasSnapshot	*bool
		Disks		*[]DiskInfo
	} { HasSnapshot: &hasSnapshot, Disks: &disks })
	return
}

// Delete a disk from a Server.
// @name:      The name of the server
// @acctAlias: The alias of the account that owns the servers (optional)
// @busId:     The SCSI bus ID of the disk.
// @devId:     The SCSI device ID of the disk.
// @force:     Set to true to override safety checks that prevent deleting typical primary
//             operating system drives, e.g. SCSI Bus ID 0, SCSI Device ID 0 on Windows (
//             typically C drive) and SCSI Bus ID 0, SCSI Device IDs 0,1,2 on Linux
//             (typically boot, swap and root disks).
func (c *Client) DeleteDisk(name, acctAlias string, busId, devId string, force bool) (reqId int, err error) {
	req := struct {
		AccountAlias, Name	string
		ScsiBusID, ScsiDeviceID	string
		OverrideFailsafes	bool
	} { acctAlias, name, busId, devId, force }
	err = c.getResponse("/Server/DeleteDisk/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

// Resize a disk on a server
// @name:      The name of the server
// @acctAlias: The alias of the account that owns the servers (optional)
// @busId:     The SCSI bus ID of the disk.
// @devId:     The SCSI device ID of the disk.
// @newSizeGB: The expanded size of the disk. Must be greater than the existing disk size.
// @expandFS:  Whether to expand the file system on the disk after the resize.
func (c *Client) ResizeDisk(name, acctAlias string, busId, devId string, newSizeGB int, expandFs bool) (reqId int, err error) {
	req := struct {
		AccountAlias, Name	string
		ScsiBusID, ScsiDeviceID	string
		ResizeGuestDisk		bool
		NewSizeGB		int
	} { acctAlias, name, busId, devId, expandFs, newSizeGB }
	err = c.getResponse("/Server/ResizeDisk/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

/*
 * Create Server
 */
type CreateServerReq struct {
	// The alias of the account to own the server.
	// If not provided it will assume the account to which the API user is mapped.
	// Providing this value gives you the ability to create servers in your sub accounts. (optional)
	AccountAlias		string

	// The alias of the data center in which to create the server (optional).
	// If not provided, will default to the API user's default data center.
	LocationAlias		string

	// The name of the template to create the server from (required)
	Template		string

	// The alias for the server. Limit 6 charcters (required)
	Alias			string

	// An optional description for the server. If none is supplied the server name will be used. (optional)
	Description		string

	// The unique identifier of the Hardware Group to add this server to (required)
	HardwareGroupUUID	string

	//The type of server to create (required)
	// 1 = Standard
	// 2 = Enterprise
	ServerType		int

	// The service level/performance for the underlying data store (required)
	// 1 = Premium
	// 2 = Standard
	ServiceLevel		int

	// The number of processors to configure the server with (required)
	Cpu			int

	// The number of GB of memory to configure the server with (required)
	MemoryGB		int

	// The size in GB of an additional drive to add to the server (required).
	// If no additional drive is needed, pass in a 0 value.
	ExtraDriveGB		int

	// The primary DNS to set on the server (optional)
	// If not supplied the default value set on the account will be used.
	PrimaryDns		string

	// The secondary DNS to set on the server (optional)
	// If not supplied the default value set on the account will be used.
	SecondaryDns		string

	// The name of the network to which to deploy the server.
	// If your account has not yet been assigned a network, leave this blank and one will be assigned automatically.
	// If one or more networks are available, the network name is required.
	Network			string

	// The desired Admin/Root password (optional).
	// Please note the password must meet the password strength policy.
	// Leave blank to have the system generate a password
	Password		string

	// A list of Custom Fields associated to this server
	// ID:    Unique identifier that is associated with the Account Custom Field.
	//        Call Account/GetCustomFields for a list of all custom fields set at the account level.
	// Value:
	//        - For Text: Any value;
	//        - For Option values, call Account/GetCustomFields to see possible values to pass in.
	//        - Checkbox values should be "true" or "false".
	CustomFields            []struct { ID, Value string }
}

// Create a new Server
func (c *Client) CreateServer(req *CreateServerReq) (reqId int, err error) {
	err = c.getResponse("/Server/CreateServer/JSON", req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

/*
 * Configure Server
 */
type ConfigureServerReq struct {
	// The name of the server.
	Name			string

	// The unique identifier of the Hardware Group that this server belongs to (required).
	HardwareGroupUUID	string

	// The alias of the account to own the server (optional)
	AccountAlias		string

	// The number of processors to configure the server with (required).
	Cpu			int

	// The number of GB of memory to configure the server with (required).
	MemoryGB		int

	// The size in GB of an additional drive to add to the server (required).
	// If no additional drive is needed, pass in a 0 value.
	AdditionalStorageGB	int

	// A list of Custom Fields associated to this server.
	// See field of identical name in above CreateServerReq
	CustomFields            []struct { ID, Value string }
}

// Configure the CPU, Memory, Group and additional storage for a Server.
func (c *Client) ConfigureServer(req *ConfigureServerReq) (reqId int, err error) {
	err = c.getResponse("/Server/ConfigureServer/JSON", req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

// Power server on (or resume from a paused state).
// @name:      The name of the Server to power on.
// @acctAlias: The alias of the account that owns the server.
func (c *Client) PowerOnServer(name, acctAlias string) (reqId int, err error) {
	req := struct { Name, AccountAlias string } { name, acctAlias }
	err = c.getResponse("/Server/PowerOnServer/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

// Pause the server.
// @name:      The name of the Server to pause.
// @acctAlias: The alias of the account that owns the server.
func (c *Client) PauseServer(name, acctAlias string) (reqId int, err error) {
	req := struct { Name, AccountAlias string } { name, acctAlias }
	err = c.getResponse("/Server/PauseServer/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

// Power server off.
// @name:      The name of the Server to power off.
// @acctAlias: The alias of the account that owns the server.
func (c *Client) PowerOffServer(name, acctAlias string) (reqId int, err error) {
	req := struct { Name, AccountAlias string } { name, acctAlias }
	err = c.getResponse("/Server/PowerOffServer/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

// Shut down the operating system and then power off server.
// @name:      The name of the Server to shut down.
// @acctAlias: The alias of the account that owns the server.
func (c *Client) ShutdownServer(name, acctAlias string) (reqId int, err error) {
	req := struct { Name, AccountAlias string } { name, acctAlias }
	err = c.getResponse("/Server/ShutdownServer/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

// Reboot the server (OS reboot)
// @name:      The name of the Server to reboot.
// @acctAlias: The alias of the account that owns the server.
func (c *Client) RebootServer(name, acctAlias string) (reqId int, err error) {
	req := struct { Name, AccountAlias string } { name, acctAlias }
	err = c.getResponse("/Server/RebootServer/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

// Reset server (forced power-cycle).
// @name:      The name of the Server to reset.
// @acctAlias: The alias of the account that owns the server.
func (c *Client) ResetServer(name, acctAlias string) (reqId int, err error) {
	req := struct { Name, AccountAlias string } { name, acctAlias }
	err = c.getResponse("/Server/ResetServer/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

// Enable or disable maintenance mode on a Server.
// @enable:    Turn maintenance mode on or off.
// @name:      The name of the Server.
// @acctAlias: The alias of the account that owns the server (optional).
func (c *Client) ServerMaintenance(enable bool, name, acctAlias string) (reqId int, err error) {
	req := struct {
		Name, AccountAlias	string
		Enable			bool
	} { name, acctAlias, enable }
	err = c.getResponse("/Server/ServerMaintenance/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

// Delete the machine and release all associated resources.
// @name:      The name of the Server to delete.
// @acctAlias: The alias of the account that owns the server.
func (c *Client) DeleteServer(name, acctAlias string) (reqId int, err error) {
	req := struct { Name, AccountAlias string } { name, acctAlias }
	err = c.getResponse("/Server/DeleteServer/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

/*
 * Snapshots
 */
type SnapshotAttribute struct {
	// The full name of the Snapshot.
	Name		string

	// The description of the Snapshot.
	Description	string

	// The time (in UTC) when the Snapshot was created.
	DateCreated	microsoft.Timestamp
}

// Get the list of Snapshots associated with the server.
// @name:      The name of the Server to list snapshots of.
// @acctAlias: The alias of the account that owns the server (optional).
func (c *Client) GetSnapshots(name, acctAlias string) (snaps []SnapshotAttribute, err error) {
	req := struct { Name, AccountAlias string } { name, acctAlias }
	err = c.getResponse("/Server/GetSnapshots/JSON", &req, &struct {
		BaseResponse
		Snapshots	*[]SnapshotAttribute
	} { Snapshots: &snaps })
	return
}


// Take a server snapshot.
// @name:      The name of the Server to snapshot.
// @acctAlias: The alias of the account that owns the server (optional).
func (c *Client) SnapshotServer(name, acctAlias string) (reqId int, err error) {
	req := struct { Name, AccountAlias string } { name, acctAlias }
	err = c.getResponse("/Server/SnapshotServer/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

// Revert to a named snapshot for a specified server.
// @name:      The name of the Server to revert.
// @snapName:  The name of the Snapshot to revert to.
// @acctAlias: The alias of the account that owns the server (optional).
// FIXME: CLC only uses 1 snapshot currently. Use GetSnapshots() to get the single snapshot
//        name (if any), and drop the @snapName argument.
func (c *Client) RevertToSnapshot(name, snapName, acctAlias string) error {
	req := struct { SnapshotName, Name, AccountAlias string } { snapName, name, acctAlias }
	return c.getResponse("/Server/RevertToSnapshot/JSON", &req, new(BaseResponse))
}

// Delete a named snapshot for a specified server
// @snapName:  The name of the Snapshot to delete.
// @name:      The name of the Server to delete snapshot from.
// @acctAlias: The alias of the account that owns the server (optional).
func (c *Client) DeleteSnapshot(snapName, name, acctAlias string) error {
	req := struct { SnapshotName, Name, AccountAlias string } { snapName, name, acctAlias }
	return c.getResponse("/Server/DeleteSnapshot/JSON", &req, new(BaseResponse))
}


/*
 * Archiving Servers
 */
type ArchivedServer struct {
	// The ID of the Server.
	// Deprecated. Value is -1.
	ID		int

	// The full name of the Server.
	Name		string

	// The description of the Server as provided on creation.
	Description	string
}

// List archived servers
// @acctAlias: The alias of the account that owns the server (optional).
// @location:  The data center of the servers (optional).
// Note: dropped the GetArchiveServers() call, it seems deprecated (always returned an empty list).
func (c *Client) ListArchiveServers(acctAlias, location string) (servers []ArchivedServer, err error) {
	req := struct { AccountAlias, Location string } { acctAlias, location }
	err = c.getResponse("/Server/ListArchivedServers/JSON", &req, &struct {
		BaseResponse
		Servers		*[]ArchivedServer
	} { Servers: &servers })
	return
}

// Archive a server
// @name:      The name of the Server to archive.
// @acctAlias: The alias of the account that owns the server (optional).
func (c *Client) ArchiveServer(name, acctAlias string) (reqId int, err error) {
	req := struct { Name, AccountAlias string } { name, acctAlias }
	err = c.getResponse("/Server/ArchiveServer/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

// Restore an archived server.
// @name:      The name of the archived Server.
// @acctAlias: The alias of the account that owns the server (optional).
// @hwGrpUUID: The unique identifier of the hardware group to the restore the server to.
func (c *Client) RestoreServer(name, acctAlias, hwGrpUUID string) (reqId int, err error) {
	 req := struct {
		 Name, AccountAlias	string
		 HardwareGroupUUID	string
	 } { name, acctAlias, hwGrpUUID }
	err = c.getResponse("/Server/RestoreServer/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}
