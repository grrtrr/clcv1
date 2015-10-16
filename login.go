/*
 * CLC API v1 login/logout routines
 */
package clcv1

import (
	"github.com/grrtrr/clcv1/utils/terminal"
	"flag"
	"os"
)


// Global variables
var g_apikey, g_pass string   /* Command-line API-Key/Password */

func init() {
	flag.StringVar(&g_apikey, "k", "",    "CLC v1 API Key (if not set via CLC_V1_API_KEY)")
	flag.StringVar(&g_pass,   "p", "",    "CLC v1 API Password (if not set via CLC_V1_API_PASS)")
}

// This method is required to be called prior to calling any other method exposed by the CenturyLink Cloud API.
// This method validates your credentials and writes the Encrypted cookie required to be present for all
// subsequent calls into the API.
func (c *Client) Logon(api_key, password string) (err error) {
	var credentials struct {
		APIKey	 string
		Password string
	}

	credentials.APIKey, credentials.Password, err = resolveApiCredentials(api_key, password)
	if err != nil {
		return err
	}

	return c.getResponse("/Auth/Logon/", &credentials, new(BaseResponse))
}


// This method will log you out of the API. The Logon method must be called again prior to accessing the API again.
func (c *Client) Logout() (err error) {
	/* URL has to end in JSON, otherwise it will produce XML output */
	return c.getResponse("/Auth/Logout/JSON", nil, new(BaseResponse))
}

// Support multiple ways of resolving the v1 API Key and Password
// 1. directly (pass-through),
// 2. command-line flags (g_apikey, g_pass),
// 3. environment variables (CLC_V1_API_KEY, CLC_V1_API_PASS)
// 4. prompt for values
func resolveApiCredentials(api_key, password string) (res_key, res_pass string, err error) {
	if api_key == "" {
		api_key = g_apikey
	}
	if api_key == "" {
		api_key = os.Getenv("CLC_V1_API_KEY")
	}
	if api_key == "" {
		if api_key, err = terminal.PromptInput("API Key"); err != nil {
			return "", "", err
		}
	}

	if password == "" {
		password = g_pass
	}
	if password == "" {
		password = os.Getenv("CLC_V1_API_PASS")
	}
	if password == "" {
		if password, err = terminal.GetPass("API Password"); err != nil {
			return "", "", err
		}
	}
	return api_key, password, nil
}
