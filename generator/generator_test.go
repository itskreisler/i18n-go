package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExportedName(t *testing.T) {
	tests := map[string]string{
		"hello":          "Hello",
		"settings.title": "SettingsTitle",
		"user-name":      "UserName",
		"123test":        "Key123test",
	}

	for input, want := range tests {
		if got := ExportedName(input); got != want {
			t.Fatalf("%q: got %q, want %q", input, got, want)
		}
	}
}

// TestParseDirPluralNamedNamespace is a regression test: a namespace whose
// sub-keys are plural categories (e.g. [ram.zero]) must not be mistaken for a
// message itself, which used to produce a spurious "ram" key.
func TestParseDirPluralNamedNamespace(t *testing.T) {
	dir := writeLocales(t, map[string]string{
		"active.en.toml": `[hello]
other = "Hello!"

[ram.zero]
other = "0 MB"

[ram.all_off]
other = "all sub-options OFF"

[settings.title]
other = "Settings"
`,
	})

	messages, err := ParseDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	ids := map[string]bool{}
	for _, m := range messages {
		ids[m.ID] = true
	}

	if ids["ram"] {
		t.Error("ParseDir produced a spurious root key for a namespace with a plural-named sub-key")
	}
	for _, want := range []string{"hello", "ram.zero", "ram.all_off", "settings.title"} {
		if !ids[want] {
			t.Errorf("ParseDir missing expected key %q", want)
		}
	}
}

// TestParseDirIgnoresNonActiveToml ensures a stray non-active TOML (e.g. a
// config file) never pollutes the generated keys.
func TestParseDirIgnoresNonActiveToml(t *testing.T) {
	dir := writeLocales(t, map[string]string{
		"active.en.toml": `[hello]
other = "Hello!"
`,
		"config.toml": `[db]
other = "not a locale"
`,
	})

	messages, err := ParseDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(messages) != 1 || messages[0].ID != "hello" {
		t.Errorf("ParseDir = %v, want only [hello]", messages)
	}
}

func writeLocales(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}
