package models

// Room represents room with name and email
type Room struct{
	Name string
	EmailAddress string
}
type ExchangeRoom struct {
	ID struct {
		Name string `xml:"Name"`
		EmailAddress string `xml:"EmailAddress"`
	} `xml:"Id"`
}

func (e ExchangeRoom)ToRoom()Room {
	return Room{Name:e.ID.Name,EmailAddress:e.ID.EmailAddress}
}
type CalendarEvent struct {
	StartTime string
	EndTime   string
	BusyType  string
}