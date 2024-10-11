package ui

import (
	"MapPhotoMD/internal/constants"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// showAbout 显示关于界面
func showAbout(ap fyne.App, win fyne.Window) {
	contactButton := widget.NewButton("我的邮箱", func() {
		win.Clipboard().SetContent("1165011707@qq.com")
		ap.SendNotification(&fyne.Notification{
			Title:   "提示",
			Content: "已复制邮箱地址",
		})
	})
	blogButton := widget.NewButton("我的博客", func() {
		u, _ := url.Parse("https://www.huangoo.top")
		_ = ap.OpenURL(u)
	})
	githubButton := widget.NewButton("跳转仓库", func() {
		u, _ := url.Parse(constants.GithubURL)
		_ = ap.OpenURL(u)
	})

	content := container.NewGridWithColumns(2,
		widget.NewLabel("版本号"), widget.NewLabel(constants.AppVersion),
		widget.NewLabel("版权信息"), widget.NewLabel("Copyright © 2024 黄嚄嚄."),
		widget.NewLabel("联系开发者"), container.NewHBox(contactButton, blogButton),
		widget.NewLabel("开源协议"), container.NewHBox(widget.NewLabel("本软件使用MIT协议发行"), githubButton),
		widget.NewLabel("鸣谢"), container.NewHBox(
			widget.NewButton("StarAire", func() {
				u, _ := url.Parse(("https://sspai.com/post/80578"))
				_ = ap.OpenURL(u)
			}),
			widget.NewButton("Go", func() {
				u, _ := url.Parse(("https://github.com/golang/go"))
				_ = ap.OpenURL(u)
			}),
			widget.NewButton("Fyne", func() {
				u, _ := url.Parse("https://github.com/fyne-io/fyne")
				_ = ap.OpenURL(u)
			}),
			widget.NewButton("Go-EXIF", func() {
				u, _ := url.Parse("https://github.com/rwcarlsen/goexif")
				_ = ap.OpenURL(u)
			}),
		),
	)
	dialog.ShowCustom("关于", "关闭", content, win)
}
