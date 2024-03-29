package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aipetto/go-aipetto-utils/src/rest_errors"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strconv"
	"time"
)

const (
	headerXPublic 		= "X-Public"
	headerXClientId		= "X-Client-Id"
	headerXUserId		= "X-User-Id"
	paramAccessToken	= "access_token_id"
)


type accessToken struct {
	Id 			string	`json:"id"`
	UserId 		int64	`json:"user_id"`
	ClientId	int64	`json:"client_id"`
}

func IsPublic(request *http.Request) bool{
	if request == nil {
		return true
	}
	return request.Header.Get(headerXPublic) == "true"
}

func GetUserId(request *http.Request) int64 {
	if request == nil {
		return 0
	}
	userId, err := strconv.ParseInt(request.Header.Get(headerXUserId), 10, 64)
	if err != nil {
		return 0
	}
	return userId
}

func GetClientId(request *http.Request) int64 {
	if request == nil {
		return 0
	}
	clientId, err := strconv.ParseInt(request.Header.Get(headerXClientId), 10, 64)
	if err != nil {
		return 0
	}
	return clientId
}

func AuthenticateRequest(request *http.Request) *rest_errors.RestErr{
	if request == nil {
		return nil
	}
	// https://aipetto.com/resource?access_token=abc123
	accessTokenId := request.URL.Query().Get(paramAccessToken)
	if accessTokenId == "" {
		return nil
	}

	cleanRequest(request)

	at, err := getAccessToken(accessTokenId)
	if err != nil {
		if err.Status == http.StatusNotFound{
			return nil
		}
		return err
	}

	request.Header.Add(headerXClientId, fmt.Sprintf("%v", at.ClientId))
	request.Header.Add(headerXUserId, fmt.Sprintf("%v", at.UserId))

	return nil
}

func cleanRequest(request *http.Request) {
	if request == nil {
		return
	}
	request.Header.Del(headerXClientId)
	request.Header.Del(headerXUserId)
}

func getAccessToken(accessTokenId string) (*accessToken, *rest_errors.RestErr) {
	client := resty.New().SetHostURL("http://localhost:8082").SetTimeout(1 * time.Minute)
	resp, err := client.R().Get(fmt.Sprintf("/oauth/access_token/%s", accessTokenId))

	if err != nil {
		return nil, rest_errors.NewInternalServerError("invalid oauth rest response when trying to obtain access token", errors.New("OAuth Rest Error"))
	}

	if resp.StatusCode() > 299 {
		var restErr rest_errors.RestErr
		if err := json.Unmarshal(resp.Body(), &restErr); err != nil {
			return nil, rest_errors.NewInternalServerError("invalid error interface when trying to get the token", errors.New("Invalid Token"))
		}
		return nil, &restErr
	}

	var at accessToken
	if err := json.Unmarshal(resp.Body(), &at); err != nil {
		return nil, rest_errors.NewInternalServerError("error when trying to unmarshal token response", errors.New("Decode UnMarshal Error"))
	}
	return &at, nil
}