package mywidget

import (
	"errors"
	"regexp"
	"time"

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
var type2NameMap = map[string]string{
	ProType_Tag:     "tags",
	ProType_Aliases: "aliases",
	ProType_css:     "cssclasses",
	ProType_Text:    "",
	ProType_List:    "",
	ProType_Num:     "",
	ProType_Check:   "",
	ProType_Date:    "",
}

// 属性类型与属性值提示词的映射表
var type2PromptMap = map[string]string{
	ProType_Tag:     "tag1,tag2...",
	ProType_Aliases: "aliases1,aliases2...",
	ProType_css:     "css1,css2...",
	ProType_Text:    "任意文本",
	ProType_List:    "list1,list2...",
	ProType_Num:     "只能是数字",
	ProType_Check:   "true/false",
	ProType_Date:    "YYYY-MM-DD",
}

// 属性类型与属性值文本框检查器的映射表
var type2ValidatorMap = map[string]func(s string) error{
	ProType_Tag:     validator_default,
	ProType_Aliases: validator_default,
	ProType_css:     validator_default,
	ProType_Text:    func(s string) error { return nil },
	ProType_List:    validator_default,
	ProType_Num:     validator_num,
	ProType_Check:   validator_bool,
	ProType_Date:    validator_data,
}

type Property struct {
	widget.BaseWidget
	ProContainer *fyne.Container
	ProType      *widget.Select //属性类型下拉菜单
	ProName      *widget.Entry  //属性名称文本框
	ProValue     *widget.Entry  //属性值文本框
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

// Layout 属性类型、名称、值控件依次水平放置。属性类型固定宽度，属性名称与属性值按1:2分配剩余宽度
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

// validator_default 默认检查器。文本格式不是“item1,item2,item3”时返回error
func validator_default(s string) error {
	pat := "^[^,]+(,[^,]+)*$"
	re := regexp.MustCompile(pat)
	if re.MatchString(s) {
		return nil
	}
	return errors.New("")
}

// validator_num 数字检查器。文本格式不是连续数字时返回error
func validator_num(s string) error {
	pat := "^\\d+$"
	re := regexp.MustCompile(pat)
	if re.MatchString(s) {
		return nil
	}
	return errors.New("")
}

// validator_data 日期检查器。文本格式不是合法日期时返回error
func validator_data(s string) error {
	pat := "^\\d{4}-\\d{2}-\\d{2}$"
	re := regexp.MustCompile(pat)
	_, err := time.Parse("2006-01-02", s)
	if re.MatchString(s) && err == nil {
		return nil
	}
	return errors.New("")
}

// validator_bool 布尔值检查器。文本不为true或false时返回error
func validator_bool(s string) error {
	if s == "true" || s == "false" {
		return nil
	}
	return errors.New("")
}

// NewProperty 创建属性控件。
// 传入参数：types 属性类型下拉菜单可选项；其余为子控件的默认值。
func NewProperty(types []string, defaultType string, defaultName string, defaultValue string) *Property {
	t := &Property{}
	t.ExtendBaseWidget(t)

	t.ProName = widget.NewEntry()
	t.ProName.SetPlaceHolder("属性名称")
	t.ProName.SetText(defaultName)

	t.ProValue = widget.NewEntry()
	t.ProValue.SetText(defaultValue)

	t.ProType = widget.NewSelect(types, func(s string) {
		//选定属性类型时，自动更改属性名称，默认属性名称为空时不更改
		if type2NameMap[s] != "" {
			t.ProName.SetText(type2NameMap[s])
		}
		//根据选定的属性类型，修改属性值输入框的提示词和检查器
		t.ProValue.SetPlaceHolder(type2PromptMap[s])
		t.ProValue.Validator = type2ValidatorMap[s]
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

// GetPropertyName 获取属性名称
func (t *Property) GetPropertyName() string {
	return t.ProName.Text
}

// GetPropertyValue 获取属性值
func (t *Property) GetPropertyValue() string {
	return t.ProValue.Text
}

// GetValid 获取属性值文本框检查状态。为空或格式不合法时返回false
func (t *Property) GetValid() bool {
	return t.ProValue.Validate() == nil
}
