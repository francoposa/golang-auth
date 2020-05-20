package resources

import (
	"errors"
	"fmt"
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
	return "Invalid ExampleCom Auth User Role"
}

type ResourceOperation string

const (
	Create ResourceOperation = "create"
	Read   ResourceOperation = "read"
	Update ResourceOperation = "update"
	Delete ResourceOperation = "delete"
)

func (op ResourceOperation) IsValid() bool {
	switch op {
	case Create, Read, Update, Delete:
		return true
	}
	return false
}

type InvalidResourceOperationError struct{}

func (e *InvalidResourceOperationError) Error() string {
	return "Invalid ExampleCom Resource Operation"
}

// Resource is an ExampleCom resource which may be operated upon by a Client on behalf of the Resource Owner
type Resource string

// Valid Resources in the ExampleCom system today
const (
	UserResource        Resource = "user"
	UserAccountResource Resource = "user.account"
	UserProfileResource Resource = "user.profile"
)

func (res Resource) IsValid() bool {
	switch res {
	case UserResource, UserAccountResource, UserProfileResource:
		return true
	}
	return false
}

type InvalidResourceError struct {
	msg string
}

const invalidResourceErrorMsgPrefix string = "Invalid ExampleCom Resource"

func (e *InvalidResourceError) Error() string {
	if e.msg != "" {
		return fmt.Sprintf("%s: %s", invalidResourceErrorMsgPrefix, e.msg)
	}
	return invalidResourceErrorMsgPrefix
}

type ResourceOperationScope struct {
	Role      AuthUserRole
	Operation ResourceOperation
	Resource  Resource
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

	var operation = ResourceOperation(parts[0])
	if !operation.IsValid() {
		return nil, &InvalidResourceOperationError{}
	}

	var resource = Resource(parts[2])
	if !resource.IsValid() {
		return nil, &InvalidResourceError{}
	}

	return &ResourceOperationScope{
		Role:      role,
		Operation: operation,
		Resource:  resource,
	}, nil
}
