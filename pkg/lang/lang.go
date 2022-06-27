package lang

import (
	"context"

	"google.golang.org/grpc/metadata"
)

type Language string

var (
	EN Language = "en-US"
	FA Language = "fa-IR"
)

// Get func, with this method we check accept-language
// if context is grpc-context we check grpcgateway-accept-language
// if context is normal context we check accept-language
func Get(ctx context.Context) Language {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		val := ctx.Value("accept-language")
		if val == nil {
			return FA
		}
		return Language(val.(string))
	}

	if len(headers.Get("grpcgateway-accept-language")) == 0 {
		return FA
	}
	return Language(headers.Get("grpcgateway-accept-language")[0])
}
