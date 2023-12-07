package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Generates the point pairs from an area dividing it to `differentalRad`
// Useful for querying the Places API or any other API that restricts result, as with these point pairs you can iterate
// throw all the area and obtain all the places data for each piece of area
func generateAllPointsPairs(pointA LatLong, pointB LatLong, differentialRad float64) (*[]PointsPair, error) {
	var points = []PointsPair{}
	// Longitude loop
	currentLong := pointA.Long
	var nextLong float64
	limitLong := pointB.Long
	shouldSubstractLong := currentLong > limitLong
	latPieces := amountOfPointsLat(pointA, pointB, differentialRad)
	longPieces := amountOfPointsLong(pointA, pointB, differentialRad)
	for i := 0; i < longPieces; i++ {

		if shouldSubstractLong {
			nextLong = currentLong - differentialRad
		} else {
			nextLong = currentLong + differentialRad
		}

		// Latitude loop
		currentLat := pointA.Lat
		var nextLat float64
		limitLat := pointB.Lat
		shouldSubstractLat := currentLat > limitLat

		for i := 0; i < latPieces; i++ {

			// Increment for next iteration (Lat)
			if shouldSubstractLat {
				nextLat = currentLat - differentialRad
			} else {
				nextLat = currentLat + differentialRad
			}
			// The point is added to the list
			points = append(points, PointsPair{Start: LatLong{Lat: currentLat, Long: currentLong}, End: LatLong{Lat: nextLat, Long: nextLong}})
			currentLat = nextLat
		}
		// Increment for next iteration (Long)
		currentLong = nextLong
	}

	if expected := expectedResultingPointsLength(pointA, pointB, differentialRad); expected != len(points) {
		fmt.Println(len(points))
		return nil, errors.New(strings.Join([]string{"Expected amount of points does not equal to resulting amount of points. Should be:", strconv.Itoa(expected), " , is:", strconv.Itoa(len(points))}, ""))
	}
	return &points, nil
}

// Calculates the amount of points that result in the division of the whole area to `differentalRad`
func expectedResultingPointsLength(pointA LatLong, pointB LatLong, differential float64) int {
	return amountOfPointsLong(pointA, pointB, differential) * amountOfPointsLat(pointA, pointB, differential)
}

// Calculates the amount of points that can result from dividing the area to `differentalRad`, keeping the longitude constant
func amountOfPointsLat(pointA LatLong, pointB LatLong, differentialRad float64) int {
	height := math.Abs(math.Abs(pointA.Lat) - math.Abs(pointB.Lat))
	return int(math.Ceil(height / differentialRad))
}

// Calculates the amount of points that can result from dividing the area to `differentalRad`, keeping the latitude constant
func amountOfPointsLong(pointA LatLong, pointB LatLong, differentialRad float64) int {
	width := math.Abs(math.Abs(pointA.Long) - math.Abs(pointB.Long))
	return int(math.Ceil(width / differentialRad))
}

// Example of how it should be used for an area that roughly covers Mendoza province
func exampleUsageMendoza() {
	points, err := generateAllPointsPairs(PointAMendoza, PointBMendoza, differentialRad)
	if err != nil {
		panic(strings.Join([]string{"Error generating points", err.Error()}, ""))
	}
	for i := 0; i < len(*points); i++ {
		point := (*points)[i]

		fmt.Println("LatLng(", point.Start.Lat, ",", point.Start.Long, "),")
		fmt.Println("LatLng(", point.End.Lat, ",", point.End.Long, "),")
	}
}
