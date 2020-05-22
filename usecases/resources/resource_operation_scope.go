package resources

import (
	"errors"
	"github.com/google/uuid"
	"strings"
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
// AuthUserRole disambiguates the applicable scope of ResourceOperationScope
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
	return "Invalid ExampleCom Auth User AuthUserRole"
}

type ResourceOperation string

const (
	Create ResourceOperation = "create"
	Read   ResourceOperation = "read"
	Update ResourceOperation = "update"
	Delete ResourceOperation = "delete"
)

func (v ResourceOperation) IsValid() bool {
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

type ResourceOperationScope struct {
	ID           uuid.UUID
	AuthUserRole AuthUserRole
	Operation    ResourceOperation
	Resource     *Resource
}

func NewResourceOperationScope(
	role AuthUserRole,
	operation ResourceOperation,
	resource *Resource,
) *ResourceOperationScope {
	id := uuid.New()
	return &ResourceOperationScope{
		ID:           id,
		AuthUserRole: role,
		Operation:    operation,
		Resource:     resource,
	}
}

func ParseResourceOperationScope(s string) (*ResourceOperationScope, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return nil, errors.New(
			"Invalid Resource Operation Scope string. Must be in format `role:operation:resource`",
		)
	}

	var role = AuthUserRole(parts[0])
	if !role.IsValid() {
		return nil, &InvalidAuthUserRoleError{}
	}

	var verb = ResourceOperation(parts[1])
	if !verb.IsValid() {
		return nil, &InvalidResourceOperationError{}
	}

	var resource = NewResource(parts[2])

	return NewResourceOperationScope(role, verb, resource), nil
}
