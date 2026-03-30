package mdmcmd

// BuildInstallApplication creates the raw plist bytes for the InstallApplication command.
func (b *CommandBuilder) BuildInstallApplication(manifestURL string, bundleID string) ([]byte, string, error) {
	payload := map[string]any{
		"ManagementFlags": 1, // Remove app when MDM profile is removed
	}

	if manifestURL != "" {
		payload["ManifestURL"] = manifestURL
	} else if bundleID != "" {
		payload["Identifier"] = bundleID
	}

	return buildCommand("InstallApplication", payload)
}
