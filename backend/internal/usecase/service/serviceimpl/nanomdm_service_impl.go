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
	baseURL  string
	username string
	password string
}

func NewNanoMDMService(baseURL, username, password string) service.NanoMDMService {
	return &nanomdmServiceImpl{
		baseURL:  strings.TrimSuffix(baseURL, "/"),
		username: username,
		password: password,
	}
}

func (s *nanomdmServiceImpl) doRequest(ctx context.Context, method, path string, body interface{}, query url.Values) (*http.Response, error) {
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

func (s *nanomdmServiceImpl) DefineDEPProfile(ctx context.Context, depName string, profile interface{}) (string, error) {
	resp, err := s.doRequest(ctx, http.MethodPost, "/v1/dep/profiles", profile, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

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
	resp, err := s.doRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/dep/profiles/%s", profileUUID), nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
	resp, err := s.doRequest(ctx, http.MethodPost, "/v1/pushcert", certData, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("nanomdm error: status %d", resp.StatusCode)
	}
	return nil
}

func (s *nanomdmServiceImpl) GetPushCert(ctx context.Context) (interface{}, error) {
	resp, err := s.doRequest(ctx, http.MethodGet, "/v1/pushcert", nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nanomdm error: status %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}
