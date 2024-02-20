package userid

import (
	"context"
	"fmt"
)

// ContextKey is a string type
type ContextKey string

// ContextUserKey constant UserID
const ContextUserKey ContextKey = "UserID"

// Get returns the string representation of the value associated with the ContextUserKey in the given context.
//
// ctx: context.Context
// string
func Get(ctx context.Context) string {
	return fmt.Sprint(ctx.Value(ContextUserKey))
}
