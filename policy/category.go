package policy

import (
	"context"

	"gitag.ir/cookthepot/services/vault/services/role"
)

func CanAccessCategory(ctx context.Context) bool {
	return CanCreateCategory(ctx)
}

func CanCreateCategory(ctx context.Context) bool {
	roles := ExtractRolesClaim(ctx)
	for _, r := range roles {
		if r == role.Admin || r == role.Expert {
			return true
		}
	}
	return false
}

func CanUpdateCategory(ctx context.Context) bool {
	return CanCreateCategory(ctx)
}

func CanDeleteCategory(ctx context.Context) bool {
	return CanCreateCategory(ctx)
}
