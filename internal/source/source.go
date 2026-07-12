package source

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Options struct {
	Input    string
	Platform string
	Progress io.Writer
}

type Resolved struct {
	DisplayName string
	LocalPath   string
	Cached      bool
}

func Resolve(ctx context.Context, opts Options) (*Resolved, error) {
	if opts.Input == "" {
		return nil, errors.New("input is required")
	}
	if strings.HasPrefix(opts.Input, "oci:") {
		if opts.Platform != "" {
			return nil, errors.New("--platform is only supported with docker:// sources")
		}
		p := strings.TrimPrefix(opts.Input, "oci:")
		return resolveLocal(opts.Input, p)
	}
	if strings.HasPrefix(opts.Input, "oci-archive:") {
		if opts.Platform != "" {
			return nil, errors.New("--platform is only supported with docker:// sources")
		}
		p := strings.TrimPrefix(opts.Input, "oci-archive:")
		return resolveLocal(opts.Input, p)
	}
	if _, err := os.Stat(opts.Input); err == nil {
		if opts.Platform != "" {
			return nil, errors.New("--platform is only supported with docker:// sources")
		}
		return resolveLocal(opts.Input, opts.Input)
	}
	if strings.HasPrefix(opts.Input, "docker://") {
		platform, err := ParsePlatform(opts.Platform)
		if err != nil {
			return nil, err
		}
		return copyToCache(ctx, opts.Input, platform, opts.Progress)
	}
	if strings.HasPrefix(opts.Input, "docker-daemon:") {
		if opts.Platform != "" {
			return nil, errors.New("--platform is only supported with docker:// sources")
		}
		return copyToCache(ctx, opts.Input, Platform{}, opts.Progress)
	}
	return nil, fmt.Errorf("input is not a path; image sources must use an explicit transport prefix such as docker:// or docker-daemon:")
}

func resolveLocal(display, p string) (*Resolved, error) {
	if p == "" {
		return nil, errors.New("local path is empty")
	}
	if _, err := os.Stat(p); err != nil {
		return nil, err
	}
	return &Resolved{DisplayName: display, LocalPath: p}, nil
}

func makeTempLayoutDir() (string, error) {
	root, err := cacheRoot()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Join(root, "tmp"), 0o755); err != nil {
		return "", err
	}
	return os.MkdirTemp(filepath.Join(root, "tmp"), "image-*")
}
