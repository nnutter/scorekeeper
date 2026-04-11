package scorebook

import "strings"

func Validate(meta GameMeta, ctx GameContext, draft EventDraft) []string {
	var issues []string
	batterEvent := strings.TrimSpace(draft.BatterEvent)
	runnerEvent := strings.TrimSpace(draft.RunnerEvent)

	if strings.TrimSpace(meta.AwayTeam) == "" {
		issues = append(issues, "Away team is required.")
	}
	if strings.TrimSpace(meta.HomeTeam) == "" {
		issues = append(issues, "Home team is required.")
	}
	if strings.TrimSpace(meta.GameDate) == "" {
		issues = append(issues, "Game date is required.")
	}
	if ctx.Inning < 1 {
		issues = append(issues, "Inning must be 1 or higher.")
	}
	if strings.TrimSpace(ctx.Pitcher) == "" {
		issues = append(issues, "Pitcher is required.")
	}
	if strings.TrimSpace(draft.Batter) == "" {
		issues = append(issues, "Batter is required.")
	}

	if batterEvent == "" && runnerEvent == "" {
		issues = append(issues, "Enter a batter event or a base-running event.")
	}
	if batterEvent != "" && runnerEvent != "" && !validCombinedEvent(batterEvent, runnerEvent) {
		issues = append(issues, "Only K, W, and IW may be combined with SB, CS, OA, PO, PB, WP, or E base-running events.")
	}
	if strings.TrimSpace(draft.Pitches) != "" && !ValidPitchString(draft.Pitches) {
		issues = append(issues, "Pitches contains a character outside simplified Retrosheet pitch syntax.")
	}

	return issues
}

func validCombinedEvent(batterEvent, runnerEvent string) bool {
	switch batterEvent {
	case "K", "W", "IW":
		return validCombinedRunnerEvent(runnerEvent)
	default:
		return false
	}
}

func validCombinedRunnerEvent(runnerEvent string) bool {
	parts := strings.Split(runnerEvent, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if !isAllowedCombinedRunnerToken(part) {
			return false
		}
	}
	return true
}

func isAllowedCombinedRunnerToken(runnerEvent string) bool {
	return strings.HasPrefix(runnerEvent, "SB") ||
		strings.HasPrefix(runnerEvent, "CS") ||
		runnerEvent == "OA" ||
		strings.HasPrefix(runnerEvent, "PO") ||
		runnerEvent == "PB" ||
		runnerEvent == "WP" ||
		strings.HasPrefix(runnerEvent, "E")
}
