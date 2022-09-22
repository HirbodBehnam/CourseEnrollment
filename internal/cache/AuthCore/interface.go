package AuthCore

import "context"

type Interface interface {
	// Set must set the token for user
	Set(ctx context.Context, token string) error
	// Delete must revoke a token of user
	Delete(ctx context.Context, token string) error
	// IsValid must check if a token is valid or not
	IsValid(ctx context.Context, token string) (bool, error)
	// Replace must do something like token refresh. Delete followed by Set
	Replace(ctx context.Context, old, new string) error
}
