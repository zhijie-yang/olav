package tui

import (
	"strings"
	"testing"

	"github.com/canonical/olav/internal/oci"
	"github.com/charmbracelet/x/ansi"
)

func TestViewDoesNotExceedWindow(t *testing.T) {
	root := &oci.Node{Name: "/", Path: "/", IsDir: true}
	longJSON := []byte(`{"this-is-a-very-long-key-that-would-wrap-without-truncation":"this is a very long value that would otherwise force the pane to exceed its allocated dimensions and scroll the terminal","items":[1,2,3,4,5,6,7,8,9,10]}`)
	index := &oci.Node{Name: "index.json", Path: "/index.json", Data: longJSON, Parent: root}
	root.Children = []*oci.Node{index}
	layout := &oci.Layout{InputPath: strings.Repeat("very-long-input-path/", 20), Root: root, Files: map[string]*oci.Node{"/": root, "/index.json": index}}

	m := New(layout)
	m.width = 72
	m.height = 20
	m.selectOCI(1)
	if m.preview == nil {
		t.Fatal("expected preview")
	}
	m.preview.ScrollBy(100, m.previewHeight())

	view := m.View()
	lines := strings.Split(view, "\n")
	if len(lines) > m.height {
		t.Fatalf("view height = %d, want <= %d\n%s", len(lines), m.height, view)
	}
	for i, line := range lines {
		if w := ansi.StringWidth(line); w > m.width {
			t.Fatalf("line %d width = %d, want <= %d: %q", i+1, w, m.width, line)
		}
	}
}
