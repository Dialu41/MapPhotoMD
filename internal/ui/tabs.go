package ui

import (
	"MapPhotoMD/internal/config"
	"MapPhotoMD/internal/service"
	"MapPhotoMD/mywidget"
	"errors"
	"regexp"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	xWidget "fyne.io/x/fyne/widget"
)

// tabs 所有选项卡的指针，用于跳转选项卡
var (
	travelData    *service.TravelData
	tabs          *container.AppTabs
	travelTab     *container.TabItem
	IOputTab      *container.TabItem
	propertiesTab *container.TabItem
)

var (
	travelName *widget.Entry
	travelDate *widget.Entry
	inputPath  *mywidget.FolderOpenWithEntry
	outputPath *mywidget.FolderOpenWithEntry
)

// makeTabs 创建主窗口的选项卡
func MakeTabs(ap fyne.App, win fyne.Window, cfg *config.UserConfig) *container.AppTabs {
	travelData = service.NewTravelData()

	travelTab = container.NewTabItem("旅行信息", makeTravelTabContent(win))
	IOputTab = container.NewTabItem("导入导出设置", makeIOputTabContent(win))
	propertiesTab = container.NewTabItem("添加属性", makePropertiesTabContent(ap, win, cfg))
	tabs = container.NewAppTabs(
		travelTab,
		IOputTab,
		propertiesTab,
	)

	//选项卡靠左
	tabs.SetTabLocation(container.TabLocationLeading)

	return tabs
}

// makeTravelTabContent 创建旅行名称和旅行日期输入选项卡的内容
func makeTravelTabContent(win fyne.Window) *fyne.Container {
	travelName = widget.NewEntry()
	travelName.OnChanged = func(s string) {
		travelData.TravelName = s
	}
	travelName.SetPlaceHolder("例：故宫一日游")
	travelName.Validator = func(s string) error {
		if s == "" {
			return errors.New("")
		} else {
			return nil
		}
	}

	travelDate = widget.NewEntry()
	travelDate.OnChanged = func(s string) {
		travelData.TravelDate = s
	}
	travelDate.SetPlaceHolder("点击日历，选择旅行开始的第一天")
	travelDate.Validator = func(s string) error {
		pat := "^\\d{4}-\\d{2}-\\d{2}$"
		re := regexp.MustCompile(pat)
		_, err := time.Parse("2006-01-02", s)
		if re.MatchString(s) && err == nil {
			return nil
		}
		return errors.New("")
	}

	//日历，点击日期时将日期赋值给输入框
	datePicker := xWidget.NewCalendar(time.Now(), func(t time.Time) {
		travelDate.SetText(t.Format("2006-01-02"))
	})

	//跳转下一个选项卡
	travelNextButton := widget.NewButton("下一步", func() {
		if travelDate.Validate() == nil && travelName.Validate() == nil {
			tabs.Select(IOputTab)
		} else {
			dialog.ShowError(errors.New("旅行信息未填写或格式错误"), win)
		}
	})
	travelNextButton.Importance = widget.HighImportance

	return container.NewVBox(
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
}

// makeIOputTabContent 创建导入导出选项卡的内容
func makeIOputTabContent(win fyne.Window) *fyne.Container {
	inputPath = mywidget.NewFolderOpenWithEntry(func(s string) {
		travelData.InputPath = s
	}, "", win)
	outputPath = mywidget.NewFolderOpenWithEntry(func(s string) {
		travelData.OutputPath = s
	}, "", win)

	IOputNextButton := widget.NewButton("下一步", func() {
		if inputPath.GetValid() && outputPath.GetValid() {
			tabs.Select(propertiesTab)
		} else {
			dialog.ShowError(errors.New("导入导出路径不存在或未填写"), win)
		}
	})
	IOputNextButton.Importance = widget.HighImportance

	IOputBackButton := widget.NewButton("上一步", func() {
		tabs.Select(travelTab)
	})

	return container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("导入照片", inputPath),
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
}

func makePropertiesTabContent(ap fyne.App, win fyne.Window, cfg *config.UserConfig) *fyne.Container {
	//选择的属性类型
	types := []string{
		mywidget.ProType_Tag,
		mywidget.ProType_Aliases,
		mywidget.ProType_css,
		mywidget.ProType_Text,
		mywidget.ProType_List,
		mywidget.ProType_Num,
		mywidget.ProType_Check,
		mywidget.ProType_Date,
	}

	//默认属性类型
	defaultType := "文本"

	//所有属性控件纵向排列
	proContainer := container.NewVBox()

	//还原保存的属性设置
	for _, pro := range cfg.Properties {
		travelData.ProIndex = append(travelData.ProIndex, mywidget.NewProperty(types, pro.Type, pro.Name, pro.Value))
		proContainer.Add(travelData.ProIndex[len(travelData.ProIndex)-1])
	}

	//点击开始生成旅行记录文件及文件夹，如设置保存属性，则与设置项一并保存到config.json
	proNextButton := widget.NewButton("开始生成", func() {
		if travelName.Validate() != nil || travelDate.Validate() != nil || !inputPath.GetValid() || !outputPath.GetValid() {
			dialog.ShowError(errors.New("旅行信息、导入导出设置错误或未填写"), win)
			return
		}
		//检查填写的属性是否正确
		for _, pIndex := range travelData.ProIndex {
			if !pIndex.GetValid() {
				dialog.ShowError(errors.New("属性格式错误，无法进行下一步"), win)
				return
			}
		}
		//清空之前保存的属性
		cfg.Properties = cfg.Properties[:0]
		//按照用户设置，选择是否保存属性
		if cfg.SaveProperties {
			for _, pIndex := range travelData.ProIndex {
				proData := pIndex.GetPropertyData()
				if proData.Name != "" { //未指定属性名称的不保存
					cfg.Properties = append(cfg.Properties, proData)
				}
			}
		}
		cfg.SaveConfigFile(ap)
		invalidPhotos := travelData.GenerateMD(cfg)
		if len(invalidPhotos) != 0 {
			str := "无法转换的照片如下，请检查其是否存在经纬度信息：\n"
			for _, p := range invalidPhotos {
				str = str + p + "\n"
			}
			dialog.ShowInformation("提示", str, win)
		} else {
			dialog.ShowInformation("提示", "生成成功！", win)
		}
	})
	proNextButton.Importance = widget.DangerImportance

	//跳转上一个选项卡
	proBackButton := widget.NewButton("上一步", func() {
		tabs.Select(IOputTab)
	})

	//点击添加一条属性
	addProButton := widget.NewButton("添加属性", func() {
		travelData.ProIndex = append(travelData.ProIndex, mywidget.NewProperty(types, defaultType, "", ""))
		proContainer.Add(travelData.ProIndex[len(travelData.ProIndex)-1])
	})
	addProButton.Importance = widget.HighImportance

	//点击删除一条属性，少于一条时无效
	deleteProButton := widget.NewButton("删除属性", func() {
		length := len(travelData.ProIndex)
		if length > 0 {
			proContainer.Remove(travelData.ProIndex[length-1])
			travelData.ProIndex = travelData.ProIndex[:length-1]
		}
	})

	return container.NewVBox(
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
}
