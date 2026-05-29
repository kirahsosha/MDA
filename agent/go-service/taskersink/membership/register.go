package membership

import maa "github.com/MaaXYZ/maa-framework-go/v4"

// Register keeps the legacy MembershipCheck action name for pipeline compatibility.
// RuntimeQuotaCheck is also registered as an alias so upstream pipelines work without membership logic.
func Register() {
	maa.AgentServerRegisterCustomAction("MembershipCheck", &MembershipCheckAction{})
	maa.AgentServerRegisterCustomAction("RuntimeQuotaCheck", &MembershipCheckAction{})
}
