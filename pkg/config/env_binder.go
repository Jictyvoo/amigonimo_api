package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Binder func() error

func BindFieldErr[T any](ptr *T, field string, parserCallback func(string) (T, error)) Binder {
	return func() (err error) {
		if resultStr := os.Getenv(field); resultStr != "" {
			if parserCallback != nil {
				*ptr, err = parserCallback(resultStr)
				return err
			}

			return tryTypeParse(ptr, resultStr)
		}

		return nil
	}
}

func BindField[T any](ptr *T, field string, parserCallback func(string) T) Binder {
	subCallback := func(v string) (T, error) {
		return parserCallback(v), nil
	}

	if parserCallback == nil {
		subCallback = nil
	}
	return BindFieldErr[T](ptr, field, subCallback)
}

func tryTypeParse(ptr any, resultStr string) (err error) {
	switch targetPtr := ptr.(type) {
	case *bool:
		*targetPtr, err = strconv.ParseBool(resultStr)
	case *float64:
		*targetPtr, err = strconv.ParseFloat(resultStr, 64)
	case *int:
		*targetPtr, err = strconv.Atoi(resultStr)
	case *int64:
		*targetPtr, err = strconv.ParseInt(resultStr, 10, 64)
	case *string:
		*targetPtr = resultStr
	case *uint:
		var asUint64 uint64
		asUint64, err = strconv.ParseUint(resultStr, 10, 32)
		*targetPtr = uint(asUint64)
	case *uint16:
		var asUint64 uint64
		asUint64, err = strconv.ParseUint(resultStr, 10, 16)
		*targetPtr = uint16(asUint64)
	case *uint64:
		*targetPtr, err = strconv.ParseUint(resultStr, 10, 64)
	case *time.Duration:
		*targetPtr, err = time.ParseDuration(resultStr)
	default:
		err = fmt.Errorf("invalid target type: %T", targetPtr)
	}
	return err
}

// BindEnv executes any number of Binder functions and returns the first error.
func BindEnv(binders ...Binder) error {
	for _, b := range binders {
		if b == nil {
			continue
		}
		if err := b(); err != nil {
			return err
		}
	}
	return nil
}
