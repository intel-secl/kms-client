package kms

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	commonTls "intel/isecl/lib/common/tls"
)

type authToken struct {
	token             string
	authorizationDate time.Time
	notAfter          time.Time
	faults            interface{}
}

// A Client is defines parameters to connect and Authenticate with a KMS
type Client struct {
	// BaseURL specifies the URL base for the KMS, for example https://keymanagement.server/
	BaseURL string
	// Username used to authenticate with the KMS. Username is only used for obtaining an authorization token, which is automatically used for requests.
	Username string
	// Password to supply for the Username
	Password string
	// CertSha256 is a pointer to a 32 byte array that specifies the fingerprint of the immediate TLS certificate to trust.
	// If the value is a non nil pointer to a 32 byte array, custom TLS verification will be used, where any valid chain of X509 certificates
	// with a self signed CA at the root will be accepted as long as the Host Certificates Fingerprint matches what is provided here
	// If the value is a nil pointer, then system standard TLS verification will be used.
	CertSha256 *[32]byte
	// A reference to the underlying http Client.
	// If the value is nil, a default client will be created and used.
	HTTPClient *http.Client
	authToken  authToken
}

func (c *Client) httpClient() *http.Client {
	if c.HTTPClient == nil {
		// init http client
		tlsConfig := tls.Config{}
		if c.CertSha256 != nil {
			// set explicit verification
			tlsConfig.InsecureSkipVerify = true
			tlsConfig.VerifyPeerCertificate = commonTls.VerifyCertBySha256(*c.CertSha256)
		}
		transport := http.Transport{
			TLSClientConfig: &tlsConfig,
		}
		c.HTTPClient = &http.Client{Transport: &transport}
	}
	return c.HTTPClient
}

func (c *Client) Key(uuid string) *KeyID {
	return &KeyID{Client: c, ID: uuid}
}

// Keys returns a sub client that operates on KMS /keys endpoints, such as creating a new key
func (c *Client) Keys() *Keys {
	return &Keys{Client: c}
}

func (c *Client) refreshAuthToken() error {
	loginForm := []byte(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, c.Username, c.Password))
	req, err := http.NewRequest(http.MethodPost, c.BaseURL, bytes.NewBuffer(loginForm))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	rsp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}
	if err := json.NewDecoder(rsp.Body).Decode(&c.authToken); err != nil {
		return err
	}
	return nil
}

func (c *Client) dispatchRequest(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept", "Application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.authToken.token)
	rsp, err := c.httpClient().Do(req)
	if err != nil && rsp.StatusCode == http.StatusUnauthorized || rsp.StatusCode == http.StatusForbidden {
		err := c.refreshAuthToken()
		if err != nil {
			// failed to refresh token, Create request is a failure
			return nil, err
		}
		// retry the request
		req.Header.Set("Authorization", "Token "+c.authToken.token)
		rsp, err = c.httpClient().Do(req)
	}
	return rsp, err
}
