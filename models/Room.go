package models

import (
	"encoding/xml"
)

// Room represents room.
type Room struct {
	Name         string
	EmailAddress string
}

// RoomList represents roomlist.
type RoomList struct {
	Name         string `xml:",Name"`
	EmailAddress string `xml:",EmailAddress"`
}

// UnmarshalXML unmarshals xml from exchange to Room.
func (r *Room) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type ExchangeRoom struct {
		ID struct {
			Name         string `xml:"Name"`
			EmailAddress string `xml:"EmailAddress"`
		} `xml:"Id"`
	}
	exroom := ExchangeRoom{}
	err := d.DecodeElement(&exroom, &start)
	if err != nil {
		return err
	}
	*r = Room{Name: exroom.ID.Name, EmailAddress: exroom.ID.EmailAddress}
	return nil
}
