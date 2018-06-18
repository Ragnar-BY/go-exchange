package exchange

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"text/template"

	"github.com/Ragnar-BY/go-exchange/models"
)

// GetMailBox returns mailBox, if can find unique by resolveName.
func (e Exchange2006) GetMailBox(resolveName string) (*models.Mailbox, error) {
	t, err := template.New("GetMailBox").Parse(resolveNamesRequest)
	if err != nil {
		return nil, fmt.Errorf("[GetMailBox] cannot create template %v", err)
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, resolveName)
	if err != nil {
		return nil, fmt.Errorf("[GetMailBox] cannot parse template: %v", err)
	}
	response, err := e.Post(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("[GetMailBox] cannot post %v", err)
	}
	return parseGetEmailResponse(response)
}

var resolveNamesRequest = `<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
               xmlns:xsd="http://www.w3.org/2001/XMLSchema"
               xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"
               xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types">
  <soap:Body>
    <ResolveNames xmlns="http://schemas.microsoft.com/exchange/services/2006/messages"
                  xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types"
                  ReturnFullContactData="true">
      <UnresolvedEntry>{{.}}</UnresolvedEntry>
    </ResolveNames>
  </soap:Body>
</soap:Envelope>`

func parseGetEmailResponse(response string) (*models.Mailbox, error) {
	type ResolveNamesResponseMessage struct {
		ResponseClass string `xml:"ResponseClass,attr"`
		ResponseCode  string `xml:"ResponseCode"`
		MessageText   string `xml:"MessageText"`
		ResolutionSet *struct {
			Resolution []*struct {
				Mailbox *models.Mailbox `xml:"Mailbox"`
			} `xml:"Resolution"`
		} `xml:"ResolutionSet"`
	}
	var rnrm ResolveNamesResponseMessage
	decoder := xml.NewDecoder(bytes.NewBufferString(response))
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "ResolveNamesResponseMessage" {
				err := decoder.DecodeElement(&rnrm, &se)
				if err != nil {
					return nil, err
				}
				if rnrm.ResponseClass == "Error" || rnrm.ResponseClass == "Warning" {
					return nil, errors.New(rnrm.MessageText)
				}
				if rnrm.ResolutionSet == nil {
					return nil, errors.New("cannot parse resolutions")
				}
				if len(rnrm.ResolutionSet.Resolution) == 0 {
					return nil, errors.New("cannot find any resolutions")
				}
				if rnrm.ResolutionSet.Resolution[0].Mailbox == nil {
					return nil, errors.New("cannot find any mailboxes")
				}
				return rnrm.ResolutionSet.Resolution[0].Mailbox, nil
			}
		}
	}
	return nil, errors.New("cannot find ResolveNamesResponseMessage in response")
}

// FindInContactList finds unique user in contactList by resolveName, returns error, if can not.
// ResolveName can be as "name.surname", "name surname", "surname".
func (e Exchange2006) FindInContactList(resolveName string) error {
	_, err := e.GetMailBox(resolveName)
	return err
}
