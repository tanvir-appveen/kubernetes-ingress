package configs

import (
	"reflect"
	"sort"
	"testing"

	networking "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestParseRewrites(t *testing.T) {
	serviceName := "coffee-svc"
	serviceNamePart := "serviceName=" + serviceName
	rewritePath := "/beans/"
	rewritePathPart := "rewrite=" + rewritePath
	rewriteService := serviceNamePart + " " + rewritePathPart

	serviceNameActual, rewritePathActual, err := parseRewrites(rewriteService)
	if serviceName != serviceNameActual || rewritePath != rewritePathActual || err != nil {
		t.Errorf("parseRewrites(%s) should return %q, %q, nil; got %q, %q, %v", rewriteService, serviceName, rewritePath, serviceNameActual, rewritePathActual, err)
	}
}

func TestParseRewritesWithLeadingAndTrailingWhitespace(t *testing.T) {
	serviceName := "coffee-svc"
	serviceNamePart := "serviceName=" + serviceName
	rewritePath := "/beans/"
	rewritePathPart := "rewrite=" + rewritePath
	rewriteService := "\t\n " + serviceNamePart + " " + rewritePathPart + " \t\n"

	serviceNameActual, rewritePathActual, err := parseRewrites(rewriteService)
	if serviceName != serviceNameActual || rewritePath != rewritePathActual || err != nil {
		t.Errorf("parseRewrites(%s) should return %q, %q, nil; got %q, %q, %v", rewriteService, serviceName, rewritePath, serviceNameActual, rewritePathActual, err)
	}
}

func TestParseRewritesInvalidFormat(t *testing.T) {
	rewriteService := "serviceNamecoffee-svc rewrite=/"

	_, _, err := parseRewrites(rewriteService)
	if err == nil {
		t.Errorf("parseRewrites(%s) should return error, got nil", rewriteService)
	}
}

func TestParseStickyService(t *testing.T) {
	serviceName := "coffee-svc"
	serviceNamePart := "serviceName=" + serviceName
	stickyCookie := "srv_id expires=1h domain=.example.com path=/"
	stickyService := serviceNamePart + " " + stickyCookie

	serviceNameActual, stickyCookieActual, err := parseStickyService(stickyService)
	if serviceName != serviceNameActual || stickyCookie != stickyCookieActual || err != nil {
		t.Errorf("parseStickyService(%s) should return %q, %q, nil; got %q, %q, %v", stickyService, serviceName, stickyCookie, serviceNameActual, stickyCookieActual, err)
	}
}

func TestParseStickyServiceInvalidFormat(t *testing.T) {
	stickyService := "serviceNamecoffee-svc srv_id expires=1h domain=.example.com path=/"

	_, _, err := parseStickyService(stickyService)
	if err == nil {
		t.Errorf("parseStickyService(%s) should return error, got nil", stickyService)
	}
}

func TestFilterMasterAnnotations(t *testing.T) {
	masterAnnotations := map[string]string{
		"nginx.org/rewrites":                "serviceName=service1 rewrite=rewrite1",
		"nginx.org/ssl-services":            "service1",
		"nginx.org/hsts":                    "True",
		"nginx.org/hsts-max-age":            "2700000",
		"nginx.org/hsts-include-subdomains": "True",
	}
	removedAnnotations := filterMasterAnnotations(masterAnnotations)

	expectedfilteredMasterAnnotations := map[string]string{
		"nginx.org/hsts":                    "True",
		"nginx.org/hsts-max-age":            "2700000",
		"nginx.org/hsts-include-subdomains": "True",
	}
	expectedRemovedAnnotations := []string{
		"nginx.org/rewrites",
		"nginx.org/ssl-services",
	}

	sort.Strings(removedAnnotations)
	sort.Strings(expectedRemovedAnnotations)

	if !reflect.DeepEqual(expectedfilteredMasterAnnotations, masterAnnotations) {
		t.Errorf("filterMasterAnnotations returned %v, but expected %v", masterAnnotations, expectedfilteredMasterAnnotations)
	}
	if !reflect.DeepEqual(expectedRemovedAnnotations, removedAnnotations) {
		t.Errorf("filterMasterAnnotations returned %v, but expected %v", removedAnnotations, expectedRemovedAnnotations)
	}
}

func TestFilterMinionAnnotations(t *testing.T) {
	minionAnnotations := map[string]string{
		"nginx.org/rewrites":                "serviceName=service1 rewrite=rewrite1",
		"nginx.org/ssl-services":            "service1",
		"nginx.org/hsts":                    "True",
		"nginx.org/hsts-max-age":            "2700000",
		"nginx.org/hsts-include-subdomains": "True",
	}
	removedAnnotations := filterMinionAnnotations(minionAnnotations)

	expectedfilteredMinionAnnotations := map[string]string{
		"nginx.org/rewrites":     "serviceName=service1 rewrite=rewrite1",
		"nginx.org/ssl-services": "service1",
	}
	expectedRemovedAnnotations := []string{
		"nginx.org/hsts",
		"nginx.org/hsts-max-age",
		"nginx.org/hsts-include-subdomains",
	}

	sort.Strings(removedAnnotations)
	sort.Strings(expectedRemovedAnnotations)

	if !reflect.DeepEqual(expectedfilteredMinionAnnotations, minionAnnotations) {
		t.Errorf("filterMinionAnnotations returned %v, but expected %v", minionAnnotations, expectedfilteredMinionAnnotations)
	}
	if !reflect.DeepEqual(expectedRemovedAnnotations, removedAnnotations) {
		t.Errorf("filterMinionAnnotations returned %v, but expected %v", removedAnnotations, expectedRemovedAnnotations)
	}
}

func TestMergeMasterAnnotationsIntoMinion(t *testing.T) {
	masterAnnotations := map[string]string{
		"nginx.org/proxy-buffering":       "True",
		"nginx.org/proxy-buffers":         "2",
		"nginx.org/proxy-buffer-size":     "8k",
		"nginx.org/hsts":                  "True",
		"nginx.org/hsts-max-age":          "2700000",
		"nginx.org/proxy-connect-timeout": "50s",
		"nginx.com/jwt-token":             "$cookie_auth_token",
	}
	minionAnnotations := map[string]string{
		"nginx.org/client-max-body-size":  "2m",
		"nginx.org/proxy-connect-timeout": "20s",
	}
	mergeMasterAnnotationsIntoMinion(minionAnnotations, masterAnnotations)

	expectedMergedAnnotations := map[string]string{
		"nginx.org/proxy-buffering":       "True",
		"nginx.org/proxy-buffers":         "2",
		"nginx.org/proxy-buffer-size":     "8k",
		"nginx.org/client-max-body-size":  "2m",
		"nginx.org/proxy-connect-timeout": "20s",
	}
	if !reflect.DeepEqual(expectedMergedAnnotations, minionAnnotations) {
		t.Errorf("mergeMasterAnnotationsIntoMinion returned %v, but expected %v", minionAnnotations, expectedMergedAnnotations)
	}
}

func TestParseProxyConnectTimeoutAnnotation(t *testing.T) {
	baseCfg := NewDefaultConfigParams()
	defaultTimeout := baseCfg.ProxyConnectTimeout
	tests := []struct {
		ing      *IngressEx
		baseCfg  *ConfigParams
		expected string
	}{
		{
			ing: &IngressEx{
				Ingress: &networking.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{},
					},
				},
			},
			baseCfg:  baseCfg,
			expected: defaultTimeout,
		},
		{
			ing:      createTestIngress("nginx.org/proxy-connect-timeout", "30s"),
			baseCfg:  baseCfg,
			expected: "30s",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-connect-timeout", "1m 30s"),
			baseCfg:  baseCfg,
			expected: "1m 30s",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-connect-timeout", "1m30s"),
			baseCfg:  baseCfg,
			expected: "1m30s",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-connect-timeout", "30s 2m"),
			baseCfg:  baseCfg,
			expected: "30s 2m",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-connect-timeout", "10s"),
			baseCfg:  baseCfg,
			expected: "10s",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-connect-timeout", "60secs"),
			baseCfg:  baseCfg,
			expected: defaultTimeout,
		},
		{
			ing:      createTestIngress("nginx.org/proxy-connect-timeout", "invalid_time_string"),
			baseCfg:  baseCfg,
			expected: defaultTimeout,
		},
	}

	for _, test := range tests {
		cfg := parseAnnotations(test.ing, test.baseCfg, true, false, false)
		if test.expected != cfg.ProxyConnectTimeout {
			t.Errorf("parseAnnotations() returned cfg.ProxyConnectTimeout with %v, but expected %v", cfg.ProxyConnectTimeout, test.expected)
		}
	}
}

func TestParseProxyReadTimeoutAnnotation(t *testing.T) {
	baseCfg := NewDefaultConfigParams()
	defaultTimeout := baseCfg.ProxyReadTimeout
	tests := []struct {
		ing      *IngressEx
		baseCfg  *ConfigParams
		expected string
	}{
		{
			ing: &IngressEx{
				Ingress: &networking.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{},
					},
				},
			},
			baseCfg:  baseCfg,
			expected: defaultTimeout,
		},
		{
			ing:      createTestIngress("nginx.org/proxy-read-timeout", "30s"),
			baseCfg:  baseCfg,
			expected: "30s",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-read-timeout", "1m 30s"),
			baseCfg:  baseCfg,
			expected: "1m 30s",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-read-timeout", "1m30s"),
			baseCfg:  baseCfg,
			expected: "1m30s",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-read-timeout", "30s 2m"),
			baseCfg:  baseCfg,
			expected: "30s 2m",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-read-timeout", "10s"),
			baseCfg:  baseCfg,
			expected: "10s",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-read-timeout", "60secs"),
			baseCfg:  baseCfg,
			expected: defaultTimeout,
		},
		{
			ing:      createTestIngress("nginx.org/proxy-read-timeout", "invalid_time_string"),
			baseCfg:  baseCfg,
			expected: defaultTimeout,
		},
	}

	for _, test := range tests {
		cfg := parseAnnotations(test.ing, test.baseCfg, true, false, false)
		if test.expected != cfg.ProxyReadTimeout {
			t.Errorf("parseAnnotations() returned cfg.ProxyReadTimeout with %v, but expected %v", cfg.ProxyReadTimeout, test.expected)
		}
	}
}

func TestParseProxySendTimeoutAnnotation(t *testing.T) {
	baseCfg := NewDefaultConfigParams()
	defaultTimeout := baseCfg.ProxySendTimeout
	tests := []struct {
		ing      *IngressEx
		baseCfg  *ConfigParams
		expected string
	}{
		{
			ing: &IngressEx{
				Ingress: &networking.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{},
					},
				},
			},
			baseCfg:  baseCfg,
			expected: defaultTimeout,
		},
		{
			ing:      createTestIngress("nginx.org/proxy-send-timeout", "30s"),
			baseCfg:  baseCfg,
			expected: "30s",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-send-timeout", "1m 30s"),
			baseCfg:  baseCfg,
			expected: "1m 30s",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-send-timeout", "1m30s"),
			baseCfg:  baseCfg,
			expected: "1m30s",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-send-timeout", "30s 2m"),
			baseCfg:  baseCfg,
			expected: "30s 2m",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-send-timeout", "10s"),
			baseCfg:  baseCfg,
			expected: "10s",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-send-timeout", "60secs"),
			baseCfg:  baseCfg,
			expected: defaultTimeout,
		},
		{
			ing:      createTestIngress("nginx.org/proxy-send-timeout", "invalid_time_string"),
			baseCfg:  baseCfg,
			expected: defaultTimeout,
		},
	}

	for _, test := range tests {
		cfg := parseAnnotations(test.ing, test.baseCfg, true, false, false)
		if test.expected != cfg.ProxySendTimeout {
			t.Errorf("parseAnnotations() returned cfg.ProxySendTimeout with %v, but expected %v", cfg.ProxySendTimeout, test.expected)
		}
	}
}

func TestParseFailTimeoutAnnotation(t *testing.T) {
	baseCfg := NewDefaultConfigParams()
	defaultTimeout := baseCfg.FailTimeout
	tests := []struct {
		ing      *IngressEx
		baseCfg  *ConfigParams
		expected string
	}{
		{
			ing: &IngressEx{
				Ingress: &networking.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{},
					},
				},
			},
			baseCfg:  baseCfg,
			expected: defaultTimeout,
		},
		{
			ing:      createTestIngress("nginx.org/fail-timeout", "30s"),
			baseCfg:  baseCfg,
			expected: "30s",
		},
		{
			ing:      createTestIngress("nginx.org/fail-timeout", "1m 30s"),
			baseCfg:  baseCfg,
			expected: "1m 30s",
		},
		{
			ing:      createTestIngress("nginx.org/fail-timeout", "1m30s"),
			baseCfg:  baseCfg,
			expected: "1m30s",
		},
		{
			ing:      createTestIngress("nginx.org/fail-timeout", "30s 2m"),
			baseCfg:  baseCfg,
			expected: "30s 2m",
		},
		{
			ing:      createTestIngress("nginx.org/fail-timeout", "10s"),
			baseCfg:  baseCfg,
			expected: "10s",
		},
		{
			ing:      createTestIngress("nginx.org/fail-timeout", "60secs"),
			baseCfg:  baseCfg,
			expected: defaultTimeout,
		},
		{
			ing:      createTestIngress("nginx.org/fail-timeout", "invalid_time_string"),
			baseCfg:  baseCfg,
			expected: defaultTimeout,
		},
	}

	for _, test := range tests {
		cfg := parseAnnotations(test.ing, test.baseCfg, true, false, false)
		if test.expected != cfg.FailTimeout {
			t.Errorf("parseAnnotations() returned cfg.FailTimeout with %v, but expected %v", cfg.FailTimeout, test.expected)
		}
	}
}

func TestParseClientMaxBodySizeAnnotation(t *testing.T) {
	baseCfg := NewDefaultConfigParams()
	defaultSize := baseCfg.ClientMaxBodySize
	tests := []struct {
		ing      *IngressEx
		baseCfg  *ConfigParams
		expected string
	}{
		{
			ing: &IngressEx{
				Ingress: &networking.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{},
					},
				},
			},
			baseCfg:  baseCfg,
			expected: defaultSize,
		},
		{
			ing:      createTestIngress("nginx.org/client-max-body-size", "2k"),
			baseCfg:  baseCfg,
			expected: "2k",
		},
		{
			ing:      createTestIngress("nginx.org/client-max-body-size", "16M"),
			baseCfg:  baseCfg,
			expected: "16M",
		},
		{
			ing:      createTestIngress("nginx.org/client-max-body-size", "1g"),
			baseCfg:  baseCfg,
			expected: "1g",
		},
		{
			ing:      createTestIngress("nginx.org/client-max-body-size", "12M"),
			baseCfg:  baseCfg,
			expected: "12M",
		},
		{
			ing:      createTestIngress("nginx.org/client-max-body-size", "32Megabytes"),
			baseCfg:  baseCfg,
			expected: defaultSize,
		},
		{
			ing:      createTestIngress("nginx.org/client-max-body-size", "invalid_offset_string"),
			baseCfg:  baseCfg,
			expected: defaultSize,
		},
	}

	for _, test := range tests {
		cfg := parseAnnotations(test.ing, test.baseCfg, true, false, false)
		if test.expected != cfg.ClientMaxBodySize {
			t.Errorf("parseAnnotations() returned cfg.ClientMaxBodySize with %v, but expected %v", cfg.ClientMaxBodySize, test.expected)
		}
	}
}

func TestParseProxyBufferSizeAnnotation(t *testing.T) {
	baseCfg := NewDefaultConfigParams()
	defaultSize := baseCfg.ProxyBufferSize
	tests := []struct {
		ing      *IngressEx
		baseCfg  *ConfigParams
		expected string
	}{
		{
			ing: &IngressEx{
				Ingress: &networking.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{},
					},
				},
			},
			baseCfg:  baseCfg,
			expected: defaultSize,
		},
		{
			ing:      createTestIngress("nginx.org/proxy-buffer-size", "2k"),
			baseCfg:  baseCfg,
			expected: "2k",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-buffer-size", "16M"),
			baseCfg:  baseCfg,
			expected: "16M",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-buffer-size", "12M"),
			baseCfg:  baseCfg,
			expected: "12M",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-buffer-size", "1g"),
			baseCfg:  baseCfg,
			expected: defaultSize,
		},
		{
			ing:      createTestIngress("nginx.org/proxy-buffer-size", "32Megabytes"),
			baseCfg:  baseCfg,
			expected: defaultSize,
		},
		{
			ing:      createTestIngress("nginx.org/proxy-buffer-size", "invalid_size_string"),
			baseCfg:  baseCfg,
			expected: defaultSize,
		},
	}

	for _, test := range tests {
		cfg := parseAnnotations(test.ing, test.baseCfg, true, false, false)
		if test.expected != cfg.ProxyBufferSize {
			t.Errorf("parseAnnotations() returned cfg.ProxyBufferSize with %v, but expected %v", cfg.ProxyBufferSize, test.expected)
		}
	}
}

func TestParseProxyMaxTempFileSizeAnnotation(t *testing.T) {
	baseCfg := NewDefaultConfigParams()
	defaultSize := baseCfg.ProxyMaxTempFileSize
	tests := []struct {
		ing      *IngressEx
		baseCfg  *ConfigParams
		expected string
	}{
		{
			ing: &IngressEx{
				Ingress: &networking.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{},
					},
				},
			},
			baseCfg:  baseCfg,
			expected: defaultSize,
		},
		{
			ing:      createTestIngress("nginx.org/proxy-max-temp-file-size", "2k"),
			baseCfg:  baseCfg,
			expected: "2k",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-max-temp-file-size", "16M"),
			baseCfg:  baseCfg,
			expected: "16M",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-max-temp-file-size", "12M"),
			baseCfg:  baseCfg,
			expected: "12M",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-max-temp-file-size", "1g"),
			baseCfg:  baseCfg,
			expected: defaultSize,
		},
		{
			ing:      createTestIngress("nginx.org/proxy-max-temp-file-size", "32Megabytes"),
			baseCfg:  baseCfg,
			expected: defaultSize,
		},
		{
			ing:      createTestIngress("nginx.org/proxy-max-temp-file-size", "invalid_size_string"),
			baseCfg:  baseCfg,
			expected: defaultSize,
		},
	}

	for _, test := range tests {
		cfg := parseAnnotations(test.ing, test.baseCfg, true, false, false)
		if test.expected != cfg.ProxyMaxTempFileSize {
			t.Errorf("parseAnnotations() returned cfg.ProxyMaxTempFileSize with %v, but expected %v", cfg.ProxyMaxTempFileSize, test.expected)
		}
	}
}

func TestParseUpstreamZoneSizeAnnotation(t *testing.T) {
	baseCfg := NewDefaultConfigParams()
	defaultSize := baseCfg.UpstreamZoneSize
	tests := []struct {
		ing      *IngressEx
		baseCfg  *ConfigParams
		expected string
	}{
		{
			ing: &IngressEx{
				Ingress: &networking.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{},
					},
				},
			},
			baseCfg:  baseCfg,
			expected: defaultSize,
		},
		{
			ing:      createTestIngress("nginx.org/upstream-zone-size", "2k"),
			baseCfg:  baseCfg,
			expected: "2k",
		},
		{
			ing:      createTestIngress("nginx.org/upstream-zone-size", "16M"),
			baseCfg:  baseCfg,
			expected: "16M",
		},
		{
			ing:      createTestIngress("nginx.org/upstream-zone-size", "12M"),
			baseCfg:  baseCfg,
			expected: "12M",
		},
		{
			ing:      createTestIngress("nginx.org/upstream-zone-size", "1g"),
			baseCfg:  baseCfg,
			expected: defaultSize,
		},
		{
			ing:      createTestIngress("nginx.org/upstream-zone-size", "32Megabytes"),
			baseCfg:  baseCfg,
			expected: defaultSize,
		},
		{
			ing:      createTestIngress("nginx.org/upstream-zone-size", "invalid_size_string"),
			baseCfg:  baseCfg,
			expected: defaultSize,
		},
	}

	for _, test := range tests {
		cfg := parseAnnotations(test.ing, test.baseCfg, true, false, false)
		if test.expected != cfg.UpstreamZoneSize {
			t.Errorf("parseAnnotations() returned cfg.UpstreamZoneSize with %v, but expected %v", cfg.UpstreamZoneSize, test.expected)
		}
	}
}

func TestParseMaxFailsAnnotation(t *testing.T) {
	baseCfg := NewDefaultConfigParams()
	defaultNum := baseCfg.MaxFails
	tests := []struct {
		ing      *IngressEx
		baseCfg  *ConfigParams
		expected int
	}{
		{
			ing: &IngressEx{
				Ingress: &networking.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{},
					},
				},
			},
			baseCfg:  baseCfg,
			expected: defaultNum,
		},
		{
			ing:      createTestIngress("nginx.org/max-fails", "0"),
			baseCfg:  baseCfg,
			expected: 0,
		},
		{
			ing:      createTestIngress("nginx.org/max-fails", "1"),
			baseCfg:  baseCfg,
			expected: 1,
		},
		{
			ing:      createTestIngress("nginx.org/max-fails", "100"),
			baseCfg:  baseCfg,
			expected: 100,
		},
		{
			ing:      createTestIngress("nginx.org/max-fails", "-1"),
			baseCfg:  baseCfg,
			expected: defaultNum,
		},
		{
			ing:      createTestIngress("nginx.org/max-fails", "-100"),
			baseCfg:  baseCfg,
			expected: defaultNum,
		},
		{
			ing:      createTestIngress("nginx.org/max-fails", "invalid_non_negative_int"),
			baseCfg:  baseCfg,
			expected: defaultNum,
		},
	}

	for _, test := range tests {
		cfg := parseAnnotations(test.ing, test.baseCfg, true, false, false)
		if test.expected != cfg.MaxFails {
			t.Errorf("parseAnnotations() returned cfg.MaxFails with %v, but expected %v", cfg.MaxFails, test.expected)
		}
	}
}

func TestParseMaxConnsAnnotation(t *testing.T) {
	baseCfg := NewDefaultConfigParams()
	defaultNum := baseCfg.MaxConns
	tests := []struct {
		ing      *IngressEx
		baseCfg  *ConfigParams
		expected int
	}{
		{
			ing: &IngressEx{
				Ingress: &networking.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{},
					},
				},
			},
			baseCfg:  baseCfg,
			expected: defaultNum,
		},
		{
			ing:      createTestIngress("nginx.org/max-conns", "0"),
			baseCfg:  baseCfg,
			expected: 0,
		},
		{
			ing:      createTestIngress("nginx.org/max-conns", "1"),
			baseCfg:  baseCfg,
			expected: 1,
		},
		{
			ing:      createTestIngress("nginx.org/max-conns", "100"),
			baseCfg:  baseCfg,
			expected: 100,
		},
		{
			ing:      createTestIngress("nginx.org/max-conns", "-1"),
			baseCfg:  baseCfg,
			expected: defaultNum,
		},
		{
			ing:      createTestIngress("nginx.org/max-conns", "-100"),
			baseCfg:  baseCfg,
			expected: defaultNum,
		},
		{
			ing:      createTestIngress("nginx.org/max-conns", "invalid_non_negative_int"),
			baseCfg:  baseCfg,
			expected: defaultNum,
		},
	}

	for _, test := range tests {
		cfg := parseAnnotations(test.ing, test.baseCfg, true, false, false)
		if test.expected != cfg.MaxConns {
			t.Errorf("parseAnnotations() returned cfg.MaxConns with %v, but expected %v", cfg.MaxConns, test.expected)
		}
	}
}

func TestParseProxyBuffersAnnotation(t *testing.T) {
	baseCfg := NewDefaultConfigParams()
	defaultSetting := baseCfg.ProxyBuffers
	tests := []struct {
		ing      *IngressEx
		baseCfg  *ConfigParams
		expected string
	}{
		{
			ing: &IngressEx{
				Ingress: &networking.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{},
					},
				},
			},
			baseCfg:  baseCfg,
			expected: defaultSetting,
		},
		{
			ing:      createTestIngress("nginx.org/proxy-buffers", "8 2k"),
			baseCfg:  baseCfg,
			expected: "8 2k",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-buffers", "4 16M"),
			baseCfg:  baseCfg,
			expected: "4 16M",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-buffers", "1 12M"),
			baseCfg:  baseCfg,
			expected: "1 12M",
		},
		{
			ing:      createTestIngress("nginx.org/proxy-buffers", "2k"),
			baseCfg:  baseCfg,
			expected: defaultSetting,
		},
		{
			ing:      createTestIngress("nginx.org/proxy-buffers", "8 1g"),
			baseCfg:  baseCfg,
			expected: defaultSetting,
		},
		{
			ing:      createTestIngress("nginx.org/proxy-buffers", "6 32Megabytes"),
			baseCfg:  baseCfg,
			expected: defaultSetting,
		},
		{
			ing:      createTestIngress("nginx.org/proxy-buffers", "invalid_proxy_buffers_string"),
			baseCfg:  baseCfg,
			expected: defaultSetting,
		},
	}

	for _, test := range tests {
		cfg := parseAnnotations(test.ing, test.baseCfg, true, false, false)
		if test.expected != cfg.ProxyBuffers {
			t.Errorf("parseAnnotations() returned cfg.ProxyBuffers with %v, but expected %v", cfg.ProxyBuffers, test.expected)
		}
	}
}

func createTestIngress(name string, value string) *IngressEx {
	return &IngressEx{
		Ingress: &networking.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					name: value,
				},
			},
		},
	}
}
