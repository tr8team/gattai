package core_action

import (
	"fmt"
	"strings"
	"strconv"
)

const (
	CmpEqual string = "equal"
	CmpNotEqual		= "not_equal"
	CmpContain		= "contain"
	CmpNotContain	= "not_contain"
	CmpIntLessThan  = "int_less_than"
	CmpIntMoreThan  = "int_more_than"
)

type Comparison struct {
	Condition string
	Value string
}

func ExpectedTest(expected string, conditon string, expected_value string) (bool,error) {
	result := false
	switch conditon {
	case CmpEqual:
		result = (strings.TrimSpace(expected) == strings.TrimSpace(expected_value))
	case CmpNotEqual:
		result = (strings.TrimSpace(expected) != strings.TrimSpace(expected_value))
	case CmpContain:
		result = strings.Contains(strings.TrimSpace(expected), strings.TrimSpace(expected_value))
	case CmpNotContain:
		result =!strings.Contains(strings.TrimSpace(expected), strings.TrimSpace(expected_value))
	case CmpIntLessThan:
		exp_int, err := strconv.Atoi(expected)
		if  err != nil {
			return result, fmt.Errorf("ExpectedTest strconvAtoi error: %s error: %v",expected, err)
		}
		exp_val, err := strconv.Atoi(expected_value)
		if  err != nil {
			return result, fmt.Errorf("ExpectedTest strconvAtoi error: %s error: %v",expected_value, err)
		}
		result = (exp_int < exp_val)
	case CmpIntMoreThan:
		exp_int, err := strconv.Atoi(expected)
		if  err != nil {
			return result, fmt.Errorf("ExpectedTest strconvAtoi error: %s error: %v",expected, err)
		}
		exp_val, err := strconv.Atoi(expected_value)
		if  err != nil {
			return result, fmt.Errorf("ExpectedTest strconvAtoi error: %s error: %v",expected_value, err)
		}
		result = (exp_int > exp_val)
	default:
		return result, fmt.Errorf("ExpectedTest condition is not supported error: %s",conditon)
	}
	return result, nil
}
