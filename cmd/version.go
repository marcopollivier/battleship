package cmd

import "runtime/debug"

// develVersion is the version the Go tool stamps on a build that is not
// from a module version or a tag.
const develVersion = "(devel)"

// version returns the game version stamped in the binary, such as
// "v1.0.0", or "" when there is none (a development build).
// info and ok are what debug.ReadBuildInfo returns.
func version(info *debug.BuildInfo, ok bool) string {
	if !ok || info == nil || info.Main.Version == develVersion {
		return ""
	}
	return info.Main.Version
}

// versionText is what --versao shows: "battleship v1.0.0" or
// "battleship (desenvolvimento)".
func versionText(v string) string {
	if v == "" {
		return "battleship (desenvolvimento)"
	}
	return "battleship " + v
}
