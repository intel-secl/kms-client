package kms

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// KeyID represents a single key id on the KMS, equating to /keys/id
type KeyID struct {
	*Client
	ID string
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
func (k *KeyID) Transfer() ([]byte, error) {
	var keyObj KeyObj
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/keys/%s/transfer", k.BaseURL, k.ID), nil)
	if err != nil {
		return nil, err
	}
	rsp, err := k.dispatchRequest(req)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != http.StatusOK {
		return nil, errors.New("kms-client: failed to transfer key")
	}
	err = json.NewDecoder(rsp.Body).Decode(&keyObj)
	if err != nil {
		return nil, err
	}
	return keyObj.Key, nil
}

// Keys represents the resource collection of Keys on the KMS
type Keys struct {
	*Client
}

// Create sends a /po
func (k *Keys) Create(key KeyInfo) (*KeyInfo, error) {
	// marshal KeyInfo
	kiJSON, err := json.Marshal(&key)
	if err != nil {
		return nil, err
	}
	// this aint gonna work
	req, err := http.NewRequest("POST", k.BaseURL+"/keys", bytes.NewBuffer(kiJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Token "+k.authToken.token)
	rsp, err := k.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode == http.StatusUnauthorized || rsp.StatusCode == http.StatusForbidden {
		err := k.refreshAuthToken()
		if err != nil {
			// failed to refresh token, Create request is a failure
			return nil, err
		}
		// retry the request
		req.Header.Set("Authorization", "Token "+k.authToken.token)
		rsp, err = k.httpClient().Do(req)
		if err != nil {
			return nil, err
		}
	}
	if rsp.StatusCode != http.StatusOK {
		return nil, errors.New("kms-client: failed to create key")
	}
	var kiOut KeyInfo
	err = json.NewDecoder(rsp.Body).Decode(&kiOut)
	if err != nil {
		return nil, err
	}
	return &kiOut, nil
}
