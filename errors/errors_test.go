package errors_test

import (
	"errorHandler/errors"
	errorsBase "errors"
	"testing"
)

type errorType uint32

const (
	userError errorType = iota
	systemError
)

func TestErrorInGenerator(t *testing.T) {
	factory := errors.NewFactory(
		func(typeOfErrorArg errorType, descriptorArg string, failback string, argumentsArg []interface{}) (string, error) {
			return "", errorsBase.New("abc")
		}, systemError)
	err1 := factory.New(userError, "message", "message for failback")
	if err1.Type() != systemError {
		t.Errorf(`New(userError,"message") don't failed on generator error`)
	}
	if err1.Error() != "abc" {
		t.Errorf(`message from generator was abc, but is %s`, err1.Error())
	}
}
func TestErrorMethods(t *testing.T) {
	var typeOfError errorType
	var descriptor string
	var failback string
	var arguments []interface{}
	factory := errors.NewFactory(
		func(typeOfErrorArg errorType, descriptorArg string, failbackArg string, argumentsArg []interface{}) (string, error) {
			typeOfError = typeOfErrorArg
			descriptor = descriptorArg
			arguments = argumentsArg
			failback = failbackArg
			return "test", nil
		}, systemError)
	err1 := factory.New(userError, "message", "message for failback")
	if typeOfError != userError {
		t.Errorf(`New(userError,"message","mesage for failback") don't send userError`)
	}
	if descriptor != "message" {
		t.Errorf(`New(userError,"message","mesage for failback") don't send "message"`)
	}
	if failback != "message for failback" {
		t.Errorf(`New(userError,"message","mesage for failback") don't send "message for failback"`)
	}
	if len(arguments) != 0 {
		t.Errorf(`New(userError,"message","mesage for failback") don't send arguments`)
	}
	if err1.Error() != "test" {
		t.Errorf(`message from generator was test, but is %s`, err1.Error())
	}
	if err1.Type() != userError {
		t.Errorf(`New(userError,"message").Error()!=userError`)
	}
	err2 := factory.New(systemError, "message2", "", 5)
	if typeOfError != systemError {
		t.Errorf(`New(systemError,"message2","", 5) don't send systemError`)
	}
	if failback != "" {
		t.Errorf(`New(systemError,"message2","", 5) don't send ""`)
	}
	if descriptor != "message2" {
		t.Errorf(`New(systemError,"message2","", 5) don't send "message2"`)
	}
	if len(arguments) != 1 || arguments[0] != 5 {
		t.Errorf(`New(systemError,"message2",5) don't send 5`)
	}
	if err2.Type() != systemError {
		t.Errorf(`New(userError,"message").Error()!=userError`)
	}

}
