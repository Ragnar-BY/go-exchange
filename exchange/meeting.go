package exchange

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"text/template"
	"time"

	"github.com/Ragnar-BY/go-exchange/models"
)

// AddMeeting create meeting in room with required attendees with subject from start to end time.
// There is no error if room or attendee is busy.
func (e Exchange2006) AddMeeting(room models.Room, attendees []string, start time.Time, end time.Time, subject string) (*models.CalendarItem, error) {
	data := struct {
		Room      models.Room
		Attendees []string
		Start     string
		End       string
		Subject   string
	}{
		Room:      room,
		Attendees: attendees,
		Start:     start.Format("2006-01-02T15:04:05"),
		End:       end.Format("2006-01-02T15:04:05"),
		Subject:   subject,
	}
	t, err := template.New("meetings").Parse(addMeetingRequest())
	if err != nil {
		return nil, fmt.Errorf("[AddMeeting] cannot create template %v", err)
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return nil, fmt.Errorf("[AddMeeting] cannot parse template: %v", err)
	}
	response, err := e.Post(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("[AddMeeting] cannot post %v", err)
	}
	item, err := parseAddMeetingResponse(response)
	return item, err
}
func addMeetingRequest() string {
	return `<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:m="http://schemas.microsoft.com/exchange/services/2006/messages" 
       xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <t:RequestServerVersion Version="Exchange2007_SP1" />
    <t:TimeZoneContext>
      <t:TimeZoneDefinition Id="Belarus Standard Time" />
    </t:TimeZoneContext>
  </soap:Header>
  <soap:Body>
    <m:CreateItem SendMeetingInvitations="SendToAllAndSaveCopy">
      <m:Items>
        <t:CalendarItem>
          <t:Subject>{{.Subject}}</t:Subject>
          <t:Body BodyType="HTML">Meeting is created by the bot</t:Body>
          <t:ReminderMinutesBeforeStart>30</t:ReminderMinutesBeforeStart>
          <t:Start>{{.Start}}</t:Start>
          <t:End>{{.End}}</t:End>
          <t:Location>{{.Room.Name}}</t:Location>
          <t:RequiredAttendees>			
			<t:Attendee>
				<t:Mailbox>
					<t:EmailAddress>{{.Room.EmailAddress}}</t:EmailAddress>
				</t:Mailbox>
			</t:Attendee>{{range .Attendees}}
			<t:Attendee>
				<t:Mailbox>
					<t:EmailAddress>{{.}}</t:EmailAddress>
				</t:Mailbox>
			</t:Attendee>{{end}}
          </t:RequiredAttendees>          
          <t:MeetingTimeZone TimeZoneName="Belarus Standard Time" />
        </t:CalendarItem>
      </m:Items>
    </m:CreateItem>
  </soap:Body>
</soap:Envelope>`
}

func parseAddMeetingResponse(response string) (*models.CalendarItem, error) {

	type CreateItemResponseMessage struct {
		ResponseClass string `xml:"ResponseClass,attr"`
		ResponseCode  string `xml:"ResponseCode"`
		MessageText   string `xml:"MessageText"`
		Items         struct {
			CalendarItem models.CalendarItem `xml:"CalendarItem"`
		} `xml:"Items"`
	}
	var cirm CreateItemResponseMessage
	decoder := xml.NewDecoder(bytes.NewBufferString(response))
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "CreateItemResponseMessage" {
				err := decoder.DecodeElement(&cirm, &se)
				if err != nil {
					return nil, err
				}
				if cirm.ResponseClass == "Error" {
					return nil, errors.New(cirm.MessageText)
				}
				return &cirm.Items.CalendarItem, nil
			}
		}
	}
	return nil, errors.New("cannot find CreateItemResponseMessage in response")
}
