package types_result

import (
	types_item "app_burse_backend/internal/types/error_item"
)

type Result interface {
	Success() bool

	Data() interface{}

	ErrorItem() *types_item.ErrorItem

	Error() string
}

type ErrorResult struct {
	Exception *types_item.ErrorItem
}

func NewErrorResult(options ...OptionResult) *ErrorResult {
	e := &ErrorResult{
		Exception: types_item.NewErrorItem(types_item.WithMsg("An error occurred.")),
	}

	for _, option := range options {
		option(e)
	}

	return e
}

func (r *ErrorResult) Success() bool {
	return false
}

func (r *ErrorResult) Data() interface{} {
	return r.Exception.Data()
}

func (r *ErrorResult) Error() string {
	return r.Exception.Error()
}

func (r *ErrorResult) ErrorItem() *types_item.ErrorItem {
	return r.Exception
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

func (r *SuccessResult) ErrorItem() *types_item.ErrorItem {
	return nil
}

func (r *SuccessResult) Error() string {
	return ""
}
