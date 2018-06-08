package go_exchange

import "github.com/Ragnar-BY/go-exchange/models"

// Exchange is interface for requests to MS Exchange.
type Exchange interface {
	GetRooms() ([]models.Room,error)
}
