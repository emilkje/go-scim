package e2etest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/intility/scim/server"
)

type ListResponse struct {
	Resources    []Resource `json:"Resources"`
	Schemas      []string   `json:"schemas"`
	ItemsPerPage int        `json:"itemsPerPage"`
	StartIndex   int        `json:"startIndex"`
	TotalResults int        `json:"totalResults"`
}

type Resource struct {
	Attributes  []Attribute `json:"attributes"`
	Description string      `json:"description"`
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Schemas     []string    `json:"schemas"`
}

type Attribute struct {
	Description string `json:"description"`
	Mutability  string `json:"mutability"`
	Name        string `json:"name"`
	Returned    string `json:"returned"`
	Type        string `json:"type"`
	Uniqueness  string `json:"uniqueness"`
	CaseExact   bool   `json:"caseExact"`
	MultiValued bool   `json:"multiValued"`
	Required    bool   `json:"required"`
}

type ValidatorFunc func(res Resource) error

var schemaValidators = map[string]ValidatorFunc{
	"urn:ietf:params:scim:schemas:core:2.0:User":                 validateUserSchema,
	"urn:ietf:params:scim:schemas:core:2.0:Group":                validateGroupSchema,
	"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": validateEnterpriseUser,
}

func TestUserCreate(t *testing.T) {
	scimServer := server.NewServer(&testlogger{})
	srv := httptest.NewServer(scimServer)

	client := srv.Client()

	resp, err := client.Get(srv.URL + "/Schemas")
	if err != nil {
		t.Errorf("failed to call /Schemas: %s", err.Error())
		return
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code: %d", resp.StatusCode)
		return
	}

	var content ListResponse
	err = json.NewDecoder(resp.Body).Decode(&content)
	if err != nil {
		t.Errorf("failed to decode response: %s", err.Error())
	}

	for _, resource := range content.Resources {
		if err := validate(t, resource); err != nil {
			t.Errorf("failed to validate schema: %s", err)
		}
	}
}

func validate(t *testing.T, resource Resource) error {
	t.Helper()

	if validate, ok := schemaValidators[resource.ID]; ok {
		return validate(resource)
	}

	return errors.New("unknown schema: " + resource.ID)
}

func validateGroupSchema(resource Resource) error {
	var err error
	wellknownAttrs := []string{
		"userName",
		"phoneNumbers",
		"surname",
	}

	for _, attr := range resource.Attributes {
		if !slices.Contains(wellknownAttrs, attr.Name) {
			err = errors.Join(errors.New("unknown attr: "+attr.Name), err)
		}
	}

	return err
}

func validateUserSchema(resource Resource) error {
	var err error
	wellknownAttrs := []string{
		"name",
		"userName",
		"displayName",
		"userType",
		"title",
		"profileUrl",
		"nickName",
		"preferredLanguage",
		"phoneNumbers",
		"emails",
		"locale",
		"timezone",
		"active",
		"photos",
	}

	discoveredAttrs := []string{}

	for _, attr := range resource.Attributes {
		discoveredAttrs = append(discoveredAttrs, attr.Name)
		if !slices.Contains(wellknownAttrs, attr.Name) {
			err = errors.Join(errors.New("unknown attr: "+attr.Name), err)
		}
	}

	for _, attr := range wellknownAttrs {
		if !slices.Contains(discoveredAttrs, attr) {
			err = errors.Join(errors.New("missing attribute in user schema: "+attr), err)
		}
	}

	return err
}

func validateEnterpriseUser(resource Resource) error {
	return nil
}

type CreateUserRequest struct {
	Name       Name     `json:"name"`
	ExternalID string   `json:"externalId"`
	UserName   string   `json:"userName"`
	Meta       Meta     `json:"meta"`
	Schemas    []string `json:"schemas"`
	Emails     []Email  `json:"emails"`
	Roles      []Role   `json:"roles"`
	Active     bool     `json:"active"`
}

type Email struct {
	Type    string `json:"type"`
	Value   string `json:"value"`
	Primary bool   `json:"primary"`
}

type Meta struct {
	ResourceType string `json:"resourceType"`
}

type Name struct {
	Formatted  string `json:"formatted"`
	FamilyName string `json:"familyName"`
	GivenName  string `json:"givenName"`
}

type Role struct {
	// Define attributes for Role if needed, currently empty based on the payload
}

func TestPostUser(t *testing.T) {
	scimServer := server.NewServer(&testlogger{})
	srv := httptest.NewServer(scimServer)

	client := srv.Client()
	// Create an instance of PostPayload with sample data
	payload := CreateUserRequest{
		Schemas: []string{
			"urn:ietf:params:scim:schemas:core:2.0:User",
			"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User",
		},
		ExternalID: "0a21f0f2-8d2a-4f8e-bf98-7363c4aed4ef",
		UserName:   "mctestie",
		Active:     true,
		Emails: []Email{
			{
				Primary: true,
				Type:    "work",
				Value:   "mctesties@testuser.com",
			},
			{
				Value: "testies@foo.com",
			},
		},
		Meta: Meta{
			ResourceType: "User",
		},
		Name: Name{
			Formatted:  "Testuser McTesties",
			FamilyName: "McTesties",
			GivenName:  "Testuser",
		},
		Roles: []Role{}, // No roles provided in this example
	}

	// Create a buffer to hold the serialized JSON
	var buf bytes.Buffer

	// Serialize the payload into the buffer
	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		t.Fatalf("failed to serialize payload: %s", err.Error())
	}

	resp, err := client.Post(srv.URL+"/Users", "application/json", &buf)
	if err != nil {
		t.Fatalf("failed to post payload: %s", err.Error())
	}

	if resp.StatusCode != http.StatusCreated {
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("failed to read response body: %s", err.Error())
			return
		}

		t.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bytes))
	}
}

type testlogger struct{}

// Error implements scim.Logger.
func (t *testlogger) Error(args ...interface{}) {
	fmt.Printf("%s", args)
}
