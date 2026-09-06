package utils

import (
	"errors"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetQueryProps(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name     string
		endpoint string
		wantErr  bool
	}{
		{
			name:     "Get query props - all set with revRequired true",
			wantErr:  false,
			endpoint: "/api/v1/namespaces/komodorio/charts?name=testing&namespace=testing&revision=1",
		},
		{
			name:     "Get query props - no namespace with revRequired true",
			wantErr:  false,
			endpoint: "/api/v1/namespaces/komodorio/charts?name=testing&revision=1",
		},
		{
			name:     "Get query props - no name with revRequired true",
			wantErr:  true,
			endpoint: "/api/v1/namespaces/komodorio/charts?namespace=testing&revision=1",
		},
		{
			name:     "Get query props - with apiVersion specified",
			wantErr:  false,
			endpoint: "/api/v1/namespaces/komodorio/charts?name=testing&apiVersion=apps/v1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", tt.endpoint, nil)
			_, err := GetQueryProps(c)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetQueryProps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestChartAndVersion(t *testing.T) {
	tests := []struct {
		name      string
		params    string
		wantChart string
		wantVer   string
		wantError bool
	}{
		{
			name:      "Chart and version - successfully parsing chart and version",
			params:    "chart-1.0.0",
			wantChart: "chart",
			wantVer:   "1.0.0",
			wantError: false,
		},
		{
			name:      "Chart and version - successfully parsing chart and version",
			params:    "chart-v1.0.0",
			wantChart: "chart",
			wantVer:   "v1.0.0",
			wantError: false,
		},
		{
			name:      "Chart and version - successfully parsing chart and version",
			params:    "chart-v1.0.0-alpha",
			wantChart: "chart",
			wantVer:   "v1.0.0-alpha",
			wantError: false,
		},
		{
			name:      "Chart and version - successfully parsing chart and version",
			params:    "chart-1.0.0-alpha",
			wantChart: "chart",
			wantVer:   "1.0.0-alpha",
			wantError: false,
		},
		{
			name:      "Chart and version - parsing chart without version",
			params:    "chart",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, b, err := ChartAndVersion(tt.params)
			if (err != nil) != tt.wantError {
				t.Errorf("ChartAndVersion() error = %v, wantErr %v", err, tt.wantError)
				return
			}

			if a != tt.wantChart {
				t.Errorf("ChartAndVersion() got = %v, want %v", a, tt.wantChart)
			}

			if b != tt.wantVer {
				t.Errorf("ChartAndVersion() got1 = %v, want %v", b, tt.wantVer)
			}
		})
	}
}

func TestQualifiedKind(t *testing.T) {
	tests := []struct {
		kind       string
		apiVersion string
		want       string
	}{
		{"Widget", "new.example.com/v1alpha1", "Widget.new.example.com"},
		{"Widget", "old.example.com/v1alpha1", "Widget.old.example.com"},
		{"Deployment", "apps/v1", "Deployment"},
		{"ConfigMap", "v1", "ConfigMap"},
		{"ConfigMap", "", "ConfigMap"},
		{"Middleware", "traefik.io/v1alpha1", "Middleware.traefik.io"},
		{"Middleware", "traefik.containo.us/v1alpha1", "Middleware.traefik.containo.us"},
	}
	for _, tt := range tests {
		t.Run(tt.kind+"/"+tt.apiVersion, func(t *testing.T) {
			got := QualifiedKind(tt.kind, tt.apiVersion)
			if got != tt.want {
				t.Errorf("QualifiedKind(%q, %q) = %q, want %q", tt.kind, tt.apiVersion, got, tt.want)
			}
		})
	}
}

func TestEnvAsBool(t *testing.T) {
	// value: "true" | "1", default: false -> expect true
	t.Setenv("TEST", "true")
	want := true
	if EnvAsBool("TEST", false) != want {
		t.Errorf("Env 'TEST' value '%v' should be parsed to %v", os.Getenv("TEST"), want)
	}
	t.Setenv("TEST", "1")
	want = true
	if EnvAsBool("TEST", false) != want {
		t.Errorf("Env 'TEST' value '%v' should be parsed to %v", os.Getenv("TEST"), want)
	}

	// value: "false" | "0", default: true  -> expect false
	t.Setenv("TEST", "false")
	want = false
	if EnvAsBool("TEST", true) != want {
		t.Errorf("Env 'TEST' value '%v' should be parsed to %v", os.Getenv("TEST"), want)
	}
	t.Setenv("TEST", "0")
	want = false
	if EnvAsBool("TEST", true) != want {
		t.Errorf("Env 'TEST' value '%v' should be parsed to %v", os.Getenv("TEST"), want)
	}

	// value: "" | *, default: false -> expect false
	t.Setenv("TEST", "")
	want = false
	if EnvAsBool("TEST", false) != want {
		t.Errorf("Env 'TEST' value '%v' should be parsed to %v", os.Getenv("TEST"), want)
	}
	t.Setenv("TEST", "10random")
	want = false
	if EnvAsBool("TEST", false) != want {
		t.Errorf("Env 'TEST' value '%v' should be parsed to %v", os.Getenv("TEST"), want)
	}

	// value: "" | *, default: true -> expect true
	t.Setenv("TEST", "")
	want = true
	if EnvAsBool("TEST", true) != want {
		t.Errorf("Env 'TEST' value '%v' should be parsed to %v", os.Getenv("TEST"), want)
	}
	t.Setenv("TEST", "10random")
	want = true
	if EnvAsBool("TEST", true) != want {
		t.Errorf("Env 'TEST' value '%v' should be parsed to %v", os.Getenv("TEST"), want)
	}

	// uppercase and whitespace trimming
	t.Setenv("TEST", " TRUE ")
	want = true
	if EnvAsBool("TEST", false) != want {
		t.Errorf("Env 'TEST' value ' TRUE ' should be parsed to true")
	}
}

func TestTempFile(t *testing.T) {
	content := "foo: bar\nkey: value"
	fname, cleanup, err := TempFile(content)
	if err != nil {
		t.Fatalf("TempFile failed: %v", err)
	}
	defer cleanup()

	data, err := os.ReadFile(fname)
	if err != nil {
		t.Fatalf("failed to read created temp file: %v", err)
	}
	if string(data) != content {
		t.Errorf("content mismatch: got %q, want %q", string(data), content)
	}

	cleanup()
	if _, err := os.Stat(fname); !os.IsNotExist(err) {
		t.Errorf("expected file to be removed after cleanup")
	}
}

func TestRunCommand_Empty(t *testing.T) {
	_, err := RunCommand([]string{}, nil)
	if err == nil {
		t.Errorf("expected error for empty command")
	}
}

func TestCmdError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      CmdError
		expected string
	}{
		{
			name:     "with stderr",
			err:      CmdError{StdErr: "error from stderr"},
			expected: "error from stderr",
		},
		{
			name:     "with orig error",
			err:      CmdError{OrigError: errors.New("original error")},
			expected: "original error",
		},
		{
			name:     "with neither",
			err:      CmdError{},
			expected: "command failed with unknown error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("CmdError.Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestRunCommand_Success(t *testing.T) {
	out, err := RunCommand([]string{"echo", "test_output"}, nil)
	if err != nil {
		t.Fatalf("RunCommand failed: %v", err)
	}
	if !strings.Contains(out, "test_output") {
		t.Errorf("expected output to contain 'test_output', got %q", out)
	}
}

func TestRunCommand_WithEnv(t *testing.T) {
	out, err := RunCommand([]string{"sh", "-c", "echo $CUSTOM_VAR"}, map[string]string{
		"CUSTOM_VAR": "helm_dashboard_test",
	})
	if err != nil {
		t.Fatalf("RunCommand with env failed: %v", err)
	}
	if !strings.Contains(out, "helm_dashboard_test") {
		t.Errorf("expected output to contain 'helm_dashboard_test', got %q", out)
	}
}

func TestRunCommand_ExitError(t *testing.T) {
	_, err := RunCommand([]string{"sh", "-c", "exit 2"}, nil)
	if err == nil {
		t.Fatalf("expected command error for non-zero exit")
	}
	if _, ok := err.(CmdError); !ok {
		t.Errorf("expected error of type CmdError, got %T", err)
	}
}
