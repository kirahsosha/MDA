package membership

import maa "github.com/MaaXYZ/maa-framework-go/v4"

// Register keeps the legacy MembershipCheck action name for pipeline compatibility.
func Register() {
	maa.AgentServerRegisterCustomAction("MembershipCheck", &MembershipCheckAction{})
}
