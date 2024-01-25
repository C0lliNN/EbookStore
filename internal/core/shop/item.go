package shop

type Item struct {
	ID             string `gorm:"primaryKey"`
	Name           string
	Price          int64
	PreviewImageID string
	// OrderID it is only defined when the item is part of an order
	OrderID string `gorm:"primaryKey"`
}
