package mywidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// 内置的属性类型
const (
	ProType_Tag     = "标签"
	ProType_Aliases = "别名"
	ProType_css     = "样式"
	ProType_Text    = "文本"
	ProType_List    = "列表"
	ProType_Num     = "数字"
	ProType_Check   = "复选框"
	ProType_Date    = "日期"
)

// 属性类型与默认属性名称的映射表
var typeMap = map[string]string{
	ProType_Tag:     "tags",
	ProType_Aliases: "aliases",
	ProType_css:     "cssclasses",
	ProType_Text:    "",
	ProType_List:    "",
	ProType_Num:     "",
	ProType_Check:   "",
	ProType_Date:    "",
}

type Property struct {
	widget.BaseWidget
	ProContainer *fyne.Container
	ProType      *widget.Select
	ProName      *widget.Entry
	ProValue     *widget.Entry
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
func NewProperty(types []string, defaultType string, defaultName string, defaultValue string) *Property {
	t := &Property{}
	t.ExtendBaseWidget(t)

	t.ProName = widget.NewEntry()
	t.ProName.SetPlaceHolder("属性名称")
	t.ProName.SetText(defaultName)

	t.ProValue = widget.NewEntry()
	t.ProValue.SetPlaceHolder("属性值")
	t.ProValue.SetText(defaultValue)
	//根据属性类型，对输入内容进行检验

	t.ProType = widget.NewSelect(types, func(s string) {
		//选定属性类型时，自动更改属性名称，默认属性名称为空时不更改
		if typeMap[s] != "" {
			t.ProName.SetText(typeMap[s])
		}
		//根据选定的属性类型，修改属性值输入框的提示词

	})
	//设置默认属性类型
	t.ProType.SetSelected(defaultType)

	t.ProContainer = container.New(&PropertyLayout{}, t.ProType, t.ProName, t.ProValue)

	return t

}

func (t *Property) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(t.ProContainer)
}

// GetPropertyData 获取传入属性控件的子控件的值，以结构体指针作为返回
func (t *Property) GetPropertyData() *PropertyData {
	data := &PropertyData{}

	data.Name = t.ProName.Text
	data.Type = t.ProType.Selected
	data.Value = t.ProValue.Text

	return data
}

func (t *Property) GetPropertyName() string {
	return t.ProName.Text
}

func (t *Property) GetPropertyValue() string {
	return t.ProValue.Text
}
