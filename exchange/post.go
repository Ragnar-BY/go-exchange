package exchange

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Azure/go-ntlmssp"
)

//Post make post request with contents.
func (e *Exchange2006) Post(contents []byte) (string, error) {
	req, err := http.NewRequest("POST", e.URL, bytes.NewBuffer(contents))
	if err != nil {
		return "", err
	}
	//	req2.Header.Set("Host", e.User+"@"+host)
	req.Header.Set("Content-Type", "text/xml")
	req.SetBasicAuth(e.User, e.Password)

	client := &http.Client{
		Transport: ntlmssp.Negotiator{
			RoundTripper: &http.Transport{},
		},
	}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(body), nil
}
