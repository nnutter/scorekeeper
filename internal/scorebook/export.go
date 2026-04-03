package scorebook

import (
	"fmt"
	"net/url"
	"strings"
)

func ExportText(book Book) string {
	var lines []string
	lines = append(lines, fmt.Sprintf("game,%s,away=%s,home=%s", safe(book.Meta.GameDate), safe(book.Meta.AwayTeam), safe(book.Meta.HomeTeam)))
	lines = append(lines, "")

	for _, entry := range book.Entries {
		if entry.Mode == ModeRun {
			lines = append(lines, fmt.Sprintf("run,%d,%s,pitcher=%s,batter=%s,event=%s", entry.Inning, entry.Half, safe(entry.Pitcher), safe(entry.Batter), safe(entry.RunnerEvent)))
		} else {
			parts := []string{
				fmt.Sprintf("play,%d,%s", entry.Inning, entry.Half),
				"pitcher=" + safe(entry.Pitcher),
				"batter=" + safe(entry.Batter),
			}
			if entry.Pitches != "" {
				parts = append(parts, "pitches="+safe(entry.Pitches))
			}
			parts = append(parts, "event="+safe(entry.BatterEvent))
			if entry.Advances != "" {
				parts = append(parts, "adv="+safe(entry.Advances))
			}
			lines = append(lines, strings.Join(parts, ","))
		}
		if entry.Note != "" {
			lines = append(lines, fmt.Sprintf("note,%d,%s,%q", entry.Inning, entry.Half, entry.Note))
		}
		lines = append(lines, "")
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func MailtoLink(book Book) string {
	subject := fmt.Sprintf("Scorekeeper %s at %s %s", safe(book.Meta.AwayTeam), safe(book.Meta.HomeTeam), safe(book.Meta.GameDate))
	body := ExportText(book)
	return "mailto:?subject=" + url.QueryEscape(subject) + "&body=" + url.QueryEscape(body)
}

func safe(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), ",", " ")
}
