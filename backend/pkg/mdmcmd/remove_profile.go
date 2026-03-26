package mdmcmd

// RemoveProfile generates a RemoveProfile MDM command.
// The identifier is the PayloadIdentifier of the profile to remove.
func (b *CommandBuilder) RemoveProfile(identifier string) ([]byte, string, error) {
	payload := map[string]any{
		"Identifier": identifier,
	}

	return buildCommand("RemoveProfile", payload)
}
