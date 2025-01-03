package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/vallieres/fg-market-onboarding/model"
)

const ShopifyURL = "https://fg-test-one.myshopify.com"

type MarketingConsent struct {
	MarketingOptInLevel string `json:"marketingOptInLevel"`
	MarketingState      string `json:"marketingState"`
}

type Metafield struct {
	Value     string `json:"value"`
	Namespace string `json:"namespace"`
	Key       string `json:"key"`
	Type      string `json:"type"`
}

type CustomerInput struct {
	Email                 string           `json:"email"`
	FirstName             string           `json:"firstName"`
	LastName              string           `json:"lastName"`
	Addresses             []Address        `json:"addresses"`
	EmailMarketingConsent MarketingConsent `json:"emailMarketingConsent"`
	Metafields            Metafield        `json:"metafields"`
}

type CustomerDogMetafields struct {
	Name             string  `json:"name"`
	Species          string  `json:"species"`
	Breed            string  `json:"breed"`
	Age              float64 `json:"age"`
	WeightLbs        int     `json:"weight_lbs"`
	HealthConditions string  `json:"health_conditions"`
}

type GraphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

type ShopifyResponse struct {
	Data struct {
		CustomerCreate struct {
			UserErrors []struct {
				Field   []string `json:"field"`
				Message string   `json:"message"`
			} `json:"userErrors"`
			Customer Customer `json:"customer"`
		} `json:"customerCreate"`
	} `json:"data"`
	Extensions struct {
		Cost struct {
			RequestedQueryCost int `json:"requestedQueryCost"`
			ActualQueryCost    int `json:"actualQueryCost"`
			ThrottleStatus     struct {
				MaximumAvailable   float64 `json:"maximumAvailable"`
				CurrentlyAvailable int     `json:"currentlyAvailable"`
				RestoreRate        float64 `json:"restoreRate"`
			} `json:"throttleStatus"`
		} `json:"cost"`
	} `json:"extensions"`
}

type Address struct {
	Address1 string `json:"address1,omitempty"`
	City     string `json:"city,omitempty"`
	Country  string `json:"country,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Zip      string `json:"zip,omitempty"`
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

type CustomerService struct {
	shopifyAppToken        string
	shopifyStorefrontToken string
}

func NewCustomerService(adminAPIToken string, storefrontToken string) *CustomerService {
	return &CustomerService{
		shopifyAppToken:        adminAPIToken,
		shopifyStorefrontToken: storefrontToken,
	}
}

//nolint:funlen
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
					  metafields(first: 10) {
						edges {
						  node {
							key
							namespace
							value
						  }
						}
					  }
					  phone
					  emailMarketingConsent {
						marketingState
						marketingOptInLevel
						consentUpdatedAt
					  }
					  firstName
					  lastName
					  smsMarketingConsent {
						marketingState
						marketingOptInLevel
					  }
					  addresses {
						country
						zip
					  }
					}
				  }
				}
`

	subscribed := "SUBSCRIBED"
	if !details.MailingList {
		subscribed = "NOT_SUBSCRIBED"
	}

	// Build Metafields JSON
	metafields := []CustomerDogMetafields{{
		Name:             details.PetName,
		Species:          details.PetSpecies,
		Breed:            details.PetBreed,
		Age:              details.PetAge,
		WeightLbs:        details.PetWeight,
		HealthConditions: details.PetHealthConditions,
	}}
	metafieldsJSON, errMarshalMeta := json.Marshal(metafields)
	if errMarshalMeta != nil {
		fmt.Println("Error marshalling Metafields to JSON:", errMarshalMeta)
		return "", errMarshalMeta
	}

	// Convert to string for display
	metafieldsJSONString := string(metafieldsJSON)

	// Create request payload
	addresses := []Address{{
		Zip:     details.ZipCode,
		Country: details.Country,
	}}
	input := CustomerInput{
		Email:     details.Email,
		FirstName: details.FirstName,
		LastName:  details.LastName,
		Addresses: addresses,
		EmailMarketingConsent: MarketingConsent{
			MarketingOptInLevel: "SINGLE_OPT_IN",
			MarketingState:      subscribed,
		},
		Metafields: Metafield{
			Namespace: "custom",
			Type:      "json",
			Key:       "animals",
			Value:     metafieldsJSONString,
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
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error creating customer: %v\n", resp.Status)
		return "", errors.New("error creating customer: " + resp.Status)
	}
	defer resp.Body.Close()

	// Read response
	body, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		fmt.Printf("Error reading response: %v\n", errRead)
		return "", errRead
	}

	var response ShopifyResponse
	var errorList []string
	if errUnMarsh := json.Unmarshal(body, &response); errUnMarsh != nil {
		fmt.Printf("Error unmarshaling response: %v\n", errUnMarsh)
		return "", errUnMarsh
	}
	if len(response.Data.CustomerCreate.UserErrors) > 0 {
		for _, userError := range response.Data.CustomerCreate.UserErrors {
			errorList = append(errorList, userError.Message)
		}
		return "", errors.New(strings.Join(errorList, "\n"))
	}

	return response.Data.CustomerCreate.Customer.Email, nil
}
