package fixtures

// Builder is a generic interface for all fixtures
type Builder[T any] interface {
	Build() *T
}
