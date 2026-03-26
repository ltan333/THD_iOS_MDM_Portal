package mdmcmd

// RestartDeviceOptions contains optional parameters for the RestartDevice command.
type RestartDeviceOptions struct {
	// NotifyUser shows a notification to the user before restart (macOS only).
	NotifyUser bool
}

// RestartDevice generates a RestartDevice MDM command.
// This command restarts the device.
// Supported on iOS 10.3+ and macOS 10.13+.
func (b *CommandBuilder) RestartDevice(opts *RestartDeviceOptions) ([]byte, string, error) {
	payload := make(map[string]any)

	if opts != nil && opts.NotifyUser {
		payload["NotifyUser"] = true
	}

	return buildCommand("RestartDevice", payload)
}
