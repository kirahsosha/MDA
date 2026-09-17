package membership

import maa "github.com/MaaXYZ/maa-framework-go/v4"

// Register keeps the legacy MembershipCheck action name for pipeline compatibility.
// RuntimeQuotaCheck and QuotaDisplayAction are also registered as aliases
// so upstream pipelines work without membership logic.
func Register() {
	maa.AgentServerRegisterCustomAction("MembershipCheck", &MembershipCheckAction{})
	maa.AgentServerRegisterCustomAction("RuntimeQuotaCheck", &MembershipCheckAction{})
	maa.AgentServerRegisterCustomAction("QuotaDisplayAction", &MembershipCheckAction{})
}
