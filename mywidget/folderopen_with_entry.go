package mywidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type FolderOpenWithEntry struct {
	widget.BaseWidget
	feContainer *fyne.Container
}

type FolderOpenWithEntryLayout struct{}

func (lo *FolderOpenWithEntryLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, obj := range objects {
		minSize = minSize.Max(obj.MinSize())
	}
	return minSize
}

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

func NewFolderOpenWithEntry(entryChanged func(s string), entryPlaceHolder string, win fyne.Window) *FolderOpenWithEntry {
	entry := widget.NewEntry()
	button := widget.NewButton("打开文件夹", func() {
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
			entry.SetText(list.Path())
		}, win)
	})

	entry.OnChanged = entryChanged
	entry.SetPlaceHolder(entryPlaceHolder)

	t := &FolderOpenWithEntry{}
	t.ExtendBaseWidget(t)
	t.feContainer = container.New(&FolderOpenWithEntryLayout{}, entry, button)

	return t
}

func (t *FolderOpenWithEntry) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(t.feContainer)
}

func (t *FolderOpenWithEntry) SetEntryText(s string) {
	for _, obj := range t.feContainer.Objects {
		switch v := obj.(type) {
		case *widget.Entry:
			v.SetText(s)
		}
	}
}

func (t *FolderOpenWithEntry) Enable() {
	for _, obj := range t.feContainer.Objects {
		switch v := obj.(type) {
		case *widget.Entry:
			v.Enable()
		case *widget.Button:
			v.Enable()
		}
	}
}

func (t *FolderOpenWithEntry) Disable() {
	for _, obj := range t.feContainer.Objects {
		switch v := obj.(type) {
		case *widget.Entry:
			v.Disable()
		case *widget.Button:
			v.Disable()
		}
	}
}
