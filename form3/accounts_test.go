package form3

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
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

func TestUnit_AccountsService_Create_BadRequest(t *testing.T) {
	client, mux, _, teardown := setupClientWithStubbedApi()
	defer teardown()

	mux.HandleFunc("/organisation/accounts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", jsonApiMediaType)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error_message": "Your request is not good",
			"error_code": "I don't like it!"
		}`))
	})

	_, _, err := client.Accounts.Create(context.Background(), &Account{})
	if err == nil {
		t.Error("Create did not return error")
	}
	if !strings.Contains(err.Error(), "Your request is not good") {
		t.Errorf("Create returned error: %v, missing error_message %v", err.Error(), "Your request is not good")
	}
	if !strings.Contains(err.Error(), "I don't like it!") {
		t.Errorf("Create returned error: %v, missing error_code %v", err.Error(), "I don't like it!")
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

func TestUnit_AccountsService_List(t *testing.T) {
	client, mux, _, teardown := setupClientWithStubbedApi()
	defer teardown()

	mux.HandleFunc("/organisation/accounts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", jsonApiMediaType)
		testFormValues(t, r, values{
			"page[number]": "1",
			"page[size]":   "10",
		})
		fmt.Fprint(w, `
		{
			"data": [
				{
					"id":"ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
					"type": "accounts",
					"organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
					"version": 7,
					"attributes": {
						"country": "GB",
						"base_currency": "GBP"
					}
				}
			]
		}`)
	})

	accountsListResponse, _, err := client.Accounts.List(context.Background(), &ListOptions{PageNumber: 1, PageSize: 10})
	if err != nil {
		t.Errorf("Accounts.List returned error: %v", err)
	}

	want := &AccountDetailsListResponse{
		Data: []*Account{
			{
				ID:             String("ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"),
				Type:           String("accounts"),
				OrganisationId: String("eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"),
				Version:        Int(7),
				Attributes: &AccountAttributes{
					Country:      String("GB"),
					BaseCurrency: String("GBP"),
				},
			},
		},
	}
	if !reflect.DeepEqual(accountsListResponse, want) {
		t.Errorf("Accounts.List returned %+v, want %+v", accountsListResponse, want)
	}
}

func TestUnit_AccountsService_Delete(t *testing.T) {
	client, mux, _, teardown := setupClientWithStubbedApi()
	defer teardown()

	mux.HandleFunc("/organisation/accounts/d97a4470-299f-11eb-adc1-0242ac120002", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testFormValues(t, r, values{
			"version": "1",
		})
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Accounts.Delete(context.Background(), "d97a4470-299f-11eb-adc1-0242ac120002", 1)
	if err != nil {
		t.Errorf("Accounts.Delete returned error: %v", err)
	}
}

func TestIntegration_AccountsService(t *testing.T) {
	client, _ := setupClientWithFakedApi()

	testAccount := &Account{
		ID:             String("ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"),
		Type:           String("accounts"),
		OrganisationId: String("eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"),
		Attributes: &AccountAttributes{
			Country:      String("GB"),
			BaseCurrency: String("GBP"),
			BankId:       String("400300"),
			BankIdCode:   String("GBDSC"),
			BIC:          String("NWBKGB22"),
		},
	}

	createAccountResponse, _, err := client.Accounts.Create(
		context.Background(),
		testAccount)

	if err != nil {
		t.Errorf("Accounts.List returned error: %v", err)
	}

	accountId := createAccountResponse.ID

	fetchAccountResponse, _, err := client.Accounts.Fetch(context.Background(), *accountId)
	if err != nil {
		t.Errorf("Accounts.Fetch returned error: %v", err)
	}

	if *fetchAccountResponse.Data.ID != *testAccount.ID {
		t.Errorf("Accounts.Fetch returned %+v, want %+v", fetchAccountResponse.Data.ID, testAccount.ID)
	}

	listAccountsResponse, _, err := client.Accounts.List(context.Background(), &ListOptions{PageNumber: 0, PageSize: 1})
	if err != nil {
		t.Errorf("Accounts.List returned error: %v", err)
	}

	if len(listAccountsResponse.Data) < 1 {
		t.Errorf("Accounts.List returned %+v, want %+v", len(listAccountsResponse.Data), 1)
	}

	_, err = client.Accounts.Delete(context.Background(), *accountId, *createAccountResponse.Version)
	if err != nil {
		t.Errorf("Accounts.Delete returned error: %v", err)
	}
}
