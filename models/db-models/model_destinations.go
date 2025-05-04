package dbmodels

type Destination struct {
	ID        string   `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string   `json:"name"`
	Address   string   `json:"address"`
	Longitude string   `json:"longitude"`
	Latitude  string   `json:"latitude"`

	//foreign key - reviews
	// Reviews   []Review `gorm:"foreignKey:DestinationID"`
}