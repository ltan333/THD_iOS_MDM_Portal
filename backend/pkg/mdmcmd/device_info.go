package mdmcmd

// DeviceInformation generates a DeviceInformation MDM command.
// The queries parameter specifies which device attributes to return.
// If queries is empty, all available attributes are returned.
//
// Common queries include:
// - DeviceName, OSVersion, BuildVersion, ModelName, Model, ProductName
// - SerialNumber, UDID, IMEI, MEID, PhoneNumber
// - BatteryLevel, AvailableDeviceCapacity, DeviceCapacity
// - WiFiMAC, BluetoothMAC, EthernetMAC
// - IsSupervised, IsDeviceLocatorServiceEnabled, IsActivationLockEnabled
func (b *CommandBuilder) DeviceInformation(queries []string) ([]byte, string, error) {
	payload := make(map[string]any)

	if len(queries) > 0 {
		payload["Queries"] = queries
	}

	return buildCommand("DeviceInformation", payload)
}

// CommonDeviceQueries returns a list of commonly used device information queries.
func CommonDeviceQueries() []string {
	return []string{
		"DeviceName",
		"OSVersion",
		"BuildVersion",
		"ModelName",
		"Model",
		"ProductName",
		"SerialNumber",
		"UDID",
		"BatteryLevel",
		"AvailableDeviceCapacity",
		"DeviceCapacity",
		"WiFiMAC",
		"IsSupervised",
	}
}
