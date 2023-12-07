package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	loadDotEnv()
	discoverFlag, detailsFlag, index, differential := parseFlags()
	var discoveredIds []string

	if discoverFlag {
		discoveredIdsMap := discoverPlaces(index, differential)
		fmt.Println(len(*discoveredIdsMap), " were discovered in total!")
		if detailsFlag {
			discoveredIds = make([]string, len(*discoveredIdsMap))
			for k := range *discoveredIdsMap {
				discoveredIds = append(discoveredIds, k)
			}
		}
	}

	// once every discovered id is found, the second call to each place detail is done
	// TODO: read from discoveredPlaces.csv file if no discover flag = false
	if detailsFlag {
		var ids *[]string
		if discoveredIds == nil || len(discoveredIds) == 0 {
			ids = readIdsFromFile()
		}
		getPlacesDetail(ids)
		// uploadPlaceResults(results)
	}

}

func getPlacesDetail(ids *[]string) *[]any {
	var parsedResults []any
	for i, id := range *ids {
		var dto PlaceResultDTO
		// sometimes this query fails with 404 NOT FOUND, for no reason, just query separately the ids that failed and you'll get the response
		res := makePlaceDetailCall(id)
		err := json.Unmarshal(res, &dto)
		check(err)
		place := PlaceResult{}
		place.FromPlacesResultDTO(&dto)
		parsedResults = append(parsedResults, place)
		fmt.Printf("\033[G\033[K")
		fmt.Printf("%d detail queried", i)
		// we sleep some millisecond to avoid having issues with PlacesApi
		time.Sleep(200 * time.Millisecond)
	}
	return &parsedResults
}

func discoverPlaces(startingIndex int, differential float64) *map[string]bool {
	// a map is used as a set as a convenience to not repeat any id
	discoveredIdsMap := map[string]bool{}
	lastPairIndex := 0
	pairs, err := generateAllPointsPairs(PointAMendoza, PointBMendoza, differential)
	check(err)
	if startingIndex >= len(*pairs) {
		panic("The index provided is greater than or equal to the total points generated")
	}
	fmt.Printf("%d will be queried\n", len(*pairs)-startingIndex)
	fmt.Printf("Aprox minutes: %d\n", 300*(len(*pairs)-startingIndex)/1000/60)

	// a channel is created to react to the SO signal interrupt or terminate and save the discovered ids
	// and the last coordinate queried. useful for any interruption in the process to continue with the next index
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		dumpKeysToFile(discoveredIdsMap)
		dumpLastCoordinateToFile((*pairs)[lastPairIndex])
		os.Exit(1)
	}()

	for i := startingIndex; i < len(*pairs); i++ {
		lastPairIndex = startingIndex
		pair := (*pairs)[i]
		res := makeDiscoverPlacesApiCall(pair.Start, pair.End)
		type Response struct {
			Places []PlaceResultDTO `json:"places"`
		}
		var resObj Response

		err := json.Unmarshal(res, &resObj)
		check(err)
		for _, v := range resObj.Places {
			discoveredIdsMap[v.Id] = true
		}
		// we sleep some millisecond to avoid having issues with PlacesApi
		time.Sleep(200 * time.Millisecond)
		fmt.Printf("\033[G\033[K")
		fmt.Printf("%d queried", i)
		fmt.Printf(", %d discovered", len(discoveredIdsMap))
	}
	dumpKeysToFile(discoveredIdsMap)
	return &discoveredIdsMap
}

func parseFlags() (bool, bool, int, float64) {
	discoverPtr := flag.Bool("discover", true, "discover flag for running the discover task")
	detailsPtr := flag.Bool("details", true, "details flag for running the details query task and upload the results to db")
	indexPtr := flag.Int("lastIndex", 0, "index of generated points to start the discover query from. (default: 0)")
	differentialPtr := flag.Float64("differential", differentialRad, "differential, represented in radians, squares are made diff*diff. less differential is more chances to leave a place behind in a discover query due to the hard limit in results by the places api (20 res max). (constants.go -> differentialRad)")
	flag.Parse()
	return *discoverPtr, *detailsPtr, *indexPtr, *differentialPtr
}

// reads from discoveredPlaces.csv file if no discover flag = false
func readIdsFromFile() *[]string {
	fmt.Println("Reading discovered places from discoveredPlaces.csv")
	file, err := os.Open("discoveredPlaces.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// Create a new CSV reader
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	var data []string
	if len(records) == 0 {
		panic("No ids were found in the discoveredPlaces.csv file")
	}
	for _, record := range records {
		data = append(data, record[0])
	}

	return &data
}

func dumpKeysToFile(m map[string]bool) {
	dumpDataToFile("\""+strings.Join(*(keysFromMap(&m)), "\"\n\"")+"\"\n", "discoveredPlaces.csv", true)
}

func dumpLastCoordinateToFile(pair PointsPair) {
	jsonObj, _ := json.Marshal(pair)
	dumpDataToFile(string(jsonObj), "lastPointPair.json", false)
}
