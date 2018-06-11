package exchange

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"

	"github.com/Ragnar-BY/go-exchange/models"
)

// GetRooms return rooms from roomlist as email.
func (e Exchange2006) GetRooms(roomlist string) ([]models.Room, error) {
	t, err := template.New("rooms").Parse(getRoomsRequest())
	if err != nil {
		return nil, fmt.Errorf("[GetRooms] cannot parse template %v", err)
	}
	var buf bytes.Buffer
	t.Execute(&buf, roomlist)
	content, err := e.Post(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("[GetRooms] cannot post %v", err)
	}
	return parseRooms(content)
}

func getRoomsRequest() string {
	return `<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xmlns:m="http://schemas.microsoft.com/exchange/services/2006/messages"
xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types"
xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
<soap:Header>
	<t:RequestServerVersion Version="Exchange2010" />
	</soap:Header>
	<soap:Body>
	<m:GetRooms>
	<m:RoomList>
	<t:EmailAddress>{{.}}</t:EmailAddress>
	</m:RoomList>
	</m:GetRooms>
	</soap:Body>
	</soap:Envelope>`
}

func parseRooms(soap string) ([]models.Room, error) {
	decoder := xml.NewDecoder(bytes.NewBufferString(soap))
	var rooms = make([]models.Room, 0)
	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "Room" {
				var room models.Room
				err := decoder.DecodeElement(&room, &se)
				if err != nil {
					return nil, err
				}
				rooms = append(rooms, room)
			}
		}
	}
	return rooms, nil
}
