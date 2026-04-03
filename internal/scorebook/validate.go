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
	if batterEvent != "" && runnerEvent != "" {
		issues = append(issues, "Enter either a batter event or a base-running event, not both.")
	}
	if strings.TrimSpace(draft.Pitches) != "" && !ValidPitchString(draft.Pitches) {
		issues = append(issues, "Pitches contains a character outside simplified Retrosheet pitch syntax.")
	}

	return issues
}
