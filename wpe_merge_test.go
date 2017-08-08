package main_test

import (
	"testing"
	"gopkg.in/h2non/gock.v1"
	"github.com/adamrmelnyk/csv_builder"
)

func TestAccountInfo(t *testing.T) {
	defer gock.Off()
  gock.New("http://interview.wpengine.io/v1/accounts").
    Get("/777").
    Reply(200).
		JSON(`{"account_id": 777, "status": "fake status", "created_on": "2012-01-12"}`)

	account := main.AccountInfo(777)
	if account.Account_id != 777 {
		t.Error("Expected 777, got: ", account.Account_id)
	}
	if account.Status != "fake status" {
		t.Error("Expected fake status, got: ", account.Status)
	}
	if account.Created_on != "2012-01-12" {
		t.Error("Expected 2012-01-12, got: ", account.Created_on)
	}
}