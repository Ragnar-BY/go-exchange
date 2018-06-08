package exchange

import (
	"github.com/Ragnar-BY/go-exchange/models"
	"fmt"
	"encoding/xml"
	"bytes"
)

type Exchange2006 struct {
	User string
	Password string
	Url string
}

func (e Exchange2006) GetRooms() ([]models.Room, error) {
	content,err:=e.Post([]byte(getRoomsRequest()))
	if err!= nil{
		return nil,err
	}
	fmt.Println(content)

	fmt.Println(parseRooms(content))
	return nil, nil
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
	<t:EmailAddress>Tolstogo10@itechart-group.com</t:EmailAddress>
	</m:RoomList>
	</m:GetRooms>
	</soap:Body>
	</soap:Envelope>`
}


func parseRooms(soap string)  []models.Room{
	decoder := xml.NewDecoder(bytes.NewBufferString(soap))


	var rooms = make([]models.Room,0)

	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "Room" {
				var room models.ExchangeRoom
				decoder.DecodeElement(&room, &se)
				rooms = append(rooms,models.Room{Name:room.ID.Name,EmailAddress:room.ID.EmailAddress})

			}
		}
	}
	return rooms
}
