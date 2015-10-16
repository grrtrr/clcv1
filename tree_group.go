/*
 * Directory-tree like representation of Hardware Groups and Servers.
 */
package clcv1

import (
	"fmt"
)

type GroupNode struct {
	// The Hardware Group that this tree node wraps
	*HardwareGroup

	// Link back to the upper level (nil if root node)
	Parent		*GroupNode

	// Folder/leaf elements - servers
	Servers		[]*Server

	// Children - subgroups contained within this group
	Children	[]*GroupNode
}

// Represent the hardware groups in @location as a tree.
// @location:  The data center location to query for groups.
// @acctAlias: The alias of the account that owns the groups.
// @servers:   If true, also populate the @Servers field of each group.
func (c *Client) GetGroupHierarchy(location, acctAlias string, servers bool) (root *GroupNode, err error) {
	var uuidMap = make(map[string]*GroupNode)

	if hwgroups, err := c.GetGroups(location, acctAlias); err != nil {
		return nil, err
	} else {
		for i := range hwgroups {
			uuidMap[hwgroups[i].UUID] = &GroupNode{ HardwareGroup: &hwgroups[i] }
		}
	}

	for _, grp := range uuidMap {
		if grp.ParentUUID == "" {
			root = grp
		} else if parent, ok := uuidMap[grp.ParentUUID]; !ok {
			return nil, fmt.Errorf("No parent found for non-root group %q at %s", grp.Name, location)
		} else {
			grp.Parent      = parent
			parent.Children = append(parent.Children, grp)
		}
	}
	if root == nil {
		return nil, fmt.Errorf("No root hardware group found at %s", location)
	}

	if servers {
		srv, err := c.GetAllServers(acctAlias, "", location)
		if err != nil {
			return nil, fmt.Errorf("Failed to look up servers at %s: %s", location, err)
		}
		for i := range srv {
			if group, ok := uuidMap[srv[i].HardwareGroupUUID]; !ok {
				return nil, fmt.Errorf("Failed to look up HW Group UUID for %s", srv[i].Name)
			} else {
				group.Servers = append(group.Servers, &srv[i])
			}
		}
	}
	return
}

// Do a depth-first traversal of the tree to find a specific node.
// @root:  where to start the search at
// @found: function indicating whether the passed GroupNode is the one looked for
func FindGroupNode(root *GroupNode, found func(*GroupNode) bool) *GroupNode {
	if found(root) {
		return root
	}
	for _, c := range root.Children {
		if node := FindGroupNode(c, found); node != nil {
			return node
		}
	}
	return nil
}

// Print group hierarchy starting at @g, using initial indentation @indent.
func PrintGroupHierarchy(g *GroupNode, indent string) {
	var groupLine string

	if g.IsSystemGroup && g.ParentUUID != "" {
		groupLine = fmt.Sprintf("%s[%s]/", indent, g.Name)
	} else {
		groupLine = fmt.Sprintf("%s%s/", indent, g.Name)
	}
	fmt.Printf("%-70s %s\n", groupLine, g.UUID)

	for _, s := range g.Servers {
		fmt.Printf("%s", indent + "    ")
		if s.PowerState == "Started" {
			fmt.Printf("*%s\n", s.Name)
		} else if g.IsSystemGroup && g.Name == "Templates" {
			fmt.Printf("%s\t%s\n", s.Name, s.Description)
		} else {
			fmt.Printf("%s\n", s.Name)
		}
	}

	for _, c := range g.Children {
		PrintGroupHierarchy(c, indent + "    ")
	}
}
