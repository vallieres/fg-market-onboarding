package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/vallieres/fg-market-onboarding/model"
)

const ShopifyURL = "https://fg-test-one.myshopify.com"

type CustomerService struct {
	shopifyAppToken        string
	shopifyStorefrontToken string
}

type MarketingConsent struct {
	MarketingOptInLevel string `json:"marketingOptInLevel"`
	MarketingState      string `json:"marketingState"`
}

type CustomerInput struct {
	Email                 string           `json:"email"`
	FirstName             string           `json:"firstName"`
	LastName              string           `json:"lastName"`
	EmailMarketingConsent MarketingConsent `json:"emailMarketingConsent"`
}

type GraphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

type ShopifyResponse struct {
	Data struct {
		CustomerCreate CustomerCreate `json:"customerCreate"`
	} `json:"data"`
}

type Address struct {
	Address1 string `json:"address1"`
	City     string `json:"city"`
	Country  string `json:"country"`
	Phone    string `json:"phone"`
	Zip      string `json:"zip"`
}

type Customer struct {
	ID                    string            `json:"id"`
	Email                 string            `json:"email"`
	Phone                 *string           `json:"phone"`
	TaxExempt             bool              `json:"taxExempt"`
	EmailMarketingConsent MarketingConsent  `json:"emailMarketingConsent"`
	FirstName             string            `json:"firstName"`
	LastName              string            `json:"lastName"`
	SmsMarketingConsent   *MarketingConsent `json:"smsMarketingConsent"`
	Addresses             []Address         `json:"addresses"`
}

type UserError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type CustomerCreate struct {
	UserErrors []UserError `json:"userErrors"`
	Customer   Customer    `json:"customer"`
}

func NewCustomerService(adminAPIToken string, storefrontToken string) *CustomerService {
	return &CustomerService{
		shopifyAppToken:        adminAPIToken,
		shopifyStorefrontToken: storefrontToken,
	}
}

func (c *CustomerService) Create(ctx context.Context, details model.OnboardPostBody) (string, error) {
	query := `mutation customerCreate($input: CustomerInput!) {
		customerCreate(input: $input) {
			userErrors {
				field
				message
			}
			customer {
				id
				email
				phone
				taxExempt
				emailMarketingConsent {
					marketingState
					marketingOptInLevel
					consentUpdatedAt
				}
				firstName
				lastName
				amountSpent {
					amount
					currencyCode
				}
				smsMarketingConsent {
					marketingState
					marketingOptInLevel
				}
				addresses {
					address1
					city
					country
					phone
					zip
				}
			}
		}
	}`

	// Create request payload
	input := CustomerInput{
		Email:     details.Email,
		FirstName: details.FirstName,
		LastName:  details.LastName,
		EmailMarketingConsent: MarketingConsent{
			MarketingOptInLevel: "CONFIRMED_OPT_IN",
			MarketingState:      "SUBSCRIBED",
		},
	}

	requestBody := GraphQLRequest{
		Query: query,
		Variables: map[string]any{
			"input": input,
		},
	}

	jsonData, errMarshal := json.Marshal(requestBody)
	if errMarshal != nil {
		fmt.Printf("Error marshaling JSON: %v\n", errMarshal)
		return "", errMarshal
	}

	// Create HTTP request
	req, errReq := http.NewRequestWithContext(ctx, http.MethodPost, ShopifyURL+"/admin/api/2024-10/graphql.json",
		bytes.NewBuffer(jsonData))
	if errReq != nil {
		fmt.Printf("Error creating request: %v\n", errReq)
		return "", errReq
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Shopify-Access-Token", c.shopifyAppToken)

	// Send request
	client := &http.Client{}
	resp, errDo := client.Do(req)
	if errDo != nil {
		fmt.Printf("Error sending request: %v\n", errDo)
		return "", errDo
	}
	defer resp.Body.Close()

	// Read response
	body, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		fmt.Printf("Error reading response: %v\n", errRead)
		return "", errRead
	}

	var response ShopifyResponse
	if errUnMarsh := json.Unmarshal(body, &response); errUnMarsh != nil {
		fmt.Printf("Error unmarshaling response: %v\n", errUnMarsh)
		return "", errUnMarsh
	}

	return response.Data.CustomerCreate.Customer.Email, nil
}
