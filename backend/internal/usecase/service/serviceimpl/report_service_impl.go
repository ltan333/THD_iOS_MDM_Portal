package serviceimpl

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
)

type reportServiceImpl struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) service.ReportService {
	return &reportServiceImpl{repo: repo}
}

func (s *reportServiceImpl) ExportDevicesCSV(ctx context.Context, opts query.QueryOptions) ([]byte, error) {
	devices, err := s.repo.GetDevicesForExport(ctx, opts)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất dữ liệu thiết bị").WithError(err)
	}

	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	
	// Write Header
	if err := w.Write([]string{
		"ID", "Name", "Serial Number", "Platform", "Model", "OS Version", 
		"Status", "Compliance Status", "Is Enrolled", "MAC Address", "IP Address", "Created At",
	}); err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo file CSV").WithError(err)
	}

	// Write Data
	for _, d := range devices {
		enrolledStr := "No"
		if d.IsEnrolled {
			enrolledStr = "Yes"
		}
		
		record := []string{
			fmt.Sprintf("%v", d.ID),
			d.Name,
			d.SerialNumber,
			string(d.Platform),
			d.Model,
			d.OsVersion,
			string(d.Status),
			string(d.ComplianceStatus),
			enrolledStr,
			d.MACAddress,
			d.IPAddress,
			d.CreatedAt.UTC().Format("2006-01-02 15:04:05"),
		}
		if err := w.Write(record); err != nil {
			return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi ghi dữ liệu dòng CSV").WithError(err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi flush file CSV").WithError(err)
	}

	return b.Bytes(), nil
}

func (s *reportServiceImpl) ExportAlertsCSV(ctx context.Context, opts query.QueryOptions) ([]byte, error) {
	alerts, err := s.repo.GetAlertsForExport(ctx, opts)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất dữ liệu cảnh báo").WithError(err)
	}

	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	
	if err := w.Write([]string{
		"ID", "Severity", "Title", "Type", "Status", "Device ID", "Created At", "Resolved At",
	}); err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo file CSV").WithError(err)
	}

	for _, a := range alerts {
		resolvedAt := ""
		if a.ResolvedAt != nil {
			resolvedAt = a.ResolvedAt.UTC().Format("2006-01-02 15:04:05")
		}
		
		record := []string{
			fmt.Sprintf("%v", a.ID),
			string(a.Severity),
			a.Title,
			string(a.Type),
			string(a.Status),
			a.DeviceID,
			a.CreatedAt.UTC().Format("2006-01-02 15:04:05"),
			resolvedAt,
		}
		if err := w.Write(record); err != nil {
			return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi ghi dữ liệu dòng CSV").WithError(err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi flush file CSV").WithError(err)
	}

	return b.Bytes(), nil
}

func (s *reportServiceImpl) ExportApplicationsCSV(ctx context.Context, opts query.QueryOptions) ([]byte, error) {
	apps, err := s.repo.GetApplicationsForExport(ctx, opts)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất dữ liệu ứng dụng").WithError(err)
	}

	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	
	if err := w.Write([]string{
		"ID", "Name", "Bundle ID", "Platform", "Type", "Created At",
	}); err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo file CSV").WithError(err)
	}

	for _, a := range apps {
		record := []string{
			fmt.Sprintf("%v", a.ID),
			a.Name,
			a.BundleID,
			string(a.Platform),
			string(a.Type),
			a.CreatedAt.UTC().Format("2006-01-02 15:04:05"),
		}
		if err := w.Write(record); err != nil {
			return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi ghi dữ liệu dòng CSV").WithError(err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi flush file CSV").WithError(err)
	}

	return b.Bytes(), nil
}
