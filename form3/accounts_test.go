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

func TestIntegration_AccountsService_Create(t *testing.T) {
	client, _ := setupClientWithFakedApi()

	testAccount := &Account{
		ID:             String("1d50df61-db36-483c-9975-41141d691be1"),
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

	// Create account
	createAccountResponse, _, err := client.Accounts.Create(
		context.Background(),
		testAccount)

	if err != nil {
		t.Errorf("Accounts.Create returned error: %v", err)
		return
	}

	if *createAccountResponse.ID != *testAccount.ID {
		t.Errorf("Accounts.Create returned unexpected account: %v, want %v", *createAccountResponse.ID, *testAccount.ID)
	}

	// Clean up
	_, err = client.Accounts.Delete(context.Background(), *createAccountResponse.ID, *createAccountResponse.Version)
}

func TestIntegration_AccountsService_Create_Duplicate(t *testing.T) {
	client, _ := setupClientWithFakedApi()

	testAccount1, testAccount2 := &Account{
		ID:             String("3ff062d2-803f-4078-8067-699b3a0a0ba9"),
		Type:           String("accounts"),
		OrganisationId: String("eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"),
		Attributes: &AccountAttributes{
			Country:      String("GB"),
			BaseCurrency: String("GBP"),
			BankId:       String("400300"),
			BankIdCode:   String("GBDSC"),
			BIC:          String("NWBKGB22"),
		},
	},
		&Account{
			ID:             String("3ff062d2-803f-4078-8067-699b3a0a0ba9"),
			Type:           String("accounts"),
			OrganisationId: String("58d8a2c8-29ca-11eb-adc1-0242ac120002"),
			Attributes: &AccountAttributes{
				Country:      String("GB"),
				BaseCurrency: String("GBP"),
				BankId:       String("501600"),
				BankIdCode:   String("GBDXN"),
				BIC:          String("GDTKGB88"),
			},
		}

	// Create first account
	createAccountResponse1, _, err := client.Accounts.Create(
		context.Background(),
		testAccount1)

	if err != nil {
		t.Errorf("Accounts.Create returned error: %v", err)
	}

	// Create second account
	_, _, err = client.Accounts.Create(
		context.Background(),
		testAccount2)

	if !strings.Contains(err.Error(), "409") {
		t.Error("Accounts.Create did not return expected 409 error")
	}

	// Clean up
	_, err = client.Accounts.Delete(context.Background(), *createAccountResponse1.ID, *createAccountResponse1.Version)
}

func TestIntegration_AccountsService_Fetch(t *testing.T) {
	client, _ := setupClientWithFakedApi()

	testAccount := &Account{
		ID:             String("c041bbf4-19f3-4baa-9406-37a2a4be73e3"),
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

	// Create account
	createAccountResponse, _, err := client.Accounts.Create(
		context.Background(),
		testAccount)

	if err != nil {
		t.Errorf("Accounts.Create returned error: %v", err)
	}

	// Fetch account
	fetchAccountResponse, _, err := client.Accounts.Fetch(context.Background(), *createAccountResponse.ID)
	if err != nil {
		t.Errorf("Accounts.Fetch returned error: %v", err)
	}

	if *fetchAccountResponse.Data.ID != *testAccount.ID {
		t.Errorf("Accounts.Fetch returned %+v, want %+v", fetchAccountResponse.Data.ID, testAccount.ID)
	}

	// Clean up
	_, err = client.Accounts.Delete(context.Background(), *createAccountResponse.ID, *createAccountResponse.Version)
}

func TestIntegration_AccountsService_Fetch_Missing(t *testing.T) {
	client, _ := setupClientWithFakedApi()

	accountId := "c041bbf4-19f3-4baa-9406-37a2a4be73e3"

	// Fetch account
	_, httpResponse, _ := client.Accounts.Fetch(context.Background(), accountId)

	if httpResponse.StatusCode != 404 {
		t.Errorf("Accounts.Fetch returned status %v, want %v", httpResponse.StatusCode, 404)
	}
}

func TestIntegration_AccountsService_List(t *testing.T) {
	client, _ := setupClientWithFakedApi()

	testAccount1, testAccount2 := &Account{
		ID:             String("3d0d787f-32a8-4236-b6b7-39e1ca14ec04"),
		Type:           String("accounts"),
		OrganisationId: String("eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"),
		Attributes: &AccountAttributes{
			Country:      String("GB"),
			BaseCurrency: String("GBP"),
			BankId:       String("400300"),
			BankIdCode:   String("GBDSC"),
			BIC:          String("NWBKGB22"),
		},
	},
		&Account{
			ID:             String("58b3085c-e55a-4517-aead-aa6d5b9ea05e"),
			Type:           String("accounts"),
			OrganisationId: String("58d8a2c8-29ca-11eb-adc1-0242ac120002"),
			Attributes: &AccountAttributes{
				Country:      String("GB"),
				BaseCurrency: String("GBP"),
				BankId:       String("501600"),
				BankIdCode:   String("GBDXN"),
				BIC:          String("GDTKGB88"),
			},
		}

	// Create first account
	createAccountResponse1, _, err := client.Accounts.Create(
		context.Background(),
		testAccount1)

	if err != nil {
		t.Errorf("Accounts.Create returned error: %v", err)
	}

	// Create second account
	createAccountResponse2, _, err := client.Accounts.Create(
		context.Background(),
		testAccount2)

	if err != nil {
		t.Errorf("Accounts.Create returned error: %v", err)
	}

	// Test list
	listAccountsResponse, _, err := client.Accounts.List(context.Background(), &ListOptions{PageNumber: 0, PageSize: 1})
	if err != nil {
		t.Errorf("Accounts.List returned error: %v", err)
	}

	if len(listAccountsResponse.Data) < 1 {
		t.Errorf("Accounts.List returned %+v accounts, want %+v accounts", len(listAccountsResponse.Data), 1)
	}

	// Test list with paging
	listAccountsResponse2, _, err := client.Accounts.List(context.Background(), &ListOptions{PageNumber: 1, PageSize: 1})

	if err != nil {
		t.Errorf("Accounts.List returned error: %v", err)
	}

	if len(listAccountsResponse2.Data) < 1 {
		t.Errorf("Accounts.List returned %+v accounts, want %+v accounts", len(listAccountsResponse2.Data), 1)
	}

	if *listAccountsResponse.Data[0].ID == *listAccountsResponse2.Data[0].ID {
		t.Error("Accounts.List paging did not work returned same account on both pages")
	}

	// Clean up
	_, err = client.Accounts.Delete(context.Background(), *createAccountResponse1.ID, *createAccountResponse1.Version)
	_, err = client.Accounts.Delete(context.Background(), *createAccountResponse2.ID, *createAccountResponse2.Version)
}

func TestIntegration_AccountsService_Delete(t *testing.T) {
	client, _ := setupClientWithFakedApi()

	testAccount := &Account{
		ID:             String("5b207923-d1b9-4bdd-8c7e-60b28db655e7"),
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

	// Create the account
	createAccountResponse, _, err := client.Accounts.Create(
		context.Background(),
		testAccount)

	if err != nil {
		t.Errorf("Accounts.Create returned error: %v", err)
	}

	accountId := createAccountResponse.ID

	// Delete the account
	_, err = client.Accounts.Delete(context.Background(), *accountId, *createAccountResponse.Version)
	if err != nil {
		t.Errorf("Accounts.Delete returned error: %v", err)
	}
}

func TestIntegration_AccountsService_Delete_Missing(t *testing.T) {
	client, _ := setupClientWithFakedApi()

	accountId := "c77da171-6fb3-45ac-ba57-24d94eb13ba5"

	// Delete the account
	response, _ := client.Accounts.Delete(context.Background(), accountId, 1)

	if response.StatusCode != 204 {
		t.Errorf("Accounts.Create returned status %v, want %v", response.StatusCode, 204)
	}
}
