|path string splited map infected dead graphic ds1 ds2|

path := 'C:\Users\Home\Documents\GitHub\COVID19-Simulation\covid19\files\stats.txt'.
string := FileStream readOnlyFileNamed: path
				do: [ :stream | stream contents ].
splited := string splitOn: '|'.
map := splited collect: 
	[ :line | (line splitOn: Character space) 
	collect: [ :number | number asNumber] 
].

infected := Set  new.
dead := Set  new.
map do: [ :each | infected add: ((each) at: 1) @ ((each) at: 2)].
map do: [ :each | dead add: ((each) at: 1) @ ((each) at: 3)].

graphic := RTGrapher new.

ds1 := RTData new.
ds1 label: 'infected'.
ds1 dotShape circle.
ds1 points: infected.
ds1 connectColor: Color blue.
ds1 y: #y.
ds1 x: #x.
graphic add: ds1.

ds2 := RTData new.
ds2 label: 'dead'.
ds2 dotShape circle.
ds2 points: dead.
ds2 connectColor: Color red.
ds2 y: #y.
ds2 x: #x.
graphic add: ds2.
	
graphic axisX title: 'day'; noTick; noDecimal; numberOfLabels: map size.
graphic axisY title: 'infected and dead'; noDecimal.
graphic legend below.
graphic