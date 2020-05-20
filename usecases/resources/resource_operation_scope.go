package resources

import (
	"github.com/google/uuid"
)

//// OpenID Connect Scope as defined in OIDC Core Section 5.4 - Requesting Claims using Scope Values
//// An OpenID Connect Scope is a OpenID Connect-specific value for the generalized OAuth2
//// definition of a Scope provided in RFC 6749 Section 3.3. - Access Token Scope
//type OIDCScope string
//
//// Valid OpenID Connect Scopes as defined in OIDC Core Section 5.4
//const (
//	OpenID  OIDCScope = "openid"
//	Profile OIDCScope = "profile"
//	Email   OIDCScope = "email"
//	Address OIDCScope = "address"
//	Phone   OIDCScope = "phone"
//)
//
//func (scope OIDCScope) IsValid() bool {
//	switch scope {
//	case OpenID, Profile, Email, Address, Phone:
//		return true
//	}
//	return false
//}
//
//type InvalidOIDCScopeError struct{}
//
//func (e *InvalidOIDCScopeError) Error() string {
//	return "Invalid OpenID Connect Scope: OIDC Core Section 5.4"
//}

// AuthUserRole is the role a user is assigned within ExampleCom's system
// AuthUserRole disambiguates the applicable scope of ResourceVerbScope
// where an AuthUserRole of "admin" gives permission for all resources in the system
// while an AuthUserRole of "user" limits permission for only the resources owned by the given user
type AuthUserRole string

// Valid AuthUserRoles in the ExampleCom system today
const (
	AdminRole AuthUserRole = "admin"
	UserRole  AuthUserRole = "user"
)

func (role AuthUserRole) IsValid() bool {
	switch role {
	case AdminRole, UserRole:
		return true
	}
	return false
}

type InvalidAuthUserRoleError struct{}

func (e *InvalidAuthUserRoleError) Error() string {
	return "Invalid ExampleCom Auth User Role"
}

type ResourceVerb string

const (
	Create ResourceVerb = "create"
	Read   ResourceVerb = "read"
	Update ResourceVerb = "update"
	Delete ResourceVerb = "delete"
)

func (v ResourceVerb) IsValid() bool {
	switch v {
	case Create, Read, Update, Delete:
		return true
	}
	return false
}

type InvalidResourceOperationError struct{}

func (e *InvalidResourceOperationError) Error() string {
	return "Invalid ExampleCom Resource Operation"
}

type ResourceVerbScope struct {
	ID       uuid.UUID
	Role     AuthUserRole
	Verb     ResourceVerb
	Resource Resource
}

func NewResourceVerbScope(role AuthUserRole, verb ResourceVerb, resource Resource) *ResourceVerbScope {
	id := uuid.New()
	return &ResourceVerbScope{
		ID:       id,
		Role:     role,
		Verb:     verb,
		Resource: resource,
	}
}

//
//func ParseResourceVerbScope(s string) (*ResourceVerbScope, error) {
//	parts := strings.Split(s, ":")
//	if len(parts) != 3 {
//		return nil, errors.New(
//			"Invalid Resource Verb Scope string. Must be in format `role:verb:resource`",
//		)
//	}
//
//	var role = AuthUserRole(parts[0])
//	if !role.IsValid() {
//		return nil, &InvalidAuthUserRoleError{}
//	}
//
//	var verb = ResourceVerb(parts[1])
//	if !verb.IsValid() {
//		return nil, &InvalidResourceOperationError{}
//	}
//
//	var resource = Resource(parts[2])
//	if !resource.IsValid() {
//		return nil, &InvalidResourceError{}
//	}
//
//	return NewResourceVerbScope(role, verb, resource), nil
//}
