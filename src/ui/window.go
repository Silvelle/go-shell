package ui

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Silvelle/go-shell/src/shell"
)

const (
	windowWidth  = 800
	windowHeight = 500
	prompt       = "$ "
)

// terminal connects a shell to the widgets that show the dialog.
type terminal struct {
	sh     *shell.Shell
	window fyne.Window
	output *widget.TextGrid
	scroll *container.Scroll
	input  *widget.Entry
}

// Run opens the emulator window and blocks until it is closed
func Run(sh *shell.Shell) {
	a := app.New()
	w := a.NewWindow(shell.SystemTitle())
	t := newTerminal(sh, w)

	w.SetContent(container.NewBorder(nil, t.input, nil, nil, t.scroll))
	w.Resize(fyne.NewSize(windowWidth, windowHeight))
	w.Canvas().Focus(t.input)
	w.ShowAndRun()
}

// newTerminal create teh output area and the input line.
func newTerminal(sh *shell.Shell, w fyne.Window) *terminal {
	t := &terminal{
		sh:     sh,
		window: w,
		output: widget.NewTextGrid(),
		input:  widget.NewEntry(),
	}
	t.scroll = container.NewScroll(t.output)
	t.input.SetPlaceHolder("Type a command and press Enter")
	t.input.OnSubmitted = t.submit
	return t
}

// submit runs one input line and prints the command and its result.
func (t *terminal) submit(line string) {
	t.input.SetText("")
	t.print(prompt + line)

	out, err := t.sh.Execute(line)
	if errors.Is(err, shell.ErrExit) {
		t.window.Close()
		return
	}
	if err != nil {
		t.print("error: " + err.Error())
	}
	if out != "" {
		t.print(out)
	}
}

func (t *terminal) print(text string) {
	t.output.SetText(t.output.Text() + text + "\n")
	t.scroll.ScrollToBottom()
}
