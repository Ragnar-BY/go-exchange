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
