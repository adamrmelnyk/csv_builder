package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var client = &http.Client{Timeout: 30 * time.Second}
var apiURL = "http://interview.wpengine.io/v1/accounts"

// Account is an exported struct
// that contains all the same data present in the json api
type Account struct {
	Account_id int    `json:"account_id"`
	Status     string `json:"status"`
	Created_on string `json:"created_on"`
}

// AccountInfo is an exported function
// that calls an API and returns an Account struct
func AccountInfo(id int) Account {
	accountURL := fmt.Sprintf("%s/%d", apiURL, id)
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
	jsonErr := json.Unmarshal(body, &apiResponse)
	if jsonErr != nil {
		os.Stderr.WriteString("Failed to unmarshal JSON")
		os.Exit(1)
	}
	return apiResponse
}

func main() {
	if len(os.Args) < 3 {
		os.Stderr.WriteString("Wrong number of arguements.\nCorrect format: wpe_merge <input file.csv> <output file.csv>")
		os.Exit(1)
	}
	inputFile := os.Args[1]
	outputFile := os.Args[2]

	csvFile, err := os.Open(inputFile)
	if err != nil {
		os.Stderr.WriteString("Failed to open input file\n")
		os.Exit(1)
	}
	reader := csv.NewReader(csvFile)

	outFile, err := os.Create(outputFile)
	if err != nil {
		os.Stderr.WriteString("Failed to create output file\n")
		os.Exit(1)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()
	writer.Write([]string{"Account ID", "First Name", "Created On", "Status", "Status Set on"})

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			os.Stderr.WriteString("Failed to read file\n")
			os.Exit(1)
		}
		id, err := strconv.Atoi(record[0])
		if err == nil {
			// TODO: What if nothing comes back?
			account := AccountInfo(id)
			writer.Write([]string{record[0], record[2], record[3], account.Status, account.Created_on})
		}
	}
}
