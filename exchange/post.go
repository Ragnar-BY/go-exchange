package exchange

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
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

	response, err := e.client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	err = isSOAPFault(string(body))
	if err != nil {
		return "", err
	}
	return string(body), nil
}
func isSOAPFault(soap string) error {
	decoder := xml.NewDecoder(bytes.NewBufferString(soap))
	type Fault struct {
		Faultstring string `xml:"faultstring"`
	}
	var fault Fault
	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "Fault" {
				err := decoder.DecodeElement(&fault, &se)
				if err != nil {
					return err
				}
				return errors.New(fault.Faultstring)
			}
		}
	}
	return nil
}
