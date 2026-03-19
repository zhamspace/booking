package constant

const (
	ServiceName = "booking"

	MaxPageSize = 1000
)

// Roles

const (
	RoleAdmin = "core:admin"
)

func RoleIsValid(v string) bool {
	return v == RoleAdmin
}

var RoleWeight = map[string]int{
	RoleAdmin: 3,
}
