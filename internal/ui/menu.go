package ui

import (
	"MapPhotoMD/internal/config"

	"fyne.io/fyne/v2"
)

// makeMenu 用于创建菜单栏
func MakeMenu(ap fyne.App, win fyne.Window, cfg *config.UserConfig) *fyne.MainMenu {
	settingItem := fyne.NewMenuItem("设置", func() {
		showSettings(ap, win, cfg)
	})
	helpItem := fyne.NewMenuItem("使用说明", func() {
		ShowHelp(ap)
	})
	aboutItem := fyne.NewMenuItem("关于", func() {
		showAbout(ap, win)
	})

	//添加菜单项到菜单栏
	options := fyne.NewMenu("选项",
		settingItem, //设置
		helpItem,    //使用说明
		aboutItem,   //关于
	)
	mainMenu := fyne.NewMainMenu(options)

	return mainMenu
}
