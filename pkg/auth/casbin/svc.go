package casbin

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/auth0.v5"
	"gopkg.in/auth0.v5/management"
)

type service struct {
}

//unc (s *service) AllowAccess(jwt string, account string, resource string) error {
//	e, err := casbin.NewEnforcer("path/to/model.conf", "path/to/policy.csv")
//
//	// Use our role manager.
//	// clientID is the Client ID.
//	// clientSecret is the Client Secret.
//	// tenant is your tenant name. If your domain is: abc.auth0.com, then abc is your tenant name.
//	// apiEndpoint is the base URL for your Auth0 Authorization Extension, it should
//	// be something like: "https://abc.us.webtask.io/adf6e2f2b84784b57522e3b19dfc9201", there is
//	// no "/admins", "/admins/login", "/users" or "/api" in the end.
//	rm := NewAuth0RoleManager(
//		"your_client_id",
//		"your_client_secret",
//		"your_tenant_name",
//		"your_base_url_for_auth0_authorization_extension")
//
//	e.SetRoleManager(rm)
//
//	// If our role manager relies on Casbin policy (like reading "g"
//	// policy rules), then we have to set the role manager before loading
//	// policy.
//	//
//	// Otherwise, we can set the role manager at any time, because role
//	// manager has nothing to do with the adapter.
//	e.LoadPolicy()
//
//	e.Enforce(usr, acc, rol, act)
//

func (s *service) UserInfo(accessToken string) {

	m, err := management.New("domain", management.WithClientCredentials("id", "secret"))
	if err != nil {
		// handle err
	}
	c := &management.Client{
		Name:        auth0.String("Client Name"),
		Description: auth0.String("Long description of client"),
	}

	err = m.Client.Create(c)
	if err != nil {
		// handle err
	}

	url := "http://localhost:3010/api/private"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("authorization", "Bearer YOUR_ACCESS_TOKEN")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}
