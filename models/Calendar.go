package models

// CalendarEvent represents event, busytype is "Free" or "Busy".
type CalendarEvent struct {
	StartTime string
	EndTime   string
	BusyType  string
}

// CalendarEventArray is array of CalendarEvents.
type CalendarEventArray struct {
	CalendarEvent []CalendarEvent `xml:"CalendarEvent"`
}

// CalendarItem represents an Exchange calendar item.
type CalendarItem struct {
	ItemId struct {
		Id        string `xml:"Id,attr"`
		ChangeKey string `xml:"ChangeKey,attr"`
	} `xml:"ItemId"`
}
