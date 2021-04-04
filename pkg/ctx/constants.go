package ctx

// ContextKey is the type that should be used
type ContextKey string

const (
	AccountIDContextField ContextKey = "AccountId"
	UserIDContextField    ContextKey = "UserId"
	ProjectIDContextField ContextKey = "ProjectId"
)
