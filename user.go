package kms

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	logger "github.com/sirupsen/logrus"
)

//UserInfo is a representation of key information
type UserInfo struct {
	UserID         string `json:"id"`
	Username       string `json:"username"`
	TransferKeyPem string `json:"transfer_key_pem"`
}

//Users is a representation of key information
type Users struct {
	Users []UserInfo `json:"users"`
}

type User struct {
	client *Client
}

// GetKmsUser is used to get the kms user information
func (k *Keys) GetKmsUser() (UserInfo, error) {
	var userInfo UserInfo
	var users Users
	logger.Info("Retrieving kms user information")

	baseURL, err := url.Parse(k.client.BaseURL)
	if err != nil {
		return userInfo, err
	}
	getUserURL, err := url.Parse(fmt.Sprintf("users?usernameEqualTo=%s", k.client.Username))
	if err != nil {
		return userInfo, err
	}
	reqURL := baseURL.ResolveReference(getUserURL)
	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		return userInfo, errors.New("Error creating request")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	rsp, err := k.client.dispatchRequest(req)
	if err != nil {
		return userInfo, errors.New("Error sending request")
	}
	defer rsp.Body.Close()
	err = json.NewDecoder(rsp.Body).Decode(&users)
	if err != nil {
		return userInfo, err
	}
	return users.Users[0], nil
}

func (k *Keys) RegisterUserPubKey(publicKey []byte, userID string) error {
	baseURL, err := url.Parse(k.client.BaseURL)
	if err != nil {
		return err
	}
	keyXferURL, err := url.Parse(fmt.Sprintf("users/%s/transfer-key", userID))
	if err != nil {
		return err
	}
	reqURL := baseURL.ResolveReference(keyXferURL)
	req, err := http.NewRequest("PUT", reqURL.String(), bytes.NewBuffer(publicKey))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-pem-file")
	_, err = k.client.dispatchRequest(req)
	if err != nil {
		return err
	}

	return nil
}
