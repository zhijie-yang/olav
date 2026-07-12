package source

import (
	"fmt"
	"runtime"
	"strings"
)

type Platform struct {
	All          bool
	OS           string
	Architecture string
	Variant      string
}

func ParsePlatform(s string) (Platform, error) {
	if s == "" {
		return CurrentPlatform(), nil
	}
	if s == "all" {
		return Platform{All: true}, nil
	}
	parts := strings.Split(s, "/")
	if len(parts) < 2 || len(parts) > 3 || parts[0] == "" || parts[1] == "" {
		return Platform{}, fmt.Errorf("invalid platform %q, expected os/arch, os/arch/variant, or all", s)
	}
	p := Platform{OS: parts[0], Architecture: normalizeArch(parts[1])}
	if len(parts) == 3 {
		if parts[2] == "" {
			return Platform{}, fmt.Errorf("invalid platform %q, variant is empty", s)
		}
		p.Variant = parts[2]
	}
	return p, nil
}

func CurrentPlatform() Platform {
	return Platform{OS: runtime.GOOS, Architecture: normalizeArch(runtime.GOARCH), Variant: currentVariant()}
}

func normalizeArch(arch string) string {
	switch arch {
	case "aarch64":
		return "arm64"
	case "x86_64":
		return "amd64"
	default:
		return arch
	}
}

func currentVariant() string {
	if runtime.GOARCH != "arm" {
		return ""
	}
	return "v7"
}

func (p Platform) String() string {
	if p.All {
		return "all"
	}
	if p.Variant != "" {
		return p.OS + "/" + p.Architecture + "/" + p.Variant
	}
	return p.OS + "/" + p.Architecture
}
