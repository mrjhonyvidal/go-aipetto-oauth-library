package oauth

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
fmt.Println("About to start our test cases...")
os.Exit(m.Run())
}

func TestLoginNotError(t *testing.T) {
	accessTokenId := "abc123"

	client := resty.New().SetHostURL("http://localhost:8082").SetTimeout(1 * time.Minute)
	resp, err := client.R().Get(fmt.Sprintf("/oauth/access_token/%s", accessTokenId))

	assert.NotNil(t, resp)
	assert.Nil(t, err)
}
