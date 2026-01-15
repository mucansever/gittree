package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mucansever/gittree/internal/git"
	"github.com/mucansever/gittree/internal/tree"
)

type Model struct {
	items    []tree.Item
	cursor   int
	repo     *git.Repository
	err      error
	quitting bool
	message  string
}

func NewModel(items []tree.Item, repo *git.Repository) Model {
	return Model{
		items: items,
		repo:  repo,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter":
			selected := m.items[m.cursor]
			// skip checkout if it's the root
			if selected.BranchName == "." {
				m.message = "Cannot checkout root node."
				return m, nil
			}

			branchName := strings.TrimSuffix(selected.BranchName, "*")
			err := m.repo.Checkout(branchName)
			if err != nil {
				m.err = err
				m.message = fmt.Sprintf("Error checking out %s: %v", branchName, err)
			} else {
				m.message = fmt.Sprintf("Checked out %s", branchName)
				m.quitting = true
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.quitting {
		return fmt.Sprintf("%s\n", m.message)
	}

	s := lipgloss.NewStyle().Bold(true).Render("Navigate branches (Up/Down), Enter to checkout, q to quit.") + "\n\n"

	for i, item := range m.items {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}

		line := item.Text
		if m.cursor == i {
			line = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).Render(line)
			cursor = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(cursor)
		}

		s += fmt.Sprintf("%s%s\n", cursor, line)
	}

	if m.message != "" {
		s += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(m.message) + "\n"
	}

	return s
}
