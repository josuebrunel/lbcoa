package fizzbuzz

import (
	"context"
	"errors"
	"fizzbuzz/pkg/apiresponse"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	ErrMissingParam    = errors.New("missing query parameters")
	ErrValMustBeInt    = errors.New("value must be an integer")
	ErrValMustBeGTZero = errors.New("value must be greater than 0")
	ErrValEmpty        = errors.New("value must not be empty")
)

type Args struct {
	Int1  int
	Int2  int
	Limit int
	Str1  string
	Str2  string
}

func (a Args) String() string {
	return fmt.Sprintf("int1=%d&int2=%d&limit=%d&str1=%s&str2=%s", a.Int1, a.Int2, a.Limit, a.Str1, a.Str2)
}

// parseIntNValidate parses the given string to an integer and validates it.
// It returns an error if the string is not a valid integer or if the integer is less than 1.
// It returns the integer and nil if the string is a valid integer and greater than or equal to 1.
func parseIntNValidate(s string) (int, error) {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, ErrValMustBeInt
	}
	if v < 1 {
		return 0, ErrValMustBeGTZero
	}
	return v, nil
}

// parseStrNValidate parses the given string and validates it.
// It returns an error if the string is empty.
// It returns the string and nil if the string is not empty.
func parseStrNValidate(s string) (string, error) {
	if strings.EqualFold(s, "") {
		return "", ErrValEmpty
	}
	return s, nil
}

// argsFromQuery parses the query parameters from the given url.Values and returns an Args struct.
// It returns an error if any of the query parameters are missing or invalid.
// It returns the Args struct and nil if all query parameters are valid.
func argsFromQuery(qs url.Values) (Args, error) {
	var (
		args       = Args{}
		err        error
		paramCount = 0
	)
	for k, v := range qs {
		switch k {
		case "int1":
			args.Int1, err = parseIntNValidate(v[0])
			paramCount++
		case "int2":
			args.Int2, err = parseIntNValidate(v[0])
			paramCount++
		case "limit":
			args.Limit, err = parseIntNValidate(v[0])
			paramCount++
		case "str1":
			args.Str1, err = parseStrNValidate(v[0])
			paramCount++
		case "str2":
			args.Str2, err = parseStrNValidate(v[0])
			paramCount++
		}
		if err != nil {
			return Args{}, err
		}
	}
	if paramCount < 5 {
		return args, ErrMissingParam
	}
	slog.Info("Args", "args", args, "error", err)
	return args, err
}

// Handler is the main handler for the fizzbuzz endpoint.
// It parses the query parameters, validates them, and then calls the fizzBuzz function to generate the result.
// It returns a JSON response with the result.  If there is an error, it returns a JSON response with the error.
func Handler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		args, err := argsFromQuery(r.URL.Query())
		if err != nil {
			apiresponse.New(w, http.StatusBadRequest, "", err)
			return
		}
		apiresponse.New(w, http.StatusOK, fizzBuzz(ctx, args), nil)
	}
}

// fizzBuzz generates a fizzbuzz sequence based on the given arguments.
// It returns a string containing the fizzbuzz sequence.
func fizzBuzz(ctx context.Context, args Args) string {
	res := []string{}
	for i := 1; i <= args.Limit; i++ {
		if i%args.Int1 == 0 && i%args.Int2 == 0 {
			res = append(res, args.Str1+args.Str2)
		} else if i%args.Int1 == 0 {
			res = append(res, args.Str1)
		} else if i%args.Int2 == 0 {
			res = append(res, args.Str2)
		} else {
			res = append(res, strconv.Itoa(i))
		}
	}
	slog.Info("result", "result", res)
	return strings.Join(res, ",")
}
