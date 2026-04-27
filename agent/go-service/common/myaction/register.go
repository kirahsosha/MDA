package myaction

import maa "github.com/MaaXYZ/maa-framework-go/v4"

// Register registers the sample custom action with the agent server.
func Register() {
	maa.AgentServerRegisterCustomAction("my_action_111", &MyCustomAction{})
}
