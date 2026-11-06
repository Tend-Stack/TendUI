package components

import (
	"bytes"
	"context"
	"io"
	"strconv"
	"testing"

	"github.com/a-h/templ"
)

func BenchmarkRender(b *testing.B) {
	table := func(count int) templ.Component {
		p := DataTableProps{Caption: "Benchmark projects", Columns: []TableColumn{{Label: "Project", RowHeader: true}, {Label: "Status"}, {Label: "Hours", Alignment: AlignEnd}}}
		for i := 0; i < count; i++ {
			p.Rows = append(p.Rows, TableRow{Cells: []TableCell{{Text: "Project " + strconv.Itoa(i)}, {Text: "Active"}, {Text: "24"}}})
		}
		return DataTable(p)
	}
	for name, component := range map[string]templ.Component{
		"Card":    Card(CardProps{Body: templ.Raw("<h2>A small project</h2><p>A bounded component benchmark.</p>")}),
		"Table10": table(10), "Table100": table(100),
		"PageTable100": Layout(LayoutProps{Title: "Projects", Theme: ThemeDark}, table(100)),
	} {
		b.Run(name, func(b *testing.B) {
			ctx := context.Background()
			var sample bytes.Buffer
			if err := component.Render(ctx, &sample); err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(sample.Len()))
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := component.Render(ctx, io.Discard); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
