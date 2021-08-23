package domain

import (
	"encoding/json"
	"github.com/aerostatka/banking-lib/logger"
	"net/http"
	"net/url"
)

type AuthRepository interface {
	IsAuthorized(token string, route string, vars map[string]string) bool
}

type RemoteAuthRepository struct {
	
}

func (r RemoteAuthRepository) IsAuthorized(token string, route string, vars map[string]string) bool  {
	u := buildVerifyUrl(token, route, vars)

	if response, err := http.Get(u); err != nil {
		logger.Error("Error while sending....")

		return false
	} else {
		m := map[string]bool{}

		if err = json.NewDecoder(response.Body).Decode(&m); err != nil {
			logger.Error("Error while while decoding response from Auth server: " + err.Error())

			return false
		}

		return m["isAuthorized"]
	}
}

func buildVerifyUrl(token string, route string, vars map[string]string) string {
	u := url.URL{
		Host: "localhost:8281",
		Path: "/auth/verify",
		Scheme: "http",
	}

	q := u.Query()
	q.Add("token", token)
	q.Add("routeName", route)

	for k, v := range vars {
		q.Add(k, v)
	}

	u.RawQuery = q.Encode()

	return u.String()
}

func NewAuthRepository() RemoteAuthRepository {
	return RemoteAuthRepository{}
}