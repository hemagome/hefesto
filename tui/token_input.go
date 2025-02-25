package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hemagome/hefesto/storage"
)

var (
	focusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)

	blurredStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("99")).
			Bold(true).
			MarginLeft(2)

	helpStyle = blurredStyle.
			Italic(true).
			MarginLeft(2)

	instructionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("86")).
				MarginLeft(2)
)

type TokenInput struct {
	textInput textinput.Model
	err       error
}

func NewTokenInput() *TokenInput {
	ti := textinput.New()
	ti.Placeholder = "Pega tu token de GitHub aquí"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	return &TokenInput{
		textInput: ti,
	}
}

func (m *TokenInput) Init() tea.Cmd {
	return textinput.Blink
}

func (m *TokenInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.textInput.Value() != "" {
				err := storage.StoreToken("github", m.textInput.Value())
				if err != nil {
					m.err = err
					return m, nil
				}
				return m, tea.Quit
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *TokenInput) View() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(titleStyle.Render("🔑 GitHub Token Configuration"))
	b.WriteString("\n\n")

	instructions := []string{
		"Para crear un Fine-grained personal access token:",
		"1. Dirigete a GitHub.com → Settings → Developer settings → Personal access tokens → Fine-grained tokens",
		"2. Click en 'Generate new token'",
		"3. Establece el nombre del token: 'Hefesto'",
		"4. Establece el tiempo de expiración",
		"5. Selecciona 'All repositories' para dar acceso a todos los repositorios",
		"6. En la sección de 'Repository permissions':",
		"   • Contents: Read and write",
		"   • Metadata: Read-only",
		"   • Administration: Read and write",
		"7. Click en 'Generate token'",
		"8. Copia el token generado y pégalo abajo",
	}

	for _, instruction := range instructions {
		b.WriteString(instructionStyle.Render(instruction))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")

	b.WriteString(helpStyle.Render("Presiona Enter para guardar • Esc para salir"))
	b.WriteString("\n")

	if m.err != nil {
		b.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			MarginLeft(2).
			Render(fmt.Sprintf("Error: %v", m.err)))
		b.WriteString("\n")
	}

	return b.String()
}
