package errs

import "emperror.dev/errors"

// In does what emperror.dev/errors.Is in loop
func In(err error, targets ...error) bool {
	for _, e := range targets {
		if errors.Is(err, e) {
			return true
		}
	}

	return false
}
