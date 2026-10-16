package membership

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// QuotaRefillType is retained for backward compatibility with older membership logic.
// The local fork intentionally keeps the type and wire format but does not consume the
// underlying remote quota refill metadata.
type QuotaRefillType string

const (
	QuotaRefillTypeNone      QuotaRefillType = "none"
	QuotaRefillTypeSponsor   QuotaRefillType = "sponsor"
	QuotaRefillTypeMember    QuotaRefillType = "member"
	QuotaRefillTypeFallback  QuotaRefillType = "fallback"
)

type quotaMultiplier struct {
	BasePermille  int64
	ExtraPermille int64
}

const multiplierScale int64 = 1000

func unmultipliedQuotaMultiplier() quotaMultiplier {
	return quotaMultiplier{BasePermille: multiplierScale, ExtraPermille: multiplierScale}
}

func regularQuotaMultiplier(entry string) quotaMultiplier {
	_ = entry
	return unmultipliedQuotaMultiplier()
}

func (m quotaMultiplier) totalPermille() int64 {
	return m.BasePermille + m.ExtraPermille
}

func (m quotaMultiplier) billableSecondsFromReal(realSeconds int64, flush bool) int64 {
	_ = flush
	if realSeconds <= 0 {
		return 0
	}
	return realSeconds
}

func isHighConsumptionEntry(entry string) bool {
	_ = entry
	return false
}

func replaceQuotaStateFile(src, dst string) error {
	if src == "" || dst == "" || src == dst {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return os.Rename(src, dst)
}

func acquireQuotaFileLock(path string) (func() error, error) {
	if path == "" {
		return func() error { return nil }, nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	return func() error {
		if file == nil {
			return nil
		}
		return file.Close()
	}, nil
}

// SponsorURL preserves the legacy URL contract without contacting remote services.
func SponsorURL(status *MembershipStatus) string {
	if status == nil {
		return ""
	}
	if status.DeviceCode == (DeviceCodeV7{}) {
		return "https://member.maa.plus/sponsor"
	}
	payload, err := json.Marshal(status.DeviceCode)
	if err != nil {
		return "https://member.maa.plus/sponsor"
	}
	return "https://member.maa.plus/sponsor?device=" + url.QueryEscape(base64.RawURLEncoding.EncodeToString(payload))
}

// DeviceCodeFromSponsorURL decodes the legacy sponsor URL payload back into the stable
// six-field device representation, keeping the compatibility API intact without enabling
// membership enforcement.
func DeviceCodeFromSponsorURL(rawURL string) (DeviceCodeV7, error) {
	if rawURL == "" {
		return DeviceCodeV7{}, fmt.Errorf("empty sponsor url")
	}
	if strings.Contains(rawURL, "://") {
		parsed, err := url.Parse(rawURL)
		if err == nil && parsed != nil {
			if deviceValue := parsed.Query().Get("device"); deviceValue != "" {
				rawURL = deviceValue
			}
		}
	}
	if rawURL == "" {
		return DeviceCodeV7{}, fmt.Errorf("missing device payload")
	}
	if decoded, err := base64.RawURLEncoding.DecodeString(rawURL); err == nil {
		rawURL = string(decoded)
	}
	var device DeviceCodeV7
	if err := json.Unmarshal([]byte(rawURL), &device); err != nil {
		return DeviceCodeV7{}, err
	}
	return device, nil
}
