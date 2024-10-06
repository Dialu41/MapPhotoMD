package main

import (
	"MapPhotoMD/mywidget"
	"encoding/json"
	"io"
	"net/url"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	xWidget "fyne.io/x/fyne/widget"
)

// UserConfig 用户配置数据
type UserConfig struct {
	Key            string                   `json:"key"`             //高德key
	MovePhoto      bool                     `json:"move_photo"`      //是否转存照片
	PhotoPath      string                   `json:"photo_path"`      //转存路径
	DeletePhoto    bool                     `json:"delete_Photo"`    //是否删除原照片
	SaveProperties bool                     `json:"save_properties"` //是否保存YAML属性
	Properties     []*mywidget.PropertyData `json:"properties"`      //旅行记录YAML属性
}

var config UserConfig

// proIndex 保存所有property控件的地址
var proIndex []*fyne.Container

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

	win.Resize(fyne.NewSize(640, 440))
	win.SetMaster()
	win.SetMainMenu(makeMenu(ap, win)) //设置菜单栏
	win.SetContent(makeTabs(ap, win))  //设置各选项卡的内容
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

// readConfigFile 用于读取配置文件，如成功则将数据保存到config，否则发送错误提醒
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
	readConfigFile(ap)

	//Key
	gdKeyEntry := widget.NewPasswordEntry()
	gdKeyEntry.SetText(config.Key)          //还原设置
	gdKeyEntry.OnChanged = func(s string) { //自动保存到config变量
		config.Key = s
	}

	//是否删除原照片
	deletePhotoRadio := widget.NewRadioGroup([]string{"是", "否"}, func(s string) {
		if s == "是" { //自动保存到config变量
			config.DeletePhoto = true
		} else {
			config.DeletePhoto = false
		}
	})
	switch config.DeletePhoto { //还原设置
	case true:
		deletePhotoRadio.SetSelected("是")
	case false:
		deletePhotoRadio.SetSelected("否")
	}

	//照片转存路径
	photoPathEntry := widget.NewEntry()
	photoPathEntry.OnChanged = func(s string) { //自动保存到config
		config.PhotoPath = s
	}
	photoPathEntry.Disable() //默认不转存
	photoPathEntry.SetPlaceHolder("默认为 旅行名称/pictures")
	photoPathEntry.SetText(config.PhotoPath) //还原设置

	//点击选定转存路径，并显示在输入框中
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
		//改变转存路径选择控件的状态，保存设置到config
		if s == "是" {
			photoPathEntry.Enable()
			photoPathButton.Enable()
			deletePhotoRadio.Enable()
			config.MovePhoto = true
		} else {
			photoPathEntry.Disable()
			photoPathButton.Disable()
			deletePhotoRadio.Disable()
			config.MovePhoto = false
		}
	})
	switch config.MovePhoto { //还原设置
	case true:
		movePhotoRadio.SetSelected("是")
	case false:
		movePhotoRadio.SetSelected("否")
	}

	//是否保存属性
	savePropertiesRadio := widget.NewRadioGroup([]string{"是", "否"}, func(s string) {
		if s == "是" {
			config.SaveProperties = true
		} else {
			config.SaveProperties = false
		}
	})
	switch config.SaveProperties { //还原设置
	case true:
		savePropertiesRadio.SetSelected("是")
	case false:
		savePropertiesRadio.SetSelected("否")
	}

	items := []*widget.FormItem{
		widget.NewFormItem("高德Key", gdKeyEntry),
		widget.NewFormItem("是否转存照片", movePhotoRadio),
		widget.NewFormItem("照片转存路径", photoPath),
		widget.NewFormItem("是否删除原照片", deletePhotoRadio),
		widget.NewFormItem("是否保存属性", savePropertiesRadio),
	}

	settingDialog := dialog.NewForm("设置", "保存", "取消", items, func(b bool) {
		//用户选择取消，则直接返回
		if !b {
			return
		}
		//用户选择保存，则保存输入的Key
		saveConfigFile(ap)
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
func makeTabs(ap fyne.App, win fyne.Window) *container.AppTabs {
	//tabs 所有选项卡的指针，用于跳转选项卡
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

	//日历，点击日期时将日期赋值给输入框
	datePicker := xWidget.NewCalendar(time.Now(), func(t time.Time) {
		travelDate.SetText(t.Format("2006-01-02"))
	})

	//跳转下一个选项卡
	travelNextButton := widget.NewButton("下一步", func() {
		tabs.Select(IOputTab)
	})
	travelNextButton.Importance = widget.HighImportance

	travelTabContent := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("旅行名称", travelName),
			widget.NewFormItem("旅行日期", travelDate),
			widget.NewFormItem("", datePicker)),
		//保持跳转按钮靠下
		layout.NewSpacer(),
		//保持跳转按钮居中
		container.NewHBox(
			layout.NewSpacer(),
			travelNextButton,
			layout.NewSpacer(),
		),
	)

	/*********设置导入导出选项卡********/
	inputPhotoEntry := widget.NewEntry()

	inputPhotoButton := widget.NewButton("选择文件夹", func() {
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
			inputPhotoEntry.SetText(list.Path())
		}, win)
	})
	inputPhoto := container.New(&FileSelectLayout{}, inputPhotoEntry, inputPhotoButton)

	outputPathEntry := widget.NewEntry()
	outputPathButton := widget.NewButton("选择文件夹", func() {
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
			outputPathEntry.SetText(list.Path())
		}, win)
	})
	outputPath := container.New(&FileSelectLayout{}, outputPathEntry, outputPathButton)
	IOputNextButton := widget.NewButton("下一步", func() {
		tabs.Select(propertiesTab)
	})
	IOputNextButton.Importance = widget.HighImportance
	IOputBackButton := widget.NewButton("上一步", func() {
		tabs.Select(travelTab)
	})
	IOputTabContent := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("导入照片", inputPhoto),
			widget.NewFormItem("导出到Ob库", outputPath)),
		//保持按钮靠下
		layout.NewSpacer(),
		//保持按钮居中
		container.NewHBox(
			layout.NewSpacer(),
			IOputBackButton,
			IOputNextButton,
			layout.NewSpacer(),
		),
	)

	/*********设置添加属性选项卡********/
	//属性类型及其对应的默认属性名称
	type2Name := map[string]string{
		"标签":  "tags",
		"别名":  "aliases",
		"样式":  "cssclasses",
		"文本":  "",
		"列表":  "",
		"数字":  "",
		"复选框": "",
		"日期":  "",
	}
	//默认属性类型
	defaultType := "文本"

	//所有属性控件纵向排列
	proContainer := container.NewVBox()

	//还原保存的属性设置
	for _, pro := range config.Properties {
		proIndex = append(proIndex, mywidget.NewProperty(type2Name, pro.Type, pro.Name, pro.Value))
		proContainer.Add(proIndex[len(proIndex)-1])
	}

	//点击开始生成旅行记录文件及文件夹，如设置保存属性，则与设置项一并保存到config.json
	proNextButton := widget.NewButton("开始生成", func() {
		//读取已保存的配置
		readConfigFile(ap)
		//按照用户设置，选择是否保存属性
		if config.SaveProperties {
			config.Properties = config.Properties[:0] //清空已保存的属性，以存入新属性
			for _, pIndex := range proIndex {
				proData := mywidget.GetPropertyData(pIndex)
				if proData.Name != "" { //未指定属性名称的不保存
					config.Properties = append(config.Properties, proData)
				}
			}
			saveConfigFile(ap)
		}
	})
	proNextButton.Importance = widget.DangerImportance

	//跳转上一个选项卡
	proBackButton := widget.NewButton("上一步", func() {
		tabs.Select(IOputTab)
	})

	//点击添加一条属性
	addProButton := widget.NewButton("添加属性", func() {
		proIndex = append(proIndex, mywidget.NewProperty(type2Name, defaultType, "", ""))
		proContainer.Add(proIndex[len(proIndex)-1])
	})
	addProButton.Importance = widget.HighImportance

	//点击删除一条属性，少于一条时无效
	deleteProButton := widget.NewButton("删除属性", func() {
		length := len(proIndex)
		if length > 0 {
			proContainer.Remove(proIndex[length-1])
			proIndex = proIndex[:length-1]
		}
	})

	propertiesTabContent := container.NewVBox(
		proContainer,
		//增删属性按钮居中
		container.NewHBox(
			layout.NewSpacer(),
			deleteProButton,
			addProButton,
			layout.NewSpacer(),
		),
		//跳转及生成按钮中间靠下
		layout.NewSpacer(),
		container.NewHBox(
			layout.NewSpacer(),
			proBackButton,
			proNextButton,
			layout.NewSpacer(),
		),
	)

	travelTab = container.NewTabItem("旅行信息", travelTabContent)
	IOputTab = container.NewTabItem("导入导出设置", IOputTabContent)
	propertiesTab = container.NewTabItem("添加属性", propertiesTabContent)
	tabs = container.NewAppTabs(
		travelTab,
		IOputTab,
		propertiesTab,
	)

	//选项卡靠左
	tabs.SetTabLocation(container.TabLocationLeading)

	return tabs
}

// saveConfigFile 保存设置到config.json
func saveConfigFile(ap fyne.App) {
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
}
