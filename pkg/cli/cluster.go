package cli

import (
	"crypto/x509"
	"strconv"
	"strings"
	"time"

	"github.com/dcos/dcos-cli/pkg/config"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
)

// Cluster is a subset representation of a DC/OS CLI configuration.
//
// It is a proxy struct on top of a config which provides user-friendly getters and setters for common
// configurations such as "core.dcos_url" or "core.ssl_verify". It leverages Go types as much as possible.
type Cluster struct {
	config *config.Config
}

// NewCluster returns a new cluster for a given config, if omitted it uses an empty config.
func NewCluster(conf *config.Config) *Cluster {
	if conf == nil {
		conf = config.Empty()
	}
	return &Cluster{config: conf}
}

// URL returns the public master URL of the DC/OS cluster.
func (c *Cluster) URL() string {
	url := cast.ToString(c.config.Get("core.dcos_url"))
	return strings.TrimRight(url, "/")
}

// SetURL sets the public master URL of the DC/OS cluster.
func (c *Cluster) SetURL(url string) {
	c.config.Set("core.dcos_url", url)
}

// ACSToken returns the token generated by authenticating
// to DC/OS using the Admin Router Access Control Service.
func (c *Cluster) ACSToken() string {
	return cast.ToString(c.config.Get("core.dcos_acs_token"))
}

// SetACSToken sets the token generated by authenticating
// to DC/OS using the Admin Router Access Control Service.
func (c *Cluster) SetACSToken(acsToken string) {
	c.config.Set("core.dcos_acs_token", acsToken)
}

// TLS returns the configuration for TLS clients.
func (c *Cluster) TLS() TLS {
	tlsVal := cast.ToString(c.config.Get("core.ssl_verify"))

	// Try to cast the value to a bool, true means we verify
	// server certificates, false means we skip verification.
	if verify, err := strconv.ParseBool(tlsVal); err == nil {
		return TLS{Insecure: !verify}
	}

	// The value is not a string representing a bool thus it is a path to a root CA bundle.
	rootCAsPEM, err := afero.ReadFile(c.config.Fs(), tlsVal)
	if err != nil {
		return TLS{
			Insecure:    true,
			RootCAsPath: tlsVal,
		}
	}

	// Decode the PEM root certificate(s) into a cert pool.
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(rootCAsPEM) {
		return TLS{
			Insecure:    true,
			RootCAsPath: tlsVal,
		}
	}

	// The cert pool has been successfully created, store it in the TLS config.
	return TLS{
		RootCAs:     certPool,
		RootCAsPath: tlsVal,
	}
}

// SetTLS returns the configuration for TLS clients.
func (c *Cluster) SetTLS(tls TLS) {
	c.config.Set("core.ssl_verify", tls.String())
}

// Timeout returns the HTTP request timeout once the connection is established.
func (c *Cluster) Timeout() time.Duration {
	timeout := c.config.Get("core.timeout")
	return time.Duration(cast.ToInt64(timeout)) * time.Second
}

// SetTimeout sets the HTTP request timeout once the connection is established.
func (c *Cluster) SetTimeout(timeout time.Duration) {
	c.config.Set("core.timeout", timeout.Seconds())
}

// Name returns the custom name for the cluster.
func (c *Cluster) Name() string {
	return cast.ToString(c.config.Get("cluster.name"))
}

// SetName sets a custom name for the cluster.
func (c *Cluster) SetName(name string) {
	c.config.Set("cluster.name", name)
}

// Config returns the cluster's config.
func (c *Cluster) Config() *config.Config {
	return c.config
}

// TLS holds the configuration for TLS clients.
type TLS struct {
	// Insecure specifies if server certificates should be accepted without verification.
	//
	// Skipping verification against the system's CA bundle or a cluster-specific CA is highly discouraged
	// and should only be done during testing/development.
	Insecure bool

	// Path to the root CA bundle.
	RootCAsPath string

	// A pool of root CAs to verify server certificates against.
	RootCAs *x509.CertPool
}

// String creates a string from a TLS struct.
func (tls *TLS) String() string {
	if tls.RootCAsPath != "" {
		return tls.RootCAsPath
	}
	return strconv.FormatBool(!tls.Insecure)
}
