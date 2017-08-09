package main_test

import (
	"testing"

	main "github.com/adamrmelnyk/csv_builder"
	"gopkg.in/h2non/gock.v1"
)

func TestAccountInfo(t *testing.T) {
	defer gock.Off()
	gock.New("http://interview.wpengine.io/v1/accounts").
		Get("/777").
		Reply(200).
		JSON(`{"account_id": 777, "status": "fake status", "created_on": "2012-01-12"}`)

	firstAcc, err := main.AccountInfo(777)
	if firstAcc.AccountId != 777 {
		t.Error("Expected 777, got: ", firstAcc.AccountId)
	}
	if firstAcc.Status != "fake status" {
		t.Error("Expected fake status, got: ", firstAcc.Status)
	}
	if firstAcc.CreatedOn != "2012-01-12" {
		t.Error("Expected 2012-01-12, got: ", firstAcc.CreatedOn)
	}
	if err != nil {
		t.Error("Expected err to be nil, got: ", err)
	}

	gock.New("http://interview.wpengine.io/v1/accounts").
		Get("/777").
		Reply(200).
		JSON(`{"account": 777, "not_status": "fake status", "not_created_on": "2012-01-12"}`)

	secondAcc, _ := main.AccountInfo(777)
	if secondAcc.AccountId != 0 {
		t.Error("Expected 0, got: ", secondAcc.AccountId)
	}
	if secondAcc.Status != "" {
		t.Error("Expected '', got: ", secondAcc.Status)
	}
	if secondAcc.CreatedOn != "" {
		t.Error("Expected '', got: ", secondAcc.CreatedOn)
	}

	gock.New("http://interview.wpengine.io/v1/accounts").
		Get("/777").
		Reply(404).
		JSON(`{"detail": "Not found."}`)

	_, err = main.AccountInfo(777)
	if err == nil {
		t.Error("Expected '404, account not found', got: ", err)
	}

	gock.New("http://interview.wpengine.io/v1/accounts").
		Get("/777").
		Reply(500).
		JSON(`{"Error": "Internal Server Error"}`)

	_, err = main.AccountInfo(777)
	if err == nil {
		t.Error("Expected 'Error 500: Could not access account', git: ", err)
	}
}
