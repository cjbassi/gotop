T is max height in row
S(T) is all widgets with height T
R(T) is all widgets with height < T
X is len(R) > 0 ? 1 : 0
C is len(S) + X
Make row
Make C columns
Place S
Recurse with R; place result


 1         2          3       4          5
cpu/2...............  mem/1.  6:procs/2..........
3:temp/1.  2:disk/2.........  |..................
|........  |................  |..................
|........  power/2..........  |..................
net/2...............  batt..  |..................

 1         2          3       4          5
cpu/2...............  6:procs/2........  mem/1...
2:disk/2............  |................  3:temp/1   
|...................  |................  |.......
power/2.............  |................  |.......
net/2...............  |................  batt

 1         2          3       4          5
1x2.................  3x2..............  1x1.....    221    221
2x2.................  |||||||||||||||||  3x1.....    21     2x1
||||||||||||||||||||  |||||||||||||||||  ||||||||           
1x1......  1x1......  1x2..............  1x1.....    1121
1x2.................  1x2..............  ||||||||    22     22x
1x1......  1x4...................................    14     

initial columns = initial row
fill
	pattern for row
	does pattern fit columns?
		yes: place widgets
		no: new row w/ new columns; fill

does fit
	cw < patt_c_w
