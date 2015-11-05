package clcv1

import (
	"fmt"
)

type UserRole int

const (
	ServerAdministrator  UserRole =  2
	BillingManager                =  3
	DNSManager                    =  8
	AccountAdministrator          =  9
	AccountViewer                 = 10
	NetworkManager                = 12
	SecurityManager               = 13
	ServerOperator                = 14
	ServerScheduler               = 15
)

func (r UserRole) String() string {
	switch r {
	case ServerAdministrator:
		return "Server-Administrator"
	case BillingManager:
		return "Billing-Manager"
	case DNSManager:
		return "DNS Manager"
	case AccountAdministrator:
		return "Account-Administrator"
	case AccountViewer:
		return "Account-Viewer"
	case NetworkManager:
		return "Network-Manager"
	case SecurityManager:
		return "Security-Manager"
	case ServerOperator:
		return "Server-Operator"
	case ServerScheduler:
		return "Server-Scheduler"
	default:
		panic(fmt.Errorf("Unknown UserRole %d", int(r)))
	}
}

type User struct {
	// Short code of account.
	AccountAlias		string

	// User name, which is typically the email address.
	UserName		string

	// Email address of the user.
	EmailAddress		string

	// First name of the user.
	FirstName		string

	// Last name of the user.
	LastName		string

	// Name used for single-sign-on process.
	SAMLUserName		string

	// Additional email address for the user.
	AlternateEmailAddress	string

	// Job title of the user.
	Title			string

	// Office phone number of the user.
	OfficeNumber		string

	// Mobile phone number of the user.
	MobileNumber		string

	// Flag indicating whether this user can receive SMS messages.
	AllowSMS		bool

	// Fax number for the user.
	FaxNumber		string

	// Time zone that the user resides in.
	TimeZoneID		string

	// List of values indicating the roles assigned to this user.
	Roles			[]UserRole

	// FIXME: as of Sep 29/2015, the output also contains
	//Status		string	// e.g., "Status": "Active"
}

// Get all users assigned to a given account.
// @accountAlias: account alias to use (required)
func (c *Client) GetUsers(accountAlias string) (users []User, err error) {
	req := struct { AccountAlias string } { accountAlias }
	err = c.getResponse("/User/GetUsers/JSON", &req, &struct {
		BaseResponse
		Users	*[]User
	} { Users: &users })
	return
}

// Get the details of a specific user associated with a given account.
// @userName:     user name to query (typically the email address, required)
// @accountAlias: account alias to use (required)
func (c *Client) GetUserDetails(userName, accountAlias string) (details User, err error) {
	req := struct { AccountAlias, UserName string } { accountAlias, userName }
	err = c.getResponse("/User/GetUserDetails/JSON", &req, &struct {
		BaseResponse
		UserDetails	*User
	} { UserDetails: &details })
	return
}
