package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
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

const minimumArgs = 3

// Account is an exported struct
// that contains all the same data present in the json api
type Account struct {
	AccountID int    `json:"account_id"`
	Status    string `json:"status"`
	CreatedOn string `json:"created_on"`
}

// breakAndThrowErr is an unexported function
// that writes to stderr and exits
func breakAndThrowErr(err string) {
	os.Stderr.WriteString(err)
	os.Exit(1)
}

// AccountInfo is an exported function
// that calls an API and returns an Account struct
func AccountInfo(id int) (Account, error) {
	apiResponse := Account{}
	accountURL := fmt.Sprintf("%s/%d", apiURL, id)
	response, err := client.Get(accountURL)
	if err != nil {
		return apiResponse, err
	}
	if response.StatusCode >= 400 {
		if response.StatusCode == 404 {
			statusErr := errors.New("404, Account not found")
			return apiResponse, statusErr
		}
		e := fmt.Sprintf("Error %d: Could not access account", response.StatusCode)
		statusErr := errors.New(e)
		return apiResponse, statusErr
	}
	defer response.Body.Close()

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return apiResponse, readErr
	}

	jsonErr := json.Unmarshal(body, &apiResponse)
	if jsonErr != nil {
		return apiResponse, jsonErr
	}
	return apiResponse, nil
}

// CombineDataInCSV is an unexported function
// that combines data from the input csv and the api
// and writes to the out file.
func combineDataInCSV(reader *csv.Reader, writer *csv.Writer) {
	hErr := writer.Write([]string{"Account ID", "First Name", "Created On", "Status", "Status Set on"})
	if hErr != nil {
		breakAndThrowErr("Failed to write headers to new csv")
	}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			breakAndThrowErr("Failed to read file\n")
		}
		writeToCsv(record, writer)
	}
}

// writeToCsv is an unexported function
// that takes an account record and
// combines it with data from the api
func writeToCsv(record []string, writer *csv.Writer) {
	id, err := strconv.Atoi(record[0])
	if err == nil && len(record) == 4 {
		account, err := AccountInfo(id)
		if err == nil {
			wErr := writer.Write([]string{record[0], record[2], record[3], account.Status, account.CreatedOn})
			if wErr != nil {
				breakAndThrowErr("Failed to write to new CSV")
			}
		}
	}
}

func main() {
	if len(os.Args) < minimumArgs {
		breakAndThrowErr("Wrong number of arguements.\nCorrect format: wpe_merge <input file.csv> <output file.csv>")
	}
	inputFile := os.Args[1]
	outputFile := os.Args[2]

	csvFile, err := os.Open(inputFile)
	if err != nil {
		breakAndThrowErr("Failed to open input file\n")
	}
	reader := csv.NewReader(csvFile)

	outFile, err := os.Create(outputFile)
	if err != nil {
		breakAndThrowErr("Failed to create output file\n")
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	combineDataInCSV(reader, writer)
}
