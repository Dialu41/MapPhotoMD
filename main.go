package main

import (
	"encoding/json"
	"io"
	"log"
	"net/url"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	xWidget "fyne.io/x/fyne/widget"
)

// config 从配置文件中读取到的数据
var config map[string]interface{}

// appVersion 软件版本号
const appVersion = "v1.0"

// githubURL 项目github仓库链接
const githubURL = ""

func main() {
	ap := app.NewWithID("MapPhotoMD")
	win := ap.NewWindow("MapPhotoMD")

	config = make(map[string]interface{})
	readConfig(ap) //读取配置文件

	win.SetMaster()
	win.Resize(fyne.NewSize(640, 460))
	win.SetMainMenu(makeMenu(ap, win)) //设置菜单栏
	win.SetContent(makeTabs())         //设置各选项卡的内容
	win.CenterOnScreen()               //主窗口居中显示

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

// readConfig 用于读取配置文件，如成功则将数据保存到全局变量，否则发送错误提醒
func readConfig(ap fyne.App) {
	//打开配置文件
	file, err := os.Open("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			ap.SendNotification(&fyne.Notification{
				Title:   "提示",
				Content: "未找到配置文件，请先在设置中完成相关设置并保存",
			})
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
	err = json.Unmarshal(data, &config)
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
	var photoPathEntry *widget.Entry
	gdKeyEntry := widget.NewPasswordEntry()
	keyValue, ok := config["Key"]
	if !ok {
		keyValue = ""
	}
	gdKeyEntry.SetText(keyValue.(string))

	photoPathEntry = widget.NewEntry()
	photoPathEntry.Disable()
	photoPathEntry.SetPlaceHolder("默认为 旅行名称/pictures")
	photoPathEntry.Resize(fyne.NewSize(100, photoPathEntry.MinSize().Height))
	photoPath := container.NewHBox(
		photoPathEntry,
	)
	photoPath.Resize(fyne.NewSize(200, photoPath.MinSize().Height))

	movePhotoRadio := widget.NewRadioGroup([]string{"是", "否"}, func(s string) {
		if s == "是" {
			photoPathEntry.Enable()
			config["isMovePhoto"] = "yes"
		} else {
			photoPathEntry.Disable()
			config["isMovePhoto"] = "no"
		}
	})
	moveRadioValue, ok := config["isMovePhoto"]
	if !ok {
		movePhotoRadio.SetSelected("否")
	}
	switch moveRadioValue {
	case "yes":
		movePhotoRadio.SetSelected("是")
	case "no":
		movePhotoRadio.SetSelected("否")
	}

	items := []*widget.FormItem{
		widget.NewFormItem("高德Key", gdKeyEntry),
		widget.NewFormItem("是否转存照片", movePhotoRadio),
		widget.NewFormItem("照片转存路径", photoPath),
	}

	settingDialog := dialog.NewForm("设置", "保存", "取消", items, func(b bool) {
		//用户选择取消，则直接返回
		if !b {
			return
		}
		//用户选择保存，则保存输入的Key
		config["Key"] = gdKeyEntry.Text
		jsonData, err := json.Marshal(config)
		if err != nil {
			ap.SendNotification(&fyne.Notification{
				Title:   "错误",
				Content: "序列化配置文件时出错，请联系开发者",
			})
			return
		}
		err = os.WriteFile("config.json", jsonData, 0644)
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

// makeTabs 创建主窗口的选项卡
func makeTabs() *container.AppTabs {
	var tabs *container.AppTabs
	var (
		travelTab     *container.TabItem
		IOputTab      *container.TabItem
		propertiesTab *container.TabItem
	)

	/********设置旅行信息选项卡*********/
	travelName := widget.NewEntry()
	travelName.SetPlaceHolder("例：故宫一日游")
	travelDate := widget.NewEntry()
	travelDate.SetPlaceHolder("点击日历，选择旅行开始的第一天")
	datePicker := xWidget.NewCalendar(time.Now(), func(t time.Time) {
		travelDate.SetText(t.Format("2006-01-02"))
	})
	travelTabContent := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "旅行名称", Widget: travelName},
			{Text: "旅行日期", Widget: travelDate},
			{Text: "", Widget: datePicker},
		},
		OnSubmit: func() {
			tabs.Select(IOputTab)
		},
		SubmitText: "下一步",
	}

	/*********设置导入导出选项卡********/
	IOputTabContent := &widget.Form{}

	propertiesTabContent := &widget.Form{}

	travelTab = container.NewTabItem("旅行信息", travelTabContent)
	IOputTab = container.NewTabItem("导入导出设置", IOputTabContent)
	propertiesTab = container.NewTabItem("添加属性", propertiesTabContent)
	tabs = container.NewAppTabs(
		travelTab,
		IOputTab,
		propertiesTab,
	)
	tabs.SetTabLocation(container.TabLocationLeading)
	return tabs
}

func debugPrintf(str string) {
	log.Printf("%s", str)
}
