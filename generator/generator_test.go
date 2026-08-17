package generator

import "testing"

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
