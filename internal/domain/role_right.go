package domain

type RoleRight struct {
	ID      string
	RoleID  string
	Section string
	Route   string
	RCreate int
	RRead   int
	RUpdate int
	RDelete int
}
