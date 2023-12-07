package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const urlStringDiscoverPlace = "https://places.googleapis.com/v1/places:searchText"
const urlStringPlaceDetails = "https://places.googleapis.com/v1/places/"

var apiKey string

// the details api call costs, so the amount of request is minimized by requesting only the places
// that were found on the discovery call
func makePlaceDetailCall(placeId string) []byte {
	setApiKeyFromEnv()
	fieldMasks := []string{
		"id",
		// SKU text search basic
		"displayName",
		"addressComponents",
		"formattedAddress",
		"googleMapsUri",
		"location",
		"shortFormattedAddress",
		// SKU text search advanced
		"regularOpeningHours.weekdayDescriptions",
		"nationalPhoneNumber",
		"rating",
		"userRatingCount",
		// SKU text search preferred
		// Apparently not giving results
		// "evChargeOptions",
		// "fuelOptions",
		// "servesCoffee",
	}

	req, err := http.NewRequest("GET", urlStringPlaceDetails+placeId+"?regionCode=ar&languageCode=es", nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Goog-FieldMask", strings.Join(fieldMasks, ","))
	// TODO: replace with env provided api key
	req.Header.Add("X-Goog-Api-Key", apiKey)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer res.Body.Close()
	resbody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if res.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Error %d : %s", res.StatusCode, resbody)
	}
	return resbody
}

// only the basic fields are requested, this way the cost of the request is zero
// this is done only for discovery of the places inside the square
func makeDiscoverPlacesApiCall(pointA LatLong, pointB LatLong) []byte {
	setApiKeyFromEnv()
	fieldMasks := []string{
		"places.id",
	}

	reqbody :=
		fmt.Sprintf(`{
		"textQuery": "gas station",
		"includedType": "gas_station",
		"strictTypeFiltering": true,
		"regionCode": "AR",
		"languageCode": "es",
		"locationRestriction": {
			"rectangle": {
				"low": {
			  		"latitude": %s,
			  		"longitude": %s
					},
				"high": {
			  		"latitude": %s,
			  		"longitude": %s
					}
		  		}
			}
		}`,
			strconv.FormatFloat(pointB.Lat, 'f', 8, 64),
			strconv.FormatFloat(pointA.Long, 'f', 8, 64),
			strconv.FormatFloat(pointA.Lat, 'f', 8, 64),
			strconv.FormatFloat(pointB.Long, 'f', 8, 64),
		)

	req, err := http.NewRequest("POST", urlStringDiscoverPlace, bytes.NewBuffer([]byte(reqbody)))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Goog-FieldMask", strings.Join(fieldMasks, ","))
	req.Header.Add("X-Goog-Api-Key", apiKey)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	defer res.Body.Close()
	resbody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if res.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Error %d : %s", res.StatusCode, resbody)
	}
	return resbody
}

func setApiKeyFromEnv() {
	if apiKey == "" {
		apiKey = os.Getenv("GMAPS_API_KEY")
		if apiKey == "" {
			panic("You must set your 'GMAPS_API_KEY' environment variable.")
		}
	}
}
