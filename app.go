package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type App struct {
	*gtk.Window

	ctx    context.Context
	rbox   *gtk.Box
	logger *slog.Logger
}

func NewApp(ctx context.Context) *App {
	var (
		err     error
		win     *gtk.Window
		rootBox *gtk.Box

		logger = slog.With("gtk log")

		app *App
	)

	// top level window
	{
		gtk.Init(nil)
		if win, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL); err != nil {
			logger.ErrorContext(ctx,
				"Unable to create top window",
				slog.Any("err", err))
			os.Exit(1)
		}
		win.SetTitle("forward")
		win.Connect("destroy", func() {
			gtk.MainQuit()
		})
		win.SetResizable(false) // 禁用
		win.Resize(1, 1)
	}

	// root box
	{
		if rootBox, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5); err != nil {
			logger.ErrorContext(ctx,
				"Unable to create root box",
				slog.Any("err", err))
			os.Exit(1)
		}
		rootBox.SetMarginTop(10)
		win.Add(rootBox)
	}

	// add column button
	{
		var (
			hbox   *gtk.Box
			button *gtk.Button
			icon   *gtk.Image

			err error
		)
		if hbox, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5); err != nil {
			logger.ErrorContext(ctx,
				"Unable to create horizontal box",
				slog.Any("err", err))
			os.Exit(1)
		}

		if button, err = gtk.ButtonNew(); err != nil {
			logger.ErrorContext(ctx,
				"Unable to create add button",
				slog.Any("err", err))
			os.Exit(1)
		}

		icon, _ = gtk.ImageNewFromIconName(
			"list-add",
			gtk.ICON_SIZE_BUTTON)
		button.SetImage(icon)
		button.SetAlwaysShowImage(true)

		button.Connect("clicked", func() {
			slog.InfoContext(ctx, "you clicked add column button")
			if err = app.addColumn(); err != nil {
				slog.ErrorContext(ctx, "Unable to add column in add button click", slog.Any("err", err))
			}
		})

		hbox.PackEnd(button, false, false, 5)
		hbox.ShowAll()
		rootBox.PackEnd(hbox, false, false, 5)
	}

	{
		app = &App{
			Window: win,
			ctx:    ctx,
			logger: logger,
			rbox:   rootBox,
		}
		if err = app.addColumn(); err != nil {
			logger.ErrorContext(ctx, "Unable to add column", slog.Any("err", err))
			os.Exit(1)
		}
	}

	return app
}

func (a *App) addColumn() (err error) {

	hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	if err != nil {
		a.logger.ErrorContext(a.ctx,
			"Unable to create horizontal box",
			slog.Any("err", err))
		return err
	}

	// namespace
	if err = a.namespaceComboBox(hbox); err != nil {
		return err
	}
	// types
	if err = a.typesComboBox(hbox); err != nil {
		return err
	}
	// resources
	if err = a.resourcesEntryComboBox(hbox); err != nil {
		return err
	}
	// portForwardEntry
	if err = a.portForWardEntry(hbox); err != nil {
		return err
	}
	// forward
	if err = a.forwardButton(hbox); err != nil {
		return err
	}
	// remove
	if err = a.removeButton(hbox); err != nil {
		return err
	}

	a.rbox.Add(hbox)
	a.rbox.ShowAll()
	return
}

func (a *App) removeButton(box *gtk.Box) (err error) {
	remove, err := gtk.ButtonNew()
	if err != nil {
		a.logger.ErrorContext(a.ctx,
			"Unable to create remove button:",
			slog.Any("err", err))
		return err
	}
	remove.Connect("clicked", func() {
		parent, _ := box.GetParent()
		p := parent.(*gtk.Box)
		if p.GetChildren().Length() != 2 {
			parent.(*gtk.Box).Remove(box)
		}
	})

	icon, _ := gtk.ImageNewFromIconName(
		"window-close-symbolic",
		gtk.ICON_SIZE_BUTTON)
	remove.SetImage(icon)
	remove.SetAlwaysShowImage(true)

	box.PackStart(remove, false, false, 5)

	return nil
}

func (a *App) forwardButton(box *gtk.Box) (err error) {
	button, err := gtk.ButtonNewWithLabel("forward")
	if err != nil {
		a.logger.ErrorContext(a.ctx,
			"Unable to create forward button:",
			slog.Any("err", err))
		return err
	}
	button.ToWidget().SetSizeRequest(100, -1)
	button.Connect("clicked", func() {
		label, _ := button.GetLabel()
		if label == "forward" {
			button.SetLabel("stop")
		} else {
			button.SetLabel("forward")
		}
	})
	box.PackStart(button, false, false, 5)

	return nil
}

func (a *App) resourcesEntryComboBox(box *gtk.Box) (err error) {

	r, err := gtk.ComboBoxTextNewWithEntry()
	if err != nil {
		a.logger.ErrorContext(a.ctx,
			"Unable to create resources entry combo box:",
			slog.Any("err", err))
		return err
	}

	lStore, _ := gtk.ListStoreNew(glib.TYPE_STRING)
	r.SetModel(lStore)

	// TODO
	names := []string{
		"pod/user-rpc-1nidnqin",
		"pod/user-rpc-fasdnqin",
		"pod/user-rpc-1nignidn",
		"pod/user-rpc-jfindgab",
		"pod/user-rpc-fjbgbasg",
		"pod/user-rpc-hgiabgba",
		"pod/user-rpc-jigagbst",
	}

	for k, v := range names {
		iter := lStore.Insert(k)
		err := lStore.Set(iter, []int{0}, []any{v})
		if err != nil {
			a.logger.ErrorContext(a.ctx, "Unable to set rescources item liststore", slog.Any("err", err))
			return err
		}
	}
	entry, _ := r.GetEntry()
	completion, _ := gtk.EntryCompletionNew()
	completion.SetModel(lStore)
	completion.SetTextColumn(0)
	completion.SetInlineCompletion(true)
	entry.SetCompletion(completion)

	box.PackStart(r, false, true, 5)
	return
}

func (a *App) portForWardEntry(box *gtk.Box) (err error) {
	var keyPressEvent = func(entry *gtk.Entry, event *gdk.Event) bool {
		key := gdk.EventKeyNewFromEvent(event).KeyVal()
		if (key >= '0' && key <= '9') || key == gdk.KEY_BackSpace || key == gdk.KEY_Delete {
			return false
		}
		return true
	}

	origin, err := gtk.EntryNew()
	if err != nil {
		a.logger.ErrorContext(a.ctx,
			"Unable to create origin port entry",
			slog.Any("err", err))
		return err
	}
	origin.SetWidthChars(5)
	origin.Connect("key-press-event", keyPressEvent)

	forward, err := gtk.EntryNew()
	if err != nil {
		a.logger.ErrorContext(a.ctx,
			"Unable to create forward port entry",
			slog.Any("err", err))
		return err
	}
	forward.SetWidthChars(5)
	forward.Connect("key-press-event", keyPressEvent)

	label, _ := gtk.LabelNew("=>")

	box.PackStart(origin, false, true, 5)
	box.PackStart(label, false, true, 0)
	box.PackStart(forward, false, true, 5)

	return
}

func (a *App) typesComboBox(box *gtk.Box) (err error) {

	r, err := gtk.ComboBoxTextNew()
	if err != nil {
		a.logger.ErrorContext(a.ctx,
			"Unable to create resources combo box:",
			slog.Any("err", err))
		return err
	}

	r.AppendText("pods")
	r.AppendText("deployment")
	r.AppendText("replicaset")
	r.AppendText("service")

	box.PackStart(r, false, true, 5)

	return
}

func (a *App) namespaceComboBox(box *gtk.Box) (err error) {

	namespace, err := gtk.ComboBoxTextNew()
	if err != nil {
		a.logger.ErrorContext(a.ctx,
			"Unable to create namespace combo box:",
			slog.Any("err", err))
		return err
	}

	// TODO add text
	namespace.AppendText("default")
	namespace.AppendText("backend")
	namespace.AppendText("message")

	box.PackStart(namespace, false, true, 5)

	return
}

func (a *App) Run() {
	a.ShowAll()
	gtk.Main()
}
