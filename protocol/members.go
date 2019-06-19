package protocol

import (
	"github.com/AlexandreBelling/go-boojum/identity"
)

// IDList is a wrapper for a list of identity.ID
type IDList []identity.ID

// Min returns the position of the identity.ID with the smallest bigint representation
func (i IDList) Min() (int, identity.ID) {

	bestPosition, bestValue := 0, i[0]
	bestValueBig := bestValue.Big()

	for currentPosition, currentValue := range i[1:] {
		currentValueBig := currentValue.Big()
		if bestValueBig.Cmp(currentValueBig) == 1 {

			bestPosition = currentPosition
			bestValue = currentValue
			bestValueBig = currentValueBig
		}
	}

	return bestPosition, bestValue
}

// Max returns the position of the identity.ID with the biggest bigint representation
func (i IDList) Max() (int, identity.ID) {

	bestPosition, bestValue := 0, i[0]
	bestValueBig := bestValue.Big()

	for currentPosition, currentValue := range i[1:] {
		currentValueBig := currentValue.Big()
		if bestValueBig.Cmp(currentValueBig) == -1 {

			bestPosition = currentPosition
			bestValue = currentValue
			bestValueBig = currentValueBig
		}
	}

	return bestPosition, bestValue
}

// SmallestHigherThan return the smallest identity.ID that is higher
// than a given identity.ID. Returns the minimum if no records are higher
// than the given identity.ID
func (i IDList) SmallestHigherThan(id identity.ID) (int, identity.ID) {

	bestPosition, bestValue := i.FirstHigherThan(id)

	if bestPosition == -1 {
		return i.Min()
	}

	referenceIDBig := id.Big()
	bestValueBig := bestValue.Big()
	for currentPosition, currentValue := range i[bestPosition+1:] {

		println("Hello there")

		currentValueBig := currentValue.Big()
		if referenceIDBig.Cmp(currentValueBig) == 1 && bestValueBig.Cmp(currentValueBig) == -1 {

			bestPosition = currentPosition
			bestValue = currentValue
			bestValueBig = currentValueBig
		}
	}

	return bestPosition, bestValue
}

// FirstHigherThan returns the first record in IDList higher than a given reference identity.ID.
func (i IDList) FirstHigherThan(id identity.ID) (int, identity.ID) {

	
	referenceIDBig := id.Big()
	println("Hello there")
	for currentPosition, current := range i {

		println("Hello there")

		currentBig := current.Big()
		if referenceIDBig.Cmp(currentBig) == 1 {
			return currentPosition, current
		}
	}
	return -1, ""
}
