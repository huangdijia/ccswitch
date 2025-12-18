package termui

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/term"
)

var (
	// ErrCanceled is returned when the user cancels the selection.
	ErrCanceled = errors.New("selection canceled")
	// ErrNotTerminal is returned when interactive selection is requested but stdin/stdout is not a terminal.
	ErrNotTerminal = errors.New("not a terminal")
)

type SelectConfig struct {
	In           *os.File
	Out          *os.File
	Prompt       string
	Hint         string
	Items        []string
	DefaultIndex int
}

func SelectString(cfg SelectConfig) (string, error) {
	if cfg.In == nil || cfg.Out == nil {
		return "", fmt.Errorf("invalid terminal io")
	}
	if len(cfg.Items) == 0 {
		return "", fmt.Errorf("no items to select")
	}
	if cfg.DefaultIndex < 0 || cfg.DefaultIndex >= len(cfg.Items) {
		cfg.DefaultIndex = 0
	}
	if cfg.Prompt == "" {
		cfg.Prompt = "Select:"
	}
	if cfg.Hint == "" {
		cfg.Hint = "↑/↓ to move, Enter to select, q to cancel"
	}

	inFD := int(cfg.In.Fd())
	outFD := int(cfg.Out.Fd())
	if !term.IsTerminal(inFD) || !term.IsTerminal(outFD) {
		return "", ErrNotTerminal
	}

	previousState, err := term.MakeRaw(inFD)
	if err != nil {
		return "", err
	}
	defer func() { _ = term.Restore(inFD, previousState) }()

	// Hide cursor while rendering.
	fmt.Fprint(cfg.Out, "\x1b[?25l")
	defer fmt.Fprint(cfg.Out, "\x1b[?25h")

	selected := cfg.DefaultIndex
	renderedLines := 0

	render := func() {
		if renderedLines > 0 {
			fmt.Fprintf(cfg.Out, "\x1b[%dA", renderedLines)
		}

		// Prompt line.
		fmt.Fprint(cfg.Out, "\x1b[2K\r")
		fmt.Fprintln(cfg.Out, cfg.Prompt)

		// Items.
		for i, item := range cfg.Items {
			line := "  " + item
			if i == selected {
				line = "\x1b[7m> " + item + "\x1b[0m"
			}
			fmt.Fprint(cfg.Out, "\x1b[2K\r")
			fmt.Fprintln(cfg.Out, line)
		}

		// Hint line.
		fmt.Fprint(cfg.Out, "\x1b[2K\r")
		fmt.Fprintln(cfg.Out, cfg.Hint)

		renderedLines = 1 + len(cfg.Items) + 1
	}

	render()

	// Simple escape-sequence parser (supports arrow keys).
	escPending := false
	csiPending := false
	ss3Pending := false

	buf := make([]byte, 1)
	for {
		n, err := cfg.In.Read(buf)
		if err != nil {
			return "", err
		}
		if n == 0 {
			continue
		}
		b := buf[0]

		if csiPending {
			csiPending = false
			switch b {
			case 'A':
				if selected > 0 {
					selected--
					render()
				}
			case 'B':
				if selected < len(cfg.Items)-1 {
					selected++
					render()
				}
			}
			continue
		}
		if ss3Pending {
			ss3Pending = false
			switch b {
			case 'A':
				if selected > 0 {
					selected--
					render()
				}
			case 'B':
				if selected < len(cfg.Items)-1 {
					selected++
					render()
				}
			}
			continue
		}
		if escPending {
			escPending = false
			switch b {
			case '[':
				csiPending = true
			case 'O':
				ss3Pending = true
			case 27: // ESC ESC
				return "", ErrCanceled
			}
			continue
		}

		switch b {
		case '\r', '\n':
			fmt.Fprintln(cfg.Out)
			return cfg.Items[selected], nil
		case 3: // Ctrl+C (in raw mode)
			return "", ErrCanceled
		case 'q', 'Q':
			return "", ErrCanceled
		case 'k', 'K':
			if selected > 0 {
				selected--
				render()
			}
		case 'j', 'J':
			if selected < len(cfg.Items)-1 {
				selected++
				render()
			}
		case 27: // ESC prefix (for arrows or cancellation via ESC ESC)
			escPending = true
		}
	}
}
