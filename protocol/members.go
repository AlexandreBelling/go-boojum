package protocol

// IDList is a wrapper for a list of ID
type IDList []ID

// Min returns the position of the ID with the smallest bigint representation
func (i IDList) Min() (int, ID) {

	bestPosition, bestValue := 0, i[0]
	bestValueBig := bestValue.Big()

	for currentPosition, currentValue := range i[1:] {
		currentValueBig := currentValue.Big()
		if bestValueBig.Cmp(currentValueBig) == 1 {

			bestPosition = currentPosition
			bestValue 	 = currentValue
			bestValueBig = currentValueBig
		}
	}

	return bestPosition, bestValue
}

// Max returns the position of the ID with the biggest bigint representation
func (i IDList) Max() (int, ID) {

	bestPosition, bestValue := 0, i[0]
	bestValueBig := bestValue.Big()

	for currentPosition, currentValue := range i[1:] {
		currentValueBig := currentValue.Big()
		if bestValueBig.Cmp(currentValueBig) == -1 {

			bestPosition = currentPosition
			bestValue 	 = currentValue
			bestValueBig = currentValueBig
		}
	}

	return bestPosition, bestValue
}

// SmallestHigherThan return the smallest ID that is higher 
// than a given ID. Returns the minimum if no records are higher 
// than the given ID
func (i IDList) SmallestHigherThan(id ID) (int, ID) {

	bestPosition, bestValue := i.FirstHigherThan(id)

	if bestPosition == -1 { return i.Min() }

	referenceIDBig 	:= id.Big()
	bestValueBig	:= bestValue.Big()
	for currentPosition, currentValue := range i[bestPosition+1:] {

		currentValueBig := currentValue.Big()
		if referenceIDBig.Cmp(currentValueBig) == 1 && bestValueBig.Cmp(currentValueBig) == -1 {

			bestPosition = currentPosition
			bestValue 	 = currentValue
			bestValueBig = currentValueBig 
		}
	}

	return bestPosition, bestValue
}

// FirstHigherThan returns the first record in IDList higher than a given reference ID. 
func (i IDList) FirstHigherThan(id ID) (int, ID) {

	referenceIDBig := id.Big()
	for currentPosition, current := range i {
		
		currentBig := current.Big()
		if referenceIDBig.Cmp(currentBig) == 1 {
			return currentPosition, current 
		}
	}
	return -1, ID{}
}
	

