package membership

import (
	maa "github.com/MaaXYZ/maa-framework-go/v4"
)

// MembershipCheckAction is a custom action that checks membership before executing member-only tasks.
// Membership logic has been disabled and this action now acts as a no-op.
type MembershipCheckAction struct{}

var _ maa.CustomActionRunner = &MembershipCheckAction{}

// var notifyOnce sync.Once

func (a *MembershipCheckAction) Run(_ *maa.Context, _ *maa.CustomActionArg) bool {
	// 会员等级限制、远端校验与赞助逻辑均已停用。
	// 保留 custom action 入口以兼容现有 pipeline 配置，但不再执行任何检查。
	return true
}
