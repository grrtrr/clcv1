package clcv1

import (
	"github.com/grrtrr/clcv1/microsoft"
)

/*
 * Account Billing Summary
 */
type AccountSummary struct {
	// The current estimate of total hourly charges, assuming the current hourly run rate.
	// FIXME: in the v1 API doc this is called 'MonthlyEstimates' instead of 'MonthlyEstimate'
	MonthlyEstimate		float64

	// The total of actual hourly charges incurred thus far.
	MonthToDate		float64

	// The total charges incurred during the current hour.
	CurrentHour		float64

	// The total charges incurred during the previous hour.
	PreviousHour		float64

	// The total one time charges incurred this month
	// (this would be for non-recurring charges such as domain name registration, SSL Certificates, etc.).
	OneTimeCharges		float64

	// The total charges incurred this month to date, including one-time charges.
	MonthToDateTotal	float64
}

// Get monthly and hourly charges and estimates for a given account or collection of accounts.
// @acctAlias: Short code of the account to query.
func (c *Client) GetAccountSummary(acctAlias string) (summary AccountSummary, err error) {
	req := struct { AccountAlias string } { acctAlias }
	err = c.getResponse("/Billing/GetAccountSummary/JSON", &req, &struct {
		BaseResponse
		*AccountSummary
	} { AccountSummary: &summary })
	return
}

type BillingHistory struct {
	// Mirrors @acctAlias or the default account
	AccountAlias		string

	// FIXME: the following field is mentioned in the v1 API, but does not appear:
	// The total unpaid balance associated with the account.
	// OutstandingBalance	float64

	// Array of ledger entries
	BillingHistory		[]struct{
		// Identifier for the account's invoice.
		InvoiceID		string

		// Date of the invoice.
		Date			microsoft.Timestamp

		// Descriptive text of the invoice; typically the name of the invoice.
		Description		string

		// Charges associated with the invoice period.
		Debit			float64

		// Credits applied to the account during this invoice period.
		Credit			float64

		// FIXME: the following field is mentioned in the v1 API, but does not appear:
		// Total balance due for this invoice.
		// OutstandingBalance	float64

		// FIXME: Undocumented value that appears in the output
		DisplayPrecedence	int
	}
}

// Get the entire billing history for a given account or collection of accounts. 
// @acctAlias: Short code of the account to query.
func (c *Client) GetBillingHistory(acctAlias string) (history BillingHistory, err error) {
	req := struct { AccountAlias string } { acctAlias }
	err = c.getResponse("/Billing/GetBillingHistory/JSON", &req, &struct {
		BaseResponse
		*BillingHistory
	} { BillingHistory: &history })
	return
}

type CostEstimate struct {
	// Estimated cost of this group given current rate of usage.
	MonthlyEstimate	float64

	// Current charges so far this month.
	MonthToDate	float64

	// Charges for the current hour of usage.
	CurrentHour	float64

	// Charges for the previous hour of usage.
	PreviousHour	float64
}

/* 
 * Servers
 */
// Get the estimated monthly cost for a given server.
// @serverId:  Name of the server.
// @acctAlias: Short code of the account to query.
func (c *Client) GetServerEstimate(serverId, acctAlias string) (estimate CostEstimate, err error) {
	req := struct { AccountAlias, ServerName string } { acctAlias, serverId }
	err = c.getResponse("/Billing/GetServerEstimate/JSON", &req, &struct {
		BaseResponse
		*CostEstimate
	} { CostEstimate: &estimate })
	return
}

type ServerHourlyCharges struct {
	// Name given to the server.
	ServerName	string

	// Short code for a particular account.
	AccountAlias	string

	// Start of the query period.
	StartDate	microsoft.Timestamp

	// End of the query period.
	EndDate		microsoft.Timestamp

	// Aggregation of charges for the server.
	Summary		CostEstimate

	// Per hour for the server.
	HourlyCharges	[]struct {
		// FIXME: the following appear in the output,
		//        but are not documented.
		Hour		string
		ProcessorCost	string
		MemoryCost 	string
		StorageCost	string
		OSCost		string
	}
}

// Get the server-based hourly cost for a given time period. 
// @serverId:  Name of the server.
// @start:     Start date of the date range (optional).
// @end:       End date of the date range (optional).
// @acctAlias: Short code of the account to query.
func (c *Client) GetServerHourlyCharges(serverId, start, end, acctAlias string) (hourly ServerHourlyCharges, err error) {
	req := struct { AccountAlias, ServerName, StartDate, EndDate string } { acctAlias, serverId, start, end }
	err = c.getResponse("/Billing/GetServerHourlyCharges/JSON", &req, &struct {
		BaseResponse
		*ServerHourlyCharges
	} { ServerHourlyCharges: &hourly })
	return
}

/*
 * Hardware Groups
 */
// Get estimated costs for a group of servers.
// @grpUuid:   Unique identifier of the hardware group to estimate.
// @acctAlias: Short code of the account to query.
func (c *Client) GetGroupEstimate(grpUuid, acctAlias string) (estimate CostEstimate, err error) {
	req := struct { AccountAlias, HardwareGroupUUID string } { acctAlias, grpUuid }
	err = c.getResponse("/Billing/GetGroupEstimate/JSON", &req, &struct {
		BaseResponse
		*CostEstimate
	} { CostEstimate: &estimate })
	return
}

type GroupSummaries struct {
	// Account that this group summary applied to.
	AccountAlias	string

	// Start date of the query range.
	StartDate	string

	// End date of the query range.
	EndDate		string

	// Overview of costs for the group.
	Summary		CostEstimate

	// Details of individual costs of each group.
	GroupTotals	[]struct {
		// User-given name to this group.
		GroupName	string

		// Unique identifier for this specific group.
		GroupID		int

		// Data center alias corresponding to this group.
		LocationAlias	string

		// Charges and estimated costs for this group based on current usage
		CostEstimate

		// Collection of servers that make up the group, and the individual cost of each.
		ServerTotals	[]struct{
			// Name of this server that is incurring charges.
			ServerName	string

			// Charges and estimated costs for the server based on current usage.
			CostEstimate
		}
	}
}

// Get the charges for groups and servers within a given account and date range.
// @startDate: Start date of the date range (optional).
// @endDate:   End date of the date range (optional).
// @acctAlias: Short code of the account to query.
func (c *Client) GetGroupSummaries(startDate, endDate, acctAlias string) (summaries GroupSummaries, err error) {
	req := struct { AccountAlias, StartDate, EndDate string } { acctAlias, startDate, endDate }
	err = c.getResponse("/Billing/GetGroupSummaries/JSON", &req, &struct {
		BaseResponse
		*GroupSummaries
	} { GroupSummaries: &summaries })
	return
}

/*
 * Invoices
 */
type Invoice struct {
	// Unique identifier of the invoice.
	ID			string

	// Date of the invoice.
	InvoiceDate		microsoft.Timestamp

	// Billing terms of this invoice.
	Terms			string

	// Name of the company associated with the account.
	CompanyName		string

	// Short code of the account.
	AccountAlias		string

	// FIXME: this one is not documented, but appears in the output
        PricingAccountAlias	string

	// FIXME: this one is not documented, but appears in the output
        ParentAccountAlias	string

	// Street address associated with the account.
	Address1		string

	// Secondary street address associated with the account.
	Address2		string

	// City associated with the account.
	City			string

	// State or province associated with the account.
	StateProvince		string

	// Postcal code associated with the account.
	PostalCode		string

	// Email address associated with the account.
	BillingContactEmail	string

	// Secondary email address associated with the account.
	InvoiceCCEmail		string

	// Total amount of this invoice.
	TotalAmount		float64

	// Purchase Order identifier.
	PONumber		string

	// Individual line item on the invoice.
	LineItems		[]struct {
		// Typically the name of the billed resource or container.
		Description		string

		// Data center alias associated with this resource.
		ServiceLocation		string

		// Cost of one unit of the resource.
		UnitCost		float64

		// Unit cost multiplied by the quantity.
		ItemTotal		float64

		// Count of the item that is being charged.
		Quantity		int

		// Individual line item description and cost.
		// For instance, may refer to the servers within a group.
		ItemDetails		[]struct {
			// Typically the name of the lowest level billed resource, such as a server name.
			Description	string

			// Cost of this resource.
			Cost		float64
		}
	}
}

type InvoiceDetails struct {
	// Invoice details and line items.
	Invoice				Invoice

	// Indicator of support level for the account
	// (developer, legacy, professional, enterprise)
	SupportLevel			string

	/*
         FIXME: the following fields appear in the documentation at
                https://www.ctl.io/api-docs/v1/#billing-getinvoicedetails,
	        but not in the output (Oct 2015):
	
        // Previous balance on the account.
	OpeningBalance			float64

	// New charges for the invoice period.
	NewCharges			float64

	// Amount of payments received.
	Payments			float64

	// Total billed amount for this invoice.
	EndingBalance			float64

	// Total amount owed by this account.
	CurrentOutstandingBalance	float64
	*/
}

// Get the details for a given invoice within an account.
// @invoiceID: Unique identifier for a given invoice.
// @acctAlias: Short code of the account to query (optional).
func (c *Client) GetInvoiceDetails(invoiceID, acctAlias string) (details InvoiceDetails, err error) {
	req := struct { AccountAlias, InvoiceID string } { acctAlias, invoiceID }
	err = c.getResponse("/Billing/GetInvoiceDetails/JSON", &req, &struct {
		BaseResponse
		*InvoiceDetails
	} { InvoiceDetails: &details })
	return
}
