package mdmcmd

// ShutDownDevice generates a ShutDownDevice MDM command.
// This command shuts down the device.
// Supported on iOS 10.3+ and macOS 10.13+.
func (b *CommandBuilder) ShutDownDevice() ([]byte, string, error) {
	return buildCommand("ShutDownDevice", nil)
}
