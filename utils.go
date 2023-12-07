package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func keysFromMap[K comparable, V any](m *map[K]V) *[]K {
	list := []K{}
	for k := range *m {
		list = append(list, k)
	}
	return &list
}

func dumpDataToFile(data string, fileName string, append bool) {
	// creation or opening of file
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	check(err)
	if !append {
		f.Truncate(0)
	}
	defer f.Close()
	_, err = f.WriteString(data)
	check(err)
}

// Loads godotenv library to read env vars from .env file
func loadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func writeGeneratedPointsToFile(pairs *[]PointsPair) {
	sb := strings.Builder{}
	sb.WriteString(`{ "points": {`)
	for i := 0; i < len(*pairs); i++ {
		pair := (*pairs)[i]
		sb.WriteString(fmt.Sprintf(`"%d":`, i))
		marshalledPair, err := json.Marshal(pair)
		check(err)
		sb.WriteString(string(marshalledPair))
		sb.WriteString(",")
	}
	str := sb.String()
	woLastComma := str[:len(str)-1]
	resStr := woLastComma + "} }"
	dumpDataToFile(resStr, "pointsData.json", false)
}
