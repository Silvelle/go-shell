// Package ui implements the graphical terminal window of the emulator.
package ui

import (
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/Silvelle/go-shell/src/shell"
)

const (
	windowWidth  = 800
	windowHeight = 500
)

// Options describe what the window does right after it opens.
type Options struct {
	// Messages are printed before anything else, e.g. the settings dump.
	Messages []string
	// ScriptPath is the startup script to run; empty means none.
	ScriptPath string
}

// terminal connects a shell to the widgets that show the dialog.
type terminal struct {
	sh     *shell.Shell
	window fyne.Window
	output *widget.TextGrid
	scroll *container.Scroll
	input  *widget.Entry
}

// Run opens the emulator window and blocks until it is closed.
func Run(sh *shell.Shell, opts Options) {
	a := app.New()
	w := a.NewWindow(shell.SystemTitle())
	t := newTerminal(sh, w)

	w.SetContent(container.NewBorder(nil, t.input, nil, nil, t.scroll))
	w.Resize(fyne.NewSize(windowWidth, windowHeight))
	w.Canvas().Focus(t.input)
	go fyne.Do(func() { t.start(opts) })
	w.ShowAndRun()
}

// newTerminal creates the output area and the input line.
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

// start prints the startup messages and runs the startup script.
func (t *terminal) start(opts Options) {
	for _, m := range opts.Messages {
		t.print(m)
	}
	if opts.ScriptPath == "" {
		return
	}
	t.handle(t.sh.RunScript(opts.ScriptPath, t.print))
}

// submit runs one line typed by the user.
func (t *terminal) submit(line string) {
	t.input.SetText("")
	t.handle(t.sh.Interact(line, t.print))
}

// handle closes the window on exit and prints any other error.
func (t *terminal) handle(err error) {
	if errors.Is(err, shell.ErrExit) {
		t.window.Close()
		return
	}
	if err != nil {
		t.print("error: " + err.Error())
	}
}

// print appends a line to the output area and scrolls to it.
// The line is also written to the console, so OS test scripts show it.
func (t *terminal) print(text string) {
	fmt.Println(text)
	t.output.SetText(t.output.Text() + text + "\n")
	t.scroll.ScrollToBottom()
}
