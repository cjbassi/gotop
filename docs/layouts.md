# Layouts

gotop can parse and render layouts from a specification file.  The format is
intentionally simple.  The amount of nesting levels is limited.  Some examples
are in the `layouts` directory; you can try each of these with, e.g.,
`gotop --layout-file layouts/procs`.  If you stick your layouts in
`$XDG_CONFIG_HOME/gotop`, you can reference them on the command line with the
`-l` argument, e.g. `gotop -l procs`.

The syntax for each widget in a row is:
```
(rowspan:)?widget(/weight)?
```
and these are separated by spaces.

1. Each line is a row
2. Empty lines are skipped
3. Spaces are compressed (so you can do limited visual formatting)
4. Legal widget names are: cpu, disk, mem, temp, batt, net, procs
5. Widget names are not case sensitive
4. The simplest row is a single widget, by name, e.g. `cpu`
5. **Weights**
    1. Widgets with no weights have a weight of 1.
    2. If multiple widgets are put on a row with no weights, they will all have
       the same width.
    3. Weights are integers
    4. A widget will have a width proportional to its weight divided by the
       total weight count of the row. E.g.,

       ```
       cpu      net
       disk/2   mem/4
       ```

       The first row will have two widgets: the CPU and network widgets; each
       will be 50% of the total width wide.  The second row will have two
       widgets: disk and memory; the first will be 2/6 ~= 33% wide, and the
       second will be 5/7 ~= 67% wide (or, memory will be twice as wide as disk).
9.  If prefixed by a number and colon, the widget will span that number of
    rows downward. E.g.

    ```
    mem   2:cpu
    net
    ```

    Here, memory and network will be in the same row as CPU, one over the other,
    and each half as high as CPU; it'll look like this:

    ```
     +------+------+
     | Mem  |      |
     +------+ CPU  |
     | Net  |      |
     +------+------+
    ```
     
10. Negative, 0, or non-integer weights will be recorded as "1".  Same for row
    spans. 
11. Unrecognized widget names will cause the application to abort.                          
12. In rows with multi-row spanning widgets **and** weights, weights in
    lower rows are ignored.  Put the weight on the widgets in that row, not
    in later (spanned) rows.
13. Widgets are filled in top down, left-to-right order.
14. The larges row span in a row defines the top-level row span; all smaller
    row spans constitude sub-rows in the row. For example, `cpu mem/3 net/5`
    means that net/5 will be 5 rows tall overall, and mem will compose 3 of
    them. If following rows do not have enough widgets to fill the gaps,
    spacers will be used.

Yes, you're clever enough to break the layout algorithm, but if you try to
build massive edifices, you're in for disappointment.

To list all built-in color schemes, call:

```
gotop --list layouts
```

