package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// createTestRouter: create a gin router for testing
func createTestRouter(config *IPFilterConfig) *gin.Engine {
	router := gin.New()
	router.Use(IPFilter(config))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	return router
}

func TestDefaultIPFilterConfig(t *testing.T) {
	config := newDefaultIPFilterConfig()
	config.Build()

	if config.Mode != IPFilterModeBlacklist {
		t.Errorf("Expected mode to be IPFilterModeBlacklist, got %v", config.Mode)
	}

	if config.DenyStatusCode != http.StatusForbidden {
		t.Errorf("Expected deny status code to be %d, got %d", http.StatusForbidden, config.DenyStatusCode)
	}

	if config.DenyMessage != "IP access denied" {
		t.Errorf("Expected deny message to be 'IP access denied', got '%s'", config.DenyMessage)
	}

	if config.InvalidIPAction != InvalidIPDeny {
		t.Errorf("Expected invalid IP action to be InvalidIPDeny, got %v", config.InvalidIPAction)
	}
}

func TestWhitelistIP(t *testing.T) {
	config := &IPFilterConfig{
		Mode:         IPFilterModeWhitelist,
		WhitelistIPs: []string{"127.0.0.1", "192.168.1.100"},
		BlacklistIPs: []string{"192.168.1.1"},
	}

	router := createTestRouter(config)

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for whitelisted IP, got %d", http.StatusOK, w.Code)
	}

	// Test non-whitelisted IP should be denied in whitelist mode
	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.2:12345"
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d for non-whitelisted IP in whitelist mode, got %d", http.StatusForbidden, w.Code)
	}
}

func TestWhitelistCIDR(t *testing.T) {
	config := &IPFilterConfig{
		Mode:           IPFilterModeWhitelist,
		WhitelistCIDRs: []string{"192.168.1.0/24", "10.0.0.0/8"},
		BlacklistIPs:   []string{"192.168.1.50"},
	}

	router := createTestRouter(config)

	// Test IP in whitelist CIDR but also in blacklist IP - should be denied (blacklist overrides)
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.50:12345"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d for blacklisted IP (even in whitelist CIDR), got %d", http.StatusForbidden, w.Code)
	}

	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "10.1.1.1:12345"
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for whitelisted CIDR IP, got %d", http.StatusOK, w.Code)
	}
}

func TestBlacklistIP(t *testing.T) {
	config := &IPFilterConfig{
		Mode:           IPFilterModeBlacklist,
		BlacklistIPs:   []string{"192.168.1.1", "10.0.0.1"},
		DenyStatusCode: http.StatusForbidden,
		DenyMessage:    "Access denied",
	}

	router := createTestRouter(config)

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d for blacklisted IP, got %d", http.StatusForbidden, w.Code)
	}

	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.2:12345"
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for non-blacklisted IP, got %d", http.StatusOK, w.Code)
	}
}

func TestBlacklistCIDR(t *testing.T) {
	config := &IPFilterConfig{
		Mode:           IPFilterModeBlacklist,
		BlacklistCIDRs: []string{"192.168.1.0/24", "172.16.0.0/16"},
		DenyStatusCode: http.StatusForbidden,
	}

	router := createTestRouter(config)

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.100:12345"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d for blacklisted CIDR IP, got %d", http.StatusForbidden, w.Code)
	}

	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "172.16.1.1:12345"
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d for blacklisted CIDR IP, got %d", http.StatusForbidden, w.Code)
	}

	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "10.0.0.1:12345"
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for non-blacklisted CIDR IP, got %d", http.StatusOK, w.Code)
	}
}

func TestAllowAllMode(t *testing.T) {
	config := &IPFilterConfig{
		Mode:         IPFilterModeAllowAll,
		BlacklistIPs: []string{"192.168.1.1"},
	}

	router := createTestRouter(config)

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d in AllowAll mode, got %d", http.StatusOK, w.Code)
	}
}

func TestDenyAllMode(t *testing.T) {
	config := &IPFilterConfig{
		Mode:           IPFilterModeDenyAll,
		WhitelistIPs:   []string{"127.0.0.1"},
		DenyStatusCode: http.StatusForbidden,
	}

	router := createTestRouter(config)

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d in DenyAll mode, got %d", http.StatusForbidden, w.Code)
	}
}

func TestGetClientIPFromHeaders(t *testing.T) {
	config := &IPFilterConfig{
		Mode:           IPFilterModeBlacklist,
		BlacklistIPs:   []string{"1.2.3.4"},
		DenyStatusCode: http.StatusForbidden,
	}

	router := createTestRouter(config)

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	req.Header.Set("X-Forwarded-For", "1.2.3.4")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d for X-Forwarded-For blacklisted IP, got %d", http.StatusForbidden, w.Code)
	}

	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	req.Header.Set("X-Real-IP", "1.2.3.4")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d for X-Real-IP blacklisted IP, got %d", http.StatusForbidden, w.Code)
	}

	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	req.Header.Set("X-Forwarded-For", "192.168.1.100")
	req.Header.Set("X-Real-IP", "1.2.3.4")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for X-Forwarded-For priority test, got %d", http.StatusOK, w.Code)
	}
}

func TestWhitelistModeWithBlacklistOverride(t *testing.T) {
	// In whitelist mode, blacklist overrides whitelist
	config := &IPFilterConfig{
		Mode:           IPFilterModeWhitelist,
		WhitelistIPs:   []string{"192.168.1.1"},
		BlacklistIPs:   []string{"192.168.1.1"},
		DenyStatusCode: http.StatusForbidden,
	}

	router := createTestRouter(config)

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected blacklist to override whitelist in whitelist mode, got status %d", w.Code)
	}
}

func TestBlacklistModeWithWhitelistOverride(t *testing.T) {
	// In blacklist mode, whitelist overrides blacklist
	config := &IPFilterConfig{
		Mode:         IPFilterModeBlacklist,
		WhitelistIPs: []string{"192.168.1.1"},
		BlacklistIPs: []string{"192.168.1.1"},
	}

	router := createTestRouter(config)

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected whitelist to override blacklist in blacklist mode, got status %d", w.Code)
	}
}

func TestNilConfigHandling(t *testing.T) {
	router := gin.New()
	router.Use(IPFilter(nil))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected default config to allow access, got status %d", w.Code)
	}
}

func TestInvalidIPHandling(t *testing.T) {
	configDeny := &IPFilterConfig{
		Mode:            IPFilterModeBlacklist,
		InvalidIPAction: InvalidIPDeny,
		DenyStatusCode:  http.StatusForbidden,
		DenyMessage:     "Invalid IP denied",
	}

	router := createTestRouter(configDeny)

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = ":12345"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d for empty IP with InvalidIPDeny, got %d", http.StatusForbidden, w.Code)
	}

	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	req.Header.Set("X-Real-IP", "invalid-ip-format")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d for invalid IP format with InvalidIPDeny, got %d", http.StatusForbidden, w.Code)
	}

	configAllow := &IPFilterConfig{
		Mode:            IPFilterModeBlacklist,
		InvalidIPAction: InvalidIPAllow,
	}

	router = createTestRouter(configAllow)

	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = ":12345"
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for empty IP with InvalidIPAllow, got %d", http.StatusOK, w.Code)
	}
}

func TestEmptyIPHeaders(t *testing.T) {
	config := &IPFilterConfig{
		Mode:            IPFilterModeBlacklist,
		InvalidIPAction: InvalidIPDeny,
		DenyStatusCode:  http.StatusForbidden,
	}

	router := gin.New()
	router.Use(IPFilter(config))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = ""
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d for completely invalid request, got %d", http.StatusForbidden, w.Code)
	}
}

func TestMaliciousIPSpoofing(t *testing.T) {
	config := &IPFilterConfig{
		Mode:            IPFilterModeWhitelist,
		WhitelistIPs:    []string{"127.0.0.1"},
		InvalidIPAction: InvalidIPDeny,
		DenyStatusCode:  http.StatusForbidden,
	}

	router := createTestRouter(config)

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "unknown:12345"
	req.Header.Set("X-Real-IP", "127.0.0.1.malicious")
	req.Header.Set("X-Forwarded-For", "not-an-ip")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d for malicious IP spoofing, got %d", http.StatusForbidden, w.Code)
	}

	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "malicious-host:port"
	req.Header.Set("X-Real-IP", "fake-localhost")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d for completely fake request, got %d", http.StatusForbidden, w.Code)
	}
}

func TestHotUpdate(t *testing.T) {
	config := &IPFilterConfig{
		Mode:           IPFilterModeBlacklist,
		BlacklistIPs:   []string{},
		DenyStatusCode: http.StatusForbidden,
		DenyMessage:    "Access denied",
	}

	router := createTestRouter(config)

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d before hot update, got %d", http.StatusOK, w.Code)
	}

	config.BlacklistIPs = append(config.BlacklistIPs, "192.168.1.1")

	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d after hot update, got %d", http.StatusForbidden, w.Code)
	}

	config.Mode = IPFilterModeAllowAll

	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d after mode change to AllowAll, got %d", http.StatusOK, w.Code)
	}
}

func TestConcurrentAccess(t *testing.T) {
	config := &IPFilterConfig{
		Mode:           IPFilterModeBlacklist,
		WhitelistIPs:   []string{"127.0.0.1"},
		BlacklistIPs:   []string{"192.168.1.1"},
		DenyStatusCode: http.StatusForbidden,
		DenyMessage:    "Access denied",
	}

	router := createTestRouter(config)

	done := make(chan bool)
	errors := make(chan string, 100)

	for i := 0; i < 50; i++ {
		go func(id int) {
			defer func() { done <- true }()

			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "127.0.0.1:12345"
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				errors <- fmt.Sprintf("Goroutine %d: Expected %d for whitelist IP, got %d",
					id, http.StatusOK, w.Code)
			}

			req = httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.1:12345"
			w = httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusForbidden {
				errors <- fmt.Sprintf("Goroutine %d: Expected %d for blacklist IP, got %d",
					id, http.StatusForbidden, w.Code)
			}
		}(i)
	}

	for i := 0; i < 50; i++ {
		<-done
	}

	close(errors)

	for err := range errors {
		t.Error(err)
	}
}

func TestRealWorldScenario(t *testing.T) {
	// Real world scenario: Allow internal networks but block some problematic IPs
	config := &IPFilterConfig{
		Mode:           IPFilterModeWhitelist,
		WhitelistIPs:   []string{"127.0.0.1"},
		WhitelistCIDRs: []string{"192.168.0.0/16", "10.0.0.0/8"},
		BlacklistIPs:   []string{"192.168.1.50", "10.0.0.100"}, // Block some problematic internal IPs
		DenyStatusCode: http.StatusForbidden,
		DenyMessage:    "Access denied by security policy",
	}

	router := createTestRouter(config)

	testCases := []struct {
		name        string
		ip          string
		expected    int
		description string
	}{
		{"Localhost", "127.0.0.1", http.StatusOK, "localhost should be allowed"},
		{"InternalIP", "192.168.1.100", http.StatusOK, "internal IP should be allowed"},
		{"AnotherInternalIP", "10.1.1.1", http.StatusOK, "another internal IP should be allowed"},
		{"BlacklistedInternalIP", "192.168.1.50", http.StatusForbidden, "blacklisted internal IP should be denied"},
		{"AnotherBlacklistedInternalIP", "10.0.0.100", http.StatusForbidden, "another blacklisted internal IP should be denied"},
		{"ExternalIP", "8.8.8.8", http.StatusForbidden, "external IP should be denied in whitelist mode"},
		{"NonWhitelistedIP", "172.16.1.1", http.StatusForbidden, "non-whitelisted IP should be denied"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = tc.ip + ":12345"
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expected {
				t.Errorf("%s: %s - Expected status %d, got %d",
					tc.name, tc.description, tc.expected, w.Code)
			}
		})
	}
}

func BenchmarkIPFilterMiddleware(b *testing.B) {
	config := &IPFilterConfig{
		Mode:           IPFilterModeBlacklist,
		WhitelistIPs:   []string{"127.0.0.1", "192.168.1.1"},
		BlacklistIPs:   []string{"1.2.3.4", "5.6.7.8"},
		BlacklistCIDRs: []string{"172.16.0.0/12"},
	}

	router := createTestRouter(config)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.100:12345"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func TestIPSourcePriority(t *testing.T) {
	// Test custom IP source priority: X-Real-IP first, then X-Forwarded-For
	config := &IPFilterConfig{
		Mode:         IPFilterModeBlacklist,
		IPSources:    []IPSource{IPSourceXRealIp, IPSourceXForwardedFor, IPSourceRemoteAddr},
		BlacklistIPs: []string{"1.2.3.4"}, // Blacklist X-Forwarded-For IP
	}

	router := createTestRouter(config)

	// Both X-Real-IP and X-Forwarded-For are present
	// X-Real-IP should be preferred (not in blacklist, so allowed)
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	req.Header.Set("X-Real-Ip", "127.0.0.1")     // Should be used (not blacklisted)
	req.Header.Set("X-Forwarded-For", "1.2.3.4") // Should be ignored (blacklisted)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d when X-Real-IP has priority, got %d", http.StatusOK, w.Code)
	}
}

func TestDefaultIPSourcePriority(t *testing.T) {
	// Test default IP source priority: X-Forwarded-For first
	config := &IPFilterConfig{
		Mode:         IPFilterModeBlacklist,
		BlacklistIPs: []string{"1.2.3.4"}, // Blacklist X-Forwarded-For IP
	}

	router := createTestRouter(config)

	// Both X-Forwarded-For and X-Real-IP are present
	// X-Forwarded-For should be preferred (blacklisted, so denied)
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	req.Header.Set("X-Forwarded-For", "1.2.3.4") // Should be used (blacklisted)
	req.Header.Set("X-Real-IP", "127.0.0.1")     // Should be ignored
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d when X-Forwarded-For has priority, got %d", http.StatusForbidden, w.Code)
	}
}

func TestIPSourceFallback(t *testing.T) {
	// Test IP source fallback: when preferred sources are not available
	config := &IPFilterConfig{
		Mode:         IPFilterModeBlacklist,
		IPSources:    []IPSource{IPSourceXRealIp, IPSourceRemoteAddr},
		BlacklistIPs: []string{"192.168.1.1"},
	}

	router := createTestRouter(config)

	// No X-Real-IP header, should fallback to RemoteAddr
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345" // Should be blacklisted
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d when falling back to RemoteAddr, got %d", http.StatusForbidden, w.Code)
	}
}

func BenchmarkHotUpdate(b *testing.B) {
	config := &IPFilterConfig{
		Mode:         IPFilterModeBlacklist,
		BlacklistIPs: []string{"1.2.3.4"},
	}

	router := createTestRouter(config)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if i%1000 == 0 {
			config.BlacklistIPs = append(config.BlacklistIPs[:0], "1.2.3.4", fmt.Sprintf("1.2.3.%d", i%255))
		}

		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.100:12345"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
