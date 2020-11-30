package configs

import (
	"testing"

	v1 "k8s.io/api/core/v1"
)

func TestParseProxyConnectTimeoutKey(t *testing.T) {
	defaultTimeout := NewDefaultConfigParams().ProxyConnectTimeout
	tests := []struct {
		cfgMap   *v1.ConfigMap
		expected string
	}{
		{
			cfgMap: &v1.ConfigMap{
				Data: map[string]string{},
			},
			expected: defaultTimeout,
		},
		{
			cfgMap:   createTestConfigMap("proxy-connect-timeout", "60s"),
			expected: "60s",
		},
		{
			cfgMap:   createTestConfigMap("proxy-connect-timeout", "1m 30s"),
			expected: "1m 30s",
		},
		{
			cfgMap:   createTestConfigMap("proxy-connect-timeout", "1m30s"),
			expected: "1m30s",
		},
		{
			cfgMap:   createTestConfigMap("proxy-connect-timeout", "30s 2m"),
			expected: "30s 2m",
		},
		{
			cfgMap:   createTestConfigMap("proxy-connect-timeout", "60secs"),
			expected: defaultTimeout,
		},
		{
			cfgMap:   createTestConfigMap("proxy-connect-timeout", "invalid_time_string"),
			expected: defaultTimeout,
		},
	}

	for _, test := range tests {
		cfg := ParseConfigMap(test.cfgMap, true, false)
		if test.expected != cfg.ProxyConnectTimeout {
			t.Errorf("ParseConfigMap() returned cfg.ProxyConnectTimeout with %v, but expected %v", cfg.ProxyConnectTimeout, test.expected)
		}
	}
}

func TestParseProxyReadTimeoutKey(t *testing.T) {
	defaultTimeout := NewDefaultConfigParams().ProxyReadTimeout
	tests := []struct {
		cfgMap   *v1.ConfigMap
		expected string
	}{
		{
			cfgMap: &v1.ConfigMap{
				Data: map[string]string{},
			},
			expected: defaultTimeout,
		},
		{
			cfgMap:   createTestConfigMap("proxy-read-timeout", "60s"),
			expected: "60s",
		},
		{
			cfgMap:   createTestConfigMap("proxy-read-timeout", "1m 30s"),
			expected: "1m 30s",
		},
		{
			cfgMap:   createTestConfigMap("proxy-read-timeout", "1m30s"),
			expected: "1m30s",
		},
		{
			cfgMap:   createTestConfigMap("proxy-read-timeout", "30s 2m"),
			expected: "30s 2m",
		},
		{
			cfgMap:   createTestConfigMap("proxy-read-timeout", "60secs"),
			expected: defaultTimeout,
		},
		{
			cfgMap:   createTestConfigMap("proxy-read-timeout", "invalid_time_string"),
			expected: defaultTimeout,
		},
	}

	for _, test := range tests {
		cfg := ParseConfigMap(test.cfgMap, true, false)
		if test.expected != cfg.ProxyReadTimeout {
			t.Errorf("ParseConfigMap() returned cfg.ProxyReadTimeout with %v, but expected %v", cfg.ProxyReadTimeout, test.expected)
		}
	}
}

func TestParseProxySendTimeoutKey(t *testing.T) {
	defaultTimeout := NewDefaultConfigParams().ProxySendTimeout
	tests := []struct {
		cfgMap   *v1.ConfigMap
		expected string
	}{
		{
			cfgMap: &v1.ConfigMap{
				Data: map[string]string{},
			},
			expected: defaultTimeout,
		},
		{
			cfgMap:   createTestConfigMap("proxy-send-timeout", "60s"),
			expected: "60s",
		},
		{
			cfgMap:   createTestConfigMap("proxy-send-timeout", "1m 30s"),
			expected: "1m 30s",
		},
		{
			cfgMap:   createTestConfigMap("proxy-send-timeout", "1m30s"),
			expected: "1m30s",
		},
		{
			cfgMap:   createTestConfigMap("proxy-send-timeout", "30s 2m"),
			expected: "30s 2m",
		},
		{
			cfgMap:   createTestConfigMap("proxy-send-timeout", "60secs"),
			expected: defaultTimeout,
		},
		{
			cfgMap:   createTestConfigMap("proxy-send-timeout", "invalid_time_string"),
			expected: defaultTimeout,
		},
	}

	for _, test := range tests {
		cfg := ParseConfigMap(test.cfgMap, true, false)
		if test.expected != cfg.ProxySendTimeout {
			t.Errorf("ParseConfigMap() returned cfg.ProxySendTimeout with %v, but expected %v", cfg.ProxySendTimeout, test.expected)
		}
	}
}

func TestParseFailTimeoutKey(t *testing.T) {
	defaultTimeout := NewDefaultConfigParams().FailTimeout
	tests := []struct {
		cfgMap   *v1.ConfigMap
		expected string
	}{
		{
			cfgMap: &v1.ConfigMap{
				Data: map[string]string{},
			},
			expected: defaultTimeout,
		},
		{
			cfgMap:   createTestConfigMap("fail-timeout", "60s"),
			expected: "60s",
		},
		{
			cfgMap:   createTestConfigMap("fail-timeout", "1m 30s"),
			expected: "1m 30s",
		},
		{
			cfgMap:   createTestConfigMap("fail-timeout", "1m30s"),
			expected: "1m30s",
		},
		{
			cfgMap:   createTestConfigMap("fail-timeout", "30s 2m"),
			expected: "30s 2m",
		},
		{
			cfgMap:   createTestConfigMap("fail-timeout", "60secs"),
			expected: defaultTimeout,
		},
		{
			cfgMap:   createTestConfigMap("fail-timeout", "invalid_time_string"),
			expected: defaultTimeout,
		},
	}

	for _, test := range tests {
		cfg := ParseConfigMap(test.cfgMap, true, false)
		if test.expected != cfg.FailTimeout {
			t.Errorf("ParseConfigMap() returned cfg.FailTimeout with %v, but expected %v", cfg.FailTimeout, test.expected)
		}
	}
}

func TestParseClientMaxBodySizeKey(t *testing.T) {
	defaultSize := NewDefaultConfigParams().ClientMaxBodySize
	tests := []struct {
		cfgMap   *v1.ConfigMap
		expected string
	}{
		{
			cfgMap: &v1.ConfigMap{
				Data: map[string]string{},
			},
			expected: defaultSize,
		},
		{
			cfgMap:   createTestConfigMap("client-max-body-size", "2k"),
			expected: "2k",
		},
		{
			cfgMap:   createTestConfigMap("client-max-body-size", "16M"),
			expected: "16M",
		},
		{
			cfgMap:   createTestConfigMap("client-max-body-size", "1g"),
			expected: "1g",
		},
		{
			cfgMap:   createTestConfigMap("client-max-body-size", "12M"),
			expected: "12M",
		},
		{
			cfgMap:   createTestConfigMap("client-max-body-size", "32Megabytes"),
			expected: defaultSize,
		},
		{
			cfgMap:   createTestConfigMap("client-max-body-size", "invalid_offset_string"),
			expected: defaultSize,
		},
	}

	for _, test := range tests {
		cfg := ParseConfigMap(test.cfgMap, true, false)
		if test.expected != cfg.ClientMaxBodySize {
			t.Errorf("ParseConfigMap() returned cfg.ClientMaxBodySize with %v, but expected %v", cfg.ClientMaxBodySize, test.expected)
		}
	}
}

func TestParseProxyBufferSizeKey(t *testing.T) {
	defaultSize := NewDefaultConfigParams().ProxyBufferSize
	tests := []struct {
		cfgMap   *v1.ConfigMap
		expected string
	}{
		{
			cfgMap: &v1.ConfigMap{
				Data: map[string]string{},
			},
			expected: defaultSize,
		},
		{
			cfgMap:   createTestConfigMap("proxy-buffer-size", "2k"),
			expected: "2k",
		},
		{
			cfgMap:   createTestConfigMap("proxy-buffer-size", "16M"),
			expected: "16M",
		},
		{
			cfgMap:   createTestConfigMap("proxy-buffer-size", "12M"),
			expected: "12M",
		},
		{
			cfgMap:   createTestConfigMap("proxy-buffer-size", "1g"),
			expected: defaultSize,
		},
		{
			cfgMap:   createTestConfigMap("proxy-buffer-size", "32Megabytes"),
			expected: defaultSize,
		},
		{
			cfgMap:   createTestConfigMap("proxy-buffer-size", "invalid_size_string"),
			expected: defaultSize,
		},
	}

	for _, test := range tests {
		cfg := ParseConfigMap(test.cfgMap, true, false)
		if test.expected != cfg.ProxyBufferSize {
			t.Errorf("ParseConfigMap() returned cfg.ProxyBufferSize with %v, but expected %v", cfg.ProxyBufferSize, test.expected)
		}
	}
}

func TestParseProxyMaxTempFileSizeKey(t *testing.T) {
	defaultSize := NewDefaultConfigParams().ProxyMaxTempFileSize
	tests := []struct {
		cfgMap   *v1.ConfigMap
		expected string
	}{
		{
			cfgMap: &v1.ConfigMap{
				Data: map[string]string{},
			},
			expected: defaultSize,
		},
		{
			cfgMap:   createTestConfigMap("proxy-max-temp-file-size", "2k"),
			expected: "2k",
		},
		{
			cfgMap:   createTestConfigMap("proxy-max-temp-file-size", "16M"),
			expected: "16M",
		},
		{
			cfgMap:   createTestConfigMap("proxy-max-temp-file-size", "12M"),
			expected: "12M",
		},
		{
			cfgMap:   createTestConfigMap("proxy-max-temp-file-size", "1g"),
			expected: defaultSize,
		},
		{
			cfgMap:   createTestConfigMap("proxy-max-temp-file-size", "32Megabytes"),
			expected: defaultSize,
		},
		{
			cfgMap:   createTestConfigMap("proxy-max-temp-file-size", "invalid_size_string"),
			expected: defaultSize,
		},
	}

	for _, test := range tests {
		cfg := ParseConfigMap(test.cfgMap, true, false)
		if test.expected != cfg.ProxyMaxTempFileSize {
			t.Errorf("ParseConfigMap() returned cfg.ProxyMaxTempFileSize with %v, but expected %v", cfg.ProxyMaxTempFileSize, test.expected)
		}
	}
}

func TestParseUpstreamZoneSizeKey(t *testing.T) {
	defaultSize := NewDefaultConfigParams().UpstreamZoneSize
	tests := []struct {
		cfgMap   *v1.ConfigMap
		expected string
	}{
		{
			cfgMap: &v1.ConfigMap{
				Data: map[string]string{},
			},
			expected: defaultSize,
		},
		{
			cfgMap:   createTestConfigMap("upstream-zone-size", "2k"),
			expected: "2k",
		},
		{
			cfgMap:   createTestConfigMap("upstream-zone-size", "16M"),
			expected: "16M",
		},
		{
			cfgMap:   createTestConfigMap("upstream-zone-size", "12M"),
			expected: "12M",
		},
		{
			cfgMap:   createTestConfigMap("upstream-zone-size", "1g"),
			expected: defaultSize,
		},
		{
			cfgMap:   createTestConfigMap("upstream-zone-size", "32Megabytes"),
			expected: defaultSize,
		},
		{
			cfgMap:   createTestConfigMap("upstream-zone-size", "invalid_size_string"),
			expected: defaultSize,
		},
	}

	for _, test := range tests {
		cfg := ParseConfigMap(test.cfgMap, true, false)
		if test.expected != cfg.UpstreamZoneSize {
			t.Errorf("ParseConfigMap() returned cfg.UpstreamZoneSize with %v, but expected %v", cfg.UpstreamZoneSize, test.expected)
		}
	}
}

func TestParseMaxFailsKey(t *testing.T) {
	defaultNum := NewDefaultConfigParams().MaxFails
	tests := []struct {
		cfgMap   *v1.ConfigMap
		expected int
	}{
		{
			cfgMap: &v1.ConfigMap{
				Data: map[string]string{},
			},
			expected: defaultNum,
		},
		{
			cfgMap:   createTestConfigMap("max-fails", "0"),
			expected: 0,
		},
		{
			cfgMap:   createTestConfigMap("max-fails", "1"),
			expected: 1,
		},
		{
			cfgMap:   createTestConfigMap("max-fails", "100"),
			expected: 100,
		},
		{
			cfgMap:   createTestConfigMap("max-fails", "-1"),
			expected: defaultNum,
		},
		{
			cfgMap:   createTestConfigMap("max-fails", "-100"),
			expected: defaultNum,
		},
		{
			cfgMap:   createTestConfigMap("max-fails", "invalid_non_negative_int"),
			expected: defaultNum,
		},
	}

	for _, test := range tests {
		cfg := ParseConfigMap(test.cfgMap, true, false)
		if test.expected != cfg.MaxFails {
			t.Errorf("ParseConfigMap() returned cfg.MaxFails with %v, but expected %v", cfg.MaxFails, test.expected)
		}
	}
}

func TestParseProxyBuffersKey(t *testing.T) {
	defaultSetting := NewDefaultConfigParams().ProxyBuffers
	tests := []struct {
		cfgMap   *v1.ConfigMap
		expected string
	}{
		{
			cfgMap: &v1.ConfigMap{
				Data: map[string]string{},
			},
			expected: defaultSetting,
		},
		{
			cfgMap:   createTestConfigMap("proxy-buffers", "8 2k"),
			expected: "8 2k",
		},
		{
			cfgMap:   createTestConfigMap("proxy-buffers", "4 16M"),
			expected: "4 16M",
		},
		{
			cfgMap:   createTestConfigMap("proxy-buffers", "1 12M"),
			expected: "1 12M",
		},
		{
			cfgMap:   createTestConfigMap("proxy-buffers", "2k"),
			expected: defaultSetting,
		},
		{
			cfgMap:   createTestConfigMap("proxy-buffers", "8 1g"),
			expected: defaultSetting,
		},
		{
			cfgMap:   createTestConfigMap("proxy-buffers", "6 32Megabytes"),
			expected: defaultSetting,
		},
		{
			cfgMap:   createTestConfigMap("proxy-buffers", "invalid_proxy_buffers_string"),
			expected: defaultSetting,
		},
	}

	for _, test := range tests {
		cfg := ParseConfigMap(test.cfgMap, true, false)
		if test.expected != cfg.ProxyBuffers {
			t.Errorf("ParseConfigMap() returned cfg.ProxyBuffers with %v, but expected %v", cfg.ProxyBuffers, test.expected)
		}
	}
}

func createTestConfigMap(name string, value string) *v1.ConfigMap {
	return &v1.ConfigMap{
		Data: map[string]string{
			name: value,
		},
	}
}
