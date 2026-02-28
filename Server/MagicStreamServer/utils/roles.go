package utils

const (
	RoleUser  = "USER"
	RoleAdmin = "ADMIN"
)

var RoleRank = map[string]int{
	RoleUser:  1,
	RoleAdmin: 2,
}