package mongostore

import "net/url"

// RedactURI is a connection string fit for a log: the password, if there is
// one, replaced with xxxxx. Anything that doesn't parse comes back as a
// placeholder rather than as itself, since it might still hold a password.
func RedactURI(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return "(unparseable uri)"
	}
	return u.Redacted()
}
