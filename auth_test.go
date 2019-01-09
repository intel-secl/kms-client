package kms

import (
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
	json.Unmarshal([]byte(raw), &token)
	assert.Equal(t, "BSNZJFaXojS3uXfvr6wvKU2Gzx0z1f+9PfI9VwyY1A8=", token.AuthorizationToken)
}
