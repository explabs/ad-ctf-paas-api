package models

type Service struct {
	Name string  `bson:"name"`
	Cost float64 `bson:"cost"`
	HP   float64 `bson:"hp"`
}
