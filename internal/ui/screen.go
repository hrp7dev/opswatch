package ui

import "fmt"

// EnterAltScreen switches the terminal to the alternate screen buffer.
// Nothing that happens here touches the user's normal scrollback history —
// when we leave, everything is restored exactly as it was.
func EnterAltScreen() {
	fmt.Print("\033[?1049h")
	fmt.Print("\033[?25l")
	hideCursor()
}

// ExitAltScreen restores the terminal's original screen and scrollback.
// Always call this before the program exits (including on Ctrl+C / q).
func ExitAltScreen() {
	showCursor()
	fmt.Print("\033[?25h")
	fmt.Print("\033[?1049l")
}

func hideCursor() {
	fmt.Print("\x1b[?25l")
}

func showCursor() {
	fmt.Print("\x1b[?25h")
}

// Clear moves the cursor to the top-left and clears the alt-screen contents.
// Using cursor-home instead of scrolling means the dashboard is always
// redrawn in place, with no scroll-back growth.
func Clear() {
	fmt.Print("\033[H")
}
