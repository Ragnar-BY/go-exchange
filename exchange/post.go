package exchange

import (
	"net/http"
	"bytes"
	"github.com/Azure/go-ntlmssp"
	"io/ioutil"
	"log"
)

func (e *Exchange2006)Post(contents []byte) (string, error) {

	req, err := http.NewRequest("POST", e.Url, bytes.NewBuffer(contents))
//	req2.Header.Set("Host", e.User+"@"+host)
	req.Header.Set("Content-Type", "text/xml")
	req.SetBasicAuth(e.User, e.Password)

	client := &http.Client{
		Transport: ntlmssp.Negotiator{
			RoundTripper:&http.Transport{},
		},
	}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(content), nil
}
