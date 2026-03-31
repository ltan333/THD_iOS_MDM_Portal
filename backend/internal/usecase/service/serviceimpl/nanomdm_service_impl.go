package serviceimpl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"strings"

	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/httpclient"
)

type nanomdmServiceImpl struct {
	client             *http.Client
	mdmBaseURL         string
	depBaseURL         string
	mdmUsername        string
	mdmPassword        string
	depUsername        string
	depPassword        string
	sudoPassword       string
	depSyncerContainer string
}

func NewNanoMDMService(mdmBaseURL, depBaseURL, mdmUser, mdmPass, depUser, depPass, sudoPassword, depSyncerContainer string) service.NanoMDMService {
	return &nanomdmServiceImpl{
		client:             httpclient.DefaultClient(),
		mdmBaseURL:         strings.TrimSuffix(mdmBaseURL, "/"),
		depBaseURL:         strings.TrimSuffix(depBaseURL, "/"),
		mdmUsername:        mdmUser,
		mdmPassword:        mdmPass,
		depUsername:        depUser,
		depPassword:        depPass,
		sudoPassword:       sudoPassword,
		depSyncerContainer: depSyncerContainer,
	}
}

func (s *nanomdmServiceImpl) doRequest(ctx context.Context, method, baseURL, path string, body any, query url.Values, username, password string) (*http.Response, error) {
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
	} else if method == http.MethodPost || method == http.MethodPut {
		bodyReader = bytes.NewReader([]byte("{}"))
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(username, password)
	req.Header.Set("User-Agent", "MDM-Portal/1.0")
	if body != nil {
		if _, ok := body.([]byte); ok {
			// Raw plist commands sent to NanoMDM /v1/enqueue must be text/plain per spec
			req.Header.Set("Content-Type", "text/plain")
		} else {
			req.Header.Set("Content-Type", "application/json")
		}
	} else if method == http.MethodPost || method == http.MethodPut {
		req.Header.Set("Content-Type", "application/json")
	}

	return s.client.Do(req)
}

func (s *nanomdmServiceImpl) handleResponse(resp *http.Response, target any) error {
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
	if resp.StatusCode >= 500 {
		return apperror.ErrInternalServerError.WithMessage(errMsg)
	}

	return apperror.ErrBadRequest.WithMessage(errMsg)
}

func (s *nanomdmServiceImpl) GetDEPProfile(ctx context.Context, depName, profileUUID string) (any, error) {
	query := url.Values{}
	query.Set("profile_uuid", profileUUID)
	resp, err := s.doRequest(ctx, http.MethodGet, s.depBaseURL, fmt.Sprintf("/proxy/%s/profile", depName), nil, query, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result dto.DEPProfileResponse
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *nanomdmServiceImpl) SyncDEPDevices(ctx context.Context, depName string, cursor string) (any, error) {
	var body any
	if cursor != "" {
		body = map[string]string{"cursor": cursor}
	}
	resp, err := s.doRequest(ctx, http.MethodPost, s.depBaseURL, fmt.Sprintf("/proxy/%s/devices/sync", depName), body, nil, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result any
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nanomdmServiceImpl) DisownDEPDevices(ctx context.Context, depName string, devices []string) (any, error) {
	body := struct {
		Devices []string `json:"devices"`
	}{
		Devices: devices,
	}
	resp, err := s.doRequest(ctx, http.MethodPost, s.depBaseURL, fmt.Sprintf("/proxy/%s/devices/disown", depName), body, nil, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result any
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nanomdmServiceImpl) UploadDEPToken(ctx context.Context, depName string, tokenData []byte) (any, error) {
	// Custom request to set proper Content-Type for PKCS7 data
	u, err := url.Parse(fmt.Sprintf("%s/v1/tokenpki/%s", s.depBaseURL, depName))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), bytes.NewReader(tokenData))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(s.depUsername, s.depPassword)
	req.Header.Set("User-Agent", "MDM-Portal/1.0")
	req.Header.Set("Content-Type", "application/pkcs7-mime")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	var result any
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nanomdmServiceImpl) UploadPushCert(ctx context.Context, certData []byte) (*dto.PushCertResponse, error) {
	// NanoMDM PUT /v1/pushcert expects text/plain concatenated PEM
	resp, err := s.doRequest(ctx, http.MethodPut, s.mdmBaseURL, "/v1/pushcert", certData, nil, s.mdmUsername, s.mdmPassword)
	if err != nil {
		return nil, err
	}

	var result dto.PushCertResponse
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *nanomdmServiceImpl) GetPushCert(ctx context.Context, topic string) (*dto.PushCertResponse, error) {
	query := url.Values{}
	query.Set("topic", topic)
	resp, err := s.doRequest(ctx, http.MethodGet, s.mdmBaseURL, "/v1/pushcert", nil, query, s.mdmUsername, s.mdmPassword)
	if err != nil {
		return nil, err
	}

	var result dto.PushCertResponse
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *nanomdmServiceImpl) ListDEPNames(ctx context.Context, depNames []string, limit, offset int, cursor string) (*dto.DEPNamesQueryResponse, error) {
	query := url.Values{}
	for _, name := range depNames {
		query.Add("dep_name", name)
	}
	if limit > 0 {
		query.Set("limit", fmt.Sprintf("%d", limit))
	}
	if offset > 0 {
		query.Set("offset", fmt.Sprintf("%d", offset))
	}
	if cursor != "" {
		query.Set("cursor", cursor)
	}

	resp, err := s.doRequest(ctx, http.MethodGet, s.depBaseURL, "/v1/dep_names", nil, query, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result dto.DEPNamesQueryResponse
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *nanomdmServiceImpl) GetDEPVersion(ctx context.Context) (*dto.NanoDEPVersionResponse, error) {
	resp, err := s.doRequest(ctx, http.MethodGet, s.depBaseURL, "/version", nil, nil, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result dto.NanoDEPVersionResponse
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *nanomdmServiceImpl) GetDEPConfig(ctx context.Context, depName string) (*dto.DEPConfig, error) {
	resp, err := s.doRequest(ctx, http.MethodGet, s.depBaseURL, fmt.Sprintf("/v1/config/%s", depName), nil, nil, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result dto.DEPConfig
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *nanomdmServiceImpl) SetDEPConfig(ctx context.Context, depName string, config *dto.DEPConfig) (*dto.DEPConfig, error) {
	resp, err := s.doRequest(ctx, http.MethodPut, s.depBaseURL, fmt.Sprintf("/v1/config/%s", depName), config, nil, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result dto.DEPConfig
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *nanomdmServiceImpl) GetDEPAssigner(ctx context.Context, depName string) (*dto.AssignerProfileUUID, error) {
	resp, err := s.doRequest(ctx, http.MethodGet, s.depBaseURL, fmt.Sprintf("/v1/assigner/%s", depName), nil, nil, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result dto.AssignerProfileUUID
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *nanomdmServiceImpl) SetDEPAssigner(ctx context.Context, depName string, profileUUID string) (*dto.AssignerProfileUUID, error) {
	query := url.Values{}
	query.Set("profile_uuid", profileUUID)

	resp, err := s.doRequest(ctx, http.MethodPut, s.depBaseURL, fmt.Sprintf("/v1/assigner/%s", depName), nil, query, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result dto.AssignerProfileUUID
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *nanomdmServiceImpl) GetDEPAccount(ctx context.Context, depName string) (any, error) {
	resp, err := s.doRequest(ctx, http.MethodGet, s.depBaseURL, fmt.Sprintf("/proxy/%s/account", depName), nil, nil, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result any
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nanomdmServiceImpl) GetDEPDevices(ctx context.Context, depName string, devices []string, cursor string) (any, error) {
	var query url.Values
	if cursor != "" {
		query = url.Values{}
		query.Set("cursor", cursor)
	}

	method := http.MethodPost
	var body any
	if len(devices) > 0 {
		body = map[string][]string{"devices": devices}
	} else {
		// If no devices, use GET to list all devices from Apple
		method = http.MethodGet
	}

	resp, err := s.doRequest(ctx, method, s.depBaseURL, fmt.Sprintf("/proxy/%s/devices", depName), body, query, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result any
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nanomdmServiceImpl) GetDEPTokens(ctx context.Context, depName string) (*dto.OAuth1Tokens, error) {
	resp, err := s.doRequest(ctx, http.MethodGet, s.depBaseURL, fmt.Sprintf("/v1/tokens/%s", depName), nil, nil, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result dto.OAuth1Tokens
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *nanomdmServiceImpl) UpdateDEPTokens(ctx context.Context, depName string, tokens *dto.OAuth1Tokens) (*dto.OAuth1Tokens, error) {
	resp, err := s.doRequest(ctx, http.MethodPut, s.depBaseURL, fmt.Sprintf("/v1/tokens/%s", depName), tokens, nil, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result dto.OAuth1Tokens
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *nanomdmServiceImpl) GetDEPTokenPKI(ctx context.Context, depName string, cn string, validityDays int) ([]byte, string, error) {
	query := url.Values{}
	if cn != "" {
		query.Set("cn", cn)
	}
	if validityDays > 0 {
		query.Set("validity_days", fmt.Sprintf("%d", validityDays))
	}

	resp, err := s.doRequest(ctx, http.MethodGet, s.depBaseURL, fmt.Sprintf("/v1/tokenpki/%s", depName), nil, query, s.depUsername, s.depPassword)
	if err != nil {
		return nil, "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("nanodep error: status %d, body: %s", resp.StatusCode, string(body))
	}

	contentDisp := resp.Header.Get("Content-Disposition")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return body, contentDisp, nil
}

func (s *nanomdmServiceImpl) GetMAIDJWT(ctx context.Context, depName string, serverUUID string) (string, string, string, error) {
	query := url.Values{}
	if serverUUID != "" {
		query.Set("server_uuid", serverUUID)
	}

	resp, err := s.doRequest(ctx, http.MethodGet, s.depBaseURL, fmt.Sprintf("/v1/maidjwt/%s", depName), nil, query, s.depUsername, s.depPassword)
	if err != nil {
		return "", "", "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", "", fmt.Errorf("nanodep error: status %d, body: %s", resp.StatusCode, string(body))
	}

	jwt, _ := io.ReadAll(resp.Body)
	serverUuidHeader := resp.Header.Get("X-Server-Uuid")
	jtiHeader := resp.Header.Get("X-Jwt-Jti")

	return string(jwt), serverUuidHeader, jtiHeader, nil
}

func (s *nanomdmServiceImpl) GetBypassCode(ctx context.Context, code, raw string) (*dto.BypassCodeResponse, error) {
	query := url.Values{}
	if code != "" {
		query.Set("code", code)
	}
	if raw != "" {
		query.Set("raw", raw)
	}

	resp, err := s.doRequest(ctx, http.MethodGet, s.depBaseURL, "/v1/bypasscode", nil, query, s.depUsername, s.depPassword)
	if err != nil {
		return nil, err
	}

	var result dto.BypassCodeResponse
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *nanomdmServiceImpl) EnqueueCommand(ctx context.Context, udid string, cmdData []byte) (*dto.APIResult, error) {
	resp, err := s.doRequest(ctx, http.MethodPut, s.mdmBaseURL, fmt.Sprintf("/v1/enqueue/%s", udid), cmdData, nil, s.mdmUsername, s.mdmPassword)
	if err != nil {
		return nil, err
	}

	var result dto.APIResult
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *nanomdmServiceImpl) Push(ctx context.Context, enrollments []string) (*dto.APIResult, error) {
	path := fmt.Sprintf("/v1/push/%s", strings.Join(enrollments, ","))
	resp, err := s.doRequest(ctx, http.MethodGet, s.mdmBaseURL, path, nil, nil, s.mdmUsername, s.mdmPassword)
	if err != nil {
		return nil, err
	}

	var result dto.APIResult
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *nanomdmServiceImpl) EscrowKeyUnlock(ctx context.Context, req *dto.EscrowKeyUnlockRequest) ([]byte, http.Header, int, error) {
	form := url.Values{}
	form.Set("topic", req.Topic)
	form.Set("serial", req.Serial)
	form.Set("productType", req.ProductType)
	form.Set("escrowKey", req.EscrowKey)
	form.Set("orgName", req.OrgName)
	form.Set("guid", req.Guid)
	if req.IMEI != "" {
		form.Set("imei", req.IMEI)
	}
	if req.IMEI2 != "" {
		form.Set("imei2", req.IMEI2)
	}
	if req.MEID != "" {
		form.Set("meid", req.MEID)
	}

	u, err := url.Parse(fmt.Sprintf("%s/v1/escrowkeyunlock", s.mdmBaseURL))
	if err != nil {
		return nil, nil, 0, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, nil, 0, err
	}

	httpReq.SetBasicAuth(s.mdmUsername, s.mdmPassword)
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("User-Agent", "MDM-Portal/1.0")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, nil, 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, 0, err
	}

	return body, resp.Header, resp.StatusCode, nil
}

func (s *nanomdmServiceImpl) GetVersion(ctx context.Context) (*dto.NanoMDMVersionResponse, error) {
	resp, err := s.doRequest(ctx, http.MethodGet, s.mdmBaseURL, "/version", nil, nil, s.mdmUsername, s.mdmPassword)
	if err != nil {
		return nil, err
	}

	var result dto.NanoMDMVersionResponse
	if err := s.handleResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ReloadDEPSyncer sends SIGHUP to the DEP syncer container to reload its configuration.
// This is useful after updating DEP tokens or configurations.
func (s *nanomdmServiceImpl) ReloadDEPSyncer(ctx context.Context) error {
	containerName := s.depSyncerContainer
	if containerName == "" {
		containerName = "mdm-nanodep-syncer-1"
	}

	var cmd *exec.Cmd
	if s.sudoPassword != "" {
		// Use sudo with password from stdin
		cmd = exec.CommandContext(ctx, "sudo", "-S", "docker", "kill", "-s", "SIGHUP", containerName)
		cmd.Stdin = strings.NewReader(s.sudoPassword + "\n")
	} else {
		// Try without sudo (assumes user has docker permissions or running as root)
		cmd = exec.CommandContext(ctx, "docker", "kill", "-s", "SIGHUP", containerName)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return apperror.ErrInternalServerError.WithMessage(
			fmt.Sprintf("failed to reload DEP syncer: %v, stderr: %s", err, stderr.String()),
		)
	}

	return nil
}
