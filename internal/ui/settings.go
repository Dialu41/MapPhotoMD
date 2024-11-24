package ui

import (
	"MapPhotoMD/internal/config"
	"MapPhotoMD/mywidget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// showSettings 显示设置
func showSettings(ap fyne.App, win fyne.Window, config *config.UserConfig) {
	//读取配置文件
	config.ReadConfigFile(ap)
	//临时保存设置，点击取消则不保存，反之则保存到config和文件中
	temp := struct {
		Key            string
		NotePath       string
		SaveIOpath     bool
		MovePhoto      bool
		PhotoPath      string
		DeletePhoto    bool
		PhotoQuality   int
		SaveProperties bool
	}{
		Key:            config.Key,
		NotePath:       config.NotePath,
		SaveIOpath:     config.SaveIOPath,
		MovePhoto:      config.MovePhoto,
		PhotoPath:      config.PhotoPath,
		DeletePhoto:    config.DeletePhoto,
		PhotoQuality:   config.PhotoQuality,
		SaveProperties: config.SaveProperties,
	}

	//Key
	gdKeyEntry := widget.NewPasswordEntry()
	gdKeyEntry.SetText(config.Key)          //还原设置
	gdKeyEntry.OnChanged = func(s string) { //自动保存
		temp.Key = s
	}

	//ob库路径
	notePathEntry := widget.NewEntry()
	notePathEntry.SetText(config.NotePath) //还原设置
	notePathEntry.OnChanged = func(s string) {
		temp.NotePath = s
	}
	notePathEntry.SetPlaceHolder("旅行记录在Ob库下的路径，例：生活/游记")

	//是否保存导入导出路径
	saveIOpathRadio := widget.NewRadioGroup([]string{"是", "否"}, func(s string) {
		if s == "是" { //自动保存
			temp.SaveIOpath = true
		} else {
			temp.SaveIOpath = false
		}
	})
	saveIOpathRadio.Horizontal = true
	switch config.SaveIOPath { //还原设置
	case true:
		saveIOpathRadio.SetSelected("是")
	case false:
		saveIOpathRadio.SetSelected("否")
	}

	//转存路径
	photoPath := mywidget.NewFolderOpenWithEntry(func(s string) {
		temp.PhotoPath = s
	}, "默认为 旅行名称/pictures", win)
	photoPath.SetEntryText(config.PhotoPath) //还原设置

	//是否删除原照片
	deletePhotoRadio := widget.NewRadioGroup([]string{"是", "否"}, func(s string) {
		if s == "是" { //自动保存
			temp.DeletePhoto = true
		} else {
			temp.DeletePhoto = false
		}
	})
	deletePhotoRadio.Horizontal = true
	switch config.DeletePhoto { //还原设置
	case true:
		deletePhotoRadio.SetSelected("是")
	case false:
		deletePhotoRadio.SetSelected("否")
	}

	//是否转存
	movePhotoRadio := widget.NewRadioGroup([]string{"是", "否"}, func(s string) {
		//改变转存路径选择控件的状态，临时保存设置
		if s == "是" {
			photoPath.Enable()
			deletePhotoRadio.Enable()
			temp.MovePhoto = true
		} else {
			photoPath.Disable()
			deletePhotoRadio.Disable()
			temp.MovePhoto = false
		}
	})
	movePhotoRadio.Horizontal = true
	switch config.MovePhoto { //还原设置
	case true:
		movePhotoRadio.SetSelected("是")
	case false:
		movePhotoRadio.SetSelected("是")
		//fyne的bug，不能在启动时将radio设置为disable
		//需要先enable，再disable，否则文本会一直为灰色
		defer movePhotoRadio.SetSelected("否")
	}

	//照片质量，0~100
	//数值越大照片质量越高，体积越大，100时不压缩、只复制
	photoQData := binding.BindInt(&temp.PhotoQuality)
	photoQLabel := widget.NewLabelWithData(binding.IntToString(photoQData))
	photoQLabel.Resize(photoQLabel.MinSize())
	photoQSlide := widget.NewSliderWithData(0, 100, binding.IntToFloat(photoQData))
	photoQSlide.Step = 1
	photoQualityContent := container.NewAdaptiveGrid(2,
		photoQSlide, photoQLabel,
	)

	//是否保存属性
	savePropertiesRadio := widget.NewRadioGroup([]string{"是", "否"}, func(s string) {
		if s == "是" {
			temp.SaveProperties = true
		} else {
			temp.SaveProperties = false
		}
	})
	savePropertiesRadio.Horizontal = true
	switch config.SaveProperties { //还原设置
	case true:
		savePropertiesRadio.SetSelected("是")
	case false:
		savePropertiesRadio.SetSelected("否")
	}

	items := []*widget.FormItem{
		widget.NewFormItem("高德Key", gdKeyEntry),
		widget.NewFormItem("Ob库路径", notePathEntry),
		widget.NewFormItem("是否保存导入导出设置", saveIOpathRadio),
		widget.NewFormItem("是否转存照片", movePhotoRadio),
		widget.NewFormItem("照片转存路径", photoPath),
		widget.NewFormItem("是否删除原照片", deletePhotoRadio),
		widget.NewFormItem("照片质量", photoQualityContent),
		widget.NewFormItem("是否保存属性", savePropertiesRadio),
	}

	settingDialog := dialog.NewForm("设置", "保存", "取消", items, func(b bool) {
		//用户选择取消，则直接返回
		if !b {
			return
		}
		//检查转存路径是否设置正确
		if temp.MovePhoto && photoPath.GetEntryText() != "" && !photoPath.GetValid() {
			ap.SendNotification(&fyne.Notification{
				Title:   "错误",
				Content: "设置保存失败\n转存路径错误，请重新设置并保存",
			})
			temp.PhotoPath = ""
		}
		//保存设置到config.json
		config.Key = temp.Key
		config.NotePath = temp.NotePath
		config.SaveIOPath = temp.SaveIOpath
		config.MovePhoto = temp.MovePhoto
		config.PhotoPath = temp.PhotoPath
		config.DeletePhoto = temp.DeletePhoto
		config.SaveProperties = temp.SaveProperties
		config.PhotoQuality = temp.PhotoQuality
		config.SaveConfigFile(ap)

	}, win)
	settingDialog.Resize(fyne.NewSize(500, settingDialog.MinSize().Height))
	settingDialog.Show()
}
