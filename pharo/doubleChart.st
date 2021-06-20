|path string splited map b|

path := 'C:\Users\Home\Documents\GitHub\COVID19-Simulation\covid19\files\stats.txt'.

string := FileStream readOnlyFileNamed: path
				do: [ :stream | stream contents ].
splited := string splitOn: '|'.
map := splited collect: 
	[ :line | (line splitOn: Character space) 
	collect: [ :number | number asNumber ] 
].

b := RTDoubleBarBuilder new.
b pointName: [ :row | row at: 1 ].
b points: map. 
b bottomValue: [ :row | row at: 2 ] 
		titled: 'Infected'.
b topValue: [ :row | row at: 3] 
		titled: 'Dead'.
	
b