package clcv1

/*
 * Relay Aliases
 */
type RelayAlias struct {
	// The new Relay Alias created for your account.
	Alias		string

	// The password associated with the new Relay Alias.
	Password	string

	// The current status of the Relay Alias (expected values are Active, Deleted, and Disabled)
	Status		string
}

// List all SMTP Relay aliases associated with account.
func (c *Client) ListAliases() (aliases []RelayAlias, err error) {
	err = c.getResponse("/SMTPRelay/ListAliases/JSON", nil, &struct {
		BaseResponse
		SMTPRelayAliases	*[]RelayAlias
	} { SMTPRelayAliases: &aliases })
	return
}

// Ceate a new SMTP Relay alias.
// Returns the new alias name and associated password, or error on failure.
func (c *Client) CreateAlias() (alias, password string, err error) {
	err = c.getResponse("/SMTPRelay/CreateAlias/JSON", nil, &struct{
		BaseResponse
		// The new Relay Alias created for your account
		RelayAlias	*string
		// The password associated with the new Relay Alias.
		Password	*string
	} { RelayAlias: &alias, Password: &password })
	return
}

// Disable an existing SMTP Relay alias.
// @relayAlias: The SMTP relay alias to disable.
func (c *Client) DisableAlias(relayAlias string)  error {
	req := struct { RelayAlias string } { relayAlias }
	return c.getResponse("/SMTPRelay/DisableAlias/JSON", &req, new(BaseResponse))
}

// Remove an existing SMTP Relay alias.
// @relayAlias: The SMTP relay alias to remove.
func (c *Client) RemoveAlias(relayAlias string)  error {
	req := struct { RelayAlias string } { relayAlias }
	return c.getResponse("/SMTPRelay/RemoveAlias/JSON", &req, new(BaseResponse))
}





