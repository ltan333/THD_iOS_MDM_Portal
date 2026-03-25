package serviceimpl

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/alert"
	"github.com/thienel/go-backend-template/internal/ent/application"
	"github.com/thienel/go-backend-template/internal/ent/device"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
)

type reportServiceImpl struct {
	client *ent.Client
}

func NewReportService(client *ent.Client) service.ReportService {
	return &reportServiceImpl{client: client}
}

func (s *reportServiceImpl) ExportDevicesCSV(ctx context.Context, opts query.QueryOptions) ([]byte, error) {
	q := s.client.Device.Query()

	// Apply search
	if searchFilter, ok := opts.Filters["search"]; ok {
		searchStr, _ := searchFilter.Value.(string)
		if searchStr != "" {
			q = q.Where(
				device.Or(
					device.NameContainsFold(searchStr),
					device.SerialNumberContainsFold(searchStr),
					device.ModelContainsFold(searchStr),
				),
			)
		}
	}

	devices, err := q.Order(ent.Desc(device.FieldCreatedAt)).All(ctx)
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
			d.CreatedAt.Format("2006-01-02 15:04:05"),
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
	q := s.client.Alert.Query()

	if searchFilter, ok := opts.Filters["search"]; ok {
		searchStr, _ := searchFilter.Value.(string)
		if searchStr != "" {
			q = q.Where(alert.TitleContainsFold(searchStr))
		}
	}

	alerts, err := q.Order(ent.Desc(alert.FieldCreatedAt)).All(ctx)
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
			resolvedAt = a.ResolvedAt.Format("2006-01-02 15:04:05")
		}
		
		record := []string{
			fmt.Sprintf("%v", a.ID),
			string(a.Severity),
			a.Title,
			string(a.Type),
			string(a.Status),
			a.DeviceID,
			a.CreatedAt.Format("2006-01-02 15:04:05"),
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
	q := s.client.Application.Query()

	if searchFilter, ok := opts.Filters["search"]; ok {
		searchStr, _ := searchFilter.Value.(string)
		if searchStr != "" {
			q = q.Where(
				application.Or(
					application.NameContainsFold(searchStr),
					application.BundleIDContainsFold(searchStr),
				),
			)
		}
	}

	apps, err := q.Order(ent.Desc(application.FieldCreatedAt)).All(ctx)
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
			a.CreatedAt.Format("2006-01-02 15:04:05"),
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
