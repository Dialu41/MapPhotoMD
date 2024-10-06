package mywidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Property struct {
	widget.BaseWidget
	proContainer *fyne.Container
}

// PropertyData 属性控件的数据结构体
type PropertyData struct {
	Type  string `json:"type"`  //属性类型
	Name  string `json:"name"`  //属性名称
	Value string `json:"value"` //属性值
}

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
// 传入参数：type2Name 属性类型及其对应的默认属性名称；其余为子控件的默认值。
// 布局：水平排列。属性类型控件固定宽度，属性名称控件占剩余宽度的1/3，属性值控件占剩余宽度2/3。控件间固定间隔10，控件到边缘固定距离5。
func NewProperty(type2Name map[string]string, defaultType string, defaultName string, defaultValue string) *Property {
	proName := widget.NewEntry()
	proName.SetPlaceHolder("属性名称")
	proName.SetText(defaultName)

	var types []string
	for t := range type2Name {
		//提取出属性类型
		types = append(types, t)
	}
	proType := widget.NewSelect(types, func(s string) {
		//选定属性类型时，自动更改属性名称，默认属性名称为空时不更改
		if type2Name[s] != "" {
			proName.SetText(type2Name[s])
		}
	})
	//设置默认属性类型
	proType.SetSelected(defaultType)

	proValue := widget.NewEntry()
	proValue.SetPlaceHolder("属性值")
	proValue.SetText(defaultValue)

	t := &Property{}
	t.ExtendBaseWidget(t)
	t.proContainer = container.New(&PropertyLayout{}, proType, proName, proValue)

	return t

}

func (t *Property) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(t.proContainer)
}

// GetPropertyData 获取传入属性控件的子控件的值，以结构体指针作为返回
func (t *Property) GetPropertyData() *PropertyData {
	data := PropertyData{}
	for i, obj := range t.proContainer.Objects {
		switch v := obj.(type) {
		case *widget.Select: //获取属性类型
			data.Type = v.Selected
		case *widget.Entry:
			switch i {
			case 1: //获取属性名称
				data.Name = v.Text
			case 2: //获取属性值
				data.Value = v.Text
			}
		}
	}
	return &data
}
