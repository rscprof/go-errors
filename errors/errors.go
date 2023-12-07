package errors

// Package gives an opportunity to work with different types of errors such as
// system errors, user errors etc. It also has an opportunity to produce error
// messages by user code (for example for localization)

// errorWithTypeAndArguments realise error for storing
// type of error (for example user error, system error etc)
// and message (that produced by MessageGenerator)
type errorWithTypeAndArguments[T any] struct {
	typeOfError T
	message     string
}

// Error is an interface that realise methods to get
// type of error by method Type() for example user error or system error
// and stander error interface function Error() to get message about error
type Error[T any] interface {
	Type() T
	Error() string
}

// MessageGenerator converts type of error, descriptor, failback and arguments to error string
// If mesage generator fails it returns standard error
type MessageGenerator[T any] func(typeOfError T, descriptor string, failback string, arguments []interface{}) (string, error)

// ErrorFactory has a method New to produce struct that has an interface Error
// Method New gets type of error, descriptor (for example localization key), failback (message if complex method failed) and arguments (for example for method Sprintf)
type ErrorFactory[T any] interface {
	New(typeOfError T, descriptor string, failback string, arguments ...interface{}) errorWithTypeAndArguments[T]
}

// NewFactory gets a factory for produce error structs with types and arguments
// generator is a producer of message string from descriptor and arguments
// errors in generator will produce defaultErrorType error
func NewFactory[T any](generator MessageGenerator[T], defaultErrorType T) ErrorFactory[T] {
	return errorFactory[T]{generator, defaultErrorType}
}

// errorFactory stores messageGenerator, that convert type of error, descriptor and arguiments to message
// and defaultErrorType for errors from messageGenerator
type errorFactory[T any] struct {
	messageGenerator MessageGenerator[T]
	defaultErrorType T
}

// New creates object that has an interface Error and store
// type of error (for example system or user error)
// descriptor (for example message or localization key)
// failback (for example message that returns if localization is not exists)
// arguments (for example arguments for Sprintf)
//
// # If the generator fails, then New returns error from generator with default error type
//
// We can't save arguments so we have to using it during calling
func (factory errorFactory[T]) New(typeOfError T, descriptor string, failback string, arguments ...interface{}) errorWithTypeAndArguments[T] {
	message, err := factory.messageGenerator(typeOfError, descriptor, failback, arguments)
	if err == nil {
		return errorWithTypeAndArguments[T]{
			typeOfError,
			message,
		}
	} else {
		return errorWithTypeAndArguments[T]{
			factory.defaultErrorType,
			err.Error(),
		}
	}
}

// Type get type of error
func (e errorWithTypeAndArguments[T]) Type() T {
	return e.typeOfError
}

// Error get message that describe an error
func (e errorWithTypeAndArguments[T]) Error() string {
	return e.message
}
