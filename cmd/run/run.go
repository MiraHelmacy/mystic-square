/*
Copyright © 2024 Alex Helmacy

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package run

import (
	"container/heap"
	"errors"
	"fmt"
	"math"
	"slices"

	"mysticsquare/datastructures"
	"mysticsquare/square"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type AlgorithmSelection int
type SquareDifficulty int

// difficulty constants
const (
	EASY_DIFFICULTY SquareDifficulty = 1
	HARD_DIFFICULTY SquareDifficulty = 2
	NO_PATH         SquareDifficulty = 3
)

// algorithm constants
const (
	A_STAR_SEARCH        AlgorithmSelection = 1
	DIJKSTRAS_ALGORITHM  AlgorithmSelection = 2
	BREADTH_FIRST_SEARCH AlgorithmSelection = 3
)

// options constants
const (
	ALGORITHM_LONG_OPTION   = "algorithm"
	ALGORITHM_SHORT_OPTION  = "a"
	DIFFICULTY_LONG_OPTION  = "difficulty"
	DIFFICULTY_SHORT_OPTION = "d"
)

// cli args
type CliArgs struct {
	difficulty SquareDifficulty
	algorithm  AlgorithmSelection
}

// create a new set of Cli Args
func NewRunCliArgs() (args *CliArgs, valid bool) {
	args = &CliArgs{}
	valid = true
	difficulty := SquareDifficulty(viper.GetInt(DIFFICULTY_LONG_OPTION))
	switch difficulty {
	case EASY_DIFFICULTY, HARD_DIFFICULTY, NO_PATH:
		args.difficulty = difficulty
	default:
		args = nil
		valid = false
		return
	}

	algorithm := AlgorithmSelection(viper.GetInt(ALGORITHM_LONG_OPTION))

	switch algorithm {
	case A_STAR_SEARCH, DIJKSTRAS_ALGORITHM, BREADTH_FIRST_SEARCH:
		args.algorithm = algorithm
	default:
		args = nil
		valid = false
		return
	}

	return
}

// easy mode
func easy() (initialState, targetState map[int]int) {
	initialState = map[int]int{
		1: 9, 2: 1, 3: 2,
		4: 4, 5: 6, 6: 3,
		7: 7, 8: 5, 9: 8,
	}

	targetState = map[int]int{
		1: 1, 2: 2, 3: 3,
		4: 4, 5: 5, 6: 6,
		7: 7, 8: 8, 9: 9,
	}

	return
}

// hard mode
func hard() (initialState, targetState map[int]int) {
	initialState = map[int]int{
		1: 8, 2: 7, 3: 6,
		4: 5, 5: 9, 6: 4,
		7: 3, 8: 2, 9: 1,
	}

	targetState = map[int]int{
		1: 1, 2: 2, 3: 3,
		4: 4, 5: 9, 6: 5,
		7: 6, 8: 7, 9: 8,
	}

	return
}

// no path mode
func noPath() (initialState, targetState map[int]int) {
	initialState = map[int]int{
		1: 2, 2: 1, 3: 3,
		4: 4, 5: 5, 6: 6,
		7: 7, 8: 8, 9: 9,
	}

	targetState = map[int]int{
		1: 1, 2: 2, 3: 3,
		4: 4, 5: 5, 6: 6,
		7: 7, 8: 8, 9: 9,
	}

	return
}

// description of algorithm parameter
func algorithmDescription() (description string) {
	description = fmt.Sprintf("Algorithm to use. A star: %v, Dijkstras: %v, BFS: %v", int(A_STAR_SEARCH), int(DIJKSTRAS_ALGORITHM), int(BREADTH_FIRST_SEARCH))
	return
}

// description of difficulty parameter
func difficultyDescription() (description string) {
	description = fmt.Sprintf("Difficulty of the puzzle. Easy: %v, Hard: %v, No Path: %v", int(EASY_DIFFICULTY), int(HARD_DIFFICULTY), int(NO_PATH))
	return
}

// from CliArgs create the initial state
func (args CliArgs) squaresForDifficulty() (initialState, targetState map[int]int) {
	switch args.difficulty {
	case EASY_DIFFICULTY:
		initialState, targetState = easy()
	case HARD_DIFFICULTY:
		initialState, targetState = hard()
	case NO_PATH:
		initialState, targetState = noPath()
	default:
		initialState, targetState = make(map[int]int), make(map[int]int)
	}

	return
}

// from the CliArgs return the algorithm to use
func (args CliArgs) realAlgorithm() (algorithm func(is, ts square.MysticSquare) (paths map[string]square.MysticSquare, pathFound bool)) {
	switch args.algorithm {
	case A_STAR_SEARCH:
		algorithm = func(is, ts square.MysticSquare) (paths map[string]square.MysticSquare, pathFound bool) {
			paths, pathFound = aStar(is, ts, func(current square.MysticSquare) int { return manhattanDistance(current, ts) })
			return
		}
	case DIJKSTRAS_ALGORITHM:
		algorithm = dijkstrasAlgorithm
	case BREADTH_FIRST_SEARCH:
		algorithm = bfs
	default:
		algorithm = func(is, ts square.MysticSquare) (paths map[string]square.MysticSquare, pathFound bool) {
			paths = make(map[string]square.MysticSquare)
			pathFound = false
			return
		}
	}
	return
}

// implementation of manhattan distance
func manhattanDistance(current, target square.MysticSquare) (distance int) {

	if currentValid := current.ValidateState(); !currentValid {
		panic(fmt.Sprintf("passed an invalid state to manhattan distance function: %v", current.RealState()))
	}

	if targetValid := target.ValidateState(); !targetValid {
		panic(fmt.Sprintf("passed an invalid state to manhattan distance function: %v", target.RealState()))
	}

	if bothSquaresSameSize := len(current.RealState()) == len(target.RealState()); !bothSquaresSameSize {
		panic(fmt.Sprintf("State size missmatch: %v, %v", current.RealState(), target.RealState()))
	}

	distance = 0

	currentState := make([][]int, 0)
	targetState := make([][]int, 0)
	for i := 0; i < 3; i++ {
		c1, c2, c3 := current.RealState()[1+3*i], current.RealState()[2+3*i], current.RealState()[3+3*i]
		currentState = append(currentState, []int{c1, c2, c3})
		t1, t2, t3 := target.RealState()[1+3*i], target.RealState()[2+3*i], target.RealState()[3+3*i]
		targetState = append(targetState, []int{t1, t2, t3})
	}

	for currentRow := range currentState {
		for currentColumn, currentValue := range currentState[currentRow] {
			if targetValue := targetState[currentRow][currentColumn]; targetValue != currentValue && currentValue != 9 {
				valueFound := false
				for targetRow := range targetState {
					if valueFound {
						break
					}
					for targetColumn, searchTargetValue := range targetState[targetRow] {
						if searchTargetValue == currentValue {
							distance += int(math.Abs(float64(currentRow)-float64(targetRow)) + math.Abs(float64(currentColumn)-float64(targetColumn)))
							valueFound = true
							break
						}
					}

				}
			}
		}
	}
	return
}

// a* search implementation
func aStar(initialState, targetState square.MysticSquare, h func(square.MysticSquare) int) (paths map[string]square.MysticSquare, pathFound bool) {
	if h == nil {
		panic("Invalid heuristic function")
	}

	q := datastructures.NewMysticSquarePriorityQueue()

	paths = make(map[string]square.MysticSquare)
	paths[initialState.State()] = nil

	distance := make(map[string]int)
	distance[initialState.State()] = 0

	g := func(state square.MysticSquare) (g int) {
		g = distance[state.State()]
		return
	}

	f := func(state square.MysticSquare) (f int) {
		detectOverflow := func(a, b int) (err error) {
			ErrOverflow := errors.New("integer overflow detected")
			if b > 0 {
				if a > math.MaxInt-b {
					err = ErrOverflow
					return
				}
			}
			err = nil
			return
		}
		gCurrent := g(state)
		hCurrent := h(state)
		if err := detectOverflow(gCurrent, hCurrent); err == nil {
			f = gCurrent + hCurrent
		} else {
			f = math.MaxInt
		}

		return
	}

	initialItem := datastructures.NewMysticSquareItem(initialState, f(initialState))
	itemsMap := make(map[string]*datastructures.MysticSquareItem)
	itemsMap[initialItem.Msquare.State()] = initialItem
	visited := make(map[string]bool)
	q.Push(initialItem)

	heap.Init(q)

	pathFound = false
	targetStateString := targetState.State()

	for currentItem, itemExists := q.Process(); itemExists; currentItem, itemExists = q.Process() {
		current := currentItem.Msquare
		currentStateString := current.State()
		if currentStateString == targetStateString {
			pathFound = true
			break
		}
		adjacent := make([]square.MysticSquare, 0)

		if leftState := current.MoveLeft(); leftState != nil {
			if leftSquare, err := square.NewMysticSquare(leftState); err == nil {
				adjacent = append(adjacent, leftSquare)
			}
		}

		if rightState := current.MoveRight(); rightState != nil {
			if rightSquare, err := square.NewMysticSquare(rightState); err == nil {
				adjacent = append(adjacent, rightSquare)
			}
		}

		if upState := current.MoveUp(); upState != nil {
			if upSquare, err := square.NewMysticSquare(upState); err == nil {
				adjacent = append(adjacent, upSquare)
			}
		}

		if downState := current.MoveDown(); downState != nil {
			if downSquare, err := square.NewMysticSquare(downState); err == nil {
				adjacent = append(adjacent, downSquare)
			}
		}

		for _, neighbor := range adjacent {
			neighborStateString := neighbor.State()

			if _, distanceForNeighborExists := distance[neighborStateString]; !distanceForNeighborExists {
				distance[neighborStateString] = math.MaxInt
				newItem := datastructures.NewMysticSquareItem(neighbor, f(neighbor))
				itemsMap[neighborStateString] = newItem
				heap.Push(q, newItem)
			}

			tentativeDistance := distance[currentStateString] + 1
			neighborDistance := distance[neighborStateString]
			if _, neighborVisited := visited[neighborStateString]; tentativeDistance < neighborDistance && !neighborVisited {
				paths[neighborStateString] = current
				distance[neighborStateString] = tentativeDistance
				item := itemsMap[neighborStateString]
				q.Update(item, f(neighbor))
			}
		}
		visited[currentStateString] = true
	}
	return
}

// dijkstras algorithm implementation
func dijkstrasAlgorithm(initialState, targetState square.MysticSquare) (paths map[string]square.MysticSquare, pathFound bool) {
	q := datastructures.NewMysticSquarePriorityQueue()

	paths = make(map[string]square.MysticSquare)
	paths[initialState.State()] = nil

	distance := make(map[string]int)
	distance[initialState.State()] = 0

	initialItem := datastructures.NewMysticSquareItem(initialState, distance[initialState.State()])
	itemsMap := make(map[string]*datastructures.MysticSquareItem)
	itemsMap[initialItem.Msquare.State()] = initialItem
	q.Push(initialItem)

	heap.Init(q)

	pathFound = false
	targetStateString := targetState.State()

	for currentItem, itemExists := q.Process(); itemExists; currentItem, itemExists = q.Process() {

		current := currentItem.Msquare
		currentString := current.State()

		if currentString == targetStateString {
			pathFound = true
			break
		}

		adjacent := make([]square.MysticSquare, 0)
		if leftState := current.MoveLeft(); leftState != nil {
			if leftSquare, err := square.NewMysticSquare(leftState); err == nil {
				adjacent = append(adjacent, leftSquare)
			}
		}

		if rightState := current.MoveRight(); rightState != nil {
			if rightSquare, err := square.NewMysticSquare(rightState); err == nil {
				adjacent = append(adjacent, rightSquare)
			}
		}

		if upState := current.MoveUp(); upState != nil {
			if upSquare, err := square.NewMysticSquare(upState); err == nil {
				adjacent = append(adjacent, upSquare)
			}
		}

		if downState := current.MoveDown(); downState != nil {
			if downSquare, err := square.NewMysticSquare(downState); err == nil {
				adjacent = append(adjacent, downSquare)
			}
		}

		for _, value := range adjacent {
			valueString := value.State()
			if _, distanceExists := distance[valueString]; !distanceExists {
				distance[valueString] = math.MaxInt
				newItem := datastructures.NewMysticSquareItem(value, distance[valueString])
				itemsMap[valueString] = newItem
				heap.Push(q, newItem)
			}

			currentDistanceForValue := distance[valueString]
			altDistance := distance[currentString] + 1
			if altDistance < currentDistanceForValue {
				paths[valueString] = current
				distance[valueString] = altDistance
				item := itemsMap[valueString]
				q.Update(item, altDistance)
			}
		}
	}

	return
}

// bfs implementation
func bfs(initialState square.MysticSquare, targetState square.MysticSquare) (paths map[string]square.MysticSquare, pathFound bool) {
	q := datastructures.NewMysticSquareQueue()
	visited := make(map[string]bool)
	paths = make(map[string]square.MysticSquare)
	pathFound = false
	q.Push(initialState)
	paths[initialState.State()] = nil
	visited[initialState.State()] = true

	for current, hasItem := q.Process(); hasItem; current, hasItem = q.Process() {

		if pathFound = (current.State() == targetState.State()); pathFound {
			break
		}

		adjacent := make([]square.MysticSquare, 0)

		if stateLeft := current.MoveLeft(); stateLeft != nil {
			if squareStateLeft, err := square.NewMysticSquare(stateLeft); err == nil {
				adjacent = append(adjacent, squareStateLeft)
			}
		}

		if stateRight := current.MoveRight(); stateRight != nil {
			if squareStateRight, err := square.NewMysticSquare(stateRight); err == nil {
				adjacent = append(adjacent, squareStateRight)
			}
		}

		if stateUp := current.MoveUp(); stateUp != nil {
			if squareStateUp, err := square.NewMysticSquare(stateUp); err == nil {
				adjacent = append(adjacent, squareStateUp)
			}
		}

		if stateDown := current.MoveDown(); stateDown != nil {
			if squareStateDown, err := square.NewMysticSquare(stateDown); err == nil {
				adjacent = append(adjacent, squareStateDown)
			}
		}

		for _, newSquare := range adjacent {
			if _, newSquareVisited := visited[newSquare.State()]; !newSquareVisited {
				q.Push(newSquare)
				paths[newSquare.State()] = current
				visited[newSquare.State()] = true
			}
		}
	}
	if !pathFound {
		paths = nil
	}

	return
}

// work horse of the entire command
func executeRun(args *CliArgs) (err error) {
	if args == nil {
		err = fmt.Errorf("args not provided")
		return
	}

	initialState, targetState := args.squaresForDifficulty()
	initialMysticSquare, initialErr := square.NewMysticSquare(initialState)
	targetMysticSquare, targetErr := square.NewMysticSquare(targetState)

	if initialErr == nil && targetErr == nil {
		algorithm := args.realAlgorithm()
		if paths, pathFound := algorithm(initialMysticSquare, targetMysticSquare); pathFound {
			path := make([]square.MysticSquare, 0)
			for current := targetMysticSquare; paths[current.State()] != nil; current = paths[current.State()] {
				path = append(path, current)
			}
			path = append(path, initialMysticSquare)
			slices.Reverse(path)
			fmt.Println("START")
			for _, square := range path {
				fmt.Println(square.State())
				fmt.Println()
			}
		} else {
			fmt.Println("No Path")
		}
	} else {
		err = fmt.Errorf("initial or target states invalid: %v, %v", initialErr, targetErr)
		return
	}
	return
}

// RunCmd represents the run command
var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Solve a 3x3 mystic square",
	Long:  `Chose a difficulty and an algorithm to solve the problem`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if cliArgs, argsValid := NewRunCliArgs(); argsValid {
			executeRun(cliArgs)
		} else {
			err = fmt.Errorf("args not valid")
		}
		return
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	RunCmd.Flags().IntP(ALGORITHM_LONG_OPTION, ALGORITHM_SHORT_OPTION, 0, algorithmDescription())
	RunCmd.Flags().IntP(DIFFICULTY_LONG_OPTION, DIFFICULTY_SHORT_OPTION, 0, difficultyDescription())
}

func initConfig() {
	viper.BindPFlags(RunCmd.InheritedFlags())
	viper.BindPFlags(RunCmd.LocalFlags())
}
