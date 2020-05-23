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

// AuthZUserRole is the role a user is assigned within ExampleCom's system
// AuthZUserRole disambiguates the applicable scope of ResourceOperationScope
// where an AuthZUserRole of "admin" gives permission for all resources in the system
// while an AuthZUserRole of "user" limits permission for only the resources owned by the given user
type AuthZUserRole string

// Valid AuthNUserRoles in the ExampleCom system today
const (
	AdminRole AuthZUserRole = "admin"
	UserRole  AuthZUserRole = "user"
)

func (role AuthZUserRole) IsValid() bool {
	switch role {
	case AdminRole, UserRole:
		return true
	}
	return false
}

type InvalidAuthZUserRoleError struct{}

func (e *InvalidAuthZUserRoleError) Error() string {
	return "Invalid ExampleCom Authorization User Role"
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
	ID            uuid.UUID
	AuthNUserRole AuthZUserRole
	Operation     ResourceOperation
	Resource      *Resource
}

func NewResourceOperationScope(
	role AuthZUserRole,
	operation ResourceOperation,
	resource *Resource,
) *ResourceOperationScope {
	id := uuid.New()
	return &ResourceOperationScope{
		ID:           id,
		AuthNUserRole: role,
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

	var role = AuthZUserRole(parts[0])
	if !role.IsValid() {
		return nil, &InvalidAuthZUserRoleError{}
	}

	var verb = ResourceOperation(parts[1])
	if !verb.IsValid() {
		return nil, &InvalidResourceOperationError{}
	}

	var resource = NewResource(parts[2])

	return NewResourceOperationScope(role, verb, resource), nil
}
