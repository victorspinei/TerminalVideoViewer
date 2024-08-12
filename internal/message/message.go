package message

import (
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
}

func InitialModel() model {
	return model{
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// i added it just so i could use msg
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
		return m, tea.Quit
	}
	return m, nil
}

func (m model) View() string {
	s := "Welcome to Terminal Video Viewer!\n\n"
    s += "Before playing the video, please:\n"
    s += "1. Make your terminal fullscreen.\n"
    s += "2. Zoom out to fit the video properly.\n\n"
    s += "Controls:\n"
    s += "  Space: Pause/Play\n"
    s += "  Z: Seek backward\n"
    s += "  X: Seek forward\n"
    s += "  M: Mute/Unmute volume\n"
    s += "  Q: Quit the application\n\n"
    s += "Press any key to start...\n"
	return s
}