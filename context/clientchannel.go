package context

type ClientChannel interface {
	Emit(data any) error
}
