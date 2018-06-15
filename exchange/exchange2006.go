package exchange

import (
	"net/http"

	"github.com/Azure/go-ntlmssp"
)

//Exchange2006 is main struct for requests.
type Exchange2006 struct {
	User     string
	Password string
	URL      string
	client   *http.Client
}

func NewExchange(user string, password string, url string) Exchange2006 {
	client := &http.Client{
		Transport: ntlmssp.Negotiator{
			RoundTripper: &http.Transport{},
		},
	}
	return Exchange2006{User: user, Password: password, URL: url, client: client}
}
