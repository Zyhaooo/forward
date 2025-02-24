package main

import (
	"errors"
	"log"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var errChan chan error

func main() {
	gtk.Init(nil)

	errChan = make(chan error)

	// 新建一个窗口
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("kubernetes port forward")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	go showErrorDialog(win, errChan)

	// 垂直布局
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		log.Fatal("Unable to create vertical box:", err)
	}
	box.SetMarginTop(10)

	// 添加按钮
	newAddButton(box)
	addRow(box, errChan)

	win.SetResizable(false) // 禁用
	win.Resize(1, 1)
	win.Add(box)
	win.ShowAll()
	gtk.Main()
}

func newAddButton(w *gtk.Box) {

	hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	if err != nil {
		log.Fatal("Unable to create horizontal box:", err)
	}

	addButton, err := gtk.ButtonNewWithLabel("+")
	if err != nil {
		log.Fatal("Unable to create addButton:", err)
	}

	addButton.Connect("clicked", func() {
		addRow(w, errChan)
	})

	addButton.SetSizeRequest(12, -1)
	hbox.PackEnd(addButton, false, true, 5)
	hbox.ShowAll()
	w.PackEnd(hbox, false, true, 5)
}

func addRow(box *gtk.Box, eChan chan error) {
	hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	if err != nil {
		log.Fatal("Unable to create horizontal box:", err)
	}

	{
		namespace, err := gtk.ComboBoxTextNew()
		if err != nil {
			log.Fatal("Unable to create combo box:", err)
		}
		namespace.AppendText("default")
		namespace.AppendText("backend")
		namespace.AppendText("message")

		namespace.SetActive(0)

		namespace.Connect("changed", func() {
			if namespace.GetActiveText() == "--select namespace--" {
				namespace.SetActive(-1)
			} else {
				log.Println("select =>", namespace.GetActiveText())
			}
		})
		hbox.PackStart(namespace, false, true, 5)

	}

	{
		sources := []string{
			"pods",
			"deployment",
			"replicaset",
			"service",
		}

		sComboBox, err := gtk.ComboBoxTextNew()
		if err != nil {
			log.Fatal("Unable to create sources comboBox", err)
		}

		for _, v := range sources {
			sComboBox.AppendText(v)
		}
		sComboBox.SetActive(0)

		hbox.PackStart(sComboBox, false, true, 5)
	}

	{

		names := []string{
			"pod/user-rpc-1nidnqin",
			"pod/user-rpc-fasdnqin",
			"pod/user-rpc-1nignidn",
			"pod/user-rpc-jfindgab",
			"pod/user-rpc-fjbgbasg",
			"pod/user-rpc-hgiabgba",
			"pod/user-rpc-jigagbst",
		}

		nComboBox, err := gtk.ComboBoxTextNewWithEntry()
		if err != nil {
			log.Fatal("Unable to create names comboBox", err)
		}
		listStore, _ := gtk.ListStoreNew(glib.TYPE_STRING)
		nComboBox.SetModel(listStore)

		for k, v := range names {
			iter := listStore.Insert(k)
			err := listStore.Set(iter, []int{0}, []any{v})
			if err != nil {
				log.Fatal(err)
			}
		}

		entry, _ := nComboBox.GetEntry()
		completion, _ := gtk.EntryCompletionNew()
		completion.SetModel(listStore)
		completion.SetTextColumn(0)
		completion.SetInlineCompletion(true)

		entry.SetCompletion(completion)

		hbox.PackStart(nComboBox, false, true, 5)
	}

	{
		var keyPressEvent = func(entry *gtk.Entry, event *gdk.Event) bool {
			key := gdk.EventKeyNewFromEvent(event).KeyVal()
			if (key >= '0' && key <= '9') || key == gdk.KEY_BackSpace || key == gdk.KEY_Delete {
				return false
			}
			return true
		}
		e1, err := gtk.EntryNew()
		if err != nil {
			log.Fatal("Unable to create entry:", err)
		}

		e2, err := gtk.EntryNew()
		if err != nil {
			log.Fatal("Unable to create entry:", err)
		}
		e1.SetWidthChars(5)
		e2.SetWidthChars(5)
		e1.Connect("key-press-event", keyPressEvent)
		e2.Connect("key-press-event", keyPressEvent)

		label, _ := gtk.LabelNew("=>")

		hbox.PackStart(e1, false, true, 1)
		hbox.PackStart(label, false, true, 0)
		hbox.PackStart(e2, false, true, 1)
	}

	{
		// 按钮
		button, err := gtk.ButtonNewWithLabel("forward")
		if err != nil {
			log.Fatal("Unable to create button:", err)
		}
		button.ToWidget().SetSizeRequest(100, -1)
		button.Connect("clicked", func() {
			label, _ := button.GetLabel()
			if label == "forward" {
				eChan <- errors.New("forward start")
				button.SetLabel("stop")
			} else {
				button.SetLabel("forward")
			}
		})
		hbox.PackStart(button, false, false, 5)
	}

	{
		remove, err := gtk.ButtonNew()
		if err != nil {
			log.Fatal("Unable to create remove button:", err)
		}
		remove.Connect("clicked", func() {
			if box.GetChildren().Length() != 2 {
				box.Remove(hbox)
			}
		})

		icon, _ := gtk.ImageNewFromIconName(
			"window-close-symbolic", // GNOME标准删除图标
			gtk.ICON_SIZE_BUTTON)
		remove.SetImage(icon)
		remove.SetAlwaysShowImage(true)
		remove.SetImagePosition(gtk.POS_LEFT) // 图标在文字左侧

		// 添加辅助提示
		remove.SetTooltipText("永久删除选定资源（不可恢复）")

		hbox.PackStart(remove, false, false, 5)

	}

	box.Add(hbox)
	box.ShowAll()
}
func showErrorDialog(parent *gtk.Window, err chan error) {

	for {
		select {
		case e := <-err:
			// 创建错误对话框
			dialog := gtk.MessageDialogNew(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, e.Error())

			dialog.SetModal(true)

			dialog.Connect("response", func() {
				dialog.Destroy() // 关闭后销毁对话框
			})

			dialog.ShowAll()
		default:
			continue
		}
	}
}
