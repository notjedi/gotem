package listview

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

const (
	ellipsis = "…"
)

type Item interface {
	// list.Item
	Title() string
	Description() string
	FilterValue() string
}

type CustomDelegate struct {
	Styles        list.DefaultItemStyles
	UpdateFunc    func(tea.Msg, *list.Model) tea.Cmd
	ShortHelpFunc func() []key.Binding
	FullHelpFunc  func() [][]key.Binding
	height        int
	spacing       int
}

func NewCustomDelegate() CustomDelegate {
	return CustomDelegate{
		Styles:  list.NewDefaultItemStyles(),
		height:  2,
		spacing: 1,
	}
}

func (d *CustomDelegate) SetHeight(i int) {
	d.height = i
}

func (d CustomDelegate) Height() int {
	return d.height
}

func (d *CustomDelegate) SetSpacing(i int) {
	d.spacing = i
}

func (d CustomDelegate) Spacing() int {
	return d.spacing
}

func (d CustomDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	if d.UpdateFunc == nil {
		return nil
	}
	return d.UpdateFunc(msg, m)
}

func (d CustomDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var (
		title, desc  string
		matchedRunes []int
		s            = &d.Styles
	)

	if i, ok := item.(Item); ok {
		title = i.Title()
		desc = i.Description()
	} else {
		return
	}

	if m.Width() <= 0 {
		return
	}

	// Prevent text from exceeding list width
	textwidth := uint(m.Width() - s.NormalTitle.GetPaddingLeft() - s.NormalTitle.GetPaddingRight())
	title = truncate.StringWithTail(title, textwidth, ellipsis)
	var lines []string
	for i, line := range strings.Split(desc, "\n") {
		if i >= d.height-1 {
			break
		}
		lines = append(lines, truncate.StringWithTail(line, textwidth, ellipsis))
	}
	desc = strings.Join(lines, "\n")

	// Conditions
	var (
		isSelected  = index == m.Index()
		emptyFilter = m.FilterState() == list.Filtering && m.FilterValue() == ""
		isFiltered  = m.FilterState() == list.Filtering || m.FilterState() == list.FilterApplied
	)

	// does not work see https://github.com/charmbracelet/bubbles/issues/86
	// if isFiltered && index < len(m.itemsAsFilterItems()) {
	if isFiltered {
		matchedRunes = m.MatchesForItem(index)
	}

	if emptyFilter {
		title = s.DimmedTitle.Render(title)
		desc = s.DimmedDesc.Render(desc)
	} else if isSelected && m.FilterState() != list.Filtering {
		if isFiltered {
			// Highlight matches
			unmatched := s.SelectedTitle.Inline(true)
			matched := unmatched.Copy().Inherit(s.FilterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = s.SelectedTitle.Render(title)
		desc = s.SelectedDesc.Render(desc)
	} else {
		if isFiltered {
			// Highlight matches
			unmatched := s.NormalTitle.Inline(true)
			matched := unmatched.Copy().Inherit(s.FilterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = s.NormalTitle.Render(title)
		desc = s.NormalDesc.Render(desc)
	}

	fmt.Fprintf(w, "%s\n%s", title, desc)
	return
}

func (d CustomDelegate) ShortHelp() []key.Binding {
	if d.ShortHelpFunc != nil {
		return d.ShortHelpFunc()
	}
	return nil
}

func (d CustomDelegate) FullHelp() [][]key.Binding {
	if d.FullHelpFunc != nil {
		return d.FullHelpFunc()
	}
	return nil
}
