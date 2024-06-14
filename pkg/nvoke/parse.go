package nvoke

type Parser[T any] interface {
	Parse(location string) []*T
}
