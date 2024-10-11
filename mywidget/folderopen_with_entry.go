package mywidget

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type FolderOpenWithEntry struct {
	widget.BaseWidget
	feContainer *fyne.Container
	entry       *widget.Entry
	button      *widget.Button
}

type FolderOpenWithEntryLayout struct{}

func (lo *FolderOpenWithEntryLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, obj := range objects {
		minSize = minSize.Max(obj.MinSize())
	}
	return minSize
}

// Layout 按钮在左，最小宽度；文本框在右，占用剩余宽度
func (lo *FolderOpenWithEntryLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
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

// NewFolderOpenWithEntry 创建带文本框的文件夹打开按钮。
// 传入参数：entryChanged 文本框输入改变时触发的回调函数；entryPlaceHolder 文本框占位符，输入为空时显示；win 父窗口
func NewFolderOpenWithEntry(entryChanged func(s string), entryPlaceHolder string, win fyne.Window) *FolderOpenWithEntry {
	t := &FolderOpenWithEntry{}
	t.ExtendBaseWidget(t)

	t.entry = widget.NewEntry()
	t.button = widget.NewButton("打开文件夹", func() {
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
			//选择的文件夹路径显示在输入框中
			t.entry.SetText(list.Path())
		}, win)
	})

	t.entry.OnChanged = entryChanged
	t.entry.SetPlaceHolder(entryPlaceHolder)
	//检查输入的路径是否存在
	t.entry.Validator = func(s string) error {
		_, err := os.Stat(s)
		if err == nil {
			return nil
		}
		return err
	}

	t.feContainer = container.New(&FolderOpenWithEntryLayout{}, t.entry, t.button)

	return t
}

func (t *FolderOpenWithEntry) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(t.feContainer)
}

// SetEntryText 设置文本框内容
func (t *FolderOpenWithEntry) SetEntryText(s string) {
	t.entry.SetText(s)
}

// GetEntryText 获取文本框内容
func (t *FolderOpenWithEntry) GetEntryText() string {
	return t.entry.Text
}

// Enable 使文本框和按钮都可用
func (t *FolderOpenWithEntry) Enable() {
	t.button.Enable()
	t.entry.Enable()
}

// Disable 使文本框和按钮都不可用
func (t *FolderOpenWithEntry) Disable() {
	t.button.Disable()
	t.entry.Disable()
}

// GetValid 获取文本框的检查状态，文本框为空或指向路径不存在时返回false
func (t *FolderOpenWithEntry) GetValid() bool {
	return t.entry.Validate() == nil
}
