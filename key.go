package kms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// KeyID represents a single key id on the KMS, equating to /keys/id
type KeyID struct {
	client *Client
	ID     string
}

// KeyObj is a represenation of the actual key
type KeyObj struct {
	Key []byte `json:"key"`
}

// KeyInfo is a representation of key information
type KeyInfo struct {
	KeyID           string `json:"id"`
	CipherMode      string `json:"mode"`
	Algorithm       string `json:"algorithm"`
	KeyLength       string `json:"key_length"`
	PaddingMode     string `json:"padding_mode"`
	TransferPolicy  string `json:"transfer_policy"`
	TransferLink    string `json:"transfer_link"`
	DigestAlgorithm string `json:"digest_algorithm"`
}

// Transfer performs a POST to /key/{id}/transfer to retrieve the actual key data from the KMS
func (k *KeyID) Transfer(saml []byte) ([]byte, error) {
	var samlBuf io.Reader
	if saml != nil {
		samlBuf = bytes.NewBuffer(saml)
	}
	baseURL, err := url.Parse(k.client.BaseURL)
	if err != nil {
		return nil, err
	}
	keyXferURL, err := url.Parse(fmt.Sprintf("keys/%s/transfer", k.ID))
	if err != nil {
		return nil, err
	}
	reqURL := baseURL.ResolveReference(keyXferURL)
	req, err := http.NewRequest("POST", reqURL.String(), samlBuf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/octet-stream")
	req.Header.Set("Content-Type", "application/json")
	rsp, err := k.client.dispatchRequest(req)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kms-client: failed to transfer key (HTTP Status Code: %d)", rsp.StatusCode)
	}
	bytes, _ := ioutil.ReadAll(rsp.Body)
	return bytes, nil
}

// Keys represents the resource collection of Keys on the KMS
type Keys struct {
	client *Client
}

// Create sends a POST to /keys to create a new Key with the specified parameters
func (k *Keys) Create(key KeyInfo) (*KeyInfo, error) {
	// marshal KeyInfo
	kiJSON, err := json.Marshal(&key)
	if err != nil {
		return nil, err
	}
	baseURL, err := url.Parse(k.client.BaseURL)
	if err != nil {
		return nil, err
	}
	keysURL, _ := url.Parse("keys")
	reqURL := baseURL.ResolveReference(keysURL)
	req, err := http.NewRequest("POST", reqURL.String(), bytes.NewBuffer(kiJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	rsp, err := k.client.dispatchRequest(req)
	defer rsp.Body.Close()
	var kiOut KeyInfo
	err = json.NewDecoder(rsp.Body).Decode(&kiOut)
	if err != nil {
		return nil, err
	}
	return &kiOut, nil
}
