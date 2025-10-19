package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var blockStyles = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Width(70).Foreground(lipgloss.Color("88")).Align(lipgloss.Center)

func printQuestion(q *Question) {
	if q == nil {
		style := blockStyles.Bold(true)
		fmt.Println(style.Render("No Question"))
		return
	}
	fmt.Println(blockStyles.Render(lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.NewStyle().Bold(true).Render(q.Name),
		q.Group,
		lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Render(q.Url),
	)))

}

func printGreen(s string) {
	greenStyle := blockStyles.Foreground(lipgloss.Color("34")).Bold(true)
	fmt.Println(greenStyle.Render(s))
}

func printRed(s string) {
	redStyle := blockStyles.Foreground(lipgloss.Color("160")).Bold(true)
	fmt.Println(redStyle.Render(s))
}

func printNormal(s string) {
	normalStyle := blockStyles.Foreground(lipgloss.Color("88")).Bold(true)
	fmt.Println(normalStyle.Render(s))
}

func printError(header string, output string) {
	if len(output) > 0 {
		t := table.New().
			Border(lipgloss.NormalBorder()).
			Width(72).Headers(header).
			Row(output).StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("88")).Align(lipgloss.Center)
			default:
				return lipgloss.NewStyle().Foreground(lipgloss.Color("160")) // red for errors
			}
		})
		fmt.Println(t)
	}
}
