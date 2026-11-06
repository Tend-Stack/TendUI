package components

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func TestDataTableAlignmentAcrossHeaderAndDataCells(t *testing.T) {
	for _, tc := range []struct {
		value TableAlignment
		want  string
	}{{AlignStart, "start"}, {AlignCenter, "center"}, {AlignEnd, "end"}, {"", "start"}, {`bad" onmouseover="x`, "start"}} {
		p := DataTableProps{Caption: "Alignment", Columns: []TableColumn{{Label: "Name", RowHeader: true, Alignment: tc.value}, {Label: "Value", Alignment: tc.value}}, Rows: []TableRow{{Cells: []TableCell{{Text: "Row"}, {Text: "Value"}}}}}
		html := render(t, DataTable(p))
		if strings.Count(html, `data-align="`+tc.want+`"`) != 4 || strings.Count(html, "data-align=") != 4 || strings.Contains(html, " onmouseover=") || strings.Contains(html, " else") {
			t.Fatalf("alignment %q produced malformed or duplicate attributes: %s", tc.value, html)
		}
	}
}

func TestDataTablePreservesCellContextAndErrorPropagation(t *testing.T) {
	type contextKey struct{}
	ctx := templ.WithNonce(context.WithValue(context.Background(), contextKey{}, "kept"), "nonce-kept")
	ctx = templ.WithChildren(ctx, templ.Raw("unrelated children must not leak"))
	slot := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if ctx.Value(contextKey{}) != "kept" || templ.GetNonce(ctx) != "nonce-kept" {
			return errors.New("cell context lost")
		}
		if err := templ.GetChildren(ctx).Render(ctx, w); err != nil {
			return err
		}
		return Badge(BadgeProps{Text: "<&>"}).Render(ctx, w)
	})
	p := DataTableProps{Caption: "Slots", Columns: []TableColumn{{Label: "Header", RowHeader: true}, {Label: "Data"}}, Rows: []TableRow{{Cells: []TableCell{{Text: "ignored fallback", Content: slot}, {Text: "ignored fallback", Content: slot}}}}}
	var b bytes.Buffer
	if err := DataTable(p).Render(ctx, &b); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(b.String(), "ignored fallback") || strings.Contains(b.String(), "unrelated children") || strings.Count(b.String(), "&lt;&amp;&gt;") != 2 {
		t.Fatal(b.String())
	}
	sentinel := errors.New("cell render failed")
	for i := range p.Rows[0].Cells {
		original := p.Rows[0].Cells[i]
		p.Rows[0].Cells[i].Content = templ.ComponentFunc(func(context.Context, io.Writer) error { return sentinel })
		if err := DataTable(p).Render(ctx, io.Discard); !errors.Is(err, sentinel) {
			t.Fatalf("cell %d lost its error: %v", i, err)
		}
		p.Rows[0].Cells[i] = original
	}
}

func TestDataTableHonorsCancellationBetweenCells(t *testing.T) {
	for _, rowHeader := range []bool{false, true} {
		ctx, cancel := context.WithCancel(context.Background())
		stop := templ.ComponentFunc(func(context.Context, io.Writer) error { cancel(); return nil })
		p := DataTableProps{Caption: "Cancellation", Columns: []TableColumn{{Label: "Cancel"}, {Label: "Next", RowHeader: rowHeader}}, Rows: []TableRow{{Cells: []TableCell{{Content: stop}, {Text: "after-cancel"}}}}}
		var b bytes.Buffer
		err := DataTable(p).Render(ctx, &b)
		cancel()
		if !errors.Is(err, context.Canceled) || strings.Contains(b.String(), "after-cancel") {
			t.Fatalf("render continued after cancellation: %v %s", err, b.String())
		}
	}
}
