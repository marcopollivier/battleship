package cmd

import (
	"runtime/debug"
	"testing"
)

// The version comes from the build info (spec 024 CA-4): a tag version is
// kept, while a development build has none.
func TestVersion(t *testing.T) {
	tests := []struct {
		name string
		info *debug.BuildInfo
		ok   bool
		want string
	}{
		{"tag", &debug.BuildInfo{Main: debug.Module{Version: "v1.0.0"}}, true, "v1.0.0"},
		{"devel", &debug.BuildInfo{Main: debug.Module{Version: "(devel)"}}, true, ""},
		{"empty", &debug.BuildInfo{}, true, ""},
		{"no build info", nil, false, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := version(tt.info, tt.ok); got != tt.want {
				t.Errorf("version() = %q, want %q", got, tt.want)
			}
		})
	}
}

// versionText names the game and its version, or says it is a development build.
func TestVersionText(t *testing.T) {
	tests := []struct {
		v, want string
	}{
		{"v1.0.0", "battleship v1.0.0"},
		{"", "battleship (desenvolvimento)"},
	}
	for _, tt := range tests {
		if got := versionText(tt.v); got != tt.want {
			t.Errorf("versionText(%q) = %q, want %q", tt.v, got, tt.want)
		}
	}
}

// --versao shows the version and ends with code 0, without starting a game
// and before checking the other options (spec 024 CA-4). go test stamps no
// version, so it is a development build.
func TestRun_Version(t *testing.T) {
	for _, args := range [][]string{
		{"--versao"},
		{"--versao", "--modo", "xyz"},
	} {
		code, stdout, stderr := run("", args...)
		if code != 0 {
			t.Errorf("Run(%q) code = %d, want 0", args, code)
		}
		if want := "battleship (desenvolvimento)\n"; stdout != want {
			t.Errorf("Run(%q) stdout = %q, want %q", args, stdout, want)
		}
		if stderr != "" {
			t.Errorf("Run(%q) stderr = %q, want empty", args, stderr)
		}
	}
}
