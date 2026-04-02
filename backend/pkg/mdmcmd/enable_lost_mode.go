package mdmcmd

// EnableLostModeOptions contains parameters for the EnableLostMode command.
type EnableLostModeOptions struct {
	// Message is the message to display on the lock screen.
	Message string
	// PhoneNumber is a phone number to display on the lock screen.
	PhoneNumber string
	// Footnote is a footnote string to display on the lock screen.
	Footnote string
}

// EnableLostMode generates an EnableLostMode MDM command.
// This command puts the device into Lost Mode, preventing access even with the passcode.
func (b *CommandBuilder) EnableLostMode(opts *EnableLostModeOptions) ([]byte, string, error) {
	payload := make(map[string]any)

	if opts != nil {
		if opts.Message != "" {
			payload["Message"] = opts.Message
		}
		if opts.PhoneNumber != "" {
			payload["PhoneNumber"] = opts.PhoneNumber
		}
		if opts.Footnote != "" {
			payload["Footnote"] = opts.Footnote
		}
	}

	return buildCommand("EnableLostMode", payload)
}
