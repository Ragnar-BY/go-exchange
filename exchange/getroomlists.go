package exchange

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/Ragnar-BY/go-exchange/models"
)

// GetRoomLists make request to exchange and return roomlists.
func (e Exchange2006) GetRoomLists() ([]models.RoomList, error) {
	content, err := e.Post([]byte(roomListRequest()))
	if err != nil {
		return nil, fmt.Errorf("[GetRoomLists] cannot post: %v", err)
	}
	return parseRoomLists(content)
}

func roomListRequest() string {
	return `<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:m="http://schemas.microsoft.com/exchange/services/2006/messages" 
               xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <t:RequestServerVersion Version="Exchange2010" />
  </soap:Header>
  <soap:Body>
    <m:GetRoomLists />
  </soap:Body>
</soap:Envelope>`
}

func parseRoomLists(soap string) ([]models.RoomList, error) {
	decoder := xml.NewDecoder(bytes.NewBufferString(soap))
	var roomlists struct {
		XMLName   xml.Name
		Addresses []models.RoomList `xml:"Address"`
	}
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "RoomLists" {
				err := decoder.DecodeElement(&roomlists, &se)
				if err != nil {
					return nil, err
				}
				return roomlists.Addresses, nil
			}
		}
	}
	return nil, nil
}
