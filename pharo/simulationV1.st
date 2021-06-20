|iter mode path fileLines dim cells infected dead notInfected nothing cellSize cellValue grapher noCellData infectedData notInfectedData deadData|

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

cellSize := 15.

0 to: (dim-1) do: [ :i |
		0 to: (dim-1) do: [ :j |
				cellValue := (cells at: i+1) at: j+1.
				cellValue = 0
				ifTrue: [ nothing add: j@(dim - i -1)]
				ifFalse: [ 
					cellValue = -1
					ifTrue: [ dead add: j@(dim - i -1)]
					ifFalse: [ 
						cellValue = 1
						ifTrue: [ notInfected add: j@(dim - i -1) ]
						ifFalse: [ infected add: j@(dim - i -1) ]
					]
				]
			 ]].

noCellData := RTData new.
noCellData dotShape color: Color gray .
noCellData dotShape size: cellSize.
noCellData points: nothing; x: #x; y: #y.

infectedData := RTData new.
infectedData dotShape color: Color red.
infectedData dotShape size: cellSize.
infectedData points: infected; x: #x; y: #y.


notInfectedData := RTData new .
notInfectedData dotShape color: Color green.
notInfectedData dotShape size: cellSize.
notInfectedData points: notInfected; x: #x; y: #y.

deadData := RTData new .
deadData dotShape color: Color black.
deadData dotShape size: cellSize.
deadData points: dead; x: #x; y: #y.

grapher := RTGrapher new
    add: infectedData;
    add: noCellData;
    add: notInfectedData;
    add: deadData;
    build;
    view.
.

grapher view.
