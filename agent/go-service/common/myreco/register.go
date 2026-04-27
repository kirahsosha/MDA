package myreco

import maa "github.com/MaaXYZ/maa-framework-go/v4"

// Register registers the sample custom recognition with the agent server.
func Register() {
	maa.AgentServerRegisterCustomRecognition("my_reco_222", &MyCustomRecognition{})
}
