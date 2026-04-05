package ui

func CSS() string {
	return `
:root {
  color-scheme: light;
  --bg: #f7f3eb;
  --panel: rgba(255,255,255,0.82);
  --panel-strong: #fffdf8;
  --ink: #1e1d1a;
  --muted: #675f55;
  --line: rgba(76, 60, 41, 0.18);
  --accent: #0d5c63;
  --accent-soft: #d9efef;
  --warm: #bd632f;
  --danger: #a33636;
  --shadow: 0 24px 50px rgba(56, 40, 20, 0.12);
}
* { box-sizing: border-box; }
body {
  margin: 0;
  font-family: Georgia, "Times New Roman", serif;
  color: var(--ink);
  overflow-x: hidden;
  background:
    radial-gradient(circle at top left, rgba(189,99,47,0.16), transparent 32%),
    radial-gradient(circle at top right, rgba(13,92,99,0.16), transparent 28%),
    linear-gradient(180deg, #f3ede2 0%, var(--bg) 100%);
}
a { color: var(--accent); }
.mobile-only { display: none; }
button, input, textarea {
  font: inherit;
}
button, a, input, textarea, summary {
  touch-action: manipulation;
}
.page {
  width: 100%;
  max-width: 1100px;
  margin: 0 auto;
  padding: 10px;
}
.pull-refresh-indicator {
  position: fixed;
  top: calc(max(12px, env(safe-area-inset-top)) + 10px);
  left: 50%;
  transform: translate(-50%, -24px);
  opacity: 0;
  pointer-events: none;
  z-index: 40;
  min-width: 156px;
  padding: 10px 16px;
  border-radius: 999px;
  border: 1px solid rgba(76, 60, 41, 0.3);
  background: rgba(255, 253, 248, 0.99);
  color: #4a3d31;
  box-shadow: 0 14px 28px rgba(56, 40, 20, 0.22);
  font-size: 0.92rem;
  font-weight: 600;
  letter-spacing: 0.01em;
  text-align: center;
  transition: opacity 120ms ease, transform 120ms ease, color 120ms ease, border-color 120ms ease, background 120ms ease, box-shadow 120ms ease;
}
.pull-refresh-indicator.visible {
  opacity: 1;
}
.pull-refresh-indicator.ready {
  color: #f7fffe;
  border-color: rgba(13,92,99,0.7);
  background: rgba(13, 92, 99, 0.96);
  box-shadow: 0 16px 30px rgba(13, 92, 99, 0.28);
}
.stack {
  display: grid;
  gap: 8px;
}
.main-stack {
  margin-top: 0;
}
.event-layout {
  display: grid;
  grid-template-columns: 116px minmax(0, 1fr);
  gap: 8px;
  align-items: start;
}
.panel {
  min-width: 0;
  padding: 8px 10px;
  border: 1px solid var(--line);
  border-radius: 14px;
  background: var(--panel);
  backdrop-filter: blur(16px);
  box-shadow: var(--shadow);
}
.panel h2, .panel h3 {
  display: none;
}
.grid, .entry-grid {
  display: grid;
  gap: 8px;
}
.grid.two { grid-template-columns: repeat(2, minmax(0, 1fr)); }
.grid.three { grid-template-columns: repeat(3, minmax(0, 1fr)); }
.game-info-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 144px minmax(0, 1fr) auto auto auto;
  gap: 8px;
  align-items: end;
  grid-template-areas: "away date home new copy email";
}
.game-away { grid-area: away; }
.game-date { grid-area: date; }
.game-home { grid-area: home; }
.game-new { grid-area: new; }
.game-copy { grid-area: copy; }
.game-email { grid-area: email; }
.game-info-layout .btn {
  justify-self: end;
}
.game-info-layout .input,
.game-info-layout .btn {
  height: 35px;
}
.game-info-layout input[type="date"].input {
  appearance: none;
  -webkit-appearance: none;
  display: block;
  line-height: 1.2;
  min-height: 0;
}
.game-info-layout .btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
.action-icon-btn {
  width: 35px;
  min-width: 35px;
  padding: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
.save-action-btn {
  background: #fffdf9;
  border-color: var(--line);
}
.action-icon {
  width: 18px;
  height: 18px;
  display: block;
}
.context-row {
  display: grid;
  grid-template-columns: 130px minmax(0, 1fr) 140px;
  gap: 6px;
  align-items: end;
}
.context-panel {
  align-self: start;
}
.context-layout {
  margin-top: 0;
}
.combined-grid {
  grid-template-columns: 0.95fr 0.9fr 1.1fr 1fr 1.1fr;
  align-items: end;
}
.combined-grid .field:last-child {
  grid-column: 1 / -1;
}
.field {
  display: grid;
  gap: 3px;
}
.field label, .mini-label {
  font-size: 0.82rem;
  color: var(--muted);
}
.field-label {
  line-height: 1;
}
.input, .textarea, .context-chip {
  width: 100%;
  border: 1px solid var(--line);
  border-radius: 10px;
  padding: 7px 9px;
  background: rgba(255,255,255,0.95);
}
.textarea { min-height: 54px; resize: vertical; }
.toolbar, .mode-toggle, .keyboard-row, .actions-row {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}
.actions-row {
  margin-top: 7px;
}
.btn {
  border: 1px solid var(--line);
  border-radius: 999px;
  padding: 6px 10px;
  background: #fffdf9;
  cursor: pointer;
}
.btn.primary {
  background: var(--accent);
  border-color: var(--accent);
  color: #f7fffe;
}
.btn.warm {
  background: var(--warm);
  border-color: var(--warm);
  color: #fff8f3;
}
.btn.ghost.active {
  background: var(--accent-soft);
  border-color: rgba(13,92,99,0.38);
}
.btn.danger {
  color: var(--danger);
}
.context-chip strong {
  display: block;
  margin-top: 2px;
  font-size: 0.95rem;
}
.context-chip.compact {
  padding-top: 6px;
  padding-bottom: 6px;
}
.action-field label, .context-actions label {
  visibility: hidden;
}
.full-width {
  width: 100%;
}
.context-action-row {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 6px;
}
.context-step {
  width: 100%;
  padding-left: 0;
  padding-right: 0;
}
.keyboard-group {
  display: grid;
  gap: 6px;
  padding: 6px;
  border-radius: 12px;
  border: 1px solid var(--line);
  background: rgba(255, 251, 244, 0.82);
}
.keyboard-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
}
.keyboard-mobile {
  display: none;
}
.keyboard-mobile-rail {
  display: grid;
  grid-template-rows: repeat(2, 96px);
  align-content: start;
  gap: 6px;
}
.keyboard-mobile-main {
  min-width: 0;
}
.keyboard-mobile-pane {
  display: none;
}
.keyboard-mobile-pane.active {
  display: block;
}
.keyboard-mobile-main .keyboard-group {
  height: 234px;
  align-content: start;
  overflow: hidden;
}
.keyboard-switch {
  width: 32px;
  min-width: 32px;
  height: 96px;
  min-height: 0;
  padding: 0;
  border-radius: 12px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
.keyboard-switch.active {
  background: var(--accent-soft);
  border-color: rgba(13,92,99,0.38);
}
.keyboard-row {
  gap: 6px;
}
.keyboard-panel {
  position: sticky;
  bottom: 12px;
  z-index: 20;
}
.token {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 40px;
  height: 32px;
  padding: 5px 8px;
  text-align: center;
  white-space: nowrap;
}
.notice {
  padding: 4px 8px;
  border-radius: 999px;
  background: #fff1e8;
  border: 1px solid rgba(189,99,47,0.25);
  color: #7b431d;
  font-size: 0.78rem;
  justify-self: start;
}
.notice.status {
  background: #e7f6f0;
  border-color: rgba(13,92,99,0.22);
  color: var(--accent);
}
.entry-list {
  display: grid;
  gap: 0;
}
.log-entry {
  padding: 4px 0;
  border-bottom: 1px solid var(--line);
}
.log-entry:last-child {
  border-bottom: 0;
}
.log-table {
  display: grid;
  gap: 4px;
}
.log-row {
  display: grid;
  grid-template-columns: 56px 56px 70px 84px minmax(0, 1.4fr) 110px;
  gap: 6px;
  align-items: center;
}
.log-row span {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.log-header {
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: var(--muted);
  padding-top: 0;
}
.log-note {
  color: var(--muted);
  white-space: normal;
  overflow: visible;
  text-overflow: clip;
}
.log-note-row {
  display: grid;
  grid-template-columns: 44px minmax(0, 1fr);
  gap: 8px;
  margin-top: 3px;
  font-size: 0.84rem;
}
.log-note-label {
  color: var(--muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
  font-size: 0.72rem;
  padding-top: 1px;
}
.log-actions {
  display: flex;
  gap: 4px;
  justify-content: flex-end;
}
.log-actions .btn {
  padding: 4px 8px;
  font-size: 0.78rem;
}
.icon-btn {
  width: 30px;
  height: 30px;
  padding: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
.icon-btn img {
  width: 15px;
  height: 15px;
  display: block;
}
.log-empty {
  padding: 8px 0 2px;
  color: var(--muted);
  font-size: 0.88rem;
}
.meta-line {
  color: var(--muted);
  font-size: 0.78rem;
}
.export-box {
  min-height: 200px;
  white-space: pre-wrap;
}
.export-details summary {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  color: var(--accent);
  font-size: 0.9rem;
}
.export-details summary::-webkit-details-marker {
  display: none;
}
.export-details summary::marker {
  content: "";
}
.export-details summary::after {
  content: "v";
  font-size: 0.82rem;
}
.export-details:not([open]) summary::after {
  content: "<";
}
.game-info-row + .export-details {
  margin-top: 8px;
}
@media (max-width: 900px) {
  .combined-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
  .log-row {
    grid-template-columns: 50px 50px 64px 70px minmax(0, 1.2fr) 96px;
  }
}
@media (max-width: 1100px) {
  .page {
    padding-bottom: 520px;
  }
  .keyboard-panel {
    position: fixed;
    left: 12px;
    right: 12px;
    bottom: 10px;
    margin: 0 auto;
    max-width: 1076px;
    background: rgba(247, 243, 235, 0.92);
    backdrop-filter: blur(18px);
    box-shadow: 0 18px 32px rgba(56, 40, 20, 0.18);
  }
  .keyboard-panel .stack {
    gap: 6px;
  }
  .keyboard-panel h2,
  .keyboard-panel .meta-line {
    display: none;
  }
}
@media (max-width: 720px) {
  .page {
    padding: 10px;
    padding-bottom: 272px;
  }
  .main-stack,
  .event-layout,
  .log-table,
  .log-row,
  .keyboard-mobile,
  .keyboard-mobile-main {
    min-width: 0;
  }
  .keyboard-panel {
    left: 10px;
    right: 10px;
    bottom: 10px;
    overflow: hidden;
    border-bottom-left-radius: 62px;
    border-bottom-right-radius: 62px;
  }
  .event-layout {
    grid-template-columns: 1fr;
  }
  .game-info-layout {
    grid-template-columns: minmax(0, 1fr) auto;
    grid-template-areas:
      "away new"
      "date copy"
      "home email";
    gap: 7px 8px;
  }
  .game-info-layout input[type="date"].input {
    text-align: left;
    color: var(--ink);
    -webkit-text-fill-color: var(--ink);
    justify-content: flex-start;
  }
  .game-info-layout input[type="date"].input::-webkit-date-and-time-value {
    text-align: left;
  }
  .game-info-layout input[type="date"].input::-webkit-datetime-edit {
    color: var(--ink);
    padding: 0;
  }
  .game-info-layout input[type="date"].input::-webkit-datetime-edit-fields-wrapper {
    justify-content: flex-start;
  }
  .game-info-layout .btn {
    justify-self: stretch;
    align-self: end;
  }
  .context-layout {
    display: grid;
    grid-template-columns: 70px 34px 34px 10px minmax(0, 1fr);
    grid-template-areas: "inning actions actions . pitcher";
    gap: 6px;
    align-items: end;
  }
  .context-inning {
    grid-area: inning;
    max-width: 72px;
  }
  .context-inning .field-label,
  .context-actions label,
  .context-pitcher .field-label {
    display: none;
  }
  .context-actions {
    grid-area: actions;
  }
  .context-action-row {
    grid-template-columns: repeat(2, 1fr);
  }
  .context-pitcher {
    grid-area: pitcher;
    min-width: 0;
    justify-self: end;
    width: 100%;
    max-width: 96px;
  }
  .context-chip,
  .context-actions .btn,
  .context-pitcher .input {
    height: 35px;
  }
  .context-chip.compact {
    display: flex;
    align-items: center;
    padding-top: 0;
    padding-bottom: 0;
  }
  .combined-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
  .combined-grid .field:nth-child(1) {
    grid-column: 1;
    grid-row: 1;
  }
  .combined-grid .field:nth-child(2) {
    grid-column: 2;
    grid-row: 1;
  }
  .combined-grid .field:nth-child(3) {
    grid-column: 1;
    grid-row: 2;
  }
  .combined-grid .field:nth-child(5) {
    grid-column: 2;
    grid-row: 2;
  }
  .combined-grid .field:nth-child(4) {
    grid-column: 1 / -1;
    grid-row: 3;
  }
  .combined-grid .field:last-child {
    grid-column: 1 / -1;
    grid-row: 4;
  }
  .grid.two, .grid.three, .context-row { grid-template-columns: 1fr; }
  .desktop-only { display: none; }
  .mobile-only { display: inline; }
  .log-table {
    width: 100%;
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
  }
  .log-row {
    min-width: 0;
  }
  .keyboard-grid {
    display: none;
  }
  .keyboard-mobile {
    display: grid;
    grid-template-columns: 32px minmax(0, 1fr) 32px;
    gap: 6px;
    align-items: stretch;
  }
  .log-row {
    grid-template-columns: 44px 44px 56px 64px minmax(0, 1fr) 88px;
    font-size: 0.82rem;
  }
}
`
}
