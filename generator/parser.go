package generator

import (
    "fmt"
    "os"
    "path/filepath"
    "sort"
    "strings"

    "github.com/BurntSushi/toml"
)

type Message struct {
    ID          string
    Description string
}

func ParseDir(dir string) ([]Message, error) {
    entries, err := os.ReadDir(dir)
    if err != nil {
        return nil, err
    }

    seen := map[string]bool{}
    var messages []Message

    for _, entry := range entries {
        if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".toml") {
            continue
        }

        path := filepath.Join(dir, entry.Name())
        var tree map[string]any
        if _, err := toml.DecodeFile(path, &tree); err != nil {
            return nil, fmt.Errorf("%s: %w", path, err)
        }

        for key, value := range tree {
            flatten(key, value, "", &messages, seen)
        }
    }

    sort.Slice(messages, func(i, j int) bool {
        return messages[i].ID < messages[j].ID
    })
    return messages, nil
}

func flatten(key string, value any, prefix string, out *[]Message, seen map[string]bool) {
    id := key
    if prefix != "" {
        id = prefix + "." + key
    }

    table, ok := value.(map[string]any)
    if !ok {
        addMessage(id, "", out, seen)
        return
    }

    // go-i18n message tables use fields such as "other", "one",
    // "zero", "two", "few", "many", and optionally "description".
    // The table itself is the message, not another namespace level.
    if isMessageTable(table) {
        description, _ := table["description"].(string)
        addMessage(id, description, out, seen)
        return
    }

    for child, childValue := range table {
        flatten(child, childValue, id, out, seen)
    }
}

func isMessageTable(table map[string]any) bool {
    // A message table's keys are only plural categories (plus an optional
    // description). Requiring ALL keys to be categories - not just one -
    // stops a namespace like [ram.zero] from being mistaken for a message
    // just because one of its sub-keys is named "zero".
    for key := range table {
        switch key {
        case "description", "zero", "one", "two", "few", "many", "other":
        default:
            return false
        }
    }
    return len(table) > 0
}

func addMessage(id, description string, out *[]Message, seen map[string]bool) {
    if seen[id] {
        return
    }
    seen[id] = true
    *out = append(*out, Message{ID: id, Description: description})
}
