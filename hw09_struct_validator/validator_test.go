package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	InvalidIntRule struct {
		Code int `validate:"len:1"`
	}

	InvalidStrRule struct {
		Body string `validate:"min:1"`
	}

	InvalidData struct {
		Err error `validate:"len:2"`
	}

	Regexp struct {
		Body string `validate:"regexp:^\\d+$|len:2"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "test",
				Name:   "test",
				Age:    1,
				Email:  "my@email.com",
				Role:   "stuff",
				Phones: []string{"79999999999"},
			},
			expectedErr: ValidationErrors{
				{
					Field: "ID",
					Err:   ErrValidateLen,
				},
				{
					Field: "Age",
					Err:   ErrValidateMin,
				},
			},
		},
		{
			in:          App{"test1"},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "test",
				Name:   "test",
				Age:    100,
				Email:  "my@email.com",
				Role:   "test",
				Phones: []string{"7999999999"},
			},
			expectedErr: ValidationErrors{
				{
					Field: "ID",
					Err:   ErrValidateLen,
				},
				{
					Field: "Age",
					Err:   ErrValidateMax,
				},
				{
					Field: "Role",
					Err:   ErrValidateStringIn,
				},
				{
					Field: "Phones",
					Err:   ErrValidateLen,
				},
			},
		},
		{
			in: Token{
				Header: []byte("test"),
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 100,
				Body: "test",
			},
			expectedErr: ValidationErrors{
				{
					Field: "Code",
					Err:   ErrValidateIntInt,
				},
			},
		},
		{
			in: Response{
				Code: 200,
				Body: "test",
			},
			expectedErr: nil,
		},
		{
			in: InvalidIntRule{
				Code: 100,
			},
			expectedErr: SysErr{ErrUnsupportedRuleType},
		},
		{
			in: InvalidStrRule{
				Body: "test",
			},
			expectedErr: SysErr{ErrUnsupportedRuleType},
		},
		{
			in:          InvalidData{errors.New("test")},
			expectedErr: SysErr{ErrUnsupportedType},
		},
		{
			in:          "test",
			expectedErr: ErrNotStructData,
		},
		{
			in: Regexp{
				Body: "10",
			},
			expectedErr: nil,
		},
		{
			in: Regexp{"1e"},
			expectedErr: ValidationErrors{
				{
					Field: "Body",
					Err:   ErrValidateRegexp,
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			var valErrs ValidationErrors
			var valExceptedErrs ValidationErrors
			// Проверка ошибок валидации
			if errors.As(err, &valErrs) {
				if errors.As(tt.expectedErr, &valExceptedErrs) {
					for k, e := range valErrs {
						require.Truef(t, errors.Is(e.Err, valExceptedErrs[k].Err), "actual error %q, excepted %q", err, tt.expectedErr)
					}
				} else {
					require.Fail(t, "actual error %q, excepted %q", err, tt.expectedErr)
				}
			}
			// Проверка системных ошибок
			if errors.As(err, &SysErr{}) {
				require.Truef(t, errors.Is(err, tt.expectedErr), "actual error %q, excepted %q", err, tt.expectedErr)
			}
			// Проверка отсутствия ошибок
			if tt.expectedErr == nil {
				require.NoError(t, err)
			}

			_ = tt
		})
	}
}
