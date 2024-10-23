package server

import (
	"log"

	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/optional"
	"github.com/elimity-com/scim/schema"
	"github.com/intility/scim/handlers"
)

func NewServer(logger scim.Logger) *scim.Server {
	spCfg := scim.ServiceProviderConfig{
		DocumentationURI: optional.NewString("https://docs.intility.com/identity/scim"),
	}

	sc := schema.Schema{
		ID:          "urn:ietf:params:scim:schemas:core:2.0:User",
		Name:        optional.NewString("User"),
		Description: optional.NewString("User Account"),
		Attributes: []schema.CoreAttribute{
			schema.ComplexCoreAttribute(schema.ComplexParams{
				Name:     "name",
				Required: true,
				SubAttributes: []schema.SimpleParams{
					schema.SimpleStringParams(schema.StringParams{
						Name:     "formatted",
						Required: false,
					}),
					schema.SimpleStringParams(schema.StringParams{
						Name:     "familyName",
						Required: false,
					}),
					schema.SimpleStringParams(schema.StringParams{
						Name:     "givenName",
						Required: false,
					}),
					schema.SimpleStringParams(schema.StringParams{
						Name:     "middleName",
						Required: false,
					}),
					schema.SimpleStringParams(schema.StringParams{
						Name:     "honorificPrefix",
						Required: false,
					}),
					schema.SimpleStringParams(schema.StringParams{
						Name:     "honorificSuffix",
						Required: false,
					}),
				},
			}),
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:       "displayName",
				Required:   false,
				Uniqueness: schema.AttributeUniquenessNone(),
			})),
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:       "nickName",
				Required:   false,
				Uniqueness: schema.AttributeUniquenessNone(),
			})),
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:       "userName",
				Required:   true,
				Uniqueness: schema.AttributeUniquenessServer(),
			})),
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:       "profileUrl",
				Required:   false,
				Uniqueness: schema.AttributeUniquenessNone(),
			})),
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:       "title",
				Required:   false,
				Uniqueness: schema.AttributeUniquenessNone(),
			})),
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:       "userType",
				Required:   false,
				Uniqueness: schema.AttributeUniquenessNone(),
			})),
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:       "preferredLanguage",
				Required:   false,
				Uniqueness: schema.AttributeUniquenessNone(),
			})),
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:       "locale",
				Required:   false,
				Uniqueness: schema.AttributeUniquenessNone(),
			})),
			schema.SimpleCoreAttribute(schema.SimpleDateTimeParams(schema.DateTimeParams{
				Name:     "timezone",
				Required: false,
			})),
			schema.SimpleCoreAttribute(schema.SimpleBooleanParams(schema.BooleanParams{
				Name:        "active",
				Description: optional.NewString("A Boolean value indicating the user's administrative status."),
				Required:    false,
			})),
			schema.ComplexCoreAttribute(schema.ComplexParams{
				Name:        "emails",
				Required:    false,
				MultiValued: true,
				SubAttributes: []schema.SimpleParams{
					schema.SimpleBooleanParams(schema.BooleanParams{Name: "primary"}),
					schema.SimpleStringParams(schema.StringParams{Name: "value", Required: true}),
					schema.SimpleStringParams(schema.StringParams{Name: "display", Required: false}),
					schema.SimpleStringParams(schema.StringParams{Name: "type", Required: true}),
				},
			}),
			schema.ComplexCoreAttribute(schema.ComplexParams{
				Name:        "phoneNumbers",
				Required:    false,
				MultiValued: true,
				SubAttributes: []schema.SimpleParams{
					schema.SimpleBooleanParams(schema.BooleanParams{Name: "primary"}),
					schema.SimpleStringParams(schema.StringParams{Name: "value", Required: true}),
					schema.SimpleStringParams(schema.StringParams{Name: "display", Required: false}),
					schema.SimpleStringParams(schema.StringParams{Name: "type", Required: true}),
				},
			}),
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:        "photos",
				Required:    false,
				MultiValued: true,
			})),
		},
	}

	extension := schema.Schema{
		ID:          "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User",
		Name:        optional.NewString("EnterpriseUser"),
		Description: optional.NewString("Enterprise User"),
		Attributes:  []schema.CoreAttribute{
			// schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
			// 	Name: "employeeNumber",
			// })),
			// schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
			// 	Name: "organization",
			// })),
			// schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
			// 	Name:        "roles",
			// 	MultiValued: true,
			// })),
		},
	}

	userResourceHandler := handlers.NewUserResourceHandler()

	resourceTypes := []scim.ResourceType{
		{
			ID:          optional.NewString("User"),
			Name:        "User",
			Endpoint:    "/Users",
			Description: optional.NewString("User Account"),
			Schema:      sc,
			SchemaExtensions: []scim.SchemaExtension{
				{Schema: extension},
			},
			Handler: userResourceHandler,
		},
	}

	serverArgs := &scim.ServerArgs{
		ServiceProviderConfig: &spCfg,
		ResourceTypes:         resourceTypes,
	}

	serverOpts := scim.ServerOption(
		scim.WithLogger(logger),
	)

	server, err := scim.NewServer(serverArgs, serverOpts)
	if err != nil {
		log.Fatal("failed to bootstrap server: " + err.Error())
	}

	return &server
}
