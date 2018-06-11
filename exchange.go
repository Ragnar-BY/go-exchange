package go_exchange

import (
	"time"

	"github.com/Ragnar-BY/go-exchange/models"
)

// Exchange is interface for requests to MS Exchange.
type Exchange interface {
	GetRooms(roomlist string) ([]models.Room, error)
	GetRoomLists() ([]models.RoomList, error)
	GetRoomsAvailabilityByTime(rooms []models.Room, start time.Time, end time.Time) ([]models.CalendarEventArray, error)
}
