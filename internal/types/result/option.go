package types_result

import types_item "app_burse_backend/internal/types/error_item"

type OptionResult func(Result)

func ErrorWithErrorItem(exception *types_item.ErrorItem) OptionResult {
	return func(r Result) {
		r.(*ErrorResult).Exception = exception
	}
}

func WithError(err error) OptionResult {
	return func(r Result) {
		if err != nil {
			r.(*ErrorResult).Exception = types_item.NewErrorItem(types_item.WithMsg(err.Error()))
		}
	}
}

func ErrorWithMsg(msg string) OptionResult {
	return func(r Result) {
		r.(*ErrorResult).Exception = types_item.NewErrorItem(types_item.WithMsg(msg))
	}
}
