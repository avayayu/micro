package http

import (
	"fmt"

	criticalityPkg "github.com/avayayu/micro/net/criticality"
	"github.com/avayayu/micro/net/metadata"
)

// Criticality is
func Criticality(pathCriticality criticalityPkg.Criticality) HandlerFunc {
	if !criticalityPkg.Exist(pathCriticality) {
		panic(fmt.Errorf("This criticality is not exist: %s", pathCriticality))
	}
	return func(ctx *Context) {
		md, ok := metadata.FromContext(ctx)
		if ok {
			md[metadata.Criticality] = string(pathCriticality)
		}
	}
}
