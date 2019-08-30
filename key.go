/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
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
	Key []byte `json:"key,omitempty"`
}

// KeyInfo is a representation of key information
type KeyInfo struct {
	KeyID           string `json:"id,omitempty"`
	CipherMode      string `json:"mode,omitempty"`
	Algorithm       string `json:"algorithm,omitempty"`
	KeyLength       int    `json:"key_length,omitempty"`
	PaddingMode     string `json:"padding_mode,omitempty"`
	TransferPolicy  string `json:"transfer_policy,omitempty"`
	TransferLink    string `json:"transfer_link,omitempty"`
	DigestAlgorithm string `json:"digest_algorithm,omitempty"`
}

// Error is a error struct that contains error information thrown by the actual KMS
type Error struct {
	StatusCode int
	Message    string
}

func (k Error) Error() string {
	return fmt.Sprintf("kms-client: failed (HTTP Status Code: %d)\nMessage: %s", k.StatusCode, k.Message)
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
	req.Header.Set("Content-Type", "application/samlassertion+xml")
	rsp, err := k.client.dispatchRequest(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		errMsgBytes, _ := ioutil.ReadAll(rsp.Body)
		return nil, &Error{StatusCode: rsp.StatusCode, Message: fmt.Sprintf("Failed to transfer key %s: %s", k.ID, string(errMsgBytes))}
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

// Retrieve performs a POST to /key/{id}/transfer to retrieve the actual key data from the KMS
func (k *KeyID) Retrieve() ([]byte, error) {
	var keyValue KeyObj
	baseURL, err := url.Parse(k.client.BaseURL)
	if err != nil {
		return nil, err
	}
	keyXferURL, err := url.Parse(fmt.Sprintf("keys/%s/transfer", k.ID))
	if err != nil {
		return nil, err
	}
	reqURL := baseURL.ResolveReference(keyXferURL)

	req, err := http.NewRequest("POST", reqURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	rsp, err := k.client.dispatchRequest(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		errMsgBytes, _ := ioutil.ReadAll(rsp.Body)
		return nil, &Error{StatusCode: rsp.StatusCode, Message: fmt.Sprintf("Failed to transfer key %s: %s", k.ID, string(errMsgBytes))}
	}

	err = json.NewDecoder(rsp.Body).Decode(&keyValue)
	if err != nil {
		return nil, err
	}
	return keyValue.Key, nil

}
