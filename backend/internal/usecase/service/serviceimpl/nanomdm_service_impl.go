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
)

type nanomdmServiceImpl struct {
	mdmBaseURL  string
	depBaseURL  string
	mdmUsername string
	mdmPassword string
	depUsername string
	depPassword string
}

func NewNanoMDMService(mdmBaseURL, depBaseURL, mdmUser, mdmPass, depUser, depPass string) service.NanoMDMService {
	return &nanomdmServiceImpl{
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

	req.SetBasicAuth(username, password)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	return client.Do(req)
}

func (s *nanomdmServiceImpl) DefineDEPProfile(ctx context.Context, depName string, profile interface{}) (string, error) {
	resp, err := s.doRequest(ctx, http.MethodPost, s.depBaseURL, "/v1/dep/profiles", profile, nil, s.depUsername, s.depPassword)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("nanomdm error: status %d", resp.StatusCode)
	}

	var result struct {
		ProfileUUID string `json:"profile_uuid"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.ProfileUUID, nil
}

func (s *nanomdmServiceImpl) GetDEPProfile(ctx context.Context, depName, profileUUID string) (interface{}, error) {
	resp, err := s.doRequest(ctx, http.MethodGet, s.depBaseURL, fmt.Sprintf("/v1/dep/profiles/%s", profileUUID), nil, nil, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nanomdm error: status %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nanomdmServiceImpl) ListDEPProfiles(ctx context.Context, depName string) (interface{}, error) {
	resp, err := s.doRequest(ctx, http.MethodGet, s.depBaseURL, "/v1/dep/profiles", nil, nil, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nanomdm error: status %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nanomdmServiceImpl) SyncDEPDevices(ctx context.Context, depName string) (interface{}, error) {
	resp, err := s.doRequest(ctx, http.MethodPost, s.depBaseURL, fmt.Sprintf("/proxy/%s/devices/sync", depName), nil, nil, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nanomdm error: status %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
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
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nanomdm error: status %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nanomdmServiceImpl) UploadDEPToken(ctx context.Context, depName string, tokenData []byte) (interface{}, error) {
	u, err := url.Parse(fmt.Sprintf("%s/v1/tokenpki/%s", s.depBaseURL, depName))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), bytes.NewReader(tokenData))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(s.depUsername, s.depPassword)
	req.Header.Set("Content-Type", "application/x-apple-aspen-config")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nanomdm error: status %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nanomdmServiceImpl) UploadPushCert(ctx context.Context, certData []byte) error {
	resp, err := s.doRequest(ctx, http.MethodPost, s.mdmBaseURL, "/v1/pushcert", certData, nil, s.mdmUsername, s.mdmPassword)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("nanomdm error: status %d", resp.StatusCode)
	}
	return nil
}

func (s *nanomdmServiceImpl) GetPushCert(ctx context.Context) (interface{}, error) {
	resp, err := s.doRequest(ctx, http.MethodGet, s.mdmBaseURL, "/v1/pushcert", nil, nil, s.mdmUsername, s.mdmPassword)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nanomdm error: status %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}
