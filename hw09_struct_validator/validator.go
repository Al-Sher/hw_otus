package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrNotStructData       = errors.New("элемент не является структурой")
	ErrValidateMin         = errors.New("элемент меньше минимально допустимого значения")
	ErrValidateMax         = errors.New("элемент больше максимально допустимого значения")
	ErrUnsupportedRuleType = errors.New("неподдерживаемый тип правила")
	ErrUnsupportedType     = errors.New("неподдерживаемый тип данных")
	ErrValidateIntInt      = errors.New("число находится за пределами допустимых значений")
	ErrValidateLen         = errors.New("длина строки не соответствует допустимому значению")
	ErrValidateRegexp      = errors.New("строка не соответствует регулярному выражению")
	ErrValidateStringIn    = errors.New("строка находится за пределами допустимых значений")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errs := make([]string, 0, len(v))
	for _, err := range v {
		errs = append(errs, fmt.Sprintf("%s = %s", err.Field, err.Err))
	}

	return strings.Join(errs, ";")
}

type SysErr struct {
	Err error
}

func (s SysErr) Error() string {
	return s.Err.Error()
}

func Validate(v interface{}) error {
	var validationErrors ValidationErrors

	valuesByStruct := reflect.ValueOf(v)
	typeByStruct := valuesByStruct.Type()

	if typeByStruct.Kind() != reflect.Struct {
		return ErrNotStructData
	}

	for i := 0; i < valuesByStruct.NumField(); i++ {
		fieldType := typeByStruct.Field(i)
		value := valuesByStruct.Field(i)

		validateTag := fieldType.Tag
		rules := validateTag.Get("validate")

		if rules == "" || !fieldType.IsExported() {
			continue
		}

		for _, rule := range strings.Split(rules, "|") {
			keyWithValue := strings.Split(rule, ":")
			if err := validateRule(value, fieldType.Type, keyWithValue[0], keyWithValue[1]); err != nil {
				var e SysErr
				if errors.As(err, &e) {
					return e.Err
				}

				validationError := ValidationError{
					Field: fieldType.Name,
					Err:   err,
				}
				validationErrors = append(validationErrors, validationError)
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateRule(value reflect.Value, valueType reflect.Type, keyRule string, valueRule string) error {
	switch valueType.Kind() { //nolint:exhaustive
	case reflect.Int, reflect.Int32, reflect.Int64:
		switch keyRule {
		case "min":
			return validateMinRule(int(value.Int()), valueRule)
		case "max":
			return validateMaxRule(int(value.Int()), valueRule)
		case "in":
			return validateIntInRule(int(value.Int()), valueRule)
		default:
			return SysErr{ErrUnsupportedRuleType}
		}
	case reflect.String:
		switch keyRule {
		case "len":
			return validateLenRule(value.String(), valueRule)
		case "regexp":
			return validateRegexpRule(value.String(), valueRule)
		case "in":
			return validateStringInRule(value.String(), valueRule)
		default:
			return SysErr{ErrUnsupportedRuleType}
		}
	case reflect.Slice, reflect.Array:
		if value.IsNil() {
			return nil
		}

		for i := 0; i < value.Len(); i++ {
			if err := validateRule(value.Index(i), valueType.Elem(), keyRule, valueRule); err != nil {
				return err
			}
		}

		return nil
	default:
		return SysErr{ErrUnsupportedType}
	}
}

func validateMinRule(value int, valueRule string) error {
	valueIntRule, err := strconv.Atoi(valueRule)
	if err != nil {
		return SysErr{err}
	}

	if value < valueIntRule {
		return ErrValidateMin
	}

	return nil
}

func validateMaxRule(value int, valueRule string) error {
	valueIntRule, err := strconv.Atoi(valueRule)
	if err != nil {
		return SysErr{err}
	}

	if value > valueIntRule {
		return ErrValidateMax
	}

	return nil
}

func validateIntInRule(value int, valueRule string) error {
	in := strings.Split(valueRule, ",")

	for _, v := range in {
		i, err := strconv.Atoi(v)
		if err != nil {
			return err
		}

		if value == i {
			return nil
		}
	}

	return ErrValidateIntInt
}

func validateLenRule(value string, valueRule string) error {
	l, err := strconv.Atoi(valueRule)
	if err != nil {
		return SysErr{err}
	}

	if len(value) != l {
		return ErrValidateLen
	}

	return nil
}

func validateRegexpRule(value string, valueRule string) error {
	re, err := regexp.Compile(valueRule)
	if err != nil {
		return SysErr{err}
	}

	if !re.MatchString(value) {
		return ErrValidateRegexp
	}

	return nil
}

func validateStringInRule(value string, valueRule string) error {
	in := strings.Split(valueRule, ",")

	for _, v := range in {
		if value == v {
			return nil
		}
	}

	return ErrValidateStringIn
}
