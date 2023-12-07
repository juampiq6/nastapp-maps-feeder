package main

type LatLong struct {
	Lat  float64 `json:"latitude"`
	Long float64 `json:"longitude"`
}

// For this program purpose,\n only point type is used
type LatLongGEOJson struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}

// This represents a square's diagonal line
type PointsPair struct {
	Start LatLong `json:"start"`
	End   LatLong `json:"end"`
}
