package square

import (
	"fmt"
	"strconv"
)

type MysticSquare interface {
	State() string
	MoveUp() map[int]int
	MoveDown() map[int]int
	MoveLeft() map[int]int
	MoveRight() map[int]int
	ValidateState() bool
	FindEmptySpace() int
	MapKeyToNewKey() map[int]map[string]int
	RealState() map[int]int
}

type MysticSquare3 struct {
	state    map[int]int
	strState string
}

// convert the state map to a state string
func buildMysticSquare3StateString(state map[int]int) (strState string) {
	buildRow := func(squareState map[int]int, row int) (rowString string) {
		c1 := strconv.Itoa(squareState[1+3*(row-1)])
		c2 := strconv.Itoa(squareState[2+3*(row-1)])
		c3 := strconv.Itoa(squareState[3+3*(row-1)])

		if c1 == "9" {
			c1 = " "
		}

		if c2 == "9" {
			c2 = " "
		}

		if c3 == "9" {
			c3 = " "
		}

		rowString = fmt.Sprintf("%s %s %s", c1, c2, c3)
		return
	}

	if validState := (MysticSquare3{state: state}).ValidateState(); validState {
		squareState := state
		row1 := buildRow(squareState, 1)
		row2 := buildRow(squareState, 2)
		row3 := buildRow(squareState, 3)
		strState = fmt.Sprintf("%s\n%s\n%s", row1, row2, row3)
	} else {
		strState = "Invalid State"
	}
	return
}

// create a new mystic square
func NewMysticSquare(state map[int]int) (newSquare MysticSquare, err error) {
	newSquare, err = NewMysticSquare3(state)
	return
}

// create a new 3x3 mystic square
func NewMysticSquare3(state map[int]int) (newSquare *MysticSquare3, err error) {
	newSquare = &MysticSquare3{state: state, strState: buildMysticSquare3StateString(state)}
	err = nil
	if validState := newSquare.ValidateState(); !validState {
		newSquare = nil
		err = fmt.Errorf("invalid state %v", state)
	}
	return
}

// move the empty space up
func (square MysticSquare3) MoveUp() (newSquare map[int]int) {
	if validState := square.ValidateState(); validState {
		blankSpace := square.FindEmptySpace()
		boardMapping := square.MapKeyToNewKey()
		if spaceMapping, spaceMappingExists := boardMapping[blankSpace]; spaceMappingExists {
			if newSpace, upExists := spaceMapping["up"]; upExists {
				newSquare = make(map[int]int)
				oldSquare := square.state
				for key, val := range oldSquare {
					newSquare[key] = val
				}

				temp := newSquare[blankSpace]
				newSquare[blankSpace] = newSquare[newSpace]
				newSquare[newSpace] = temp
			}
		}
	}
	return
}

// move the empty space down
func (square MysticSquare3) MoveDown() (newSquare map[int]int) {
	if validState := square.ValidateState(); validState {
		blankSpace := square.FindEmptySpace()
		boardMapping := square.MapKeyToNewKey()
		if spaceMapping, spaceMappingExists := boardMapping[blankSpace]; spaceMappingExists {
			if newSpace, downExists := spaceMapping["down"]; downExists {
				newSquare = make(map[int]int)
				oldSquare := square.state
				for key, val := range oldSquare {
					newSquare[key] = val
				}

				temp := newSquare[blankSpace]
				newSquare[blankSpace] = newSquare[newSpace]
				newSquare[newSpace] = temp
			}
		}
	}
	return
}

// move the empty space left
func (square MysticSquare3) MoveLeft() (newSquare map[int]int) {
	if validState := square.ValidateState(); validState {
		blankSpace := square.FindEmptySpace()
		boardMapping := square.MapKeyToNewKey()
		if spaceMapping, spaceMappingExists := boardMapping[blankSpace]; spaceMappingExists {
			if newSpace, leftExists := spaceMapping["left"]; leftExists {
				newSquare = make(map[int]int)
				oldSquare := square.state
				for key, val := range oldSquare {
					newSquare[key] = val
				}

				temp := newSquare[blankSpace]
				newSquare[blankSpace] = newSquare[newSpace]
				newSquare[newSpace] = temp
			}
		}
	}
	return
}

// move the empty space right
func (square MysticSquare3) MoveRight() (newSquare map[int]int) {
	if validState := square.ValidateState(); validState {
		blankSpace := square.FindEmptySpace()
		boardMapping := square.MapKeyToNewKey()
		if spaceMapping, spaceMappingExists := boardMapping[blankSpace]; spaceMappingExists {
			if newSpace, rightExists := spaceMapping["right"]; rightExists {
				newSquare = make(map[int]int)
				oldSquare := square.state
				for key, val := range oldSquare {
					newSquare[key] = val
				}

				temp := newSquare[blankSpace]
				newSquare[blankSpace] = newSquare[newSpace]
				newSquare[newSpace] = temp
			}
		}
	}
	return
}

// map for each square to where it would be if moved up, down, left or right
func (square MysticSquare3) MapKeyToNewKey() (mapping map[int]map[string]int) {
	mapping = map[int]map[string]int{
		1: {"down": 4, "right": 2},
		2: {"down": 5, "left": 1, "right": 3},
		3: {"down": 6, "left": 2},
		4: {"up": 1, "down": 7, "right": 5},
		5: {"up": 2, "down": 8, "left": 4, "right": 6},
		6: {"up": 3, "down": 9, "left": 5},
		7: {"up": 4, "right": 8},
		8: {"up": 5, "left": 7, "right": 9},
		9: {"up": 6, "left": 8},
	}

	return
}

// return the state string
func (square MysticSquare3) State() (state string) {

	state = square.strState
	return
}

// find the empty space
func (square MysticSquare3) FindEmptySpace() (emptySpace int) {
	emptySpace = -1
	state := square.state
	for key, val := range state {
		if val == 9 {
			emptySpace = key
			break
		}
	}
	return
}

// ensure the mystic square is valid
func (square MysticSquare3) ValidateState() (validState bool) {
	validTilePositions := make([]int, 9)
	validTileValues := make([]int, 9)
	board := square.state
	if len(board) != len(validTilePositions) {
		validState = false
		return
	}
	for i := 1; i < 10; i++ {
		validTilePositions[i-1] = i
		validTileValues[i-1] = i
	}
	find := func(A []int, value int) (index int) {
		for idx, val := range A {
			if value == val {
				index = idx
				return
			}
		}
		return -1
	}
	remove := func(A []int, value int) (arr []int) {
		if idx := find(A, value); idx >= 0 && idx < len(A) {
			arr = append(A[:idx], A[idx+1:]...)
		} else {
			arr = A
		}

		return
	}
	for key, value := range board {
		validTilePositions = remove(validTilePositions, key)
		validTileValues = remove(validTileValues, value)
	}

	allKeysFound := len(validTilePositions) == 0
	allValuesFound := len(validTileValues) == 0
	validState = allKeysFound && allValuesFound
	return
}

// copy the state from the square to a new map.
func (square MysticSquare3) RealState() (copy map[int]int) {
	copy = make(map[int]int)
	for key, value := range square.state {
		copy[key] = value

	}
	return
}
