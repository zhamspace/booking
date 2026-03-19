package session

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
)

type contextKey string

const (
	contextKeySession = contextKey("session")
)

func ContextWithSession(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, contextKeySession, sess)
}

func ContextWithToken(ctx context.Context, token string) context.Context {
	return ContextWithSession(ctx, parseJwt(token))
}

func ExtractFromContext(ctx context.Context) *Session {
	ctxValue := ctx.Value(contextKeySession)
	if ctxValue != nil {
		if sess, ok := ctxValue.(*Session); ok {
			return sess
		}
	}

	return &Session{}
}

func parseJwt(token string) *Session {
	sess := &Session{}

	tokenParts := strings.Split(token, ".")
	if len(tokenParts) == 3 {
		if claimsRaw, err := base64.RawURLEncoding.DecodeString(tokenParts[1]); err == nil {
			if json.Unmarshal(claimsRaw, sess) != nil {
				*sess = Session{}
			}
		}
	}

	return sess
}
