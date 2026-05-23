package membership

import (
	maa "github.com/MaaXYZ/maa-framework-go/v4"
	"github.com/rs/zerolog/log"
)

// MembershipCheckAction is a custom action that checks membership before executing member-only tasks.
// It runs synchronously in the pipeline, blocking execution for non-members.
type MembershipCheckAction struct{}

var _ maa.CustomActionRunner = &MembershipCheckAction{}

// var notifyOnce sync.Once

func (a *MembershipCheckAction) Run(_ *maa.Context, _ *maa.CustomActionArg) bool {
	status := GetMembershipStatus()

	if status.UnsupportedTier {
		log.Warn().
			Str("tier", status.MembershipType).
			Msg("MembershipCheck: unsupported tier")
	}

	if status.IsMember {
		log.Info().
			Str("tier", status.MembershipType).
			Int("level", status.UserLevel).
			Str("expiry", status.VirtualExpiry).
			Msg("MembershipCheck: member verified, allowing")

		// 赞助提示已移除。
		// notifyOnce.Do(func() {
		// 	maafocus.Print(ctx, fmt.Sprintf(
		// 		i18n.T("tasker.membership_check.verified"),
		// 		status.MembershipType, status.VirtualExpiry,
		// 	))
		// 	maafocus.Print(ctx, fmt.Sprintf(
		// 		i18n.T("tasker.membership_check.sponsor"),
		// 		sponsorURL,
		// 	))
		// })
		return true
	}

	// 赞助提示已移除。
	// maafocus.Print(ctx, fmt.Sprintf(
	// 	i18n.T("tasker.membership_check.denied"),
	// 	sponsorURL,
	// ))

	return true
}
