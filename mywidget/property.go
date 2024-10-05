package mywidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type PropertyLayout struct{}

func (lo *PropertyLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, obj := range objects {
		minSize = minSize.Max(obj.MinSize())
	}
	return minSize
}

func (lo *PropertyLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) != 3 {
		return
	}
	proType := objects[0]
	proName := objects[1]
	proValue := objects[2]

	proType.Resize(fyne.NewSize(120, proName.MinSize().Height))
	proType.Move(fyne.NewPos(5, 0))

	proName.Resize(fyne.NewSize((containerSize.Width-145)/3, proName.MinSize().Height))
	proName.Move(fyne.NewPos(135, 0))

	proValue.Resize(fyne.NewSize((containerSize.Width-140)*2/3, proValue.MinSize().Height))
	proValue.Move(fyne.NewPos(140+(containerSize.Width-140)/3, 0))
}

// NewProperty 创建组合控件Property，用于输入旅行记录的单条YAML属性。
// 传入参数：type2Name 属性类型及其对应的默认属性名称；其余为控件的占位符，未输入时显示。
// 布局：水平排列。属性类型控件固定宽度，属性名称控件占剩余宽度的1/3，属性值控件占剩余宽度2/3。控件间固定间隔10，控件到边缘固定距离5
func NewProperty(type2Name map[string]string, typePlaceHolder string, namePlaceHolder string, valuePlaceHolder string) *fyne.Container {
	proName := widget.NewEntry()
	proName.SetPlaceHolder(namePlaceHolder)

	var types []string
	for t := range type2Name {
		types = append(types, t)
	}
	proType := widget.NewSelect(types, func(s string) {
		proName.SetText(type2Name[s])

	})
	proType.Selected = typePlaceHolder

	proValue := widget.NewEntry()
	proValue.SetPlaceHolder(valuePlaceHolder)

	return container.New(&PropertyLayout{}, proType, proName, proValue)

}
