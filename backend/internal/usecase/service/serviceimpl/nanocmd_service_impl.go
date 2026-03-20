package serviceimpl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
)

type nanocmdServiceImpl struct {
	baseURL  string
	username string
	password string
}

func NewNanoCMDService(baseURL, username, password string) service.NanoCMDService {
	return &nanocmdServiceImpl{
		baseURL:  strings.TrimSuffix(baseURL, "/"),
		username: username,
		password: password,
	}
}

func (s *nanocmdServiceImpl) doRequest(ctx context.Context, method, path string, body interface{}, query url.Values) (*http.Response, error) {
	u, err := url.Parse(fmt.Sprintf("%s%s", s.baseURL, path))
	if err != nil {
		return nil, err
	}
	if query != nil {
		u.RawQuery = query.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(s.username, s.password)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	return client.Do(req)
}

func (s *nanocmdServiceImpl) GetVersion(ctx context.Context) (*dto.NanoCMDVersionResponse, error) {
	resp, err := s.doRequest(ctx, http.MethodGet, "/version", nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nanocmd error: status %d", resp.StatusCode)
	}

	var result dto.NanoCMDVersionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *nanocmdServiceImpl) StartWorkflow(ctx context.Context, name string, enrollmentIDs []string, contextStr string) (*dto.NanoCMDWorkflowStartResponse, error) {
	query := url.Values{}
	for _, id := range enrollmentIDs {
		query.Add("id", id)
	}
	if contextStr != "" {
		query.Set("context", contextStr)
	}

	resp, err := s.doRequest(ctx, http.MethodPost, fmt.Sprintf("/v1/workflow/%s/start", name), nil, query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nanocmd error: status %d", resp.StatusCode)
	}

	var result dto.NanoCMDWorkflowStartResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *nanocmdServiceImpl) GetEvent(ctx context.Context, name string) (*dto.EventSubscription, error) {
	resp, err := s.doRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/event/%s", name), nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nanocmd error: status %d", resp.StatusCode)
	}

	var result dto.EventSubscription
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *nanocmdServiceImpl) PutEvent(ctx context.Context, name string, subscription *dto.EventSubscription) error {
	resp, err := s.doRequest(ctx, http.MethodPut, fmt.Sprintf("/v1/event/%s", name), subscription, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("nanocmd error: status %d", resp.StatusCode)
	}
	return nil
}

func (s *nanocmdServiceImpl) GetFVEnableProfileTemplate(ctx context.Context) ([]byte, error) {
	resp, err := s.doRequest(ctx, http.MethodGet, "/v1/fvenable/profiletemplate", nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nanocmd error: status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func (s *nanocmdServiceImpl) GetProfile(ctx context.Context, name string) ([]byte, error) {
	resp, err := s.doRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/profile/%s", name), nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nanocmd error: status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func (s *nanocmdServiceImpl) PutProfile(ctx context.Context, name string, profileData []byte) error {
	u, err := url.Parse(fmt.Sprintf("%s/v1/profile/%s", s.baseURL, name))
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), bytes.NewReader(profileData))
	if err != nil {
		return err
	}

	req.SetBasicAuth(s.username, s.password)
	req.Header.Set("Content-Type", "application/x-apple-aspen-config")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("nanocmd error: status %d", resp.StatusCode)
	}
	return nil
}

func (s *nanocmdServiceImpl) DeleteProfile(ctx context.Context, name string) error {
	resp, err := s.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/v1/profile/%s", name), nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("nanocmd error: status %d", resp.StatusCode)
	}
	return nil
}

func (s *nanocmdServiceImpl) GetProfiles(ctx context.Context, names []string) (map[string]dto.NanoCMDProfile, error) {
	query := url.Values{}
	for _, name := range names {
		query.Add("name", name)
	}

	resp, err := s.doRequest(ctx, http.MethodGet, "/v1/profiles", nil, query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nanocmd error: status %d", resp.StatusCode)
	}

	var result map[string]dto.NanoCMDProfile
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nanocmdServiceImpl) GetCMDPlan(ctx context.Context, name string) (*dto.CMDPlan, error) {
	resp, err := s.doRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/cmdplan/%s", name), nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nanocmd error: status %d", resp.StatusCode)
	}

	var result dto.CMDPlan
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *nanocmdServiceImpl) PutCMDPlan(ctx context.Context, name string, plan *dto.CMDPlan) error {
	resp, err := s.doRequest(ctx, http.MethodPut, fmt.Sprintf("/v1/cmdplan/%s", name), plan, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("nanocmd error: status %d", resp.StatusCode)
	}
	return nil
}

func (s *nanocmdServiceImpl) GetInventory(ctx context.Context, enrollmentIDs []string) (dto.NanoCMDInventoryResponse, error) {
	query := url.Values{}
	for _, id := range enrollmentIDs {
		query.Add("id", id)
	}

	resp, err := s.doRequest(ctx, http.MethodGet, "/v1/inventory", nil, query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nanocmd error: status %d", resp.StatusCode)
	}

	var result dto.NanoCMDInventoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}
