package flags

import "net/http"

var (
	AutoDiscoveryKey     string
	DisableAutoUpdate   bool
	DisableWebSsh       bool
	MemoryModeAvailable bool
	Token               string
	Endpoint            string
	Interval            float64
	IgnoreUnsafeCert    bool
	MaxRetries          int
	ReconnectInterval   int
	InfoReportInterval  int
	IncludeNics         string
	ExcludeNics         string
	IncludeMountpoints  string
	MonthRotate         int
	CFAccessClientID    string
	CFAccessClientSecret string
)

// AddCloudflareAccessHeaders adds Cloudflare Access headers to HTTP request if configured
func AddCloudflareAccessHeaders(req *http.Request) {
	if CFAccessClientID != "" && CFAccessClientSecret != "" {
		req.Header.Set("CF-Access-Client-Id", CFAccessClientID)
		req.Header.Set("CF-Access-Client-Secret", CFAccessClientSecret)
	}
}
