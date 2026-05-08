package linkoerr

import "log/slog"

type errWithAttrs struct {
	err   error
	attrs []slog.Attr
}

func (e *errWithAttrs) Error() string {
	return e.err.Error()
}

func (e *errWithAttrs) Unwrap() error {
	return e.err
}

// WithAttrs wraps an error with structured slog attributes.
// args are key-value pairs like slog.Logger methods accept.
func WithAttrs(err error, args ...any) error {
	var attrs []slog.Attr
	for i := 0; i < len(args)-1; i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}
		attrs = append(attrs, slog.Any(key, args[i+1]))
	}
	return &errWithAttrs{err: err, attrs: attrs}
}

// Attrs extracts all slog attributes from the error chain.
func Attrs(err error) []slog.Attr {
	var attrs []slog.Attr
	for err != nil {
		if e, ok := err.(*errWithAttrs); ok {
			attrs = append(attrs, e.attrs...)
		}
		unwrapper, ok := err.(interface{ Unwrap() error })
		if !ok {
			break
		}
		err = unwrapper.Unwrap()
	}
	return attrs
}
