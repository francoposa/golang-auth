package resources

import (
	"net/url"

	"github.com/google/uuid"
)

type ClientRequiresRedirectURIError struct{}

func (e *ClientRequiresRedirectURIError) Error() string {
	return "AuthZClient Registration Requires Redirect URI: RFC 6749 Section 3.1.2.2"
}

// AuthZClient represents an authorization-provider client as defined in RFC 6749 Section 1.1 - Roles
type AuthZClient struct {
	ID          uuid.UUID
	Secret      *uuid.UUID
	RedirectURI *url.URL
	Public      bool
	FirstParty  bool
}

// Create a AuthZClient consistent with RFC 6749 Section 2 - AuthZClient Registration
func NewClient(redirectURI string, public bool, firstParty bool) (*AuthZClient, error) {

	// Require clients to register a Redirect URI as specified in
	// RFC 6749 Section 3.1.2.2 - Redirection Endpoint Registration Requirements
	//
	// The authorization-provider server MUST require the following clients to
	// register their redirection endpoint:
	//    - Public clients.
	//    - Confidential clients utilizing the implicit grant type.
	//
	// The authorization-provider server SHOULD require all clients to register their
	// redirection endpoint prior to utilizing the authorization-provider endpoint.
	var uri *url.URL
	var err error
	if redirectURI == "" {
		// TODO validate URI as defined in:
		//    - RFC 6749 Section 3.1.2 - Redirection Endpoint
		//    - OpenID Connect Core Section 3.1.2.1. - Authentication Request
		return nil, &ClientRequiresRedirectURIError{}
	} else {
		uri, err = url.Parse(redirectURI)
		if err != nil {
			return nil, err
		}
	}

	// We will not issue public clients (as defined in RFC 5749 Section 2.1 - AuthZClient Types)
	// a secret. Denying public clients a secret is not specifically required by
	// the RFC, but not issuing a secret is a good way to avoid leaking one
	var secret *uuid.UUID
	if public {
		secret = nil
	} else {
		newSecret := uuid.New()
		secret = &newSecret
	}

	id := uuid.New()
	return &AuthZClient{
		ID:          id,
		Secret:      secret,
		RedirectURI: uri,
		Public:      public,
		FirstParty:  firstParty,
	}, nil
}
