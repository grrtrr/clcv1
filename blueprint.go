package clcv1

import (
	"github.com/grrtrr/clcv1/microsoft"
	"github.com/dustin/go-humanize"
	"strings"
	"time"
	"fmt"
)

/*
 * List Blueprints
 */
type SearchBlueprintReq struct {
	// The target company size the of the Blueprint
	// 1 = 1 - 100
	// 2 = 101 - 1,000
	// 3 = 1,001 - 5,000
	// 4 = 5,000+
	CompanySize		int

	// A list of the operating systems that a Blueprint contains (optional)
	// e.g. Cent OS - 32 bit = 6,  Cent OS - 64 bit = 7
	OperatingSystems	[]int

	// A keyword search within the Name and Description of the Blueprint (optional)
	Search			string

	// The visibility level of the Blueprint.
	// 1 = Public
	// 2 = Private
	// 3 = Private Shared
	Visibility		int
}

type Blueprint struct {
	// The ID of the Blueprint.
	ID	int
	// The name of the Blueprint.
	Name	string
}

// Get a list of all Blueprints with the specified search criteria.
func (c *Client) GetBlueprints(req *SearchBlueprintReq) (blueprints []Blueprint, err error) {
	err = c.getResponse("/Blueprint/GetBlueprints/JSON", req, &struct {
		BaseResponse
		Blueprints	*[]Blueprint
	} { Blueprints: &blueprints })
	return
}

/*
 * Deployment Status
 */
type DeploymentStatus struct {
	// The ID of the Blueprint deployment request. (Obtained from DeployBlueprint.)
	RequestID	int

	// The status of the Blueprint deployment. Valid values are: NotStarted, Executing, Succeeded, Failed, and Resumed.
	CurrentStatus	string

	// A detailed description of the current status.
	Description	string

	// The percentage of the work that has been completed on the request.
	PercentComplete	int

	// The timestamp (UTC) that the most recent status was recorded on the Request.
	StatusDate	microsoft.Timestamp

	// The current step being executed.
	Step		string

	// An array of strings with the names of any servers built as part of the Blueprint.
	// This element is only populated when the Request was originally for server creation.
	Servers		[]string
}

// Get the status of the specified Blueprint deployment.
// @reqId:     The ID of the Blueprint Deployment to retrieve status for.
// @acctAlias: ID of the account (required)
// @location:  The location of the Blueprint Deployment to retrieve status for.
//             If not provided, will default to the API user's default data center.
func (c *Client) GetDeploymentStatus(reqId int, acctAlias, location string) (s *DeploymentStatus, err error) {
	req := struct {
		RequestID	int
		AccountAlias	string
		LocationAlias	string

	} { reqId, acctAlias, location }
	s = new(DeploymentStatus)
	err = c.getResponse("/Blueprint/GetDeploymentStatus/JSON", &req, &struct {
		BaseResponse
		*DeploymentStatus
	} { DeploymentStatus: s })
	return
}

// Poll the queue status of @reqId each @interval seconds, until it reaches 100%
// @reqId, @location, @acctAlias: As per GetDeploymentStatus().
// @pollInterval:                 Poll interval in seconds; use 0 for one-off.
func (c *Client) PollDeploymentStatus(reqId int, location, acctAlias string, pollInterval int) error {
	for {
		status, err := c.GetDeploymentStatus(reqId, acctAlias, location)
		if err != nil {
			return fmt.Errorf("Failed to query status of request ID %d: %s", reqId, err)
		}

		/* Clear line */
		fmt.Printf("\r\033[2K")
		fmt.Printf("#%d %d%% at %s (%s)", status.RequestID, status.PercentComplete,
			   status.StatusDate.Format(time.Kitchen), humanize.Time(status.StatusDate.Time))
		if len(status.Servers) == 1 {
			fmt.Printf(", %s", status.Servers[0])
		} else if len(status.Servers) > 1 {
			fmt.Printf(", %s", strings.Join(status.Servers, ", "))
		}
		fmt.Printf(", %s (%s)", status.CurrentStatus, status.Description)

		if status.PercentComplete == 100 || pollInterval == 0 {
			fmt.Printf("\n")
			break
		}
		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
	return nil
}
