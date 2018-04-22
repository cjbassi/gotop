# termui

A fork of [termui](https://github.com/gizak/termui) with a lot of code cleanup and (frequently asked for) improvements.

You can see an implementation/example usage of this library [here](https://github.com/cjbassi/gotop).

Some usage improvements include:
* better event/key-combo names
* more convenient event handling function
* 256 colors
* better grid system
* linegraph uses [drawille-go](https://github.com/exrook/drawille-go)
    * no longer have to choose between dot mode and braille mode; uses a superior braille mode
* table supports mouse and keyboard navigation
* table is scrollable
* more powerful table column width sizing
* visual improvements to linegraph and table

TODO:
* readd widgets that were removed like the list and bargraph
* focusable widgets
