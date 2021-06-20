|day path stringsplited map labeledData pie lb|
day := 1.

path := 'C:\Users\Home\Documents\GitHub\COVID19-Simulation\covid19\files\stats.txt'.

string := FileStream readOnlyFileNamed: path
				do: [ :stream | stream contents ].
splited := string splitOn: '|'.
map := splited collect: 
	[ :line | (line splitOn: Character space) 
	collect: [ :number | number asNumber ] 
].

labeledData := {'Tested positive' -> ((map at: day) at: 6). 'Tested negative' -> (((map  at: day) at: 5) - ((map  at: day) at: 6))}.

pie := RTPieBuilder new.
pie interaction popup.
pie objects: labeledData.
pie slice: #value.
pie labeled.
pie normalizer distinctColor.

"And add the legend"
lb := RTLegendBuilder new.
lb view: pie view.
lb addText: 'Ratio of the number of tested positive and  tested negative for day ', day asString .
lb build.

pie