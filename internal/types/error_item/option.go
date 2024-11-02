package types_item

type OptionError func(e *ErrorItem)

func WithDetails(details map[string]interface{}) OptionError {
	return func(e *ErrorItem) {
		e.Details = details
	}
}

func WithErrors(errors [][]string) OptionError {
	return func(e *ErrorItem) {
		for _, err := range errors {
			e.Errors = append(e.Errors, ErrorField{Field: err[0], Message: err[1]})
		}
	}
}

func WithStatus(status int) OptionError {
	return func(e *ErrorItem) {
		e.Status = status
	}
}

func WithMsg(message string) OptionError {
	return func(e *ErrorItem) {
		e.Msg = message
	}
}
