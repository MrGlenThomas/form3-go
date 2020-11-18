package form3

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestUnit_AccountsService_Create(t *testing.T) {
	client, mux, _, teardown := setupClientWithStubbedApi()
	defer teardown()

	mux.HandleFunc("/organisation/accounts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", jsonApiMediaType)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
		{
			"data": {
				"id": "d97a4470-299f-11eb-adc1-0242ac120002",
				"type": "accounts",
				"organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
				"version": 7,
				"attributes": {
					"country": "GB",
					"base_currency": "GBP"
				}
			}
		}`))
	})

	createdAccount, _, err := client.Accounts.Create(context.Background(), &Account{})
	if err != nil {
		t.Errorf("Create returned error: %v", err)
	}

	want := &Account{
		ID:             String("d97a4470-299f-11eb-adc1-0242ac120002"),
		Type:           String("accounts"),
		OrganisationId: String("eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"),
		Version:        Int(7),
		Attributes: &AccountAttributes{
			Country:      String("GB"),
			BaseCurrency: String("GBP"),
		},
	}
	if !reflect.DeepEqual(createdAccount, want) {
		t.Errorf("Create = %+v, want %+v", createdAccount, want)
	}
}

func TestUnit_AccountsService_Fetch(t *testing.T) {
	client, mux, _, teardown := setupClientWithStubbedApi()
	defer teardown()

	mux.HandleFunc("/organisation/accounts/d97a4470-299f-11eb-adc1-0242ac120002", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", jsonApiMediaType)
		fmt.Fprint(w, `
		{
			"data": {
				"id":"d97a4470-299f-11eb-adc1-0242ac120002",
				"type": "accounts",
				"organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
				"version": 7,
				"attributes": {
					"country": "GB",
					"base_currency": "GBP"
				}
			}
		}`)
	})

	accountResponse, _, err := client.Accounts.Fetch(context.Background(), "d97a4470-299f-11eb-adc1-0242ac120002")
	if err != nil {
		t.Errorf("Accounts.Fetch returned error: %v", err)
	}

	want := &AccountDetailsResponse{
		Data: &Account{
			ID:             String("d97a4470-299f-11eb-adc1-0242ac120002"),
			Type:           String("accounts"),
			OrganisationId: String("eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"),
			Version:        Int(7),
			Attributes: &AccountAttributes{
				Country:      String("GB"),
				BaseCurrency: String("GBP"),
			},
		},
	}

	if !reflect.DeepEqual(accountResponse, want) {
		t.Errorf("Accounts.Fetch returned %+v, want %+v", accountResponse, want)
	}
}