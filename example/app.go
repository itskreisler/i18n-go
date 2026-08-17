package main

import (
	"fmt"

	ki18n "github.com/itskreisler/i18n-go"
	"github.com/itskreisler/i18n-go/example/locales"
)

func main() {
	loc, err := ki18n.New(locales.FS, "es")
	if err != nil {
		panic(err)
	}

	fmt.Println(loc.T(locales.Hello))
	fmt.Println(loc.T(locales.Welcome))
	fmt.Println(loc.TD(locales.WelcomeUser, map[string]any{
		"Name": "Kreisler",
	}))
	fmt.Println(loc.TP(locales.Items, 3, map[string]any{
		"Count": 3,
	}))
	fmt.Println(loc.T(locales.SettingsTitle))
}
