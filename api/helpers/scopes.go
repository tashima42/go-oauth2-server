package helpers

type Scope string

const (
	AdminScope        Scope = "admin"
	ClientCreateScope Scope = "client:create"
	ClientListScope   Scope = "client:list"
)

func ScopesFromInterface(scopes []interface{}) []Scope {
	var s []Scope
	for _, scope := range scopes {
		s = append(s, Scope(scope.(string)))
	}
	return s
}
