package errutil

import (
	"errors"
	"fmt"
	"log/slog"
	"order-bot-mgmt-svc/internal/util"
	"strings"
)

func FormatErrChain(err error) string {
	var b strings.Builder
	level := 0

	tmpstr := err.Error()
	for err != nil {
		_, fprintErr := fmt.Fprintf(
			&b,
			util.If[string](level == 0, "%s\n", "↳ %s\n"),
			err.Error(),
		)
		if fprintErr != nil {
			slog.Error("errutil.FormatErrChain", "err", fprintErr)
			return tmpstr
		}
		err = errors.Unwrap(err)
		level++
	}

	return b.String()
}
