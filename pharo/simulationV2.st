|iter mode path fileLines dim cells infected dead notInfected nothing cellSize cellValue view deadCells infectedCells noCell notInfectedCells|

iter:=1.
mode:='serial'.
path := 'C:\Users\Home\Documents\GitHub\COVID19-Simulation\covid19\files\', mode, (iter asString ), '.txt'.

fileLines := path asFileReference contents lines.
dim := fileLines size.
cells := fileLines collect: 
	[ :line | (line splitOn: Character space) 
	collect: [ :number | number asNumber ] 
].

infected := Set new.
dead := Set new.
notInfected := Set new.
nothing := Set new.

cellSize := 1.

0 to: (dim-1) do: [ :i |
		0 to: (dim-1) do: [ :j |
				cellValue := (cells at: i+1) at: j+1.
				cellValue = 0
				ifTrue: [ nothing add: j@i]
				ifFalse: [ 
					cellValue = -1
					ifTrue: [ dead add: j@i]
					ifFalse: [ 
						cellValue = 1
						ifTrue: [ notInfected add: j@i ]
						ifFalse: [ infected add: j@i ]
					]
				]
			 ]].

noCell := (RTEllipse new size: cellSize; color: Color gray)
	elementsOn: (1 to: nothing size).
	
infectedCells := (RTEllipse new size: cellSize; color: Color red)
	elementsOn: (1 to: infected size).
	
notInfectedCells := (RTEllipse new size: cellSize; color: Color green)
	elementsOn: (1 to: notInfected size).
	
deadCells := (RTEllipse new size: cellSize; color: Color black)
	elementsOn: (1 to: dead size).
	

view := RTView new.


1 to: nothing size do: [ :index | (noCell at: index) 
						translateTo: 
						(nothing asArray at: index)
					].
				
1 to: infected size do: [ :index | (infectedCells at: index) 
						translateTo: 
						(infected asArray at: index)
					].
				
1 to: notInfected size do: [ :index | (notInfectedCells at: index) 
						translateTo: 
						(notInfected asArray at: index)
					].
				
1 to: dead size do: [ :index | (deadCells at: index) 
						translateTo: 
						(dead asArray at: index)
					].
				
view add: noCell;
add: infectedCells;
add: notInfectedCells;
add: deadCells;

view
