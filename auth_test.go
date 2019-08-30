/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package kms

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthTokenDecode(t *testing.T) {
	raw := `{
		"authorization_token": "BSNZJFaXojS3uXfvr6wvKU2Gzx0z1f+9PfI9VwyY1A8=",
		"authorization_date": "2019-01-09T12:36:19-0800",
		"not_after": "2019-01-09T13:06:19-0800",
		"faults": []
	  }`
	var token AuthToken
	err := json.NewDecoder(bytes.NewBufferString(raw)).Decode(&token)
	assert.NoError(t, err)
	assert.Equal(t, "BSNZJFaXojS3uXfvr6wvKU2Gzx0z1f+9PfI9VwyY1A8=", token.AuthorizationToken)
	assert.Equal(t, 2019, token.NotAfter.Year())
	assert.Equal(t, 2019, token.AuthorizationDate.Year())
}
