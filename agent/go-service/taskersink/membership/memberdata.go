package membership

import "sync"

// MembershipStatus keeps the historical status shape for compatibility.
// 会员等级限制、设备绑定和远端订阅状态均已移除，当前只返回本地默认值。
type MembershipStatus struct {
	Tier          string
	PlanCode      string
	PlanName      string
	StartsOn      string
	ExpiresOn     string
	RemainingDays int
	IsMember      bool
	UserID        string
	DeviceCode    DeviceCodeV7
}

var (
	cachedStatus   *MembershipStatus
	cachedStatusMu sync.RWMutex
)

// GetMembershipStatus returns a local default status.
// 保留该函数仅为了兼容历史调用面，当前不再执行会员校验、远端请求或机器码采集。
func GetMembershipStatus() *MembershipStatus {
	cachedStatusMu.RLock()
	if cachedStatus != nil {
		status := cachedStatus
		cachedStatusMu.RUnlock()
		return status
	}
	cachedStatusMu.RUnlock()

	status := &MembershipStatus{
		Tier:        "普通用户",
		IsMember:    false,
		DeviceCode:  DeviceCodeV7{},
	}

	cachedStatusMu.Lock()
	cachedStatus = status
	cachedStatusMu.Unlock()

	return status
}
