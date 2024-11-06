package types_item

type ErrorField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorItem struct {
	Msg    string `json:"message"`
	Status int    `json:"status"`

	Details map[string]interface{} `json:"details,omitempty"`

	Errors []ErrorField `json:"errors,omitempty"`
}

func NewErrorItem(options ...OptionError) *ErrorItem {
	e := &ErrorItem{
		Status:  0,
		Details: make(map[string]interface{}),
		Errors:  make([]ErrorField, 0),
		Msg:     "",
	}

	for _, option := range options {
		option(e)
	}

	return e
}

func (e *ErrorItem) Error() string {
	return e.Msg
}

func (e *ErrorItem) FieldErrors() []ErrorField {
	return e.Errors
}

func (e *ErrorItem) Data() map[string]interface{} {
	return e.Details
}

func (e *ErrorItem) Message() string {
	return e.Msg
}

func (e *ErrorItem) Code() int {
	return e.Status
}
