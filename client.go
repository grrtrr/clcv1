/*
 * Client API. Some ideas inspired by centurylink_sdk
 */
package clcv1

import (
	"net/http/cookiejar"
	"net/http/httputil"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"bytes"
	"time"
	"flag"
	"log"
	"fmt"
	"io"
)

const (
	BaseURL = "https://api.ctl.io/REST"
)

// Global variables
var g_debug  bool     /* Command-line debug flag */

func init() {
	flag.BoolVar(&g_debug,    "d", false, "Produce debug output")
}

// Client wraps http.Client, with logging added
type Client struct {
	*http.Client
	Log *log.Logger
}

// Return new v1 Client
func NewClient(logger *log.Logger) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	if logger == nil {
		log.New(ioutil.Discard, "", log.LstdFlags)
	}
	return &Client{ &http.Client{ Jar: jar }, logger }, nil
}

// Set the transport timeout for the client
func (c *Client) SetTimeout(timeout time.Duration) {
	c.Client.Timeout = timeout
}

// Change the logger
func (c *Client) SetLogger(logger *log.Logger) {
	if logger == nil {
		panic(fmt.Errorf("Client: logger argument must not be nil."))
	}
	c.Log = logger
}

// POST a v1 API request to @path relative to BaseURL.
// @reqModel: request model to serialize, or nil
// @resModel: result model to deserialize, must be a pointer to the expected result
// Evaluates the StatusCode of the BaseResponse (embedded) in @inModel and sets @err accordingly.
// If @err == nil, fills in @resModel, else returns error.
func (c *Client) getResponse(path string, reqModel interface{}, resModel interface{}) (err error) {
	var reqBody io.Reader

	if reqModel != nil {
		if g_debug {
			c.Log.Printf("reqModel %T %+v\n", reqModel, reqModel)
		}

		jsonReq, err := json.Marshal(reqModel)
		if err != nil {
			return fmt.Errorf("Failed to encode request model %T %+v: %s", reqModel, reqModel, err)
		}
		reqBody = bytes.NewBuffer(jsonReq)
	}

	/* resModel must be a pointer type (call-by-value) */
	if resModel == nil {
		return fmt.Errorf("Result model can not be nil")
	} else if resType := reflect.TypeOf(resModel); resType.Kind() != reflect.Ptr {
		return fmt.Errorf("Expecting pointer to result model %T", resModel)
	} else if g_debug {
		c.Log.Printf("resModel %T %+v\n", resModel, resModel)
	}

	req, err := http.NewRequest("POST", BaseURL + path, reqBody)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept",       "application/json")

	if g_debug {
		reqDump, _ := httputil.DumpRequest(req, true)
		c.Log.Printf("%s", reqDump)
	}

	res, err := c.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if g_debug {
		resDump, _ := httputil.DumpResponse(res, true)
		c.Log.Printf("%s", resDump)
	}

	/* StatusCode is used instead of the HTTP status code (which is 200 even if there was an error) */
	if res.StatusCode != 200 {
		return fmt.Errorf("POST request at %s failed with status: %q", path, res.Status)
	}

	if err = json.NewDecoder(res.Body).Decode(resModel); err != nil {
		return
	}

	br, err := ExtractBaseResponse(resModel)
	if err != nil {
		return
	}

	return br.Evaluate()
}
