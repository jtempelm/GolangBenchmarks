package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

const arraySize = 30_000 //Golang exports params with a leading capital letter so no ARRAY_SIZE convention
const testIterations = 10
const numberOfGoRoutines = 8

var largestNumberGenerated = 0
var twoDArrayRef *[arraySize][arraySize]int

func TestMain(m *testing.M) {
	fmt.Println("Main Test Method Executing")

	var twoDArray [arraySize][arraySize]int //int type should default to int32
	twoDArrayRef = &twoDArray               //Go passes arrays by value, slices by ref. Initializing 2d slices looped with make is tedious
	largestNumberGenerated = generate2DArray(twoDArrayRef)

	m.Run()
}

func TestFindTheLargestNumberInA2dArray_Serial(*testing.T) {
	for i := 0; i < testIterations; i++ {
		fmt.Printf("Golang Serial Implementation findTheLargestNumberInA2dArray - test run %d - ", i)
		findLargestNumberInArraySerial(twoDArrayRef, largestNumberGenerated)
	}
}

func TestFindTheLargestNumberInA2dArray_Parallel(*testing.T) {
	for i := 0; i < testIterations; i++ {
		fmt.Printf("Golang Parallel Implementation findTheLargestNumberInA2dArray - test run %d - ", i)
		findLargestNumberInArrayParallel(twoDArrayRef, largestNumberGenerated)
	}
}

func generate2DArray(emptyTwoDArrayRef *[arraySize][arraySize]int) int {
	var largestNumberGenerated = math.MinInt32

	for x := 0; x < len(emptyTwoDArrayRef); x++ {
		for y := 0; y < len(emptyTwoDArrayRef[x]); y++ {
			emptyTwoDArrayRef[x][y] = rand.Intn(math.MaxInt32-math.MinInt32) + math.MinInt32
			if largestNumberGenerated < emptyTwoDArrayRef[x][y] {
				largestNumberGenerated = emptyTwoDArrayRef[x][y]
			}
		}
	}

	return largestNumberGenerated
}

func findLargestNumberInArraySerial(array *[arraySize][arraySize]int, largestNumberGenerated int) {
	var startTime = getCurrentSystemTimeMillis()

	var largestNumberFound = findLargestNumberInArrayRange(array, 0, len(array))

	if largestNumberFound != largestNumberGenerated {
		panic("largestNumberFound != largestNumberGenerated benchmark is broken")
	}
	var endTime = getCurrentSystemTimeMillis()
	fmt.Println(strconv.FormatInt(endTime-startTime, 10) + "ms")
}

func findLargestNumberInArrayParallel(array *[arraySize][arraySize]int, largestNumberGenerated int) {
	var startTime = getCurrentSystemTimeMillis()

	var largestNumberFound = findLargestNumberInArrayParallelImplementation(array)

	if largestNumberFound != largestNumberGenerated {
		panic("largestNumberFound != largestNumberGenerated benchmark is broken")
	}

	var endTime = getCurrentSystemTimeMillis()
	fmt.Println(strconv.FormatInt(endTime-startTime, 10) + "ms")
}

func findLargestNumberInArrayParallelImplementation(array *[arraySize][arraySize]int) int {
	var largestNumbersFound [numberOfGoRoutines]chan int
	for i := range largestNumbersFound {
		largestNumbersFound[i] = make(chan int)
	}

	var scanRangeSize = len(array) / numberOfGoRoutines
	for i := 0; i < numberOfGoRoutines; i++ { //split the problem up into several smaller scans in the range of x, but scan the full length in vertical slices
		var startIndex = i * scanRangeSize
		var endIndex = (i * scanRangeSize) + scanRangeSize
		go findLargestNumberInArrayRangeParallel(array, startIndex, endIndex, largestNumbersFound[i])
	}

	var largestNumberFromAllGoRoutines = math.MinInt32
	var ithLargestNumber int
	for i := range largestNumbersFound {
		ithLargestNumber = <-largestNumbersFound[i]
		if ithLargestNumber > largestNumberFromAllGoRoutines {
			largestNumberFromAllGoRoutines = ithLargestNumber
		}
	}

	return largestNumberFromAllGoRoutines
}

func findLargestNumberInArrayRangeParallel(array *[arraySize][arraySize]int, startIndex int, endIndex int, largestNumber chan int) {
	largestNumber <- findLargestNumberInArrayRange(array, startIndex, endIndex)
}

func findLargestNumberInArrayRange(array *[arraySize][arraySize]int, startIndex int, endIndex int) int {
	var largestNumberInRange = math.MinInt32
	for x := startIndex; x < endIndex; x++ {
		for y := 0; y < len(array[x]); y++ {
			if largestNumberInRange < array[x][y] {
				largestNumberInRange = array[x][y]
			}
		}
	}

	return largestNumberInRange
}

//https://stackoverflow.com/questions/24122821/go-time-now-unixnano-convert-to-milliseconds
func getCurrentSystemTimeMillis() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func print2DArray(twoDArrayRef *[arraySize][arraySize]int) {
	for x := 0; x < len(twoDArrayRef); x++ {
		for y := 0; y < len(twoDArrayRef[x]); y++ {
			fmt.Printf("%d ", twoDArrayRef[x][y])
		}
		fmt.Println("")
	}
}
