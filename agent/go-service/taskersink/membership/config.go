package membership

// unsupportedTiers are membership tiers no longer supported in MDA.
// When detected, a log warning is issued but the user is still processed at their actual level.
var unsupportedTiers = map[string]bool{
	"铜Doro会员": true,
	"银Doro会员": true,
}

// membershipLevels maps tier names to their numeric user level.
var membershipLevels = map[string]int{
	"普通用户":      0,
	"铜Doro会员":    1,
	"银Doro会员":    2,
	"金Doro会员":    3,
	"金Doro企业版": 4,
}

// monthlyCost maps tier names to their monthly cost in ORANGE units.
var monthlyCost = map[string]float64{
	"普通用户":      0,
	"铜Doro会员":    1,
	"银Doro会员":    3,
	"金Doro会员":    5,
	"金Doro企业版": 100,
}

// minMemberLevel controls level-based gating in membership status calculation.
// 会员等级限制已移除：原值为 3（Gold tier），现改为 0 以允许所有等级通过。
const minMemberLevel = 0

// enableRemoteMembershipCheck toggles remote membership verification.
const enableRemoteMembershipCheck = false

// MemberDataURL external endpoint has been disabled.
const MemberDataURL = ""
