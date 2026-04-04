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
  background:
    radial-gradient(circle at top left, rgba(189,99,47,0.16), transparent 32%),
    radial-gradient(circle at top right, rgba(13,92,99,0.16), transparent 28%),
    linear-gradient(180deg, #f3ede2 0%, var(--bg) 100%);
}
a { color: var(--accent); }
button, input, textarea {
  font: inherit;
}
.page {
  max-width: 1100px;
  margin: 0 auto;
  padding: 10px;
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
.game-info-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 144px minmax(0, 1fr);
  gap: 6px;
  align-items: end;
}
.game-info-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 8px;
  align-items: end;
}
.game-info-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  justify-content: flex-end;
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
.log-table {
  display: grid;
  gap: 4px;
}
.log-row {
  display: grid;
  grid-template-columns: 56px 56px 70px 84px minmax(0, 1.4fr) minmax(90px, 1fr) 110px;
  gap: 6px;
  align-items: center;
  padding: 4px 0;
  border-bottom: 1px solid var(--line);
}
.log-row:last-child {
  border-bottom: 0;
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
    grid-template-columns: 50px 50px 64px 70px minmax(0, 1.2fr) minmax(80px, 0.9fr) 96px;
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
    padding-bottom: 780px;
  }
  .event-layout {
    grid-template-columns: 1fr;
  }
  .grid.two, .grid.three, .combined-grid, .game-info-grid, .context-row, .game-info-row { grid-template-columns: 1fr; }
  .keyboard-grid {
    grid-template-columns: 1fr;
  }
  .log-row {
    grid-template-columns: 44px 44px 56px 64px minmax(0, 1fr) 80px 88px;
    font-size: 0.82rem;
  }
}
`
}
