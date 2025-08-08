package env

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var ErrNotSet = errors.New("the env variable is not specified")

func String(
	name string, validationFunc func(v string) error,
) (string, error) {
	v, set := os.LookupEnv(name)
	if !set {
		return "", ErrNotSet
	}
	if validationFunc != nil {
		if err := validationFunc(v); err != nil {
			return "", err
		}
	}

	return v, nil
}

func Int(name string, validationFunc func(v int) error) (int, error) {
	vStr, set := os.LookupEnv(name)
	if !set {
		return 0, ErrNotSet
	}

	v, err := strconv.Atoi(vStr)
	if err != nil {
		return 0, fmt.Errorf("env %s invalid 'int' value: %w", name, err)
	}

	if validationFunc != nil {
		if err := validationFunc(v); err != nil {
			return 0, err
		}
	}

	return v, nil
}

func StringS(
	name string, validationFunc func(v []string) error,
) ([]string, error) {
	v, set := os.LookupEnv(name)
	if !set {
		return nil, ErrNotSet
	}

	s := strings.Split(v, ",")

	if validationFunc != nil {
		if err := validationFunc(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}
