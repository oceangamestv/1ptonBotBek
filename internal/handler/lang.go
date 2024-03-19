package handler

import "fmt"

func langf(formats map[string]string, key string, args ...any) string {
	f, ok := formats[key]
	if !ok {
		f = formats["en"]
	}
	return fmt.Sprintf(f, args...)
}
