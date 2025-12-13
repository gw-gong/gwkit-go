package middleware

import (
	"net"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

type IPFilterAction int

const (
	IPFilterAllow IPFilterAction = iota
	IPFilterDeny
)

type IPFilterMode int

const (
	IPFilterModeAllowAll IPFilterMode = iota
	IPFilterModeDenyAll
	IPFilterModeWhitelist
	IPFilterModeBlacklist
)

type InvalidIPAction int

const (
	InvalidIPDeny InvalidIPAction = iota
	InvalidIPAllow
)

type IPSource int

const (
	IPSourceXForwardedFor IPSource = iota // X-Forwarded-For
	IPSourceXRealIp                       // X-Real-Ip
	IPSourceXRemoteAddr                   // X-Remote-Addr
	IPSourceRemoteAddr                    // RemoteAddr
)

type IPFilterConfig struct {
	// filter mode
	Mode IPFilterMode `json:"mode" yaml:"mode" mapstructure:"mode"`

	// whitelist ips
	WhitelistIPs []string `json:"whitelist_ips" yaml:"whitelist_ips" mapstructure:"whitelist_ips"`
	// whitelist cidrs
	WhitelistCIDRs []string `json:"whitelist_cidrs" yaml:"whitelist_cidrs" mapstructure:"whitelist_cidrs"`

	// blacklist ips
	BlacklistIPs []string `json:"blacklist_ips" yaml:"blacklist_ips" mapstructure:"blacklist_ips"`
	// blacklist cidrs
	BlacklistCIDRs []string `json:"blacklist_cidrs" yaml:"blacklist_cidrs" mapstructure:"blacklist_cidrs"`

	// ip sources, the order of the sources is the priority of the ip
	// default: X-Forwarded-For > X-Real-IP > X-Remote-Addr > RemoteAddr
	IPSources []IPSource `json:"ip_sources" yaml:"ip_sources" mapstructure:"ip_sources"`

	// deny status code, default 403
	DenyStatusCode int `json:"deny_status_code" yaml:"deny_status_code" mapstructure:"deny_status_code"`
	// deny message
	DenyMessage string `json:"deny_message" yaml:"deny_message" mapstructure:"deny_message"`

	// invalid IP action: how to handle invalid/empty IP addresses, default deny invalid IP
	InvalidIPAction InvalidIPAction `json:"invalid_ip_action" yaml:"invalid_ip_action" mapstructure:"invalid_ip_action"`
}

type IPInfo struct {
	IP   string
	From string
}

func DefaultIPFilterConfig() *IPFilterConfig {
	return &IPFilterConfig{
		Mode:            IPFilterModeBlacklist,
		WhitelistIPs:    []string{},
		WhitelistCIDRs:  []string{},
		BlacklistIPs:    []string{},
		BlacklistCIDRs:  []string{},
		IPSources:       []IPSource{IPSourceXForwardedFor, IPSourceXRealIp, IPSourceXRemoteAddr, IPSourceRemoteAddr},
		DenyStatusCode:  http.StatusForbidden,
		DenyMessage:     "IP access denied",
		InvalidIPAction: InvalidIPDeny,
	}
}

func MergeDefaultIPFilterConfig(config *IPFilterConfig) *IPFilterConfig {
	defaultConfig := DefaultIPFilterConfig()
	if config == nil {
		return defaultConfig
	}
	if config.DenyStatusCode == 0 {
		config.DenyStatusCode = defaultConfig.DenyStatusCode
	}
	if config.DenyMessage == "" {
		config.DenyMessage = defaultConfig.DenyMessage
	}
	if len(config.IPSources) == 0 {
		config.IPSources = defaultConfig.IPSources
	}
	return config
}

func isIPInList(clientIP string, ipList []string) bool {
	return slices.Contains(ipList, clientIP)
}

func isIPInCIDRList(clientIP string, cidrList []string) bool {
	ip := net.ParseIP(clientIP)
	if ip == nil {
		return false
	}

	for _, cidr := range cidrList {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if ipNet.Contains(ip) {
			return true
		}
	}
	return false
}

func getIPFromSource(c *gin.Context, source IPSource) (string, string) {
	switch source {
	case IPSourceXForwardedFor:
		if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
			if ips := strings.Split(xff, ","); len(ips) > 0 {
				return strings.TrimSpace(ips[0]), "X-Forwarded-For"
			}
		}
	case IPSourceXRealIp:
		if xri := c.GetHeader("X-Real-Ip"); xri != "" {
			return xri, "X-Real-Ip"
		}
	case IPSourceXRemoteAddr:
		if xra := c.GetHeader("X-Remote-Addr"); xra != "" {
			return xra, "X-Remote-Addr"
		}
	case IPSourceRemoteAddr:
		if ip, _, err := net.SplitHostPort(c.Request.RemoteAddr); err == nil {
			return ip, "RemoteAddr"
		}
		return c.Request.RemoteAddr, "RemoteAddr"
	}
	return "", ""
}

// getClientIP: trust the authenticity of the ip, do not verify it, this problem is handled by the gateway
func getClientIP(c *gin.Context, config *IPFilterConfig) *IPInfo {
	ipInfo := &IPInfo{}

	for _, source := range config.IPSources {
		if ip, from := getIPFromSource(c, source); ip != "" {
			ipInfo.IP = ip
			ipInfo.From = from
			return ipInfo
		}
	}
	return ipInfo
}

func checkIPAccess(clientIP string, config *IPFilterConfig) IPFilterAction {
	if config.Mode == IPFilterModeAllowAll {
		return IPFilterAllow
	}
	if config.Mode == IPFilterModeDenyAll {
		return IPFilterDeny
	}

	if clientIP == "" || net.ParseIP(clientIP) == nil {
		if config.InvalidIPAction == InvalidIPDeny {
			return IPFilterDeny
		}
		return IPFilterAllow
	}

	switch config.Mode {
	case IPFilterModeWhitelist:
		if isIPInList(clientIP, config.WhitelistIPs) || isIPInCIDRList(clientIP, config.WhitelistCIDRs) {
			if isIPInList(clientIP, config.BlacklistIPs) || isIPInCIDRList(clientIP, config.BlacklistCIDRs) {
				return IPFilterDeny
			}
			return IPFilterAllow
		}
		return IPFilterDeny

	case IPFilterModeBlacklist:
		if isIPInList(clientIP, config.BlacklistIPs) || isIPInCIDRList(clientIP, config.BlacklistCIDRs) {
			if isIPInList(clientIP, config.WhitelistIPs) || isIPInCIDRList(clientIP, config.WhitelistCIDRs) {
				return IPFilterAllow
			}
			return IPFilterDeny
		}
		return IPFilterAllow

	default:
		return IPFilterAllow
	}
}

// IPFilter
// support hot update config
// support post-process function
func IPFilter(config *IPFilterConfig, postProcessFuncs ...func(iMetadata *IPInfo, allow bool)) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		currentConfig := MergeDefaultIPFilterConfig(config)

		clientIPMetadata := getClientIP(c, currentConfig)

		allow := true
		defer func() {
			for _, f := range postProcessFuncs {
				if f != nil {
					f(clientIPMetadata, allow)
				}
			}
		}()

		action := checkIPAccess(clientIPMetadata.IP, currentConfig)

		if action == IPFilterDeny {
			allow = false
			c.JSON(currentConfig.DenyStatusCode, gin.H{
				"error": currentConfig.DenyMessage,
			})
			c.Abort()
			return
		}

		c.Next()
	})
}
