package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

type IPFilterInternalConfig struct {
	isBuilt         bool                `json:"-" yaml:"-"`
	uuid            string              `json:"-" yaml:"-"`
	whiteListIPMap  map[string]struct{} `json:"-" yaml:"-"`
	whiteListIPNets []*net.IPNet        `json:"-" yaml:"-"`
	blackListIPMap  map[string]struct{} `json:"-" yaml:"-"`
	blackListIPNets []*net.IPNet        `json:"-" yaml:"-"`
}

type IPFilterConfig struct {
	IPFilterInternalConfig

	// filter mode
	Mode IPFilterMode `json:"mode" yaml:"mode"`

	// whitelist ips
	WhitelistIPs []string `json:"whitelist_ips" yaml:"whitelist_ips"`
	// whitelist cidrs
	WhitelistCIDRs []string `json:"whitelist_cidrs" yaml:"whitelist_cidrs"`

	// blacklist ips
	BlacklistIPs []string `json:"blacklist_ips" yaml:"blacklist_ips"`
	// blacklist cidrs
	BlacklistCIDRs []string `json:"blacklist_cidrs" yaml:"blacklist_cidrs"`

	// ip sources, the order of the sources is the priority of the ip
	// default: X-Forwarded-For > X-Real-IP > X-Remote-Addr > RemoteAddr
	IPSources []IPSource `json:"ip_sources" yaml:"ip_sources"`

	// deny status code, default 403
	DenyStatusCode int `json:"deny_status_code" yaml:"deny_status_code"`
	// deny message
	DenyMessage string `json:"deny_message" yaml:"deny_message"`

	// invalid IP action: how to handle invalid/empty IP addresses, default deny invalid IP
	InvalidIPAction InvalidIPAction `json:"invalid_ip_action" yaml:"invalid_ip_action"`
}

// Build: Must be called after the config is set
func (c *IPFilterConfig) Build() {
	if c == nil {
		return
	}
	c.mergeDefaultIPFilterConfig()
	c.uuid = uuid.New().String()
	c.whiteListIPMap = make(map[string]struct{}, len(c.WhitelistIPs))
	c.whiteListIPNets = make([]*net.IPNet, 0, len(c.WhitelistCIDRs))
	c.blackListIPMap = make(map[string]struct{}, len(c.BlacklistIPs))
	c.blackListIPNets = make([]*net.IPNet, 0, len(c.BlacklistCIDRs))
	for _, ip := range c.WhitelistIPs {
		c.whiteListIPMap[ip] = struct{}{}
	}
	for _, cidr := range c.WhitelistCIDRs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		c.whiteListIPNets = append(c.whiteListIPNets, ipNet)
	}
	for _, ip := range c.BlacklistIPs {
		c.blackListIPMap[ip] = struct{}{}
	}
	for _, cidr := range c.BlacklistCIDRs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		c.blackListIPNets = append(c.blackListIPNets, ipNet)
	}
	c.isBuilt = true
}

func (c *IPFilterConfig) mergeDefaultIPFilterConfig() {
	if c == nil {
		return
	}
	defaultConfig := newDefaultIPFilterConfig()
	if c.DenyStatusCode == 0 {
		c.DenyStatusCode = defaultConfig.DenyStatusCode
	}
	if c.DenyMessage == "" {
		c.DenyMessage = defaultConfig.DenyMessage
	}
	if len(c.IPSources) == 0 {
		c.IPSources = defaultConfig.IPSources
	}
}

func newDefaultIPFilterConfig() *IPFilterConfig {
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

func checkConfigShouldUpdate(currentConfig *IPFilterConfig, newConfig *IPFilterConfig) bool {
	if newConfig == nil || !newConfig.isBuilt {
		return false
	}
	if currentConfig == nil || currentConfig.uuid != newConfig.uuid {
		return true
	}
	return false
}

type configManager struct {
	mu     *sync.Mutex
	config *IPFilterConfig
}

func newConfigManager() *configManager {
	cm := &configManager{
		mu:     &sync.Mutex{},
		config: newDefaultIPFilterConfig(),
	}
	cm.config.Build()
	return cm
}

func (cm *configManager) updateConfigOnChangeFunc(newConfig *IPFilterConfig) {
	if checkConfigShouldUpdate(cm.config, newConfig) {
		cm.mu.Lock()
		defer cm.mu.Unlock()
		if checkConfigShouldUpdate(cm.config, newConfig) {
			// Deep copy IPFilterInternalConfig
			newCfg := &IPFilterConfig{}
			*newCfg = *newConfig
			newCfg.WhitelistIPs = make([]string, len(newConfig.WhitelistIPs))
			copy(newCfg.WhitelistIPs, newConfig.WhitelistIPs)
			newCfg.WhitelistCIDRs = make([]string, len(newConfig.WhitelistCIDRs))
			copy(newCfg.WhitelistCIDRs, newConfig.WhitelistCIDRs)
			newCfg.BlacklistIPs = make([]string, len(newConfig.BlacklistIPs))
			copy(newCfg.BlacklistIPs, newConfig.BlacklistIPs)
			newCfg.BlacklistCIDRs = make([]string, len(newConfig.BlacklistCIDRs))
			copy(newCfg.BlacklistCIDRs, newConfig.BlacklistCIDRs)
			newCfg.IPSources = make([]IPSource, len(newConfig.IPSources))
			copy(newCfg.IPSources, newConfig.IPSources)
			newCfg.whiteListIPMap = make(map[string]struct{}, len(newConfig.WhitelistIPs))
			for _, ip := range newConfig.WhitelistIPs {
				newCfg.whiteListIPMap[ip] = struct{}{}
			}
			newCfg.whiteListIPNets = make([]*net.IPNet, len(newConfig.whiteListIPNets))
			copy(newCfg.whiteListIPNets, newConfig.whiteListIPNets)
			newCfg.blackListIPMap = make(map[string]struct{}, len(newConfig.BlacklistIPs))
			for _, ip := range newConfig.BlacklistIPs {
				newCfg.blackListIPMap[ip] = struct{}{}
			}
			newCfg.blackListIPNets = make([]*net.IPNet, len(newConfig.blackListIPNets))
			copy(newCfg.blackListIPNets, newConfig.blackListIPNets)
			cm.config = newCfg
		}
	}
}

type IPInfo struct {
	IP   string
	From string
}

func isIPInList(clientIP string, ipMap map[string]struct{}) bool {
	_, ok := ipMap[clientIP]
	return ok
}

func isIPInCIDRList(clientIP string, ipNets []*net.IPNet) bool {
	ip := net.ParseIP(clientIP)
	if ip == nil {
		return false
	}

	for _, ipNet := range ipNets {
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
			if ip, _, err := net.SplitHostPort(xra); err == nil {
				return ip, "X-Remote-Addr"
			}
		}
	case IPSourceRemoteAddr:
		if remoteAddr := c.Request.RemoteAddr; remoteAddr != "" {
			if ip, _, err := net.SplitHostPort(remoteAddr); err == nil {
				return ip, "RemoteAddr"
			}
		}
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
		if isIPInList(clientIP, config.whiteListIPMap) || isIPInCIDRList(clientIP, config.whiteListIPNets) {
			if isIPInList(clientIP, config.blackListIPMap) || isIPInCIDRList(clientIP, config.blackListIPNets) {
				return IPFilterDeny
			}
			return IPFilterAllow
		}
		return IPFilterDeny

	case IPFilterModeBlacklist:
		if isIPInList(clientIP, config.blackListIPMap) || isIPInCIDRList(clientIP, config.blackListIPNets) {
			if isIPInList(clientIP, config.whiteListIPMap) || isIPInCIDRList(clientIP, config.whiteListIPNets) {
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
	cm := newConfigManager()
	return gin.HandlerFunc(func(c *gin.Context) {
		cm.updateConfigOnChangeFunc(config)

		clientIPMetadata := getClientIP(c, cm.config)

		allow := true
		defer func() {
			for _, f := range postProcessFuncs {
				if f != nil {
					f(clientIPMetadata, allow)
				}
			}
		}()

		action := checkIPAccess(clientIPMetadata.IP, cm.config)

		if action == IPFilterDeny {
			allow = false
			c.JSON(cm.config.DenyStatusCode, gin.H{
				"error": cm.config.DenyMessage,
			})
			c.Abort()
			return
		}

		c.Next()
	})
}
