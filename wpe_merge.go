package main

import (
	"fmt"
	"os"
	"io"
	"io/ioutil"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"time"
	"reflect"
)

// TODO: Change var name
var client = &http.Client{ Timeout: 30 * time.Second }
var apiURL = "http://interview.wpengine.io/v1/accounts/"

type Account struct {
	Account_id int `json:"account_id"`
	Status string `json:"status"`
	Created_on string `json:"created_on"`
}

// TODO: Remove if we don't end up using it
type APIResponse struct {
	Count int `json:"count"`
	Next string `json:"next"`
	Previous string `json:"previous"`
	Results []Account `json:"results"`
}

func AccountInfo(id int) Account {
	accountURL := fmt.Sprintf("%s%d", apiURL, id)
	response, err := client.Get(accountURL)
	if err != nil {
		os.Stderr.WriteString("Failed to make request to API\n")
		os.Exit(1)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		os.Stderr.WriteString("Failed to read JSON")
		os.Exit(1)
	}

	apiResponse := Account{}
	json.Unmarshal(body, &apiResponse)
	fmt.Println(apiResponse)
	return apiResponse
}

func main() {
	if len(os.Args) != 3 {
		os.Stderr.WriteString("Wrong number of arguements.\nCorrect format: wpe_merge <input file.csv> <output file.csv>")
		os.Exit(1)
	}
	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// TODO: remove after testing
	fmt.Println(inputFile)
	fmt.Println(outputFile)
	// TODO: we should also check that they are csv

	csvFile, err := os.Open(inputFile)
	if err != nil {
		os.Stderr.WriteString("Failed to open input file\n")
		os.Exit(1)
	}
	reader := csv.NewReader(csvFile)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			os.Stderr.WriteString("Failed to read file\n")
			os.Exit(1)
		}

		fmt.Println(record[0])
		// TODO: Everything is coming back as a string. do we coerce into an int
		// or do we make the request and if nothing comes back then we don't care.
		// if we do check that we can ensure that we aren't doing work we don't need to
		fmt.Println(reflect.TypeOf(record[0]))
		AccountInfo(314159)
	}

	// TODO: combine data with csv
		// TODO: the api date format is yyyy-mm-dd but the csv is mm/dd/yy
		// TODO: I need the format the data into a hash first or something so I don't have to traverse the whole thing
	// TODO: write to new csv
}
