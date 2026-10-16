package membership

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

// 只接受正式会员映射所需的最小快照；未声明的诊断字段不会参与映射。
type androidIdentityRow struct {
	Name   string   `json:"name"`
	Status string   `json:"status"`
	Values []string `json:"values"`
}

type androidIdentitySnapshot struct {
	SchemaVersion int                  `json:"schema_version"`
	Rows          []androidIdentityRow `json:"rows"`
}

var androidHex = regexp.MustCompile(`^[0-9a-fA-F]+$`)

func androidMappingHash(field, source, value string) string {
	if value == "" {
		return ""
	}
	return hashString("mda:android:client-v1:" + field + ":" + source + "\x00" + value)
}

func normalizeAndroidID(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	// 64 位整数的十六进制形式可不补前导零；归一成 16 位，拒绝零和已知占位值。
	if len(value) == 0 || len(value) > 16 || !androidHex.MatchString(value) || strings.Trim(value, "0") == "" || value == "9774d56d682e549c" {
		return ""
	}
	return strings.Repeat("0", 16-len(value)) + value
}

func androidDeviceCodeFromJSON(raw string) DeviceCodeV7 {
	_ = raw
	log.Warn().Msg("Android membership identity is disabled in this local build")
	return DeviceCodeV7{}
}

func validAndroidHash(value string) bool {
	return len(value) == 64 && androidHex.MatchString(value) && strings.Trim(value, "0") != ""
}

func isAndroidMappedDevice(device DeviceCodeV7) bool {
	_ = device
	return false
}
