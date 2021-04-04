package auth

type Service interface {
	AllowAccess(jwt string, resource string) error
}