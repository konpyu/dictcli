// Package logger provides example integration patterns for TUI logging
package logger

// Example of how to integrate logger with TUI updates:
//
// In your TUI model (internal/tui/model.go):
//
// type Model struct {
//     state   State
//     logger  *logger.TUILogger
//     // ... other fields
// }
//
// In your TUI update function (internal/tui/update.go):
//
// func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//     // Log state at the beginning
//     currentState := m.state.String()
//
//     switch msg := msg.(type) {
//     case tea.KeyMsg:
//         // Log key presses
//         m.logger.KeyPress(msg.String(), currentState)
//
//         switch msg.String() {
//         case "r", "R":
//             if m.state == StatePlaying || m.state == StateListening {
//                 m.logger.AudioPlayback("replay", m.audioPath)
//                 // Handle replay logic
//             }
//         case "s", "S":
//             oldState := m.state
//             m.state = StateSettings
//             m.logger.StateTransition(oldState.String(), m.state.String())
//         case "enter":
//             if m.state == StateListening {
//                 m.logger.UserInput(m.input, currentState)
//                 // Transition to grading
//                 m.state = StateGrading
//                 m.logger.StateTransition(currentState, m.state.String())
//             }
//         }
//
//     case GenerateSentenceMsg:
//         // Log API calls
//         start := time.Now()
//         // ... call API
//         duration := time.Since(start).Milliseconds()
//         m.logger.APICall("OpenAI", "GenerateSentence", duration)
//
//     case GradeResultMsg:
//         // Log grading results
//         m.logger.Grade(msg.WER, msg.Score, len(msg.Mistakes))
//         m.state = StateShowingResult
//         m.logger.StateTransition(currentState, m.state.String())
//
//     case SettingsUpdateMsg:
//         // Log settings changes
//         m.logger.Settings(msg.Field, msg.OldValue, msg.NewValue)
//     }
//
//     return m, nil
// }
//
// In your main.go:
//
// func main() {
//     // Parse flags
//     debug := flag.Bool("debug", false, "Enable debug logging")
//     flag.Parse()
//
//     // Initialize logger
//     logDir := filepath.Join(".", "logs")
//     if err := logger.Init(*debug, logDir); err != nil {
//         log.Fatalf("Failed to initialize logger: %v", err)
//     }
//     defer logger.Close()
//
//     // Initialize TUI logger
//     tuiLogger := logger.InitTUILogger()
//
//     // Create TUI model with logger
//     model := NewModel(tuiLogger)
//
//     // Start TUI
//     p := tea.NewProgram(model)
//     if _, err := p.Run(); err != nil {
//         logger.Error("Failed to run TUI: %v", err)
//         os.Exit(1)
//     }
// }