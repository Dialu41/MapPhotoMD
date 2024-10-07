package ui

import (
	"MapPhotoMD/internal/config"
	"MapPhotoMD/mywidget"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	xWidget "fyne.io/x/fyne/widget"
)

// proIndex 保存所有property控件的地址
var proIndex []*mywidget.Property

// makeTabs 创建主窗口的选项卡
func MakeTabs(ap fyne.App, win fyne.Window, cfg *config.UserConfig) *container.AppTabs {
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

	inputPhoto := mywidget.NewFolderOpenWithEntry(nil, "", win)
	outputPath := mywidget.NewFolderOpenWithEntry(nil, "", win)

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
	for _, pro := range cfg.Properties {
		proIndex = append(proIndex, mywidget.NewProperty(type2Name, pro.Type, pro.Name, pro.Value))
		proContainer.Add(proIndex[len(proIndex)-1])
	}

	//点击开始生成旅行记录文件及文件夹，如设置保存属性，则与设置项一并保存到config.json
	proNextButton := widget.NewButton("开始生成", func() {
		//清空之前保存的属性
		cfg.Properties = cfg.Properties[:0]
		//按照用户设置，选择是否保存属性
		if cfg.SaveProperties {
			for _, pIndex := range proIndex {
				proData := pIndex.GetPropertyData()
				if proData.Name != "" { //未指定属性名称的不保存
					cfg.Properties = append(cfg.Properties, proData)
				}
			}
		}
		cfg.SaveConfigFile(ap)
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
