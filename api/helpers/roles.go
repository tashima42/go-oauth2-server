package helpers

type Scope string

const (
	AdminScope        Scope = "admin"
	ClientCreateScope Scope = "client:create"
	ClientListScope   Scope = "client:list"
)
