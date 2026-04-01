package serviceimpl

import (
	"bytes"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	"howett.net/plist"
)

type profileGeneratorImpl struct {
	orgName string
	prefix  string
}

func NewProfileGenerator(orgName, prefix string) service.ProfileGenerator {
	if prefix == "" {
		prefix = "com.thd.mdm"
	}
	return &profileGeneratorImpl{
		orgName: orgName,
		prefix:  prefix,
	}
}

func (g *profileGeneratorImpl) GenerateXML(ctx context.Context, p *ent.Profile) ([]byte, error) {
	profileMap := map[string]any{
		"PayloadDisplayName":       p.Name,
		"PayloadIdentifier":        fmt.Sprintf("%s.profile.%d", g.prefix, p.ID),
		"PayloadOrganization":      g.orgName,
		"PayloadRemovalDisallowed": false,
		"PayloadType":              "Configuration",
		"PayloadUUID":              uuid.New().String(),
		"PayloadVersion":           1,
	}

	var payloadContent []any

	if len(p.SecuritySettings) > 0 {
		payloadContent = append(payloadContent,
			g.createPayload("com.apple.mobiledevice.passwordpolicy", g.mapSecuritySettings(p.SecuritySettings)))
	}

	if len(p.Restrictions) > 0 {
		payloadContent = append(payloadContent,
			g.createPayload("com.apple.applicationaccess", g.mapRestrictions(p.Restrictions)))
	}

	if len(p.NetworkConfig) > 0 {
		payloadContent = append(payloadContent, g.mapNetworkPayloads(p.NetworkConfig)...)
	}

	if len(p.ContentFilter) > 0 {
		if cp := g.mapContentFilter(p.ContentFilter); cp != nil {
			payloadContent = append(payloadContent,
				g.createPayload("com.apple.webcontent-filter", cp))
		}
	}

	profileMap["PayloadContent"] = payloadContent

	var buf bytes.Buffer
	encoder := plist.NewEncoder(&buf)
	encoder.Indent("\t")
	if err := encoder.Encode(profileMap); err != nil {
		return nil, fmt.Errorf("failed to encode profile to plist: %w", err)
	}

	return buf.Bytes(), nil
}

// mapSecuritySettings converts portal security_settings to com.apple.mobiledevice.passwordpolicy keys.
// Apple MDM reference: https://developer.apple.com/documentation/devicemanagement/passwordpolicy
func (g *profileGeneratorImpl) mapSecuritySettings(s map[string]any) map[string]any {
	out := map[string]any{}

	if v, ok := boolVal(s, "passcode_required"); ok && v {
		out["forcePIN"] = true
	}
	if v, ok := intVal(s, "min_passcode_length"); ok && v > 0 {
		out["minLength"] = v
	}
	if v, ok := intVal(s, "max_failed_attempts"); ok && v > 0 {
		out["maxFailedAttempts"] = v
	}
	// screen_lock_timeout is stored in seconds; Apple maxInactivity is in minutes (1–5, 0=never).
	if v, ok := intVal(s, "screen_lock_timeout"); ok {
		minutes := v / 60
		if minutes < 1 {
			minutes = 1
		}
		out["maxInactivity"] = minutes
	}
	// encryption_enabled: on iOS, encryption is automatic when a passcode is set.
	// No separate key needed; forcePIN above achieves this.

	return out
}

// mapRestrictions converts portal restriction keys to com.apple.applicationaccess keys.
// Portal keys are "xxx_disabled: true"; Apple MDM uses "allowXxx: false" (inverted).
// Apple MDM reference: https://developer.apple.com/documentation/devicemanagement/restrictions
func (g *profileGeneratorImpl) mapRestrictions(r map[string]any) map[string]any {
	keyMap := map[string]string{
		"airdrop_disabled":   "allowAirDrop",
		"camera_disabled":    "allowCamera",
		"bluetooth_disabled": "allowBluetooth",
		"app_store_disabled": "allowAppInstallation",
	}

	out := make(map[string]any, len(keyMap))
	for portalKey, appleKey := range keyMap {
		if disabled, ok := boolVal(r, portalKey); ok {
			out[appleKey] = !disabled
		}
	}
	return out
}

// mapNetworkPayloads generates one payload per network type (WiFi, VPN).
// Each is a separate MDM payload entry in PayloadContent.
func (g *profileGeneratorImpl) mapNetworkPayloads(nc map[string]any) []any {
	var payloads []any

	if wifiRaw, ok := nc["wifi"].(map[string]any); ok {
		wifi := map[string]any{
			"AutoJoin": true,
		}
		if ssid, ok := wifiRaw["ssid"].(string); ok && ssid != "" {
			wifi["SSID_STR"] = ssid
		}
		if sec, ok := wifiRaw["security_type"].(string); ok && sec != "" {
			// Normalize to Apple-accepted values: WPA, WPA2, WPA3, WEP, None, Any
			switch sec {
			case "WPA2":
				wifi["EncryptionType"] = "WPA2"
			case "WPA3":
				wifi["EncryptionType"] = "WPA3"
			case "WPA":
				wifi["EncryptionType"] = "WPA"
			case "WEP":
				wifi["EncryptionType"] = "WEP"
			default:
				wifi["EncryptionType"] = "Any"
			}
		}
		if pw, ok := wifiRaw["password"].(string); ok && pw != "" {
			wifi["Password"] = pw
		}
		payloads = append(payloads, g.createPayload("com.apple.wifi.managed", wifi))
	}

	if vpnRaw, ok := nc["vpn"].(map[string]any); ok {
		if enabled, ok := boolVal(vpnRaw, "enabled"); ok && enabled {
			vpn := map[string]any{
				"VPNType": "IKEv2",
			}
			if server, ok := vpnRaw["server"].(string); ok {
				vpn["RemoteAddress"] = server
			}
			if user, ok := vpnRaw["username"].(string); ok {
				vpn["AuthName"] = user
			}
			if pw, ok := vpnRaw["password"].(string); ok {
				vpn["AuthPassword"] = pw
			}
			payloads = append(payloads, g.createPayload("com.apple.vpn.managed", vpn))
		}
	}

	return payloads
}

// mapContentFilter converts portal content_filter to com.apple.webcontent-filter keys.
// Returns nil if no actionable settings are present.
func (g *profileGeneratorImpl) mapContentFilter(cf map[string]any) map[string]any {
	enabled, _ := boolVal(cf, "safe_browsing_enabled")
	blockedRaw, _ := cf["blocked_websites"].([]any)

	if !enabled && len(blockedRaw) == 0 {
		return nil
	}

	out := map[string]any{
		"FilterType": "BuiltIn",
	}
	if enabled {
		out["AutoFilterEnabled"] = true
	}
	if len(blockedRaw) > 0 {
		var blocked []map[string]any
		for _, item := range blockedRaw {
			if url, ok := item.(string); ok && url != "" {
				blocked = append(blocked, map[string]any{"URL": url})
			}
		}
		if len(blocked) > 0 {
			out["BlacklistedURLs"] = blocked
		}
	}
	return out
}

func (g *profileGeneratorImpl) createPayload(payloadType string, settings map[string]any) map[string]any {
	payload := map[string]any{
		"PayloadDisplayName":  payloadType,
		"PayloadIdentifier":   fmt.Sprintf("%s.payload.%s.%s", g.prefix, payloadType, uuid.New().String()),
		"PayloadOrganization": g.orgName,
		"PayloadType":         payloadType,
		"PayloadUUID":         uuid.New().String(),
		"PayloadVersion":      1,
	}

	for k, v := range settings {
		if k == "PayloadType" || k == "PayloadUUID" || k == "PayloadIdentifier" ||
			k == "PayloadOrganization" || k == "PayloadVersion" || k == "PayloadDisplayName" {
			continue
		}
		payload[k] = v
	}

	return payload
}

// boolVal safely extracts a bool from a map[string]any.
func boolVal(m map[string]any, key string) (bool, bool) {
	v, ok := m[key]
	if !ok {
		return false, false
	}
	b, ok := v.(bool)
	return b, ok
}

// intVal safely extracts an int from a map[string]any, accepting float64 (JSON default) and int.
func intVal(m map[string]any, key string) (int, bool) {
	v, ok := m[key]
	if !ok {
		return 0, false
	}
	switch n := v.(type) {
	case int:
		return n, true
	case float64:
		return int(n), true
	case int64:
		return int(n), true
	}
	return 0, false
}
