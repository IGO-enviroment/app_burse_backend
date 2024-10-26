package service

type Result interface {
	Success() bool

	Data() interface{}

	Errors() error

	Error() string
}

type ErrorResult struct {
	DataErrors error
}

func NewErrorResult(errors error) *ErrorResult {
	return &ErrorResult{
		DataErrors: errors,
	}
}

func (r *ErrorResult) Success() bool {
	return false
}

func (r *ErrorResult) Data() interface{} {
	return nil
}

func (r *ErrorResult) Error() string {
	return r.DataErrors.Error()
}

func (r *ErrorResult) Errors() error {
	return r.DataErrors
}

type SuccessResult struct {
	DataField interface{}
}

func NewSuccessResult(data interface{}) *SuccessResult {
	return &SuccessResult{
		DataField: data,
	}
}

func (r *SuccessResult) Success() bool {
	return true
}

func (r *SuccessResult) Data() interface{} {
	return r.DataField
}

func (r *SuccessResult) Errors() error {
	return nil
}

func (r *SuccessResult) Error() string {
	return ""
}
