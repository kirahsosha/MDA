package membership

import (
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// MembershipStatus keeps the historical status shape for compatibility.
// 会员等级限制、设备绑定和远端订阅状态均已移除，当前只返回本地默认值。
type MembershipStatus struct {
	Tier                        string
	TierName                    string
	TierCode                    string
	PlanCode                    string
	PlanName                    string
	StartsOn                    string
	ExpiresOn                   string
	RemainingDays               int
	RegularDailyRuntimeMinutes  int
	DailyRuntimeMinutes         int
	SpecialPeriodRuntimeMinutes int
	UnlimitedRuntime            bool
	IsMember                    bool
	UserID                      string
	DeviceCode                  DeviceCodeV7
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
		TierName:    "普通用户",
		IsMember:    false,
		DeviceCode:  DeviceCodeV7{},
	}

	cachedStatusMu.Lock()
	cachedStatus = status
	cachedStatusMu.Unlock()

	return status
}

var localVersion = "local"

// SetVersion preserves the historical version hook used by tests and legacy callers.
// The local fork keeps the API, but membership checks and remote validation remain disabled.
func SetVersion(version string) {
	version = strings.TrimSpace(version)
	if version == "" {
		localVersion = "local"
		return
	}
	localVersion = version
}

func GetVersion() string {
	return localVersion
}

func cacheStatus(status *MembershipStatus) {
	cachedStatusMu.Lock()
	defer cachedStatusMu.Unlock()
	if status == nil {
		cachedStatus = nil
		return
	}
	cachedStatus = status
}

func testStatus(minutes int, userID string) *MembershipStatus {
	status := GetMembershipStatus()
	status.Tier = "测试用户"
	status.TierName = "测试用户"
	status.TierCode = "test"
	status.PlanCode = "test"
	status.PlanName = "测试计划"
	status.StartsOn = "2024-01-01"
	status.ExpiresOn = "2099-12-31"
	status.RegularDailyRuntimeMinutes = minutes
	status.DailyRuntimeMinutes = minutes
	status.SpecialPeriodRuntimeMinutes = 0
	status.UnlimitedRuntime = false
	status.IsMember = true
	status.UserID = userID
	status.DeviceCode = DeviceCodeV7{}
	return status
}

func isolateQuotaState(t testing.TB) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "membership-quota.json")
	t.Setenv("MDA_MEMBERSHIP_QUOTA_STATE_PATH", path)
	return path
}
