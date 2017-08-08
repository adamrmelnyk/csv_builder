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

	first_account := main.AccountInfo(777)
	if first_account.Account_id != 777 {
		t.Error("Expected 777, got: ", first_account.Account_id)
	}
	if first_account.Status != "fake status" {
		t.Error("Expected fake status, got: ", first_account.Status)
	}
	if first_account.Created_on != "2012-01-12" {
		t.Error("Expected 2012-01-12, got: ", first_account.Created_on)
	}

	gock.New("http://interview.wpengine.io/v1/accounts").
		Get("/777").
		Reply(200).
		JSON(`{"account": 777, "not_status": "fake status", "not_created_on": "2012-01-12"}`)

	second_account := main.AccountInfo(777)
	if second_account.Account_id != 0 {
		t.Error("Expected 0, got: ", second_account.Account_id)
	}
	if second_account.Status != "" {
		t.Error("Expected '', got: ", second_account.Status)
	}
	if second_account.Created_on != "" {
		t.Error("Expected '', got: ", second_account.Created_on)
	}
}
