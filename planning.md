# Plan: Interactive TUI for gittree

**Objective:** Implement an interactive terminal UI to navigate the branch tree and checkout branches.

**Libs:**
- `github.com/charmbracelet/bubbletea` (core)
- `github.com/charmbracelet/lipgloss` (styling)

**1. Git Layer (`internal/git`)**
- `repository.go`: Add `Checkout(branchName string) error`.
  - Use `worktree.Checkout(&git.CheckoutOptions{...})`.

**2. Tree Layer (`internal/tree`)**
- **Helper:** Create a function to flatten the `*Tree` (recursive structure) into a linear `[]TuiItem` list for the UI loop.
- **Structure:** `TuiItem` needs `BranchName`, `Depth`, `IsCurrent`, `DisplayString` (pre-calculated tree prefix like `└── `).

**3. TUI Layer (`internal/tui`)**
- **Model:**
  - `items`: `[]TuiItem`.
  - `cursor`: `int` (index of selected item).
  - `repo`: `*git.Repository` (for checkout action).
- **Update:**
  - `Up/k`: `cursor--`.
  - `Down/j`: `cursor++`.
  - `Enter`: Call `repo.Checkout(items[cursor].BranchName)`, print success/error, `tea.Quit`.
  - `q/Esc/Ctrl+c`: `tea.Quit`.
- **View:**
  - Iterate `items`.
  - Render line using `lipgloss`.
  - Highlight line at `cursor`.

**4. CLI Integration (`cmd`)**
- Create `cmd/ui.go`.
- Add `ui` subcommand to `rootCmd`.
- Initialize and run `tea.NewProgram`.

**5. Verify**
- Run `gittree ui`.
- Navigate with arrows.
- Press Enter to checkout.
- Verify active branch with `git status`.
