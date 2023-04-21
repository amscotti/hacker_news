package hackernews

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	HackerNewsURLFormat = "https://news.ycombinator.com/item?id=%d"
)

type Story struct {
	Title       string
	By          string
	Descendants int
	Id          int
	Kids        []int
	Score       int
	Time        int
	Type        string
	URL         string
}

func newStyleWithFGColor(color string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
}

func (s *Story) PrintStyling(showSourceUrl bool) string {
	titleStyle := newStyleWithFGColor("#FAB005").Bold(true)
	metaStyle := newStyleWithFGColor("#B8B8B8")
	linkStyle := newStyleWithFGColor("#69C9D0").Underline(true)

	output := titleStyle.Render(s.Title) + "\n"
	output += metaStyle.Render(fmt.Sprintf("score: %d\tcomments: %d\tuser: %s", s.Score, s.Descendants, s.By)) + "\n"
	output += "url: " + linkStyle.Render(fmt.Sprintf(HackerNewsURLFormat, s.Id)) + "\n"
	if showSourceUrl {
		output += newStyleWithFGColor("#D6749C").Render(s.URL) + "\n"
	}
	output += "\n"

	return output
}
