package helpers

const (
	AdminScope                   = "admin"
	ClientCreateScope            = "client:create"
	ClientListScope              = "client:list"
	ClientInfoReadScope          = "client:info:read"
	DevAccountCreateScope        = "dev_account:create"
	UserAccountUserInfoReadScope = "user_account:userinfo:read"
)

var DefaultUserAccountScopes []string = []string{
	UserAccountUserInfoReadScope,
	ClientInfoReadScope,
}
var DefaultDevAccountScopes []string = []string{
	ClientCreateScope,
	ClientListScope,
	UserAccountUserInfoReadScope,
	ClientInfoReadScope,
}

var DefaultClientScopes []string = []string{
	ClientInfoReadScope,
	UserAccountUserInfoReadScope,
}

func ScopesFromInterface(scopes []interface{}) []string {
	var s []string
	for _, scope := range scopes {
		s = append(s, string(scope.(string)))
	}
	return s
}
