package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices []string
	cursor int
}

func InitialModel() model {
	return model{
		choices: []string{"YouTube - Watch a video from YouTube", "Local File - Watch a video stored on your machine"},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

const (
	YouTubeOption = "YouTube"
	LocalOption = "Local"
	NotSelected = "NotSelected"
)

var SelectedOption string = NotSelected

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices) {
				m.cursor++
			}
		case "enter", " ":
			if m.cursor == 0 {
				SelectedOption = YouTubeOption
			} else if m.cursor == 1 {
				SelectedOption = LocalOption
			}
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "Choose a video source?\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPress q to quit.\n"

	return s
}