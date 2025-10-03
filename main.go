package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Input struct {
	Enabled bool                       `json:"enabled"`
	Streams map[string]json.RawMessage `json:"streams"`
}

type Config struct {
	PolicyIDs   []string          `json:"policy_ids"`
	Package     Package           `json:"package"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Namespace   string            `json:"namespace"`
	Inputs      map[string]Input  `json:"inputs"`
	Vars        map[string]string `json:"vars"`
}

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
	)

	finalModel, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	if m, ok := finalModel.(model); ok && m.output != "" {
		fmt.Println(m.output)
	}
}

type model struct {
	textarea     textarea.Model
	output       string
	err          error
	ready        bool
	windowWidth  int
	windowHeight int
}

func initialModel() model {
	ti := textarea.New()
	ti.Placeholder = "Paste your Elastic integration JSON here..."
	ti.ShowLineNumbers = true
	ti.Focus()

	return model{
		textarea: ti,
		err:      nil,
		ready:    false,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		headerHeight := 5
		footerHeight := 2
		textareaHeight := msg.Height - headerHeight - footerHeight

		m.textarea.SetWidth(msg.Width)
		m.textarea.SetHeight(textareaHeight)

		if !m.ready {
			m.ready = true
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "ctrl+l":
			m.textarea.Reset()
			return m, nil
		case "enter":
			if m.textarea.Value() != "" {
				var config Config
				if err := json.Unmarshal([]byte(m.textarea.Value()), &config); err != nil {
					m.err = fmt.Errorf("error parsing JSON: %v", err)
					return m, nil
				}
				m.output = generateTerraform(config)
				return m, tea.Quit
			}
		}

		// Update textarea
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	var s strings.Builder

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	s.WriteString(headerStyle.Render("╔══════════════════════════════════════════════════════════════╗"))
	s.WriteString("\n")
	s.WriteString(headerStyle.Render("║       Elastic Integration to Terraform Converter             ║"))
	s.WriteString("\n")
	s.WriteString(headerStyle.Render("╚══════════════════════════════════════════════════════════════╝"))
	s.WriteString("\n\n")

	s.WriteString("Paste your Elastic integration JSON below:\n\n")
	s.WriteString(m.textarea.View())
	s.WriteString("\n\n")
	if m.err != nil {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
		s.WriteString(errorStyle.Render(fmt.Sprintf("❌ Error: %v", m.err)))
		s.WriteString("\n\n")
	}
	s.WriteString("Press Enter to convert | Esc/Ctrl+C to quit | Ctrl+l to clear")

	return s.String()
}

func generateTerraform(config Config) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`resource "elasticstack_fleet_integration_policy" "%s" {
  name                = "%s"
  namespace           = "%s"
  description         = "%s"
  agent_policy_id     = "%s"
  integration_name    = "%s"
  integration_version = "%s"
`,
		config.Package.Name,
		config.Name,
		getNamespace(config.Namespace),
		config.Description,
		config.PolicyIDs[0],
		config.Package.Name,
		config.Package.Version,
	))

	// Generate vars_json
	if len(config.Vars) > 0 {
		sb.WriteString("  vars_json = jsonencode({\n")
		i := 0
		for k, v := range config.Vars {
			if i > 0 {
				sb.WriteString(",\n")
			}
			sb.WriteString(fmt.Sprintf(`    "%s" : "%s"`, k, v))
			i++
		}
		sb.WriteString("\n  })\n")
	}

	// Generate inputs
	for inputID, input := range config.Inputs {
		sb.WriteString("  input {\n")
		sb.WriteString(fmt.Sprintf("    input_id = \"%s\"\n", inputID))
		sb.WriteString(fmt.Sprintf("    enabled  = %t\n", input.Enabled))

		// Generate streams_json
		sb.WriteString("    streams_json = jsonencode({\n")
		streamCount := 0
		for streamName, streamData := range input.Streams {
			if streamCount > 0 {
				sb.WriteString(",\n")
			}
			var streamObj map[string]interface{}
			if err := json.Unmarshal(streamData, &streamObj); err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing stream: %v\n", err)
				continue
			}
			sb.WriteString(fmt.Sprintf("      \"%s\" : ", streamName))
			sb.WriteString(printJSONToString(streamObj, 6))
			streamCount++
		}
		sb.WriteString("\n    })\n")
		sb.WriteString("  }\n")
	}
	sb.WriteString("}\n")

	return sb.String()
}

func printJSONToString(obj map[string]interface{}, indent int) string {
	jsonBytes, _ := json.MarshalIndent(obj, strings.Repeat(" ", indent), "  ")
	return string(jsonBytes)
}

func getNamespace(ns string) string {
	if ns == "" {
		return "default"
	}
	return ns
}
