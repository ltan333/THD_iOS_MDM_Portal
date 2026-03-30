package mdmcmd

// EraseDeviceOptions contains optional parameters for the EraseDevice command.
type EraseDeviceOptions struct {
	// PIN is a 6-digit PIN for macOS devices (required for Apple Silicon Macs).
	PIN string
	// PreserveDataPlan preserves the eSIM/data plan on iOS 11+ devices.
	PreserveDataPlan bool
	// DisallowProximitySetup prevents Quick Start on iOS 11+ devices.
	DisallowProximitySetup bool
	// ObliterationBehavior controls cryptographic erase behavior.
	// Values: "Default", "DoNotObliterate", "ObliterateWithWarning", "Always"
	ObliterationBehavior string
	// ReturnToService configures automatic re-enrollment after wipe.
	ReturnToService *ReturnToServiceConfig
}

// ReturnToServiceConfig configures Return to Service functionality.
type ReturnToServiceConfig struct {
	// Enabled activates Return to Service mode.
	Enabled bool
	// MDMProfileData is base64-encoded MDM enrollment profile.
	MDMProfileData []byte
	// WiFiProfileData is base64-encoded WiFi configuration profile.
	WiFiProfileData []byte
}

// EraseDevice generates an EraseDevice MDM command.
// This command factory resets the device, removing all data.
func (b *CommandBuilder) EraseDevice(opts *EraseDeviceOptions) ([]byte, string, error) {
	payload := make(map[string]any)

	if opts != nil {
		if opts.PIN != "" {
			payload["PIN"] = opts.PIN
		}
		if opts.PreserveDataPlan {
			payload["PreserveDataPlan"] = true
		}
		if opts.DisallowProximitySetup {
			payload["DisallowProximitySetup"] = true
		}
		if opts.ObliterationBehavior != "" {
			payload["ObliterationBehavior"] = opts.ObliterationBehavior
		}
		if opts.ReturnToService != nil && opts.ReturnToService.Enabled {
			rts := map[string]any{"Enabled": true}
			if opts.ReturnToService.MDMProfileData != nil {
				rts["MDMProfileData"] = opts.ReturnToService.MDMProfileData
			}
			if opts.ReturnToService.WiFiProfileData != nil {
				rts["WiFiProfileData"] = opts.ReturnToService.WiFiProfileData
			}
			payload["ReturnToService"] = rts
		}
	}

	return buildCommand("EraseDevice", payload)
}
