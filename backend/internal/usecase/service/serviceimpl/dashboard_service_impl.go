package serviceimpl

import (
	"context"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
)

type dashboardServiceImpl struct {
	repo repository.DashboardRepository
}

func NewDashboardService(repo repository.DashboardRepository) service.DashboardService {
	return &dashboardServiceImpl{repo: repo}
}

func (s *dashboardServiceImpl) GetStats(ctx context.Context) (*dto.DashboardStatsResponse, error) {
	totalDevices, err := s.repo.CountDevices(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi đếm thiết bị").WithError(err)
	}

	activeDevices, err := s.repo.CountEnrolledDevices(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi đếm thiết bị hoạt động").WithError(err)
	}

	totalUsers, err := s.repo.CountUsers(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi đếm người dùng").WithError(err)
	}

	activeUsers, err := s.repo.CountActiveUsers(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi đếm người dùng hoạt động").WithError(err)
	}

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
	total, err := s.repo.CountDevices(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi đếm thiết bị").WithError(err)
	}

	enrolled, err := s.repo.CountEnrolledDevices(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi đếm thiết bị đã đăng ký").WithError(err)
	}

	byPlatform := map[string]int64{
		"ios":     0,
		"android": 0,
		"windows": 0,
		"macos":   0,
	}
	
	platformStats, err := s.repo.GetDevicePlatformCounts(ctx)
	if err == nil {
		for platform, count := range platformStats {
			byPlatform[platform] = int64(count)
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
			"security":      0,
			"compliance":    0,
			"connectivity":  0,
			"application":   0,
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
	total, _ := s.repo.CountDevices(ctx)
	enrolled, _ := s.repo.CountEnrolledDevices(ctx)
	notEnrolled := total - enrolled

	return &dto.ChartDataResponse{
		ChartType: "doughnut",
		Title:     "Device Status",
		Labels:    []string{"Enrolled", "Not Enrolled"},
		Datasets: []dto.ChartDataset{
			{
				Label: "Devices",
				Data:  []any{enrolled, notEnrolled},
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
				Data:  []any{100, 0},
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
				Data:  []any{0, 0, 0, 0},
				Color: "#F44336,#FF9800,#FFC107,#4CAF50",
			},
		},
	}, nil
}
