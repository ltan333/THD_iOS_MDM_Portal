package mdmcmd

// DeviceLockOptions contains optional parameters for the DeviceLock command.
type DeviceLockOptions struct {
	// PIN is a 6-digit PIN for macOS devices.
	PIN string
	// Message is the lock screen message to display.
	Message string
	// PhoneNumber is a phone number to display on the lock screen.
	PhoneNumber string
}

// DeviceLock generates a DeviceLock MDM command.
// This command locks the device immediately.
// On macOS, a PIN is required to unlock. On iOS, the existing passcode is used.
func (b *CommandBuilder) DeviceLock(opts *DeviceLockOptions) ([]byte, string, error) {
	payload := make(map[string]any)

	if opts != nil {
		if opts.PIN != "" {
			payload["PIN"] = opts.PIN
		}
		if opts.Message != "" {
			payload["Message"] = opts.Message
		}
		if opts.PhoneNumber != "" {
			payload["PhoneNumber"] = opts.PhoneNumber
		}
	}

	return buildCommand("DeviceLock", payload)
}
