// Package event provides a lightweight in-process event bus for decoupling
// services that would otherwise create circular dependencies.
package event

import "sync"

// DeviceEnrolledEvent is published when a device completes MDM enrollment
// (TokenUpdate check-in received). Consumers use it to trigger actions like
// profile deployment and device inventory sync without the device service
// needing to know about those consumers directly.
type DeviceEnrolledEvent struct {
	// DeviceID is the MDM enrollment ID (UDID) of the newly enrolled device.
	DeviceID string
	// SerialNumber is the hardware serial number, may be empty if not provided.
	SerialNumber string
}

// DeviceCheckedOutEvent is published when a device unenrolls (CheckOut check-in).
type DeviceCheckedOutEvent struct {
	DeviceID string
}

// DeviceInformationReceivedEvent is published when a DeviceInformation command
// result is received via the Acknowledge webhook.
type DeviceInformationReceivedEvent struct {
	DeviceID string
	// QueryResponses contains the raw key→value map from the MDM QueryResponses dict.
	QueryResponses map[string]any
}

// ProfileInstallAckEvent is published when the device acknowledges an
// InstallProfile MDM command (success or failure).
type ProfileInstallAckEvent struct {
	// UDID is the MDM enrollment identifier of the device.
	UDID string
	// CommandUUID is the UUID of the MDM command that was acknowledged.
	// Used to look up the corresponding deployment status record.
	CommandUUID string
	// Status is "Acknowledged" for success, "Error" or "CommandFormatError" for failure.
	Status string
	// ErrorMessage contains the device-reported error, if any.
	ErrorMessage string
}

// Bus is a simple goroutine-safe publish/subscribe event bus backed by
// buffered channels. It is intentionally minimal: one channel per event type,
// multiple subscribers share the same channel (fan-out via goroutines).
type Bus struct {
	mu             sync.RWMutex
	enrolled       []chan DeviceEnrolledEvent
	checkedOut     []chan DeviceCheckedOutEvent
	devInfo        []chan DeviceInformationReceivedEvent
	profileInstall []chan ProfileInstallAckEvent
}

// NewBus creates a new event bus.
func NewBus() *Bus {
	return &Bus{}
}

// SubscribeEnrolled registers a buffered channel for DeviceEnrolledEvent.
// bufSize controls how many unprocessed events can queue before the publisher blocks.
func (b *Bus) SubscribeEnrolled(bufSize int) <-chan DeviceEnrolledEvent {
	ch := make(chan DeviceEnrolledEvent, bufSize)
	b.mu.Lock()
	b.enrolled = append(b.enrolled, ch)
	b.mu.Unlock()
	return ch
}

// SubscribeCheckedOut registers a buffered channel for DeviceCheckedOutEvent.
func (b *Bus) SubscribeCheckedOut(bufSize int) <-chan DeviceCheckedOutEvent {
	ch := make(chan DeviceCheckedOutEvent, bufSize)
	b.mu.Lock()
	b.checkedOut = append(b.checkedOut, ch)
	b.mu.Unlock()
	return ch
}

// SubscribeDeviceInformation registers a buffered channel for DeviceInformationReceivedEvent.
func (b *Bus) SubscribeDeviceInformation(bufSize int) <-chan DeviceInformationReceivedEvent {
	ch := make(chan DeviceInformationReceivedEvent, bufSize)
	b.mu.Lock()
	b.devInfo = append(b.devInfo, ch)
	b.mu.Unlock()
	return ch
}

// PublishEnrolled broadcasts a DeviceEnrolledEvent to all subscribers.
// It is non-blocking: if a subscriber's channel is full the event is dropped
// for that subscriber (and a log entry should be added by the caller).
func (b *Bus) PublishEnrolled(ev DeviceEnrolledEvent) {
	b.mu.RLock()
	subs := b.enrolled
	b.mu.RUnlock()
	for _, ch := range subs {
		select {
		case ch <- ev:
		default:
		}
	}
}

// PublishCheckedOut broadcasts a DeviceCheckedOutEvent to all subscribers.
func (b *Bus) PublishCheckedOut(ev DeviceCheckedOutEvent) {
	b.mu.RLock()
	subs := b.checkedOut
	b.mu.RUnlock()
	for _, ch := range subs {
		select {
		case ch <- ev:
		default:
		}
	}
}

// PublishDeviceInformation broadcasts a DeviceInformationReceivedEvent to all subscribers.
func (b *Bus) PublishDeviceInformation(ev DeviceInformationReceivedEvent) {
	b.mu.RLock()
	subs := b.devInfo
	b.mu.RUnlock()
	for _, ch := range subs {
		select {
		case ch <- ev:
		default:
		}
	}
}

// SubscribeProfileInstallAck registers a buffered channel for ProfileInstallAckEvent.
func (b *Bus) SubscribeProfileInstallAck(bufSize int) <-chan ProfileInstallAckEvent {
	ch := make(chan ProfileInstallAckEvent, bufSize)
	b.mu.Lock()
	b.profileInstall = append(b.profileInstall, ch)
	b.mu.Unlock()
	return ch
}

// PublishProfileInstallAck broadcasts a ProfileInstallAckEvent to all subscribers.
func (b *Bus) PublishProfileInstallAck(ev ProfileInstallAckEvent) {
	b.mu.RLock()
	subs := b.profileInstall
	b.mu.RUnlock()
	for _, ch := range subs {
		select {
		case ch <- ev:
		default:
		}
	}
}

// Close drains and closes all subscriber channels. Call this on application shutdown.
func (b *Bus) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, ch := range b.enrolled {
		close(ch)
	}
	for _, ch := range b.checkedOut {
		close(ch)
	}
	for _, ch := range b.devInfo {
		close(ch)
	}
	for _, ch := range b.profileInstall {
		close(ch)
	}
	b.enrolled = nil
	b.checkedOut = nil
	b.devInfo = nil
	b.profileInstall = nil
}
