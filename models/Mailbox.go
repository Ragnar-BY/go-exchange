package models

// Mailbox identifies a mail-enabled object.
type Mailbox struct {
	Name         string `xml:"Name"`
	EmailAddress string `xml:"EmailAddress"`
}
