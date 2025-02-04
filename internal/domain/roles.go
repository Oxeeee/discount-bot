package domain

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleStaff UserRole = "staff"
	UserRoleAdmin UserRole = "admin"
)
