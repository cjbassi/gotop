// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

var StandardColors = []Attribute{
	ColorRed,
	ColorGreen,
	ColorYellow,
	ColorBlue,
	ColorMagenta,
	ColorCyan,
	ColorWhite,
}

type RootTheme struct {
	Default AttrPair

	Block BlockTheme

	BarChart        BarChartTheme
	Gauge           GaugeTheme
	LineChart       LineChartTheme
	List            ListTheme
	Paragraph       ParagraphTheme
	PieChart        PieChartTheme
	Sparkline       SparklineTheme
	StackedBarChart StackedBarChartTheme
	Tab             TabTheme
	Table           TableTheme
}

type BlockTheme struct {
	Title  AttrPair
	Border AttrPair
}

type BarChartTheme struct {
	Bars   []Attribute
	Nums   []Attribute
	Labels []Attribute
}

type GaugeTheme struct {
	Percent Attribute
	Bar     Attribute
}

type LineChartTheme struct {
	Lines []Attribute
	Axes  Attribute
}

type ListTheme struct {
	Text AttrPair
}

type ParagraphTheme struct {
	Text AttrPair
}

type PieChartTheme struct {
	Slices []Attribute
}

type SparklineTheme struct {
	Title AttrPair
	Line  Attribute
}

type StackedBarChartTheme struct {
	Bars   []Attribute
	Nums   []Attribute
	Labels []Attribute
}

type TabTheme struct {
	Active   AttrPair
	Inactive AttrPair
}

type TableTheme struct {
	Text AttrPair
}

var Theme = RootTheme{
	Default: AttrPair{7, -1},

	Block: BlockTheme{
		Title:  AttrPair{7, -1},
		Border: AttrPair{6, -1},
	},

	BarChart: BarChartTheme{
		Bars:   StandardColors,
		Nums:   StandardColors,
		Labels: StandardColors,
	},

	Paragraph: ParagraphTheme{
		Text: AttrPair{ColorWhite, -1},
	},

	PieChart: PieChartTheme{
		Slices: StandardColors,
	},

	List: ListTheme{
		Text: AttrPair{1, -1},
	},

	StackedBarChart: StackedBarChartTheme{
		Bars:   StandardColors,
		Nums:   StandardColors,
		Labels: StandardColors,
	},

	Gauge: GaugeTheme{
		Percent: ColorWhite,
		Bar:     ColorWhite,
	},

	Sparkline: SparklineTheme{
		Line: ColorBlack,
		Title: AttrPair{
			Fg: ColorBlue,
			Bg: ColorDefault,
		},
	},

	LineChart: LineChartTheme{
		Lines: StandardColors,
		Axes:  ColorBlue,
	},

	Table: TableTheme{
		Text: AttrPair{4, -1},
	},

	Tab: TabTheme{
		Active:   AttrPair{ColorRed, ColorDefault},
		Inactive: AttrPair{ColorWhite, ColorDefault},
	},
}
