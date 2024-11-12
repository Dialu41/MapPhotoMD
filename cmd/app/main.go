package main

import (
	"MapPhotoMD/internal/config"
	"MapPhotoMD/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	ap := app.NewWithID("MapPhotoMD")
	win := ap.NewWindow("MapPhotoMD")

	cfg := config.NewUserConfig()

	cfg.ReadConfigFile(ap) //读取配置文件

	win.Resize(fyne.NewSize(600, 450))
	win.SetMaster()
	win.SetMainMenu(ui.MakeMenu(ap, win, cfg)) //设置菜单栏
	win.SetContent(ui.MakeTabs(ap, win, cfg))  //设置各选项卡的内容
	win.CenterOnScreen()                       //主窗口居中显示

	win.ShowAndRun()
}
