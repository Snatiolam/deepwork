package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	roseColor = lipgloss.Color("#ebbcba") // Pinkish Rose
	loveColor = lipgloss.Color("#eb6f92") // Deep Rose/Pink

	titleStyle = lipgloss.NewStyle().
			Foreground(roseColor).
			Bold(true).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(loveColor)

	timerStyle = lipgloss.NewStyle().
			Foreground(loveColor).
			Bold(true).
			Underline(true)
)

type model struct {
	timer     timer.Model
	input     textinput.Model
	quitting  bool
	aborted   bool
	shameMode bool
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.timer.Init(), textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		m.quitting = true
		return m, tea.Quit
	case tea.KeyMsg:
		if m.shameMode {
			switch msg.String() {
			case "enter":
				if m.input.Value() == "I surrender to my distractions" {
					m.aborted = true
					return m, tea.Quit
				}
				m.shameMode = false
				m.input.SetValue("")
				return m, nil
			case "esc":
				m.shameMode = false
				m.input.SetValue("")
				return m, nil
			}

			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "ctrl+c", "q":
			// Don't quit! Activate Shame Mode instead.
			m.shameMode = true
			m.input.Focus()
			return m, textinput.Blink
		}
	}
	return m, nil
}

// 5. VIEW (The "Painter": turns state into a string)
func (m model) View() string {
	if m.quitting {
		return "\n  ✨ Session Complete! Network restored. ✨\n\n"
	}

	if m.aborted {
		return "\n  Session aborted. You let the distractions win.\n\n"
	}

	if m.shameMode {
		s := "\n  [!] WEAKNESS DETECTED [!]\n\n"
		s += "  To cancel, type EXACTLY:\n"
		s += "  'I surrender to my distractions'\n\n"
		s += "  > " + m.input.View() + "\n\n"
		s += "  (Press Esc to return to work, or Enter to submit)\n"
		return lipgloss.NewStyle().Foreground(loveColor).Render(s)
	}

	// Center the UI in the terminal
	content := fmt.Sprintf(
		"%s\n\n  Remaining: %s\n\n  [Stay Focused]",
		titleStyle.Render("DEEP WORK ACTIVE"),
		timerStyle.Render(m.timer.View()),
	)

	return "\n" + content + "\n"
}

// 6. THE RUNNER (The entry point called from main.go)
func runTUI(minutes int) error {
	duration := time.Duration(minutes) * time.Minute

	ti := textinput.New()
	ti.Placeholder = "Type the phrase here..."
	ti.CharLimit = 100
	ti.Width = 50

	m := model{
		// NewWithInterval takes the total time and how often to 'tick' (1s)
		timer: timer.NewWithInterval(duration, time.Second),
		input: ti,
	}

	// tea.WithAltScreen() makes it "fullscreen" like Neovim/Vim
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		return err
	}
	return nil
}
