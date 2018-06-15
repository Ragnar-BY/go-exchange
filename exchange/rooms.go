package exchange

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"
	"time"

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

// GetRoomLists make request to exchange and return roomlists.
func (e Exchange2006) GetRoomLists() ([]models.RoomList, error) {
	content, err := e.Post([]byte(roomListRequest))
	if err != nil {
		return nil, fmt.Errorf("[GetRoomLists] cannot post: %v", err)
	}
	return parseRoomLists(content)
}

var roomListRequest = `<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:m="http://schemas.microsoft.com/exchange/services/2006/messages" 
               xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <t:RequestServerVersion Version="Exchange2010" />
  </soap:Header>
  <soap:Body>
    <m:GetRoomLists />
  </soap:Body>
</soap:Envelope>`

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

// GetRoomsAvailabilityByTime for array of rooms return array of events, which consists of moments, when room is busy.
func (e Exchange2006) GetRoomsAvailabilityByTime(rooms []models.Room, start time.Time, end time.Time) ([]models.Room, error) {
	t, err := template.New("roomsav").Parse(getRoomsAvailabilityRequest)
	if err != nil {
		return nil, fmt.Errorf("[GetRoomsAvailability] cannot create template %v", err)
	}
	data := struct {
		Rooms []models.Room
		Start string
		End   string
	}{
		Rooms: rooms,
		Start: start.Format("2006-01-02T15:04:05"),
		End:   end.Format("2006-01-02T15:04:05"),
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return nil, fmt.Errorf("[GetRoomsAvailabilityByTime] cannot parse template: %v", err)
	}
	content, err := e.Post(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("[GetRoomsAvailabilityByTime] cannot post: %v", err)
	}
	eventarrays, err := parseRoomAvailability(content)
	if err != nil {
		return nil, fmt.Errorf("[GetRoomsAvailabilityByTime] cannot parse: %v", err)
	}
	newRooms := make([]models.Room, 0)
	for i, r := range rooms {
		newroom := models.NewRoom(r.Name, r.EmailAddress)
		newroom.SetCalendarEvents(&eventarrays[i])
		newRooms = append(newRooms, *newroom)
	}
	return newRooms, nil
}

var getRoomsAvailabilityRequest = `<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xmlns:m="http://schemas.microsoft.com/exchange/services/2006/messages"
xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types"
xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
<soap:Header>
	<t:RequestServerVersion Version="Exchange2010" />
	<t:TimeZoneContext>
	<t:TimeZoneDefinition Id="Belarus Standard Time">
	</t:TimeZoneDefinition>
	</t:TimeZoneContext>
	</soap:Header>
	<soap:Body>
	<m:GetUserAvailabilityRequest>
	<m:MailboxDataArray>{{range .Rooms}}	
	<t:MailboxData>
	<t:Email>
	<t:Address>{{.EmailAddress}}</t:Address>
	</t:Email>
	<t:AttendeeType>Room</t:AttendeeType>
	<t:ExcludeConflicts>false</t:ExcludeConflicts>
	</t:MailboxData>{{end}}	
	</m:MailboxDataArray>
	<t:FreeBusyViewOptions>
	<t:TimeWindow>
	<t:StartTime>{{.Start}}</t:StartTime>
	<t:EndTime>{{.End}}</t:EndTime>
	</t:TimeWindow>
	<t:MergedFreeBusyIntervalInMinutes>30</t:MergedFreeBusyIntervalInMinutes>
	<t:RequestedView>FreeBusy</t:RequestedView>
	</t:FreeBusyViewOptions>
	</m:GetUserAvailabilityRequest>
	</soap:Body>
	</soap:Envelope>`

func parseRoomAvailability(soap string) ([]models.CalendarEventArray, error) {
	decoder := xml.NewDecoder(bytes.NewBufferString(soap))
	type FreeView struct {
		CalendarEventArray models.CalendarEventArray `xml:"CalendarEventArray"`
	}
	var events = make([]models.CalendarEventArray, 0)
	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "FreeBusyView" {
				var fr FreeView
				err := decoder.DecodeElement(&fr, &se)
				if err != nil {
					return nil, err
				}
				events = append(events, fr.CalendarEventArray)
			}
		}
	}
	return events, nil
}

// GetFreeRoomsByTime return all rooms, that are available since start to end time.
func (e Exchange2006) GetFreeRoomsByTime(rooms []models.Room, start time.Time, end time.Time) ([]models.Room, error) {
	roomAv, err := e.GetRoomsAvailabilityByTime(rooms, start, end)
	if err != nil {
		return nil, err
	}
	newRooms := make([]models.Room, 0)
	for _, r := range roomAv {
		if r.CalendarEvents != nil {
			if len(r.CalendarEvents.CalendarEvent) == 0 {
				newRooms = append(newRooms, r)
			}
		}
	}
	return newRooms, nil
}
