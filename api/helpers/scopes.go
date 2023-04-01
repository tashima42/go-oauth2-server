package helpers

const (
	AdminScope               = "admin"
	ClientCreateScope        = "client:create"
	ClientListScope          = "client:list"
	DevAccountCreateScope    = "dev_account:create"
	UserAccountUserInfoScope = "user_account:userinfo:read"
)

var DefaultUserAccountScopes []string = []string{UserAccountUserInfoScope}
var DefaultDevAccountScopes []string = []string{ClientCreateScope, ClientListScope, UserAccountUserInfoScope}

var DefaultUserAcountScopesString = DefaultDevAccountScopes

func ScopesFromInterface(scopes []interface{}) []string {
	var s []string
	for _, scope := range scopes {
		s = append(s, string(scope.(string)))
	}
	return s
}
