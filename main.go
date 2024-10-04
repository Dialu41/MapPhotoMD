package main

import (
	"encoding/json"
	"io"
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

// UserConfig 用户配置数据
type UserConfig struct {
	Key         string `json:"key"`          //高德key
	MovePhoto   bool   `json:"move_photo"`   //是否转存照片
	PhotoPath   string `json:"photo_path"`   //转存路径
	DeletePhoto bool   `json:"delete_Photo"` //是否删除原照片
}

var config UserConfig

// appVersion 软件版本号
const appVersion = "v1.0"

// githubURL 项目github仓库链接
const githubURL = ""

// FileSelectLayout 文件选择器布局。包含一个输入框和一个按钮，横向排布
// 按钮两侧紧贴文本，输入框填充容器剩余空间
// 传入参数时，先输入框再按钮
// 选择文件夹也可使用
type FileSelectLayout struct{}

func main() {
	ap := app.NewWithID("MapPhotoMD")
	win := ap.NewWindow("MapPhotoMD")

	readConfigFile(ap) //读取配置文件

	win.SetMaster()
	win.Resize(fyne.NewSize(640, 460))
	win.SetMainMenu(makeMenu(ap, win)) //设置菜单栏
	win.SetContent(makeTabs())         //设置各选项卡的内容
	win.CenterOnScreen()               //主窗口居中显示

	win.ShowAndRun()
}

func (lo *FileSelectLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, obj := range objects {
		minSize = minSize.Max(obj.MinSize())
	}
	return minSize
}

func (lo *FileSelectLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) != 2 {
		return
	}
	entry := objects[0]
	button := objects[1]

	buttonMinWidth := button.MinSize().Width
	button.Resize(button.MinSize())
	button.Move(fyne.NewPos(containerSize.Width-buttonMinWidth, 0))

	entry.Resize(fyne.NewSize(containerSize.Width-buttonMinWidth-10, entry.MinSize().Height))
	entry.Move(fyne.NewPos(0, 0))
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

// readConfigFile 用于读取配置文件，如成功则将数据保存到cfg，否则发送错误提醒
func readConfigFile(ap fyne.App) {
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

	//Key
	gdKeyEntry := widget.NewPasswordEntry()
	gdKeyEntry.SetText(config.Key)

	//照片转存路径
	photoPathEntry = widget.NewEntry()
	photoPathEntry.Disable()
	photoPathEntry.SetPlaceHolder("默认为 旅行名称/pictures")
	photoPathEntry.SetText(config.PhotoPath)

	photoPathButton := widget.NewButton("打开文件夹", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			//选择文件夹时出错
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			//没有选择
			if list == nil {
				return
			}
			photoPathEntry.SetText(list.Path())
		}, win)
	})
	photoPath := container.New(&FileSelectLayout{},
		photoPathEntry,
		photoPathButton,
	)

	//是否转存
	movePhotoRadio := widget.NewRadioGroup([]string{"是", "否"}, func(s string) {
		if s == "是" {
			photoPathEntry.Enable()
			photoPathButton.Enable()
		} else {
			photoPathEntry.Disable()
			photoPathButton.Disable()
		}
	})
	switch config.MovePhoto {
	case true:
		movePhotoRadio.SetSelected("是")
	case false:
		movePhotoRadio.SetSelected("否")
	}

	//是否删除原照片
	deletePhotoRadio := widget.NewRadioGroup([]string{"是", "否"}, func(s string) {})
	switch config.DeletePhoto {
	case true:
		deletePhotoRadio.SetSelected("是")
	case false:
		deletePhotoRadio.SetSelected("否")
	}

	items := []*widget.FormItem{
		widget.NewFormItem("高德Key", gdKeyEntry),
		widget.NewFormItem("是否转存照片", movePhotoRadio),
		widget.NewFormItem("照片转存路径", photoPath),
		widget.NewFormItem("是否删除原照片", deletePhotoRadio),
	}

	settingDialog := dialog.NewForm("设置", "保存", "取消", items, func(b bool) {
		//用户选择取消，则直接返回
		if !b {
			return
		}
		//用户选择保存，则保存输入的Key
		config.Key = gdKeyEntry.Text
		config.PhotoPath = photoPathEntry.Text
		switch movePhotoRadio.Selected {
		case "是":
			config.MovePhoto = true
		case "否":
			config.MovePhoto = false
		}
		switch deletePhotoRadio.Selected {
		case "是":
			config.DeletePhoto = true
		case "否":
			config.DeletePhoto = false
		}

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
		widget.NewLabel("开源协议"), container.NewHBox(widget.NewLabel("本软件使用MIT协议发行"), githubButton),
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
