package models

// ItemID element contains the unique identifier and change key of an item in the Exchange store.
type ItemID struct {
	ID        string `xml:"Id,attr"`
	ChangeKey string `xml:"ChangeKey,attr"`
}
