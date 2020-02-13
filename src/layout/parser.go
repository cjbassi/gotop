package layout

import (
	"bufio"
	"io"
	"log"
	"strconv"
	"strings"
)

// The syntax for the layout specification is:
// ```
// (rowspan:)?widget(/weight)?
// ```
// 1. Each line is a row
// 2. Empty lines are skipped
// 3. Spaces are compressed
// 4. Legal widget names are: cpu, disk, mem, temp, batt, net, procs
// 5. Names are not case sensitive
// 4. The simplest row is a single widget, by name, e.g.
//    ```
//    cpu
//    ```
// 5. Widgets with no weights have a weight of 1.
// 6. If multiple widgets are put on a row with no weights, they will all have
//    the same width.
// 7. Weights are integers
// 8. A widget will have a width proportional to its weight divided by the
//    total weight count of the row. E.g.,
//    ```
//    cpu      net
//    disk/2   mem/4
//    ```
//    The first row will have two widgets: the CPU and network widgets; each
//    will be 50% of the total width wide.  The second row will have two
//    widgets: disk and memory; the first will be 2/6 ~= 33% wide, and the
//    second will be 5/7 ~= 67% wide (or, memory will be twice as wide as disk).
// 9. If prefixed by a number and colon, the widget will span that number of
//    rows downward. E.g.
//    ```
//    2:cpu
//    mem
//    ```
//    The CPU widget will be twice as high as the memory widget.  Similarly,
//    ```
//    mem   2:cpu
//    net
//    ```
//    memory and network will be in the same row as CPU, one over the other,
//    and each half as high as CPU.
// 10. Negative, 0, or non-integer weights will be recorded as "1".  Same for row spans.
// 11. Unrecognized widgets will cause the application to abort.
// 12. In rows with multi-row spanning widgets **and** weights, weights in
//     lower rows are ignored.  Put the weight on the widgets in that row, not
//     in later (spanned) rows.
// 13. Widgets are filled in top down, left-to-right order.
// 14. The larges row span in a row defines the top-level row span; all smaller
//     row spans constitude sub-rows in the row. For example, `cpu mem/3 net/5`
//     means that net/5 will be 5 rows tall overall, and mem will compose 3 of
//     them. If following rows do not have enough widgets to fill the gaps,
//     spacers will be used.
func ParseLayout(i io.Reader) layout {
	r := bufio.NewScanner(i)
	rv := layout{Rows: make([][]widgetRule, 0)}
	var lineNo int
	for r.Scan() {
		l := strings.TrimSpace(r.Text())
		if l == "" {
			continue
		}
		row := make([]widgetRule, 0)
		ws := strings.Fields(l)
		weightTotal := 0
		for _, w := range ws {
			wr := widgetRule{Weight: 1}
			ks := strings.Split(w, "/")
			rs := strings.Split(ks[0], ":")
			var wid string
			if len(rs) > 1 {
				v, e := strconv.Atoi(rs[0])
				if e != nil {
					log.Printf("Layout error on line %d: format must be INT:STRING/INT. Error parsing %s as a int. Word was %s. Using a row height of 1.", lineNo, rs[0], w)
					v = 1
				}
				if v < 1 {
					v = 1
				}
				wr.Height = v
				wid = rs[1]
			} else {
				wr.Height = 1
				wid = rs[0]
			}
			wr.Widget = strings.ToLower(wid)
			if len(ks) > 1 {
				weight, e := strconv.Atoi(ks[1])
				if e != nil {
					log.Printf("Layout error on line %d: format must be STRING/INT. Error parsing %s as a int. Word was %s. Using a weight of 1 for widget.", lineNo, ks[1], w)
					weight = 1
				}
				if weight < 1 {
					weight = 1
				}
				wr.Weight = float64(weight)
				if len(ks) > 2 {
					log.Printf("Layout warning on line %d: too many '/' in word %s; ignoring extra junk.", lineNo, w)
				}
				weightTotal += weight
			} else {
				weightTotal += 1
			}
			row = append(row, wr)
		}
		// Prevent tricksy users from breaking their own computers
		if weightTotal <= 1 {
			weightTotal = 1
		}
		for i, w := range row {
			row[i].Weight = w.Weight / float64(weightTotal)
		}
		rv.Rows = append(rv.Rows, row)
	}
	return rv
}
