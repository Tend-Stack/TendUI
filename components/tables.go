package components

import (
	"fmt"
	"strings"

	"github.com/a-h/templ"
)

type SortDirection string
type TableAlignment string

const (
	SortNone       SortDirection  = ""
	SortAscending  SortDirection  = "ascending"
	SortDescending SortDirection  = "descending"
	AlignStart     TableAlignment = "start"
	AlignCenter    TableAlignment = "center"
	AlignEnd       TableAlignment = "end"
)

type TableColumn struct {
	Label string
	// SortURL is application-owned navigation; the component never sorts data.
	SortURL   string
	SortLabel string
	Sort      SortDirection
	Alignment TableAlignment
	RowHeader bool
}

type TableCell struct {
	Text string
	// Content takes precedence over Text and accepts a trusted templ component.
	Content templ.Component
}

type TableRow struct {
	Cells []TableCell
	// Selected is a visual state. Render a real checkbox in a cell for selection.
	Selected bool
}

type DataTableProps struct {
	Caption      string
	Columns      []TableColumn
	Rows         []TableRow
	Striped      bool
	Compact      bool
	Loading      bool
	LoadingLabel string
	Empty        templ.Component
}

func validateTable(p DataTableProps) error {
	if strings.TrimSpace(p.Caption) == "" || len(p.Columns) == 0 {
		return fmt.Errorf("tendui: DataTable requires a caption and columns")
	}
	sorted, rowHeaders := 0, 0
	for _, c := range p.Columns {
		if strings.TrimSpace(c.Label) == "" {
			return fmt.Errorf("tendui: table columns require labels")
		}
		switch c.Sort {
		case SortNone:
		case SortAscending, SortDescending:
			sorted++
		default:
			return fmt.Errorf("tendui: invalid table sort direction")
		}
		if c.RowHeader {
			rowHeaders++
		}
	}
	if sorted > 1 || rowHeaders > 1 {
		return fmt.Errorf("tendui: table permits one sorted column and one row header column")
	}
	if !p.Loading {
		for i, row := range p.Rows {
			if len(row.Cells) != len(p.Columns) {
				return fmt.Errorf("tendui: table row %d has %d cells; expected %d", i, len(row.Cells), len(p.Columns))
			}
		}
	}
	return nil
}
