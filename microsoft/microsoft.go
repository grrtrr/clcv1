/*
 * Handle the Microsoft serialization of DateTime timestamps.
 * The format is Epoch Milliseconds since January 1, 1970 in UTC.
 * See https://msdn.microsoft.com/en-us/library/bb299886.aspx#intro_to_json_sidebarb
 */
package microsoft

import (
	"strconv"
	"regexp"
	"time"
	"fmt"
)

/* Global variables */
var ts_format = regexp.MustCompile(`^\s*"\\?/Date\((-?\d+)\)\\?/"\s*$`)

type Timestamp struct {
	time.Time
}

// Return @t in Microsoft JSON format
func (t *Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("/Date(%d)/", t.UnixNano()/1000000)), nil
}

// Deserialize @b, accept both `/Date(\d+)` and `\/Date(\d+)\/`
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	m := ts_format.FindSubmatch(b)
	if len(m) == 2 {
		epochMsec, err := strconv.ParseInt(string(m[1]), 10, 64)
		if err == nil {
			t.Time = time.Unix(epochMsec/1000, (epochMsec % 1000) * 1000000)
			return nil
		}
	}
	return fmt.Errorf("Invalid timestamp format %q, %v", string(b), m)
}

func (t Timestamp) String() string {
	return t.Time.String()
}
