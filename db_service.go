package main

import (
	"fmt"
	"nastapp-maps-feeder/db"
)

const collectionName = "gas_stations"

func uploadPlaceResults(results *[]interface{}) {
	fmt.Println("Uploading to results to mongo...")
	db.InsertManyIntoCollection(collectionName, *results)
}
