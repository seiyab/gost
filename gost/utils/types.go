package utils

import "go/types"

func IsError(t types.Type) bool {
	e := types.Universe.Lookup("error")
	return types.Identical(t, e.Type())
}
