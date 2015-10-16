package clcv1

import (
	"clcv1/microsoft"
)

type ItemStatus int

const (
	All      ItemStatus = 1
	Pending  ItemStatus = 2
	Complete ItemStatus = 3
	Error    ItemStatus = 4
)

type QueueRequest struct {
	// The ID of the Request whose details were returned.
	RequestID	int

	// The current status of the request, valid values are: Not Started, Executing, Succeeded, Failed, and Resumed.
	CurrentStatus	string

	// The percentage of the work that has been completed on the request.
	PercentComplete	int

	// A description of the progress of the request, for example a description of the step currently being executed.
	ProgressDesc	string

	// A description of what the Request was created to do.
	RequestTitle	string

	// The number of the current step being executed.
	StepNumber	int

	// The timestamp (GMT) that the most recent status was recorded on the Request.
	StatusDate	microsoft.Timestamp
}

// This method can be used to get a list of Queued requests and their current status details.
func (c *Client) ListQueueRequests(status ItemStatus) (requests []QueueRequest, err error) {
	req := struct { ItemStatusType ItemStatus } { status }
	err = c.getResponse("/Queue/ListQueueRequests/JSON", &req, &struct {
		BaseResponse
		Requests	*[]QueueRequest
	} { Requests: &requests })
	return
}

// This method can be used to check the status of any of the long running
// requests which must be performed asynchronously.
func (c *Client) GetRequestStatus(requestID int) (qreq *QueueRequest, err error) {
	qreq = new(QueueRequest)
	req := struct { RequestID int } { requestID }
	err = c.getResponse("/Queue/GetRequestStatus/JSON", &req, &struct {
		BaseResponse
		RequestDetails 	*QueueRequest
	} { RequestDetails: qreq })
	return
}
