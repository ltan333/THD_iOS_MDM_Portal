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

	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/httpclient"
)

type nanomdmServiceImpl struct {
	client      *http.Client
	mdmBaseURL  string
	depBaseURL  string
	mdmUsername string
	mdmPassword string
	depUsername string
	depPassword string
}

func NewNanoMDMService(mdmBaseURL, depBaseURL, mdmUser, mdmPass, depUser, depPass string) service.NanoMDMService {
	return &nanomdmServiceImpl{
		client:      httpclient.DefaultClient(),
		mdmBaseURL:  strings.TrimSuffix(mdmBaseURL, "/"),
		depBaseURL:  strings.TrimSuffix(depBaseURL, "/"),
		mdmUsername: mdmUser,
		mdmPassword: mdmPass,
		depUsername: depUser,
		depPassword: depPass,
	}
}

func (s *nanomdmServiceImpl) doRequest(ctx context.Context, method, baseURL, path string, body interface{}, query url.Values, username, password string) (*http.Response, error) {
	u, err := url.Parse(fmt.Sprintf("%s%s", baseURL, path))
	if err != nil {
		return nil, err
	}
	if query != nil {
		u.RawQuery = query.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		if b, ok := body.([]byte); ok {
			bodyReader = bytes.NewReader(b)
		} else {
			b, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			bodyReader = bytes.NewReader(b)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(username, password)
	if body != nil {
		if _, ok := body.([]byte); !ok {
			req.Header.Set("Content-Type", "application/json")
		}
	}

	return s.client.Do(req)
}

func (s *nanomdmServiceImpl) handleResponse(resp *http.Response, target interface{}) error {
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if target == nil {
			return nil
		}
		return json.NewDecoder(resp.Body).Decode(target)
	}

	// Capture response body for better error logging
	body, _ := io.ReadAll(resp.Body)
	errMsg := fmt.Sprintf("nanomdm error: status %d, body: %s", resp.StatusCode, string(body))

	if resp.StatusCode == http.StatusNotFound {
		return apperror.ErrNotFound.WithMessage(errMsg)
	}
	if resp.StatusCode == http.StatusBadRequest {
		return apperror.ErrBadRequest.WithMessage(errMsg)
	}

	return fmt.Errorf("%s", errMsg)
}

func (s *nanomdmServiceImpl) DefineDEPProfile(ctx context.Context, depName string, profile interface{}) (string, error) {
	resp, err := s.doRequest(ctx, http.MethodPost, s.depBaseURL, "/v1/dep/profiles", profile, nil, s.depUsername, s.depPassword)
	if err != nil {
		return "", err
	}

	var result struct {
		ProfileUUID string `json:"profile_uuid"`
	}
	if err := s.handleResponse(resp, &result); err != nil {
		return "", err
	}
	return result.ProfileUUID, nil
}

func (s *nanomdmServiceImpl) GetDEPProfile(ctx context.Context, depName, profileUUID string) (interface{}, error) {
	resp, err := s.doRequest(ctx, http.MethodGet, s.depBaseURL, fmt.Sprintf("/v1/dep/profiles/%s", profileUUID), nil, nil, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nanomdmServiceImpl) ListDEPProfiles(ctx context.Context, depName string) (interface{}, error) {
	resp, err := s.doRequest(ctx, http.MethodGet, s.depBaseURL, "/v1/dep/profiles", nil, nil, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nanomdmServiceImpl) SyncDEPDevices(ctx context.Context, depName string) (interface{}, error) {
	resp, err := s.doRequest(ctx, http.MethodPost, s.depBaseURL, fmt.Sprintf("/proxy/%s/devices/sync", depName), nil, nil, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nanomdmServiceImpl) DisownDEPDevices(ctx context.Context, depName string, devices []string) (interface{}, error) {
	body := struct {
		Devices []string `json:"devices"`
	}{
		Devices: devices,
	}
	resp, err := s.doRequest(ctx, http.MethodPost, s.depBaseURL, fmt.Sprintf("/proxy/%s/devices/disown", depName), body, nil, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nanomdmServiceImpl) UploadDEPToken(ctx context.Context, depName string, tokenData []byte) (interface{}, error) {
	resp, err := s.doRequest(ctx, http.MethodPut, s.depBaseURL, fmt.Sprintf("/v1/tokenpki/%s", depName), tokenData, nil, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	// For token upload, NanoDEP might set Content-Type header manually in doRequest,
	// but I updated doRequest to handle []byte body correctly.

	var result interface{}
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nanomdmServiceImpl) UploadPushCert(ctx context.Context, certData []byte) error {
	resp, err := s.doRequest(ctx, http.MethodPost, s.mdmBaseURL, "/v1/pushcert", certData, nil, s.mdmUsername, s.mdmPassword)
	if err != nil {
		return err
	}
	return s.handleResponse(resp, nil)
}

func (s *nanomdmServiceImpl) GetPushCert(ctx context.Context) (interface{}, error) {
	resp, err := s.doRequest(ctx, http.MethodGet, s.mdmBaseURL, "/v1/pushcert", nil, nil, s.mdmUsername, s.mdmPassword)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nanomdmServiceImpl) ListDEPNames(ctx context.Context) (interface{}, error) {
	resp, err := s.doRequest(ctx, http.MethodGet, s.depBaseURL, "/v1/dep_names", nil, nil, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nanomdmServiceImpl) GetDEPConfig(ctx context.Context, depName string) (interface{}, error) {
	resp, err := s.doRequest(ctx, http.MethodGet, s.depBaseURL, fmt.Sprintf("/v1/config/%s", depName), nil, nil, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nanomdmServiceImpl) GetDEPAssigner(ctx context.Context, depName string) (interface{}, error) {
	resp, err := s.doRequest(ctx, http.MethodGet, s.depBaseURL, fmt.Sprintf("/v1/assigner/%s", depName), nil, nil, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nanomdmServiceImpl) SetDEPAssigner(ctx context.Context, depName string, assigner interface{}) (interface{}, error) {
	resp, err := s.doRequest(ctx, http.MethodPut, s.depBaseURL, fmt.Sprintf("/v1/assigner/%s", depName), assigner, nil, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nanomdmServiceImpl) GetDEPAccount(ctx context.Context, depName string) (interface{}, error) {
	resp, err := s.doRequest(ctx, http.MethodGet, s.depBaseURL, fmt.Sprintf("/proxy/%s/account", depName), nil, nil, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nanomdmServiceImpl) GetDEPDevices(ctx context.Context, depName string, cursor string) (interface{}, error) {
	var query url.Values
	if cursor != "" {
		query = url.Values{}
		query.Set("cursor", cursor)
	}
	resp, err := s.doRequest(ctx, http.MethodPost, s.depBaseURL, fmt.Sprintf("/proxy/%s/devices", depName), nil, query, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nanomdmServiceImpl) GetDEPTokens(ctx context.Context, depName string) (interface{}, error) {
	resp, err := s.doRequest(ctx, http.MethodGet, s.depBaseURL, fmt.Sprintf("/v1/tokens/%s", depName), nil, nil, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nanomdmServiceImpl) EnqueueCommand(ctx context.Context, udid string, cmdData []byte) (interface{}, error) {
	resp, err := s.doRequest(ctx, http.MethodPut, s.mdmBaseURL, fmt.Sprintf("/v1/enqueue/%s", udid), cmdData, nil, s.mdmUsername, s.mdmPassword)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}
