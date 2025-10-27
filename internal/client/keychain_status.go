package client

// UsingKeychain reports whether the token is stored in the OS keychain backend.
func (c *Client) UsingKeychain() bool {
	return c.usingSecret
}
