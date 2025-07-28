package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Color scheme
	primaryColor   = lipgloss.Color("#00D9FF")
	secondaryColor = lipgloss.Color("#FF79C6")
	successColor   = lipgloss.Color("#50FA7B")
	errorColor     = lipgloss.Color("#FF5555")
	mutedColor     = lipgloss.Color("#6272A4")
	bgColor        = lipgloss.Color("#282A36")
	fgColor        = lipgloss.Color("#F8F8F2")

	// Container styles
	containerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)

	// Title bar style
	titleBarStyle = lipgloss.NewStyle().
			Background(primaryColor).
			Foreground(bgColor).
			Padding(0, 1).
			Bold(true)

	// Status bar style
	statusBarStyle = lipgloss.NewStyle().
			Background(mutedColor).
			Foreground(fgColor).
			Padding(0, 1)

	// Input field style
	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(primaryColor).
			Padding(0, 1).
			Width(60)

	// Welcome screen styles
	welcomeTitleStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Align(lipgloss.Center)

	welcomeSubtitleStyle = lipgloss.NewStyle().
				Foreground(secondaryColor).
				Align(lipgloss.Center)

	tipStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Padding(0, 2)

	// Result screen styles
	scoreStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	mistakeStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	correctStyle = lipgloss.NewStyle().
			Foreground(successColor)

	explanationStyle = lipgloss.NewStyle().
				Foreground(fgColor).
				Padding(0, 2)

	alternativeStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Italic(true)

	// Settings modal styles
	settingsContainerStyle = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(secondaryColor).
				Padding(1, 2).
				Width(50)

	settingsItemStyle = lipgloss.NewStyle().
				Padding(0, 1)

	selectedSettingStyle = lipgloss.NewStyle().
				Background(primaryColor).
				Foreground(bgColor).
				Bold(true).
				Padding(0, 1)

	// Loading animation styles
	spinnerStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	// Error message style
	errorMsgStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	// Keyboard hint style
	keyHintStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Padding(0, 1)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)
)