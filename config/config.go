// Package config 載入 gx 的設定到 viper，設定根目錄由 gosdk 固定於 ~/.config/gx/。
//
// 設定採扁平 key（見 config/example/settings.example.json），各領域套件自行以 viper.Get* 取用，
// 不再經由一份全域 struct 轉手。
package config

import (
	_ "embed"

	gosdkcmdconfig "github.com/bizshuk/gosdk/cmd/config"
	"github.com/bizshuk/gosdk/config"
)

//go:embed example/settings.example.json
var defaultSettingsJSON string

// init 登記內嵌範例, 供 `config default` 複製到設定目錄.
func init() {
	gosdkcmdconfig.MustRegisterDefault("settings.json", []byte(defaultSettingsJSON))
}

// Default 初始化設定；內嵌範例同時作為最低層預設值。
func Default() {
	config.Default(
		config.WithAppName("gx"),
		config.WithDefaultValue(defaultSettingsJSON),
	)
}
