package main

import (
	"encoding/json"
	"errorHandler/errors"
	errorsBase "errors"
	"fmt"
	"os"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/message"
)

type errorType uint

const (
	userError errorType = iota
	systemError
)

func readIntInRange(factory errors.ErrorFactory[errorType],
	min int,
	max int) (int, errors.Error[errorType]) {
	var num int
	_, err := fmt.Scan(&num)
	if err != nil {
		return 0, factory.New(systemError, "cannot_read", "Cannot read a message")
	}
	if (num > max) || (num < min) {
		return 0, factory.New(userError, "out_of_range", "Number is out of range", num, min, max)
	}
	return num, nil
}

func getMessageGenerator(localizerUser *i18n.Localizer, localizerSystem *i18n.Localizer) errors.MessageGenerator[errorType] {

	return func(errorType errorType, descriptor string, failback string, arguments []interface{}) (string, error) {
		switch errorType {
		case systemError:
			message, err := localizerSystem.Localize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    descriptor,
					Other: failback,
				},
			})
			if err == nil {
				return fmt.Sprintf(message, arguments...), nil
			} else {
				return "", fmt.Errorf("Unexpected error from localize: %w", err)
			}
		case userError:
			message, err := localizerUser.Localize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    descriptor,
					Other: failback,
				},
			})
			if err == nil {
				return fmt.Sprintf(message, arguments...), nil
			} else {
				return "", fmt.Errorf("Unexpected error from localize: %w", err)
			}
		}
		return "", errorsBase.New("Unexpected error type")
	}
}

func main() {

	userLocale := "ru"
	systemLocale := "en"
	bundle := i18n.NewBundle(message.MatchLanguage(systemLocale))
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile("locales/en.json")
	bundle.MustLoadMessageFile("locales/ru.json")
	localizerUser := i18n.NewLocalizer(bundle, userLocale)
	localizerSystem := i18n.NewLocalizer(bundle, systemLocale)
	messageGenerator := getMessageGenerator(localizerUser, localizerSystem)
	factory := errors.NewFactory(messageGenerator, systemError)
	num, err := readIntInRange(factory, 1, 5)
	if err != nil {
		if err.Type() == userError {
			fmt.Printf(err.Error())
		} else if err.Type() == systemError {
			fmt.Fprintf(os.Stderr, err.Error())
		} else {
			fmt.Fprintf(os.Stderr, "Unexpected error type")
		}
	} else {
		fmt.Printf("You typed %d\n", num)
	}
}

// This function is only to generate locales
func forGenerator() {
	printer := message.NewPrinter(message.MatchLanguage("ru"))
	printer.Printf("out_of_range")
	printer.Printf("cannot_read")
	printer2 := message.NewPrinter(message.MatchLanguage("en"))
	printer2.Printf("out_of_range")
	printer2.Printf("cannot_read")
}
