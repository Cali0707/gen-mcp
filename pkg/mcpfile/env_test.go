package mcpfile

import (
	"os"
	"testing"

	"github.com/genmcp/gen-mcp/pkg/observability/logging"
	"github.com/stretchr/testify/assert"
)

func TestEnvOverrides(t *testing.T) {
	tt := map[string]struct {
		initialRuntime  *ServerRuntime
		expectedRuntime *ServerRuntime
		env             map[string]string
		expectErr       bool
	}{
		"no overrides": {
			initialRuntime: &ServerRuntime{
				TransportProtocol: "streamablehttp",
				StreamableHTTPConfig: &StreamableHTTPConfig{
					Port: 8080,
				},
			},
			expectedRuntime: &ServerRuntime{
				TransportProtocol: "streamablehttp",
				StreamableHTTPConfig: &StreamableHTTPConfig{
					Port: 8080,
				},
			},
		},
		"override transport protocol": {
			initialRuntime: &ServerRuntime{
				TransportProtocol: "streamablehttp",
				StreamableHTTPConfig: &StreamableHTTPConfig{
					Port: 8080,
				},
			},
			expectedRuntime: &ServerRuntime{
				TransportProtocol: "stdio",
				StreamableHTTPConfig: &StreamableHTTPConfig{
					Port: 8080,
				},
			},
			env: map[string]string{
				"GENMCP_TRANSPORTPROTOCOL": "stdio",
			},
		},
		"override nested port": {
			initialRuntime: &ServerRuntime{
				TransportProtocol: "streamablehttp",
				StreamableHTTPConfig: &StreamableHTTPConfig{
					Port: 8080,
				},
			},
			expectedRuntime: &ServerRuntime{
				TransportProtocol: "streamablehttp",
				StreamableHTTPConfig: &StreamableHTTPConfig{
					Port: 9000,
				},
			},
			env: map[string]string{
				"GENMCP_STREAMABLEHTTPCONFIG_PORT": "9000",
			},
		},
		"surfaces error on invalid env var type": {
			initialRuntime: &ServerRuntime{
				TransportProtocol: "streamablehttp",
				StreamableHTTPConfig: &StreamableHTTPConfig{
					Port: 8080,
				},
			},
			expectErr: true,
			env: map[string]string{
				"GENMCP_STREAMABLEHTTPCONFIG_PORT": "\"9000\"",
			},
		},
		"handles maps correctly": {
			initialRuntime: &ServerRuntime{
				TransportProtocol: "streamablehttp",
				StreamableHTTPConfig: &StreamableHTTPConfig{
					Port: 8080,
				},
			},
			expectedRuntime: &ServerRuntime{
				TransportProtocol: "streamablehttp",
				StreamableHTTPConfig: &StreamableHTTPConfig{
					Port: 8080,
				},
				LoggingConfig: &logging.LoggingConfig{
					InitialFields: map[string]any{
						"service": "genmcp",
					},
				},
			},
			env: map[string]string{
				"GENMCP_LOGGINGCONFIG_INITIALFIELDS": "{\"service\": \"genmcp\"}",
			},
		},
		"handles slices correctly": {
			initialRuntime: &ServerRuntime{
				TransportProtocol: "streamablehttp",
				StreamableHTTPConfig: &StreamableHTTPConfig{
					Port: 8080,
				},
			},
			expectedRuntime: &ServerRuntime{
				TransportProtocol: "streamablehttp",
				StreamableHTTPConfig: &StreamableHTTPConfig{
					Port: 8080,
				},
				LoggingConfig: &logging.LoggingConfig{
					OutputPaths: []string{"/out/1", "/out/2"},
				},
			},
			env: map[string]string{
				"GENMCP_LOGGINGCONFIG_OUTPUTPATHS": "/out/1,/out/2",
			},
		},
	}

	for tn, tc := range tt {
		t.Run(tn, func(t *testing.T) {
			for k, v := range tc.env {
				err := os.Setenv(k, v)
				assert.NoError(t, err)
			}
			defer func() {
				for k := range tc.env {
					err := os.Unsetenv(k)
					assert.NoError(t, err)
				}
			}()

			e := NewEnvRuntimeOverrider()

			err := e.ApplyOverrides(tc.initialRuntime)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRuntime, tc.initialRuntime)
			}

		})
	}
}
