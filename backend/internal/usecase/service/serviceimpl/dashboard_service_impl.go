package serviceimpl

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/device"
	"github.com/thienel/go-backend-template/internal/ent/user"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
)

type dashboardServiceImpl struct {
	client *ent.Client
}

func NewDashboardService(client *ent.Client) service.DashboardService {
	return &dashboardServiceImpl{client: client}
}

func (s *dashboardServiceImpl) GetStats(ctx context.Context) (*dto.DashboardStatsResponse, error) {
	// Count total devices
	totalDevices, err := s.client.Device.Query().Count(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi đếm thiết bị").WithError(err)
	}

	// Count active devices (enrolled)
	activeDevices, err := s.client.Device.Query().
		Where(device.IsEnrolledEQ(true)).
		Count(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi đếm thiết bị hoạt động").WithError(err)
	}

	// Count total users
	totalUsers, err := s.client.User.Query().
		Where(user.DeletedAtIsNil()).
		Count(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi đếm người dùng").WithError(err)
	}

	// Count active users
	activeUsers, err := s.client.User.Query().
		Where(user.DeletedAtIsNil(), user.StatusEQ("ACTIVE")).
		Count(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi đếm người dùng hoạt động").WithError(err)
	}

	// Compliance rate calculation (placeholder - will be updated when compliance entities exist)
	complianceRate := 0
	nonCompliantRate := 0
	if totalDevices > 0 {
		complianceRate = int((float64(activeDevices) / float64(totalDevices)) * 100)
		nonCompliantRate = 100 - complianceRate
	}

	return &dto.DashboardStatsResponse{
		TotalDevices:     int64(totalDevices),
		ActiveDevices:    int64(activeDevices),
		TotalUsers:       int64(totalUsers),
		ActiveUsers:      int64(activeUsers),
		TotalAlerts:      0, // Will be updated when Alert entity exists
		PendingAlerts:    0, // Will be updated when Alert entity exists
		TotalApps:        0, // Will be updated when Application entity exists
		DeployedApps:     0, // Will be updated when Application entity exists
		ComplianceRate:   complianceRate,
		NonCompliantRate: nonCompliantRate,
	}, nil
}

func (s *dashboardServiceImpl) GetDeviceStats(ctx context.Context) (*dto.DeviceStatsResponse, error) {
	// Count total devices
	total, err := s.client.Device.Query().Count(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi đếm thiết bị").WithError(err)
	}

	// Count enrolled devices
	enrolled, err := s.client.Device.Query().
		Where(device.IsEnrolledEQ(true)).
		Count(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi đếm thiết bị đã đăng ký").WithError(err)
	}

	// Initialize platform counts
	byPlatform := map[string]int64{
		"ios":     0,
		"android": 0,
		"windows": 0,
		"macos":   0,
	}
	var platformStats []struct {
		Platform string `json:"platform"`
		Count    int    `json:"count"`
	}
	if err := s.client.Device.Query().
		GroupBy(device.FieldPlatform).
		Aggregate(ent.Count()).
		Scan(ctx, &platformStats); err == nil {
		for _, stat := range platformStats {
			if stat.Platform != "" {
				byPlatform[stat.Platform] = int64(stat.Count)
			}
		}
	}

	// Initialize status counts
	byStatus := map[string]int64{
		"active":   int64(enrolled),
		"inactive": int64(total - enrolled),
		"pending":  0,
	}

	return &dto.DeviceStatsResponse{
		Total:        int64(total),
		Active:       int64(enrolled),
		Inactive:     int64(total - enrolled),
		Enrolled:     int64(enrolled),
		ByPlatform:   byPlatform,
		ByStatus:     byStatus,
		Compliant:    int64(enrolled), // Placeholder
		NonCompliant: 0,               // Placeholder
	}, nil
}

func (s *dashboardServiceImpl) GetAlertsSummary(ctx context.Context) (*dto.AlertsSummaryResponse, error) {
	// Placeholder implementation - will be updated when Alert entity exists
	return &dto.AlertsSummaryResponse{
		Total:        0,
		Open:         0,
		Acknowledged: 0,
		Resolved:     0,
		BySeverity: map[string]int64{
			"critical": 0,
			"high":     0,
			"medium":   0,
			"low":      0,
		},
		ByType: map[string]int64{
			"security":     0,
			"compliance":   0,
			"connectivity": 0,
			"application":  0,
			"device_health": 0,
		},
	}, nil
}

func (s *dashboardServiceImpl) GetChartData(ctx context.Context, chartType string) (*dto.ChartDataResponse, error) {
	switch chartType {
	case "devices":
		return s.getDevicesChartData(ctx)
	case "compliance":
		return s.getComplianceChartData(ctx)
	case "alerts":
		return s.getAlertsChartData(ctx)
	default:
		return nil, apperror.ErrBadRequest.WithMessage("Loại biểu đồ không hợp lệ: " + chartType)
	}
}

func (s *dashboardServiceImpl) getDevicesChartData(ctx context.Context) (*dto.ChartDataResponse, error) {
	total, _ := s.client.Device.Query().Count(ctx)
	enrolled, _ := s.client.Device.Query().Where(device.IsEnrolledEQ(true)).Count(ctx)
	notEnrolled := total - enrolled

	return &dto.ChartDataResponse{
		ChartType: "doughnut",
		Title:     "Device Status",
		Labels:    []string{"Enrolled", "Not Enrolled"},
		Datasets: []dto.ChartDataset{
			{
				Label: "Devices",
				Data:  []interface{}{enrolled, notEnrolled},
				Color: "#4CAF50,#FF5722",
			},
		},
	}, nil
}

func (s *dashboardServiceImpl) getComplianceChartData(ctx context.Context) (*dto.ChartDataResponse, error) {
	// Placeholder - will be updated when compliance data is available
	return &dto.ChartDataResponse{
		ChartType: "pie",
		Title:     "Compliance Status",
		Labels:    []string{"Compliant", "Non-Compliant"},
		Datasets: []dto.ChartDataset{
			{
				Label: "Compliance",
				Data:  []interface{}{100, 0},
				Color: "#4CAF50,#F44336",
			},
		},
	}, nil
}

func (s *dashboardServiceImpl) getAlertsChartData(ctx context.Context) (*dto.ChartDataResponse, error) {
	// Placeholder - will be updated when Alert entity exists
	return &dto.ChartDataResponse{
		ChartType: "bar",
		Title:     "Alerts by Severity",
		Labels:    []string{"Critical", "High", "Medium", "Low"},
		Datasets: []dto.ChartDataset{
			{
				Label: "Alerts",
				Data:  []interface{}{0, 0, 0, 0},
				Color: "#F44336,#FF9800,#FFC107,#4CAF50",
			},
		},
	}, nil
}
