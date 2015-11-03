package clcv1

import (
	"fmt"
)

/*
 * Hardware Groups
 */
type HardwareGroup struct {
	// The legacy ID of the Group.
	// Deprecated. Not available after May 6, 2015. Use UUID instead.
	ID		int

	// The unique identifier of the Group.
	UUID		string

	// The legacy ID of the parent Group.
	// Deprecated. Not available after May 6, 2015. Use ParentUUID instead.
	ParentID	int

	// The unique identifier of the Parent Group.
	ParentUUID	string

	// The name of the Group.
	Name		string

	// Denotes a required system Group.
	IsSystemGroup	bool
}

// Gets a list of all groups with the specified search criteria.
// @location:  The data center location to query for groups.
// @acctAlias: The alias of the account that owns the groups.
func (c *Client) GetGroups(location, acctAlias string) (hwgroups []HardwareGroup, err error) {
	req := struct { AccountAlias, Location string } { acctAlias, location }
	err = c.getResponse("/Group/GetGroups/JSON", &req, &struct {
		BaseResponse
		HardwareGroups	*[]HardwareGroup
	} { HardwareGroups: &hwgroups })
	return
}

// Look for all Hardware Groups satisfying a given criterion
// @location:  The data center location to query for groups.
// @acctAlias: The alias of the account that owns the groups.
// @found:     Function returning true if the passed Hardware Group qualifies
// Returns pointer to HardwareGroup, nil if not found; or error on failure.
func (c *Client) GetGroupsFiltered(location, acctAlias string, found func(*HardwareGroup) bool) (res []HardwareGroup, err error) {
	if groups, err := c.GetGroups(location, acctAlias); err == nil {
		for idx := range groups {
			if found(&groups[idx]) {
				res = append(res, groups[idx])
			}
		}
	}
	return
}

// Look for a (the first) Hardware Group satisfying a given criterion
// @location:  The data center location to query for groups.
// @acctAlias: The alias of the account that owns the groups.
// @found:     Function returning true if the passed Hardware Group is the one looked for.
// Returns pointer to HardwareGroup, nil if not found; or error on failure.
func (c *Client) GetGroup(location, acctAlias string, found func(*HardwareGroup) bool) (res *HardwareGroup, err error) {
	if groups, err := c.GetGroupsFiltered(location, acctAlias, found); err == nil {
		if len(groups) == 1 {
			res = &groups[0]
		} else if len(groups) > 1 {
			return nil, fmt.Errorf("Ambiguous - %d matching groups found at %s", len(groups), location)
		}
	}
	return
}

// Retrieves the root Hardware Group for a given @location
func (c *Client) GetRootGroup(location string) (*HardwareGroup, error) {
	return c.GetGroup(location, "", func(g *HardwareGroup) bool {
		return g.IsSystemGroup && g.ParentUUID == ""
	})
}

// Look up Hardware Group by @name and @location
func (c *Client) GetGroupByName(name, location, acctAlias string) (*HardwareGroup, error) {
	return c.GetGroup(location, "", func(g *HardwareGroup) bool { return g.Name == name })
}

// Look up Hardware Group by @uuid and @location
// The @location is required, since there is no global 'resolveGroup(uuid)' function.
func (c *Client) GetGroupByUUID(uuid, location, acctAlias string) (*HardwareGroup, error) {
	return c.GetGroup(location, "", func(g *HardwareGroup) bool { return g.UUID == uuid })
}

// Create a new Hardware Group.
// @acctAlias:  The alias of the account to owns the group. Can be the parent alias or a sub-account alias.
// @uuid: The unique identifier of the parent group.
// @name: The name of the Hardware Group.
// @desc: A description of the Hardware Group. If none is supplied, the Name will be used.
func (c *Client) CreateHardwareGroup(acctAlias, uuid, name, desc string) (g HardwareGroup, err error) {
	req := struct { AccountAlias,  ParentUUID, Name, Description string } {
		acctAlias, uuid, name, desc,
	}
	err = c.getResponse("/Group/CreateHardwareGroup/JSON", &req, &struct {
		BaseResponse
		Group	     *HardwareGroup
	} { Group: &g })
	return
}

// Enable or disable Maintenance Mode on a Hardware Group.
// @enable:    Turn maintenance mode on or off.
// @uuid:      The unique identifier of the Hardware Group.
// @acctAlias: The alias of the account that owns the group (optional).
func (c *Client) HardwareGroupMaintenance(enable bool, uuid, acctAlias string) (reqId int, err error) {
	req := struct {
		AccountAlias, UUID 	string
		Enable			bool
	} { acctAlias, uuid, enable }
	err = c.getResponse("/Group/HardwareGroupMaintenance/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

// Pause the Hardware Group along with all child groups and servers.
// @uuid:      The unique identifier of the Hardware Group to pause.
// @acctAlias: The alias of the account that owns the group.
func (c *Client) PauseHardwareGroup(uuid, acctAlias string) (reqId int, err error) {
	req := struct { AccountAlias, UUID string } { acctAlias, uuid }
	err = c.getResponse("/Group/PauseHardwareGroup/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

// Power on the Hardware Group along with all child groups and servers.
// @uuid:      The unique identifier of the Hardware Group to power on.
// @acctAlias: The alias of the account that owns the group.
func (c *Client) PowerOnHardwareGroup(uuid, acctAlias string) (reqId int, err error) {
	req := struct { AccountAlias, UUID string } { acctAlias, uuid }
	err = c.getResponse("/Group/PowerOnHardwareGroup/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

// Power the Hardware Group off, along with all child groups and servers.
// @uuid:      The unique identifier of the Hardware Group to power off.
// @acctAlias: The alias of the account that owns the group.
func (c *Client) PowerOffHardwareGroup(uuid, acctAlias string) (reqId int, err error) {
	req := struct { AccountAlias, UUID string } { acctAlias, uuid }
	err = c.getResponse("/Group/PowerOffHardwareGroup/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

// Shut down the Hardware Group along with all child groups and servers.
// @uuid:      The unique identifier of the Hardware Group to shut down.
// @acctAlias: The alias of the account that owns the group.
func (c *Client) ShutdownHardwareGroup(uuid, acctAlias string) (reqId int, err error) {
	req := struct { AccountAlias, UUID string } { acctAlias, uuid }
	err = c.getResponse("/Group/ShutdownHardwareGroup/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

// Reboot the Hardware Group along with all child groups and servers.
// @uuid:      The unique identifier of the Hardware Group to reboot.
// @acctAlias: The alias of the account that owns the group.
func (c *Client) RebootHardwareGroup(uuid, acctAlias string) (reqId int, err error) {
	req := struct { AccountAlias, UUID string } { acctAlias, uuid }
	err = c.getResponse("/Group/RebootHardwareGroup/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

// Reset the Hardware Group along with all child groups and servers.
// @uuid:      The unique identifier of the Hardware Group to reboot.
// @acctAlias: The alias of the account that owns the group.
func (c *Client) ResetHardwareGroup(uuid, acctAlias string) (reqId int, err error) {
	req := struct { AccountAlias, UUID string } { acctAlias, uuid }
	err = c.getResponse("/Group/ResetHardwareGroup/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

// Archive all Servers in the Group and then archive the group.
// @uuid:      The unique identifier of the Hardware Group to archive.
// @acctAlias: The alias of the account that owns the group (optional).
func (c *Client) ArchiveHardwareGroup(uuid, acctAlias string) (reqId int, err error) {
	req := struct { AccountAlias, UUID string } { acctAlias, uuid }
	err = c.getResponse("/Group/ArchiveHardwareGroup/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

// Restore an archived Hardware Group.
// @uuid:       The unique identifier of the Hardware Group to restore.
// @parentUuid: The unique identifier of the hardware group to become the restored group's parent.
// @acctAlias:  The alias of the account that owns the server (optional).
func (c *Client) RestoreHardwareGroup(uuid, parentUuid, acctAlias string) (reqId int, err error) {
	req := struct { AccountAlias, UUID, ParentUUID string } { acctAlias, uuid, parentUuid }
	err = c.getResponse("/Group/RestoreHardwareGroup/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}

// Delete the Hardware Group along with all child groups and servers.
// @uuid:     The unique identifier of the Hardware Group to delete.
// @accAlias: The alias of the account that owns the group.
//            If not provided it will assume the account to which the API user is mapped.
//            Providing this value gives you the ability to access groups in your sub accounts.
func (c *Client) DeleteHardwareGroup(uuid, acctAlias string) (reqId int, err error) {
	req := struct { AccountAlias, UUID string } { acctAlias, uuid }
	err = c.getResponse("/Group/DeleteHardwareGroup/JSON", &req, &struct {
		BaseResponse
		RequestID	*int
	} { RequestID: &reqId })
	return
}
