package form3

import (
	"context"
	"fmt"
)

// AccountsService handles communication with the accounts related
// methods of the Form3 API.
//
// Form3 API docs: https://api-docs.form3.tech/api.html#organisation-accounts
type AccountsService service

// An Account represents a bank account that is registered with Form3.
// It is used to validate and allocate inbound payments.
type Account struct {
	Type           *string            `json:"type"`
	ID             *string            `json:"id"`
	OrganisationId *string            `json:"organisation_id"`
	Version        *int               `json:"version,omitempty"`
	Attributes     *AccountAttributes `json:"attributes"`
}

// Form3 API docs: https://api-docs.form3.tech/api.html#organisation-accounts-resource
type AccountAttributes struct {
	Country                     *string  `json:"country"`                                  // ISO 3166-1 code used to identify the domicile of the account, e.g. 'GB', 'FR'
	BaseCurrency                *string  `json:"base_currency,omitempty"`                  // ISO 4217 code used to identify the base currency of the account, e.g. 'GBP', 'EUR'
	AccountNumber               *string  `json:"account_number,omitempty"`                 // Local country bank identifier. Format depends on the country. Required for most countries.
	BankId                      *string  `json:"bank_id,omitempty"`                        // Identifies the type of bank ID being used, see here for allowed value for each country. Required value depends on country attribute.
	BankIdCode                  *string  `json:"bank_id_code,omitempty"`                   // Account number. A unique account number will automatically be generated if not provided. If provided, the account number is not validated.
	BIC                         *string  `json:"bic,omitempty"`                            // SWIFT BIC in either 8 or 11 character format e.g. 'NWBKGB22'
	IBAN                        *string  `json:"iban,omitempty"`                           // IBAN of the account. Will be calculated from other fields if not supplied.
	CustomerId                  *string  `json:"customer_id,omitempty"`                    // A free-format reference that can be used to link this account to an external system
	Name                        []string `json:"name"`                                     // Name of the account holder, up to four lines possible.
	AlternativeNames            []string `json:"alternative_names,omitempty"`              // Alternative primary account names, only used for UK Confirmation of Payee
	AccountClassification       *string  `json:"account_classification,omitempty"`         // Classification of account, only used for Confirmation of Payee (CoP)
	JointAccount                *bool    `json:"joint_account,omitempty"`                  // Flag to indicate if the account is a joint account, only used for Confirmation of Payee (CoP)
	AccountMatchingOptOut       *bool    `json:"account_matching_opt_out,omitempty"`       // Flag to indicate if the account has opted out of account matching, only used for Confirmation of Payee
	SecondaryIdentification     *string  `json:"secondary_identification,omitempty"`       // Additional information to identify the account and account holder, only used for Confirmation of Payee (CoP)
	Switched                    *bool    `json:"switched,omitempty"`                       // Flag to indicate if the account has been switched away from this organisation, only used for Confirmation of Payee (CoP)
	Status                      *string  `json:"status,omitempty"`                         // Status of the account. Inferred from the status field of the newest Account Event resource associated with the account. Always confirmed for older accounts where no Account Event resources are
	Title                       *string  `json:"title,omitempty"`                          // [Deprecated] The account holder's title, e.g. Ms, Dr, Mr. Only used for UK Confirmation of Payee (CoP). Superseded by name.
	FirstName                   *string  `json:"first_name,omitempty"`                     // [Deprecated] The account holder's first name, only used for UK Confirmation of Payee (CoP). Superseded by name.
	BankAccountName             *string  `json:"bank_account_name,omitempty"`              // [Deprecated] Primary account name, only used for UK Confirmation of Payee (CoP). Superseded by name.
	AlternativeBankAccountNames *string  `json:"alternative_bank_account_names,omitempty"` // [Deprecated] Alternative primary account names, only used for UK Confirmation of Payee. Superseded by alternative_names.
}

type AccountDetailsResponse struct {
	Data  *Account `json:"data"`
	Links *Links   `json:"links"`
}

type AccountDetailsListResponse struct {
	Data  []*Account `json:"data"`
	Links *Links     `json:"links"`
}

type AccountCreation struct {
	Data *Account `json:"data"`
}

type AccountCreationResponse struct {
	Data  *Account `json:"data"`
	Links *Links   `json:"links"`
}

// Register an existing bank account with Form3 or create a new one.
// The country attribute must be specified as a minimum. Depending on the country,
// other attributes such as bank_id and bic are mandatory.
// Form3 generates account numbers and IBANs, where appropriate, in the following cases:
// - If no account number or IBAN is provided, Form3 generates a valid account number (see below). If supported by the country, an IBAN is also generated.
// - If an account number is provided but the IBAN is empty, Form3 generates an IBAN if supported by the country.
// - If only an IBAN is provided, the account number will be left empty.
// Note that a given bank_id and bic need to be registered with Form3 and connected to your organisation ID.
// Form3 API docs: https://api-docs.form3.tech/api.html#organisation-accounts-create
func (s *AccountsService) Create(ctx context.Context, account *Account) (*Account, *Response, error) {
	u := fmt.Sprintf("/v%v/organisation/accounts", apiVersion)
	payload := &AccountCreation{Data: account}
	req, err := s.client.NewRequest("POST", u, payload)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", jsonApiMediaType)

	m := &AccountCreationResponse{}
	resp, err := s.client.Do(ctx, req, m)
	if err != nil {
		return nil, resp, err
	}

	return m.Data, resp, nil
}

// Get a single account using the account ID.
// Form3 API docs: https://api-docs.form3.tech/api.html#organisation-accounts-fetch
func (s *AccountsService) Fetch(ctx context.Context, id string) (*AccountDetailsResponse, *Response, error) {
	u := fmt.Sprintf("/v%v/organisation/accounts/%v", apiVersion, id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	accountDetails := new(AccountDetailsResponse)
	resp, err := s.client.Do(ctx, req, accountDetails)
	if err != nil {
		return nil, resp, err
	}

	return accountDetails, resp, nil
}

// List accounts with the ability to page.
// Form3 API docs: https://api-docs.form3.tech/api.html#organisation-accounts-list
func (s *AccountsService) List(ctx context.Context, options *ListOptions) (*AccountDetailsListResponse, *Response, error) {
	u, err := addOptions("organisation/accounts", options)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	accountDetailsList := new(AccountDetailsListResponse)
	resp, err := s.client.Do(ctx, req, accountDetailsList)
	if err != nil {
		return nil, resp, err
	}

	return accountDetailsList, resp, nil
}
