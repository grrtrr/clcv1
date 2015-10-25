package clcv1

import (
_	"fmt"
)

/*
 * Account List
 */
type Account struct {
	// Short name associated with the account.
	AccountAlias	string

	// Short name associated with the parent account.
	ParentAlias	string

	// Data center alias (3 letters).
	// List of data center alias is retrieved from GetLocations operation.
	Location	string

	// Full business name associated with the account.
	BusinessName	string

	// Indicator of whether the account is active or not.
	IsActive	bool

	// Note: online documentation mentions an additional field "SupportLevel" (string).
	//       As of Sep 2015, this field seems no longer to be present.
}

// Get details of the API user's account and any sub-accounts.
func (c *Client) GetAccounts() (accounts []Account, err error) {
	err = c.getResponse("/Account/GetAccounts/JSON", nil, &struct {
		BaseResponse
		Accounts	*[]Account
	} { Accounts: &accounts })
	return
}

/*
 * CustomFields of accounts
 */
type AccountCustomField struct {
	// Deprecated. Value is -1. Use UUID instead.
	ID				int
	// Unique identifier of the Account Custom Field.
	UUID				string
	// Deprecated. Value is -1. Use CustomFieldType instead.
	CustomFieldTypeID		int
	// The type of field: "Text", "Option", or "Checkbox".
	CustomFieldType			string
	//  Friendly name of the Custom Field.
	Name				string
	// Whether or not the Custom Field is required.
	IsRequired			bool
	// Options for the Account Custom Field, specific to the field type.
	AccountCustomFieldOptions	[]struct { Name, Value string }
}

// Gets the account custom field definitions.
func (c *Client) GetCustomFields(accountAlias string) (acf []AccountCustomField, err error) {
	req := struct { AccountAlias string } { accountAlias }
	err = c.getResponse("/Account/GetCustomFields/JSON", &req, &struct {
		BaseResponse
		AccountCustomFields	*[]AccountCustomField
	} { AccountCustomFields: &acf })
	return
}

/*
 * Account Details
 */
type AccountStatus int
const (
	Active AccountStatus = 1
	Disabled             = 2
	Deleted              = 3
	Demo                 = 4
)

func (s AccountStatus) String() string {
	var accountStatus = [...]string{
		"Active",
		"Disabled",
		"Deleted",
		"Demo",
	}
	return accountStatus[s-1]
}

type AccountDetails struct {
	//  Short name associated with the account.
	AccountAlias		string
	//  Short name associated with parent account of the queried account.
	ParentAlias		string
	//  Data center location alias associated with this account.
	Location		string
	//  Full name of the business that the account is registered under.
	BusinessName		string
	//  Street address of the business associated with this account.
	Address1		string
	//  Secondary street address (if any) of the business associated with this account.
	Address2		string
	//  City of the business associated with this account.
	City			string
	//  State or province of the business associated with this account.
	StateProvince		string
	//  Postal code of the business associated with this account.
	PostalCode		string
	//  Country of the business associated with this account.
	Country			string
	// Telephone number of the business associated with this account.
	Telephone		string
	// Fax number (if any) of the business associated with this account.
	Fax			string
	// Time zone of the business associated with this account.
	TimeZone		string

	// Indicator of whether the account is active or not.
	Status			AccountStatus

	// Flag indicating whether this account shares the networks of its parent.
	ShareParentNetworks	bool

	//  Indicator of support level for the account.
	// developer,
	// legacy,
	// professional,
	// enterprise
	SupportLevel		string
}

// Get all of the contact information and settings for a given account.
func (c *Client) GetAccountDetails(accountAlias string) (details *AccountDetails, err error) {
	details = new(AccountDetails)
	req := struct { AccountAlias string } { accountAlias }
	err = c.getResponse("/Account/GetAccountDetails/JSON", &req, &struct {
		BaseResponse
		AccountDetails	*AccountDetails
	} { AccountDetails: details })
	return
}

/*
 * Datacenter Locations
 */
type Location struct {
	// Short name associated with the data center.
	Alias	string

	// Full name, or friendly name, of the data center.
	Region	string
}

// Get list of all valid data center location codes that are used in subsequent Account operations.
func (c *Client) GetLocations() (loc []Location, err error) {
	err = c.getResponse("/Account/GetLocations/JSON", nil, &struct {
		BaseResponse
		Locations       *[]Location
	} {Locations: &loc})
	return
}
