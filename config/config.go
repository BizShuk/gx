// Package config 載入 gx 的設定到 viper，設定根目錄由 gosdk 固定於 ~/.config/gx/。
//
// 設定採扁平 key（見 default_settings.json），各領域套件自行以 viper.Get* 取用，
// 不再經由一份全域 struct 轉手。
package config

import (
	_ "embed"

	"github.com/bizshuk/gosdk/config"
)

//go:embed default_settings.json
var defaultSettingsJSON string

// Default 初始化設定；settings.json 不存在時以內嵌預設值建立。
func Default() {
	config.Default(
		config.WithAppName("gx"),
		config.WithDefaultValue(defaultSettingsJSON),
	)
}
