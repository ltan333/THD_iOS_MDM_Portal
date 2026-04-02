package mdmcmd

// DisableLostMode generates a DisableLostMode MDM command.
// This command takes the device out of Lost Mode, restoring normal access.
func (b *CommandBuilder) DisableLostMode() ([]byte, string, error) {
	return buildCommand("DisableLostMode", nil)
}
