// Package worker provides long-running background workers that consume
// events from the event bus and perform asynchronous MDM operations.
package worker

import (
	"context"

	"github.com/thienel/go-backend-template/internal/usecase/service"
	"github.com/thienel/go-backend-template/pkg/event"
	"github.com/thienel/go-backend-template/pkg/mdmcmd"
	"github.com/thienel/tlog"
	"go.uber.org/zap"
)

// InventorySyncWorker listens for DeviceEnrolledEvent and immediately
// enqueues a DeviceInformation MDM command so that device details
// (model, OS version, battery, storage, MAC address) are populated
// automatically after enrollment without requiring manual intervention.
//
// When nanoMDM delivers the device's response, it fires an mdm.Acknowledge
// webhook that the device service handles to apply the data.
type InventorySyncWorker struct {
	mdmService service.NanoMDMService
	eventBus   *event.Bus
	cmdBuilder *mdmcmd.CommandBuilder
}

// NewInventorySyncWorker creates an InventorySyncWorker. Call Start() to begin
// processing events.
func NewInventorySyncWorker(mdmService service.NanoMDMService, eventBus *event.Bus) *InventorySyncWorker {
	return &InventorySyncWorker{
		mdmService: mdmService,
		eventBus:   eventBus,
		cmdBuilder: mdmcmd.NewBuilder(""),
	}
}

// Start subscribes to DeviceEnrolledEvent and runs the sync loop in a
// goroutine. It returns immediately; the goroutine exits when ctx is cancelled.
func (w *InventorySyncWorker) Start(ctx context.Context) {
	enrolled := w.eventBus.SubscribeEnrolled(64)

	go func() {
		tlog.Info("InventorySyncWorker started")
		for {
			select {
			case ev, ok := <-enrolled:
				if !ok {
					return
				}
				w.syncDevice(ctx, ev.DeviceID)

			case <-ctx.Done():
				tlog.Info("InventorySyncWorker stopping")
				return
			}
		}
	}()
}

// syncDevice enqueues a DeviceInformation command and triggers an APNs push
// so the device sends back its current hardware and software details.
func (w *InventorySyncWorker) syncDevice(ctx context.Context, deviceID string) {
	tlog.Info("Requesting DeviceInformation", zap.String("udid", deviceID))

	cmdData, _, err := w.cmdBuilder.DeviceInformation(mdmcmd.CommonDeviceQueries())
	if err != nil {
		tlog.Error("Failed to build DeviceInformation command",
			zap.String("udid", deviceID), zap.Error(err))
		return
	}

	if _, err := w.mdmService.EnqueueCommand(ctx, deviceID, cmdData); err != nil {
		tlog.Error("Failed to enqueue DeviceInformation",
			zap.String("udid", deviceID), zap.Error(err))
		return
	}

	if _, err := w.mdmService.Push(ctx, []string{deviceID}); err != nil {
		tlog.Error("Failed to push after DeviceInformation enqueue",
			zap.String("udid", deviceID), zap.Error(err))
	}
}

// ProfileDeployWorker listens for DeviceEnrolledEvent and triggers automatic
// profile deployment for the newly enrolled device. This replaces the direct
// call to profileService.DeployToDevice() that previously lived inside the
// device service, removing the circular dependency.
type ProfileDeployWorker struct {
	profileService service.ProfileService
	eventBus       *event.Bus
}

// NewProfileDeployWorker creates a ProfileDeployWorker. Call Start() to begin
// processing events.
func NewProfileDeployWorker(profileService service.ProfileService, eventBus *event.Bus) *ProfileDeployWorker {
	return &ProfileDeployWorker{
		profileService: profileService,
		eventBus:       eventBus,
	}
}

// Start subscribes to DeviceEnrolledEvent and runs the deploy loop in a
// goroutine. It returns immediately; the goroutine exits when ctx is cancelled.
func (w *ProfileDeployWorker) Start(ctx context.Context) {
	enrolled := w.eventBus.SubscribeEnrolled(64)

	go func() {
		tlog.Info("ProfileDeployWorker started")
		for {
			select {
			case ev, ok := <-enrolled:
				if !ok {
					return
				}
				w.deploy(ctx, ev.DeviceID)

			case <-ctx.Done():
				tlog.Info("ProfileDeployWorker stopping")
				return
			}
		}
	}()
}

func (w *ProfileDeployWorker) deploy(ctx context.Context, deviceID string) {
	tlog.Info("Auto-deploying profiles after enrollment", zap.String("udid", deviceID))
	if err := w.profileService.DeployToDevice(ctx, deviceID); err != nil {
		tlog.Error("Auto-deploy failed", zap.String("udid", deviceID), zap.Error(err))
	}
}
