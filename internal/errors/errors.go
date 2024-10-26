package errors

type ErrorField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorItem struct {
	Message string `json:"message"`
	Code    int    `json:"code"`

	Details map[string]interface{} `json:"details,omitempty"`

	Errors []ErrorField `json:"errors,omitempty"`
}

type OptionError func(e *ErrorItem)

func WithDetails(details map[string]interface{}) OptionError {
	return func(e *ErrorItem) {
		e.Details = details
	}
}

func WithErrors(errors []ErrorField) OptionError {
	return func(e *ErrorItem) {
		e.Errors = errors
	}
}

func WithCode(code int) OptionError {
	return func(e *ErrorItem) {
		e.Code = code
	}
}

func WithMessage(message string) OptionError {
	return func(e *ErrorItem) {
		e.Message = message
	}
}

func NewErrorItem(options ...OptionError) *ErrorItem {
	e := &ErrorItem{
		Code:    0,
		Details: make(map[string]interface{}),
		Errors:  make([]ErrorField, 0),
		Message: "",
	}

	for _, option := range options {
		option(e)
	}

	return e
}

func (e *ErrorItem) Error() string {
	return e.Message
}
