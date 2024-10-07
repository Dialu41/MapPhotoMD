package ui

import (
	"net/url"

	"fyne.io/fyne/v2"
)

// ShowHelp 显示使用说明
func ShowHelp(ap fyne.App) {
	u, _ := url.Parse("https://www.huangoo.top/index.php/archives/139/")
	_ = ap.OpenURL(u)
}
