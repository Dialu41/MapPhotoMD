package main

import (
	"encoding/json"
	"io"
	"log"
	"net/url"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// userData 从配置文件中读取到的数据
var userData map[string]interface{}

// appVersion 软件版本号
const appVersion = "v1.0"

// githubURL 项目github仓库链接
const githubURL = ""

func main() {
	ap := app.NewWithID("MapPhotoMD")
	win := ap.NewWindow("MapPhotoMD")

	win.SetMainMenu(makeMenu(ap, win))
	win.SetMaster()
	win.Resize(fyne.NewSize(640, 460))

	readUserData(ap, true)

	win.ShowAndRun()
}

// makeMenu 用于创建菜单栏
func makeMenu(ap fyne.App, win fyne.Window) *fyne.MainMenu {
	settingItem := fyne.NewMenuItem("设置", func() {
		showSettings(ap, win)
	})
	helpItem := fyne.NewMenuItem("使用说明", func() {
		showHelp(ap)
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

// readUserData 用于读取配置文件，如成功则将数据保存到全局变量，否则发送错误提醒
func readUserData(ap fyne.App, isWinStart bool) {
	//打开配置文件
	file, err := os.Open("userData.json")
	if err != nil {
		if os.IsNotExist(err) {
			if isWinStart {
				ap.SendNotification(&fyne.Notification{
					Title:   "提示",
					Content: "未找到配置文件，请先在设置中填写Key并保存",
				})
			}
			return
		} else {
			ap.SendNotification(&fyne.Notification{
				Title:   "错误",
				Content: "打开配置文件时出错，请联系开发者",
			})
			return
		}
	}
	defer file.Close()
	//读取配置文件
	data, err := io.ReadAll(file)
	if err != nil {
		ap.SendNotification(&fyne.Notification{
			Title:   "错误",
			Content: "读取配置文件时出错，请联系开发者",
		})
		return
	}
	//解析为JSON
	err = json.Unmarshal(data, &userData)
	if err != nil {
		ap.SendNotification(&fyne.Notification{
			Title:   "错误",
			Content: "解析JSON时出错，请联系开发者",
		})
		return
	}
}

// showSettings 显示设置
func showSettings(ap fyne.App, win fyne.Window) {
	//Key文本框
	gdKeyEntry := widget.NewPasswordEntry()
	readUserData(ap, false)
	Key, ok := userData["Key"]
	if ok {
		gdKeyEntry.SetText(Key.(string))
	}

	items := []*widget.FormItem{
		widget.NewFormItem("高德Key:", gdKeyEntry),
	}

	settingDialog := dialog.NewForm("设置", "保存", "取消", items, func(b bool) {
		//用户选择取消，则直接返回
		if !b {
			return
		}
		//用户选择保存，则保存输入的Key
		key := map[string]string{"Key": gdKeyEntry.Text}
		jsonData, err := json.Marshal(key)
		if err != nil {
			ap.SendNotification(&fyne.Notification{
				Title:   "错误",
				Content: "序列化配置文件时出错，请联系开发者",
			})
			return
		}
		err = os.WriteFile("userData.json", jsonData, 0644)
		if err != nil {
			ap.SendNotification(&fyne.Notification{
				Title:   "错误",
				Content: "保存配置文件时出错，请联系开发者",
			})
			return
		}
	}, win)
	settingDialog.Resize(fyne.NewSize(500, settingDialog.MinSize().Height))
	settingDialog.Show()
}

// showHelp 显示使用说明
func showHelp(ap fyne.App) {
	u, _ := url.Parse("https://www.huangoo.top/index.php/archives/139/")
	_ = ap.OpenURL(u)
}

// showAbout 显示关于
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
		u, _ := url.Parse(githubURL)
		_ = ap.OpenURL(u)
	})
	content := container.NewGridWithColumns(2,
		widget.NewLabel("版本号"), widget.NewLabel(appVersion),
		widget.NewLabel("版权信息"), widget.NewLabel("Copyright © 2024 黄嚄嚄."),
		widget.NewLabel("联系开发者"), container.NewHBox(contactButton, blogButton),
		widget.NewLabel("开源协议"), container.NewHBox(widget.NewLabel("本软件使用MIT协议"), githubButton),
		widget.NewLabel("鸣谢"), container.NewGridWithColumns(2,
			widget.NewButton("Go", func() {
				u, _ := url.Parse(("https://github.com/golang/go"))
				_ = ap.OpenURL(u)
			}),
			widget.NewButton("Fyne", func() {
				u, _ := url.Parse("https://github.com/fyne-io/fyne")
				_ = ap.OpenURL(u)
			})),
	)
	dialog.ShowCustom("关于", "关闭", content, win)
}

func debugPrintf(str string) {
	log.Printf("%s", str)
}
