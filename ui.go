package main

import (
	"fmt"
	"math/rand"
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

	shamePhrases = []string{
		"I am consciously choosing to abandon my deeply focused state of work, prioritizing fleeting algorithmic dopamine hits over my long-term engineering goals and personal discipline.",
		"Despite setting a strict timer for myself, I lack the requisite willpower to sustain my attention span; therefore, I am manually aborting this session and accepting my own failure.",
		"System override initiated: I acknowledge that my urge to check notifications has completely overpowered my rational intent to complete this task. I surrender to the digital noise.",
		"By typing this exact sentence, character by character, I am proving that I would rather perform a tedious, humiliating data-entry task than simply sit with my own thoughts and do the real work.",
		"Error 403: Willpower Forbidden. I am terminating the focus protocol prematurely. I recognize that this action actively sabotages my momentum, yet I am doing it anyway.",
	}
)

type model struct {
	timer        timer.Model
	input        textinput.Model
	quitting     bool
	aborted      bool
	shameMode    bool
	width        int
	height       int
	targetPhrase string
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.timer.Init(), textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		m.quitting = true
		return m, tea.Quit
	case tea.KeyMsg:
		if msg.Paste {
			return m, nil
		}

		if m.shameMode {
			switch msg.String() {
			case "enter":

				if m.input.Value() == m.targetPhrase {
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
			m.shameMode = true
			m.targetPhrase = shamePhrases[rand.Intn(len(shamePhrases))]
			m.input.Focus()
			return m, textinput.Blink
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	if m.quitting {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, "\n  ✨ Session Complete! Network restored. ✨\n\n")
	}

	if m.aborted {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, "\n  Session aborted. You let the distractions win.\n\n")
	}

	if m.shameMode {
		phraseStyle := lipgloss.NewStyle().
			Foreground(roseColor).
			Italic(true).
			Width(70)

		s := "\n  [!] WEAKNESS DETECTED [!]\n\n"
		s += "  To cancel, type EXACTLY:\n"
		s += phraseStyle.Render("  \""+m.targetPhrase+"\"") + "\n\n"
		inputLine := lipgloss.NewStyle().Width(80).Render("> " + m.input.View())
		s += inputLine + "\n\n"
		s += "  (Press Esc to return to work, or Enter to submit)\n"

		shameBox := lipgloss.NewStyle().Width(85).Align(lipgloss.Left).Foreground(loveColor).Render(s)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, shameBox)
	}

	content := fmt.Sprintf(
		"%s\n\n  Remaining: %s\n\n  [Stay Focused]",
		titleStyle.Render("DEEP WORK ACTIVE"),
		timerStyle.Render(m.timer.View()),
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func runTUI(minutes int) error {
	duration := time.Duration(minutes) * time.Minute

	ti := textinput.New()
	ti.Placeholder = "Type the phrase here..."
	ti.CharLimit = 300
	ti.Width = 80

	m := model{
		timer: timer.NewWithInterval(duration, time.Second),
		input: ti,
	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		return err
	}
	return nil
}
