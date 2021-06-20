|path string splited map b ds|

path := 'C:\Users\Home\Documents\GitHub\COVID19-Simulation\covid19\files\stats.txt'.

string := FileStream readOnlyFileNamed: path
				do: [ :stream | stream contents ].
splited := string splitOn: '|'.
map := splited collect: 
	[ :line | (line splitOn: Character space) 
	collect: [ :number | number ] 
].

map do: [ :each | ((each) at: 4) = '1' 
	ifTrue: [each add: 'infected ', ((each) at: 2), ', dead: ', ((each) at: 3), ', tested positive: ', ((each) at: 6), '/', ((each) at: 5), ', active distancing: yes']
	ifFalse: [each add: 'infected: ', ((each) at: 2), ', dead: ', ((each) at: 3), ', tested positive: ', ((each) at: 6), '/', ((each) at: 5), ', active distancing: no']].

b := RTGrapher new.
ds := RTData new.
ds interaction fixedPopupText: [ :row | row value at: 7 ]. 
ds barShape color: Color red.
ds points: map.
ds x: [ :row| (row at: 1) asNumber].
ds y: [ :row| (row at: 2) asNumber].
b add: ds.
	
ds := RTData new.
ds interaction fixedPopupText: [ :row | row value at: 7 ]. 
ds barShape color: Color lightGray.
ds points: map.
ds x: [ :row| (row at: 1) asNumber].
ds y: [ :row| (row at: 3) asNumber].
b add: ds.
b axisX noTick; title: 'days'. 
b axisY noDecimal; title: 'number of infected and dead'.
b