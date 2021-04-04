package t

type (
	String = *string
	Int64  = *int64
)

func NewString(s string) String {
	return &s
}

func NewInt64(s int64) Int64 {
	return &s
}
