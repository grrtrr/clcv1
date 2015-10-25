package clcv1

import (
	"reflect"
	"fmt"
)

// All requests will receive a response with at least these attributes:
type BaseResponse struct {
	// True if the request was successful, otherwise False.
	Success		bool

	// A description of the result. The contents of this field does not contain
	// any actionable information, it is purely intended to provide a human
	// readable description of the result.
	Message		string

	// This value will help to identify any errors which were encountered while
	// processing the request. The value of '0' indicates success, all non-zero
	// StatusCodes indicate an error state.
	StatusCode	int
}

// Evaluate the StatusCode of @b according to cases listed in the v1 API
func (b *BaseResponse) Evaluate() (err error) {
	switch b.StatusCode {
	case 0:	/* Success */
	case 2:
		err = fmt.Errorf("Unknown application error - contact support to resolve the issue (%s).", b.Message)
	case 3:
		err = fmt.Errorf("Invalid request format (%s).", b.Message)
	case 5:
		err = fmt.Errorf("Resource not found (%s).", b.Message)
	case 6:
		err = fmt.Errorf("Invalid operation (%s).", b.Message)
	case 100:
		// FIXME: re-authenticate here unless this is Logon request
		err = fmt.Errorf("The APIKey / Password combination is valid. Verify your CenturyLink Cloud Credentials.")
	case 101:
		err = fmt.Errorf("Access denied (%s).", b.Message)
	case 400:
		err = fmt.Errorf("You have reached the SMTP Relay Alias limit set on your account (%s).", b.Message)
	case 401:
		err = fmt.Errorf("Relay alias attribute required (%s).", b.Message)
	case 402:
		err = fmt.Errorf("Relay alias was previously deleted (%s).", b.Message)
	case 403:
		err = fmt.Errorf("Relay alias has already been disabled (%s).", b.Message)
	case 500:
		err = fmt.Errorf("Invalid memory value (%s).", b.Message)
	case 501:
		err = fmt.Errorf("Invalid CPU value (%s).", b.Message)
	case 502:
		err = fmt.Errorf("Alias required (%s).", b.Message)
	case 503:
		err = fmt.Errorf("Alias length exceeded (%s).", b.Message)
	case 506:
		err = fmt.Errorf("Server name required (%s).", b.Message)
	case 514:
		err = fmt.Errorf("Server password required (%s).", b.Message)
	case 541:
		err = fmt.Errorf("Hardware Group ID required (%s).", b.Message)
	case 900:
		err = fmt.Errorf("Invalid RequestID (%s)", b.Message)
	case 1000:
		err = fmt.Errorf("The IP address provided is not configured on the server (%s)", b.Message)
	case 1201:
		err = fmt.Errorf("The visibility value is missing or invalid (%s).", b.Message)
	case 1310:
		err = fmt.Errorf("The Name attribute is missing (%s).", b.Message)
	case 1410:
		err = fmt.Errorf("Name required (%s).", b.Message)
	case 1411:
		err = fmt.Errorf("Password required (%s).", b.Message)
	case 1413:
		err = fmt.Errorf("Maximum size of additional storage exceeded (%s).", b.Message)
	case 1414:
		err = fmt.Errorf("Password does not meet strength requirements (%s).", b.Message)
	case 1415:
		err = fmt.Errorf("Unknown snapshot state (%s).", b.Message)
	case 1510:
		err = fmt.Errorf("Network required (%s).", b.Message)
	case 1511:
		err = fmt.Errorf("Public IP address required (%s).", b.Message)
	case 1600:
		err = fmt.Errorf("Account alias missing (%s).", b.Message)
	case 1700:
		err = fmt.Errorf("Email address required (%s).", b.Message)
	case 1702:
		err = fmt.Errorf("First name required (%s).", b.Message)
	case 1703:
		err = fmt.Errorf("Last name required (%s).", b.Message)
	case 1704:
		err = fmt.Errorf("Username required (%s).", b.Message)
	case 1705:
		err = fmt.Errorf("User not found (%s).", b.Message)
	case 1706:
		err = fmt.Errorf("Invalid user role(s) (%s).", b.Message)
	case 1800:
		err = fmt.Errorf("Account not found (%s).", b.Message)
	case 1801:
		err = fmt.Errorf("Invalid start date (%s).", b.Message)
	case 1802:
		err = fmt.Errorf("Invalid end date (%s).", b.Message)
	default:
		err = fmt.Errorf("%s (Status Code: %d).", b.Message, b.StatusCode)
	}
	return
}

// Extract (embedded) BaseResponse from @inModel
func ExtractBaseResponse(inModel interface{}) (*BaseResponse, error) {
	var base reflect.Value

	if inModel == nil {
		return nil, fmt.Errorf("Result model can not be nil")
	}

	resType := reflect.TypeOf(inModel)

	/* inModel must be a pointer type (call-by-value) */
	if resType.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("Expecting pointer to model %T", inModel)
	}

	/* The following is sufficient if BaseResponse is passed in directly as pointer */
	base = reflect.Indirect(reflect.ValueOf(inModel))
	if base.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Result model %T is not a pointer to struct", inModel)
	}

	if resType.Elem().Name() != "BaseResponse" {
		/* BaseResponse embedded within @inModel */
		base = base.FieldByName("BaseResponse")
	}

	if base.Kind() == reflect.Invalid {
		return nil, fmt.Errorf("Unable to extract base from result model %T", inModel)
	} else if !base.CanInterface() {
		return nil, fmt.Errorf("Unable to cast base part of %T as interface", inModel)
	}

	br, ok := base.Interface().(BaseResponse);
	if !ok {
		return nil, fmt.Errorf("Unable to cast %+v of %T as BaseResponse", base, inModel)
	}
	return &br, nil
}
