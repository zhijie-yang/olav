package source

import "testing"

func TestParsePlatform(t *testing.T) {
	tests := []struct {
		input string
		want  Platform
	}{
		{input: "all", want: Platform{All: true}},
		{input: "linux/amd64", want: Platform{OS: "linux", Architecture: "amd64"}},
		{input: "linux/arm/v7", want: Platform{OS: "linux", Architecture: "arm", Variant: "v7"}},
		{input: "linux/aarch64", want: Platform{OS: "linux", Architecture: "arm64"}},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParsePlatform(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestParsePlatformRejectsInvalid(t *testing.T) {
	for _, input := range []string{"linux", "linux/", "linux/amd64/"} {
		if _, err := ParsePlatform(input); err == nil {
			t.Fatalf("expected %q to fail", input)
		}
	}
}
