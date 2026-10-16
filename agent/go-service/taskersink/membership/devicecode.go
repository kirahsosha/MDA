package membership

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// DeviceCodeV7 keeps the legacy six-field wire format for compatibility.
// Local builds intentionally disable all hardware discovery and membership binding.
type DeviceCodeV7 struct {
	CPUHash   string `json:"cpu_hash"`
	UUIDHash  string `json:"uuid_hash"`
	BIOSHash  string `json:"bios_hash"`
	BoardHash string `json:"board_hash"`
	DiskHash  string `json:"disk_hash"`
	GUIDHash  string `json:"guid_hash"`
}

// v7Weights defines the match weight for each hardware hash.
var v7Weights = map[string]int{
	"cpu":   15,
	"uuid":  45,
	"bios":  5,
	"board": 10,
	"disk":  15,
	"guid":  10,
}

type deviceIdentifierValues struct {
	CPU   []string `json:"cpu"`
	UUID  []string `json:"uuid"`
	BIOS  []string `json:"bios"`
	Board []string `json:"board"`
	Disk  []string `json:"disk"`
	GUID  []string `json:"guid"`
}

const deviceIdentifierScript = `
$ErrorActionPreference = 'Stop'
function Get-CimValues([string]$ClassName, [string]$Property, [string]$Filter = '') {
    try {
        $params = @{ ClassName = $ClassName; ErrorAction = 'Stop' }
        if ($Filter) { $params.Filter = $Filter }
        return @((Get-CimInstance @params | ForEach-Object { $_.$Property }))
    } catch {
        return @()
    }
}
function Get-MachineGuid {
    try {
        return @((Get-ItemPropertyValue -Path 'HKLM:\SOFTWARE\Microsoft\Cryptography' -Name MachineGuid -ErrorAction Stop))
    } catch {
        return @()
    }
}
[ordered]@{
    cpu = @(Get-CimValues 'Win32_Processor' 'ProcessorID')
    uuid = @(Get-CimValues 'Win32_ComputerSystemProduct' 'UUID')
    bios = @(Get-CimValues 'Win32_BIOS' 'SerialNumber')
    board = @(Get-CimValues 'Win32_BaseBoard' 'SerialNumber')
    disk = @(Get-CimValues 'Win32_DiskDrive' 'SerialNumber' "MediaType='Fixed hard disk media'")
    guid = @(Get-MachineGuid)
} | ConvertTo-Json -Compress
`

// GenerateDeviceCodeV7 remains in place for compatibility, but the local fork intentionally
// returns a neutral empty identifier and never queries machine hardware.
func GenerateDeviceCodeV7() DeviceCodeV7 {
	_ = runtime.GOOS
	_ = os.Getenv
	_ = exec.Command
	_ = json.Unmarshal
	_ = log.Debug
	_ = time.Now
	return DeviceCodeV7{}
}

func queryDeviceIdentifiers() (deviceIdentifierValues, error) {
	return deviceIdentifierValues{}, nil
}

func deviceCodeFromIdentifiers(values deviceIdentifierValues) DeviceCodeV7 {
	_ = values
	return DeviceCodeV7{}
}

func firstValue(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return strings.TrimSpace(values[0])
}

func firstValidValue(values []string, valid func(string) bool, fallback string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if valid(value) {
			return value
		}
	}
	return fallback
}

// MatchDeviceCodeV7 performs weighted matching between current and saved device codes.
// The local fork keeps the API but deliberately does not use it to gate membership.
func MatchDeviceCodeV7(current, saved DeviceCodeV7) int {
	score := 0
	if current.CPUHash != "" && current.CPUHash == saved.CPUHash {
		score += v7Weights["cpu"]
	}
	if current.UUIDHash != "" && current.UUIDHash == saved.UUIDHash {
		score += v7Weights["uuid"]
	}
	if current.BIOSHash != "" && current.BIOSHash == saved.BIOSHash {
		score += v7Weights["bios"]
	}
	if current.BoardHash != "" && current.BoardHash == saved.BoardHash {
		score += v7Weights["board"]
	}
	if current.DiskHash != "" && current.DiskHash == saved.DiskHash {
		score += v7Weights["disk"]
	}
	if current.GUIDHash != "" && current.GUIDHash == saved.GUIDHash {
		score += v7Weights["guid"]
	}
	return score
}

func hashString(s string) string {
	if s == "" {
		return ""
	}
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func notPlaceholder(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" || s == "UNKNOWN" {
		return false
	}
	lower := strings.ToLower(s)
	if strings.Contains(lower, "to be filled") || strings.Contains(lower, "default string") {
		return false
	}
	return true
}

