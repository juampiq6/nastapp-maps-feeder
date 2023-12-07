package main

import (
	"strconv"
)

type PlaceResult struct {
	Id                    string            `bson:"id"`
	DisplayName           string            `bson:"displayName"`
	AddressComponents     AddressComponents `bson:"addressComponents"`
	FormattedAddress      string            `bson:"formattedAddress"`
	GoogleMapsUri         string            `bson:"googleMapsUri"`
	Location              LatLong           `bson:"location"`
	LocationGEOJson       LatLongGEOJson    `bson:"locationGEOJson"`
	ShortFormattedAddress string            `bson:"shortFormattedAddress"`
	WeekdayDescriptions   []string          `bson:"weekdayDescriptions"`
	PhoneNumber           string            `bson:"nationalPhoneNumber"`
	Rating                string            `bson:"rating"`
	UserRatingCount       string            `bson:"userRatingCount"`
}

type AddressComponents struct {
	Country    string `json:"country" bson:"country"`
	State      string `json:"state" bson:"state"`
	District   string `json:"district" bson:"district"`
	PostalCode string `json:"postalCode" bson:"postalCode"`
}

// DTO that comes from gmaps places api
type PlaceResultDTO struct {
	Id          string `json:"id"`
	DisplayName struct {
		Text string `json:"text"`
	} `json:"displayName"`
	AddressComponents     []map[string]any `json:"addressComponents"`
	FormattedAddress      string           `json:"formattedAddress"`
	GoogleMapsUri         string           `json:"googleMapsUri"`
	Location              LatLong          `json:"location"`
	ShortFormattedAddress string           `json:"shortFormattedAddress"`
	RegularOpeningHours   struct {
		WeekdayDescriptions []string `json:"weekdayDescriptions"`
	} `json:"regularOpeningHours"`
	PhoneNumber     string  `json:"nationalPhoneNumber"`
	Rating          float64 `json:"rating"`
	UserRatingCount int     `json:"userRatingCount"`
	// SKU text search preferred
	// Apparently not giving results
	// string "places.evChargeOptions",
	// string "places.fuelOptions",
	// string "places.servesCoffee",
}

type Response struct {
	Places []PlaceResultDTO `json:"places"`
}

// Equivalences from the Places Api
const (
	Country    string = "country"
	State             = "administrative_area_level_1"
	District          = "administrative_area_level_2"
	PostalCode        = "postal_code"
)

func parseAddressComponents(maps []map[string]any) AddressComponents {
	components := AddressComponents{}
	for _, m := range maps {
		// types is casted to an array of string
		if types, ok := m["types"].([]any); ok {
			for _, t := range types {
				if text, ok := m["longText"].(string); ok {
					switch t {
					case Country:
						components.Country = text
					case State:
						components.State = text
					case District:
						components.District = text
					case PostalCode:
						components.PostalCode = text
					}
				}
			}
		}
	}
	return components
}

func (p *PlaceResult) FromPlacesResultDTO(prdto *PlaceResultDTO) {
	p.FormattedAddress = prdto.FormattedAddress
	p.GoogleMapsUri = prdto.GoogleMapsUri
	p.Id = prdto.Id
	p.Location = prdto.Location
	p.PhoneNumber = prdto.PhoneNumber
	p.Rating = strconv.FormatFloat(prdto.Rating, 'f', 1, 64)
	p.ShortFormattedAddress = prdto.ShortFormattedAddress
	p.UserRatingCount = strconv.Itoa(prdto.UserRatingCount)

	p.LocationGEOJson = LatLongGEOJson{Type: "Point", Coordinates: []float64{prdto.Location.Long, prdto.Location.Lat}}
	p.DisplayName = prdto.DisplayName.Text
	p.WeekdayDescriptions = prdto.RegularOpeningHours.WeekdayDescriptions
	p.AddressComponents = parseAddressComponents(prdto.AddressComponents)
}
