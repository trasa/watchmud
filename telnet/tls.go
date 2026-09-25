package telnet

import (
	"crypto/tls"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// defaultHandshakeTimeout is how long a connection to the TLS port gets to
// finish its handshake. A plain telnet client there never will, and nor will
// something that connects and says nothing; both hold a connection slot until
// this runs out.
const defaultHandshakeTimeout = 10 * time.Second

// tlsConfig serves cert, and nothing older than TLS 1.2.
func tlsConfig(cert *certificate) *tls.Config {
	return &tls.Config{
		MinVersion:     tls.VersionTLS12,
		GetCertificate: cert.get,
	}
}

// certificate is the server's certificate, re-read from disk when either file
// changes. Let's Encrypt certificates last ninety days, and loading a renewed
// one by restarting would disconnect every player, so each handshake checks
// the files' modification times -- a stat, cheap next to the handshake -- and
// reloads when they've moved.
//
// A renewal can be caught half-copied: the new certificate beside the old
// key, or a file cut short. Then the old certificate keeps being served and
// the next handshake tries again, since the loaded version only advances on a
// load that worked.
type certificate struct {
	certFile, keyFile string

	mu      sync.Mutex // handshakes run concurrently, one goroutine each
	cert    *tls.Certificate
	version time.Time // the newer of the two files' mtimes when cert was loaded
}

// loadCertificate fails if the files can't be read now. TLS was asked for, and
// a server that starts without it would advertise a port that doesn't work.
func loadCertificate(certFile, keyFile string) (*certificate, error) {
	c := &certificate{certFile: certFile, keyFile: keyFile}
	version, err := c.filesVersion()
	if err != nil {
		return nil, fmt.Errorf("tls certificate: %w", err)
	}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("tls certificate: %w", err)
	}
	c.cert, c.version = &cert, version
	return c, nil
}

func (c *certificate) get(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	version, err := c.filesVersion()
	if err != nil || !version.After(c.version) {
		return c.cert, nil // unchanged, or mid-renewal: what we have still works
	}
	cert, err := tls.LoadX509KeyPair(c.certFile, c.keyFile)
	if err != nil {
		log.Warn().Err(err).Msg("tls: the certificate files changed but don't load yet; still serving the old one")
		return c.cert, nil
	}
	log.Info().Msgf("tls: loaded the renewed certificate from %s", c.certFile)
	c.cert, c.version = &cert, version
	return c.cert, nil
}

func (c *certificate) filesVersion() (time.Time, error) {
	var newest time.Time
	for _, f := range []string{c.certFile, c.keyFile} {
		info, err := os.Stat(f)
		if err != nil {
			return time.Time{}, err
		}
		if info.ModTime().After(newest) {
			newest = info.ModTime()
		}
	}
	return newest, nil
}
