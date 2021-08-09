package wildyam

import (
	"context"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gari8/wildyam/domain"
	"image/color"
	"io/ioutil"
	"log"
)

type DrawManager interface {
	Draw(ctx context.Context)
}

type drawManager struct {
	Recepter domain.Recepter
}

func NewDrawManager(recep domain.Recepter) DrawManager {
	return drawManager{Recepter: recep}
}

func (m drawManager) Draw(ctx context.Context) {
	a := app.New()
	w := a.NewWindow("WildYam")
	currentBar := canvas.NewRectangle(theme.PrimaryColor())
	ncon := container.NewVBox()

	ncon.Resize(fyne.Size{
		Width:  1000.0,
		Height: 1000.0,
	})

	vcontainer := container.NewVBox(
		currentBar,
		container.NewHBox(
			&widget.Select{
				DisableableWidget: widget.DisableableWidget{},
				Selected:          "",
				Options:           []string{"/test.go", "/a.go"},
				PlaceHolder:       "please, select file!",
				OnChanged: func(s string) {
					fmt.Println(s)
				},
			},
			SetForm(),
			dialogScreen(w),
		),
		ncon,
	)

	vcontainer.Resize(fyne.Size{
		Width:  500.0,
		Height: 400.0,
	})

	w.SetContent(vcontainer)

	w.ShowAndRun()
}

func SetForm() *widget.Form {
	return &widget.Form{
		BaseWidget: widget.BaseWidget{},
		Items:      nil,
		OnSubmit: func() {

		},
		OnCancel:   nil,
		SubmitText: "Read",
		CancelText: "Cancel",
	}
}

func confirmCallback(response bool) {
	fmt.Println("Responded with", response)
}

func colorPicked(c color.Color, w fyne.Window) {
	log.Println("Color picked:", c)
	rectangle := canvas.NewRectangle(c)
	size := 2 * theme.IconInlineSize()
	rectangle.SetMinSize(fyne.NewSize(size, size))
	dialog.ShowCustom("Color Picked", "Ok", rectangle, w)
}

func imageOpened(f fyne.URIReadCloser) {
	if f == nil {
		log.Println("Cancelled")
		return
	}
	defer f.Close()

	showImage(f)
}

func fileSaved(f fyne.URIWriteCloser, w fyne.Window) {
	defer f.Close()
	_, err := f.Write([]byte("Written by Fyne demo\n"))
	if err != nil {
		dialog.ShowError(err, w)
	}
	log.Println("Saved to...", f.URI())
}

func loadImage(f fyne.URIReadCloser) *canvas.Image {
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fyne.LogError("Failed to load image data", err)
		return nil
	}
	res := fyne.NewStaticResource(f.URI().Name(), data)

	return canvas.NewImageFromResource(res)
}

func showImage(f fyne.URIReadCloser) {
	img := loadImage(f)
	if img == nil {
		return
	}
	img.FillMode = canvas.ImageFillOriginal

	w := fyne.CurrentApp().NewWindow(f.URI().Name())
	w.SetContent(container.NewScroll(img))
	w.Resize(fyne.NewSize(320, 240))
	w.Show()
}

func dialogScreen(win fyne.Window) fyne.CanvasObject {
	return container.NewVScroll(container.NewVBox(
		widget.NewButton("Info", func() {
			dialog.ShowInformation("Information", "You should know this thing...", win)
		}),
		widget.NewButton("Error", func() {
			err := errors.New("a dummy error message")
			dialog.ShowError(err, win)
		}),
		widget.NewButton("Confirm", func() {
			cnf := dialog.NewConfirm("Confirmation", "Are you enjoying this demo?", confirmCallback, win)
			cnf.SetDismissText("Nah")
			cnf.SetConfirmText("Oh Yes!")
			cnf.Show()
		}),
		widget.NewButton("File Open With Filter (.jpg or .png)", func() {
			fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err != nil {
					dialog.ShowError(err, win)
					return
				}
				if reader == nil {
					log.Println("Cancelled")
					return
				}

				imageOpened(reader)
			}, win)
			fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
			fd.Show()
		}),
		widget.NewButton("File Save", func() {
			dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, win)
					return
				}
				if writer == nil {
					log.Println("Cancelled")
					return
				}

				fileSaved(writer, win)
			}, win)
		}),
		widget.NewButton("Folder Open", func() {
			dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
				if err != nil {
					dialog.ShowError(err, win)
					return
				}
				if list == nil {
					log.Println("Cancelled")
					return
				}

				children, err := list.List()
				if err != nil {
					dialog.ShowError(err, win)
					return
				}
				out := fmt.Sprintf("Folder %s (%d children):\n%s", list.Name(), len(children), list.String())
				dialog.ShowInformation("Folder Open", out, win)
			}, win)
		}),
		widget.NewButton("Color Picker", func() {
			picker := dialog.NewColorPicker("Pick a Color", "What is your favorite color?", func(c color.Color) {
				colorPicked(c, win)
			}, win)
			picker.Show()
		}),
		widget.NewButton("Advanced Color Picker", func() {
			picker := dialog.NewColorPicker("Pick a Color", "What is your favorite color?", func(c color.Color) {
				colorPicked(c, win)
			}, win)
			picker.Advanced = true
			picker.Show()
		}),
		widget.NewButton("Form Dialog (Login Form)", func() {
			username := widget.NewEntry()
			username.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "username can only contain letters, numbers, '_', and '-'")
			password := widget.NewPasswordEntry()
			password.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "password can only contain letters, numbers, '_', and '-'")
			remember := false
			items := []*widget.FormItem{
				widget.NewFormItem("Username", username),
				widget.NewFormItem("Password", password),
				widget.NewFormItem("Remember me", widget.NewCheck("", func(checked bool) {
					remember = checked
				})),
			}

			dialog.ShowForm("Login...", "Log In", "Cancel", items, func(b bool) {
				if !b {
					return
				}
				var rememberText string
				if remember {
					rememberText = "and remember this login"
				}

				log.Println("Please Authenticate", username.Text, password.Text, rememberText)
			}, win)
		}),
	))
}
