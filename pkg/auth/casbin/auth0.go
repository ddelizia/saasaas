package casbin

import (
	auth0rolemanager "github.com/casbin/auth0-role-manager"
	"github.com/casbin/casbin/v2/log"
	"github.com/casbin/casbin/v2/rbac"
)

type Auth0RoleManager struct {
	auth0rm *auth0rolemanager.RoleManager
}

func NewAuth0RoleManager(clientID string, clientSecret string, tenant string, apiEndpoint string) rbac.RoleManager {

	rbacRoleManager := auth0rolemanager.NewRoleManager(
		clientID,
		clientSecret,
		tenant,
		apiEndpoint)

	rm := Auth0RoleManager{
		auth0rm: rbacRoleManager.(*auth0rolemanager.RoleManager),
	}

	return &rm
}

func (rm *Auth0RoleManager) SetLogger(logger log.Logger) {

}

func (rm *Auth0RoleManager) Clear() error {
	return rm.auth0rm.Clear()
}

// AddLink adds the inheritance link between two roles. role: name1 and role: name2.
// domain is a prefix to the roles (can be used for other purposes).
func (rm *Auth0RoleManager) AddLink(name1 string, name2 string, domain ...string) error {
	return rm.auth0rm.AddLink(name1, name2, domain...)
}

// DeleteLink deletes the inheritance link between two roles. role: name1 and role: name2.
// domain is a prefix to the roles (can be used for other purposes).
func (rm *Auth0RoleManager) DeleteLink(name1 string, name2 string, domain ...string) error {
	return rm.auth0rm.DeleteLink(name1, name2, domain...)
}

// HasLink determines whether a link exists between two roles. role: name1 inherits role: name2.
// domain is a prefix to the roles (can be used for other purposes).
func (rm *Auth0RoleManager) HasLink(name1 string, name2 string, domain ...string) (bool, error) {
	return rm.auth0rm.HasLink(name1, name2, domain...)
}

// GetRoles gets the roles that a user inherits.
// domain is a prefix to the roles (can be used for other purposes).
func (rm *Auth0RoleManager) GetRoles(name string, domain ...string) ([]string, error) {
	return rm.auth0rm.GetRoles(name, domain...)
}

// GetUsers gets the users that inherits a role.
// domain is a prefix to the users (can be used for other purposes).
func (rm *Auth0RoleManager) GetUsers(name string, domain ...string) ([]string, error) {
	return rm.auth0rm.GetUsers(name, domain...)
}

// PrintRoles prints all the roles to log.
func (rm *Auth0RoleManager) PrintRoles() error {
	return rm.auth0rm.PrintRoles()
}
