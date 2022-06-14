package tui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func headerView(m *Model) *tview.Flex {
	flex := tview.NewFlex().SetDirection(tview.FlexColumn)
	flex.SetBorder(false)
	flex.SetBorderPadding(0, 0, 1, 1)

	flex.AddItem(helpView(defaultHelpKeys), 0, 3, false)
	flex.AddItem(schemaVersionView(m), 0, 1, false)
	return flex
}

func schemaVersionView(m *Model) *tview.Table {
	tbl := tview.NewTable().SetBorders(false)
	tbl.SetBorder(false)

	titleCell := func(s string) *tview.TableCell {
		return tview.NewTableCell(s).
			SetStyle(textStyle.Bold(true).Foreground(tcell.ColorDarkOrange))
	}
	tbl.SetCell(0, 0, titleCell("Database:"))
	tbl.SetCell(1, 0, titleCell("Schema Version:"))
	tbl.SetCell(2, 0, titleCell("Schema Date:"))

	valCell := func(s string) *tview.TableCell {
		return tview.NewTableCell(s).SetStyle(textStyle)
	}
	tbl.SetCell(0, 1, valCell(m.databasePath))
	tbl.SetCell(1, 1, valCell(m.schemaVersion))
	tbl.SetCell(2, 1, valCell(formatTime(m.schemaDate)))

	return tbl
}

func helpView(keys keyMap) *tview.Table {
	tbl := tview.NewTable().SetBorders(false)
	tbl.SetBorder(false)

	keyCell := func(s string) *tview.TableCell {
		return tview.NewTableCell("<" + s + ">").
			SetTextColor(tcell.ColorBlue)
	}
	textCell := func(s string) *tview.TableCell {
		return tview.NewTableCell(s).
			SetStyle(textStyle.Attributes(tcell.AttrDim))
	}
	for row, binding := range keys {
		tbl.SetCell(row, 0, keyCell(binding.Key))
		tbl.SetCell(row, 1, textCell(binding.Help))
	}
	return tbl
}

func testCasesView(m *Model) *tview.Table {
	tbl := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetSelectedStyle(tcell.Style{}.Foreground(backgroundColor).Background(textColor))
	tbl.
		SetBorder(true).
		SetBorderPadding(0, 0, 1, 1).
		SetBorderAttributes(tcell.AttrDim)
	tbl.SetTitle("Test Cases")

	headerCell := func(s string) *tview.TableCell {
		s = strings.ToUpper(s)
		return tview.NewTableCell(s).
			SetStyle(textStyle.Bold(true)).
			SetExpansion(1).
			SetSelectable(false)
	}

	for col, title := range []string{
		"ID",
		"Date",
		"Name",
		"Git Sha",
		"Chain",
		"Height",
		"Tx Total",
	} {
		tbl.SetCell(0, col, headerCell(title))
	}

	contentCell := func(s string) *tview.TableCell {
		return tview.NewTableCell(s).SetStyle(textStyle).SetExpansion(1)
	}

	for i, tc := range m.testCases {
		row := i + 1 // 1 offsets header row
		pres := testCasePresenter{tc}
		for col, content := range []string{
			pres.ID(),
			pres.Date(),
			pres.Name(),
			pres.GitSha(),
			pres.ChainID(),
			pres.Height(),
			pres.TxTotal(),
		} {
			tbl.SetCell(row, col, contentCell(content))
		}
	}
	return tbl
}