package main

import (
	"github.com/starkandwayne/goutils/ansi"
)

func printMergeCombinationError(err error) {
	ansi.Printf("%s\n\n", err.Error())
	ansi.Printf("Use this 'merge' params, as separate array entries. Example:\n%s", mergeCombination)
}

func printForEachCombinationError(err error) {
	ansi.Printf("%s\n\n", err.Error())
	ansi.Printf("Use 'for_each' either with 'files' or 'in' parameter. Example :\n%s", forEachCombination)
}

func printMergeWithCombinationError(err error) {
	ansi.Printf("%s\n\n", err.Error())
	ansi.Printf("Example:\n%s", withCombination)
}

func printForEachFilesCombinationError(err error) {
	ansi.Printf("%s\n\n", err.Error())
	ansi.Printf("Example:\n%s", forEachFilesCombination)
}

func printForEachInCombinationError(err error) {
	ansi.Printf("%s\n\n", err.Error())
	ansi.Printf("Example:\n%s", forEachFilesCombination)
}

func printForEachWalkCombinationError(err error) {
	ansi.Printf("%s\n\n", err.Error())
	ansi.Printf("Example:\n%s", forEachWalkCombination)
}

func printMergeRegexpCombinationError(err error) {
	ansi.Printf("%s\n\n", err.Error())
	ansi.Printf("Example:\n%s", mergeRegexpCombination)
}

func printMergeExceptCombinationError(err error) {
	ansi.Printf("%s\n\n", err.Error())
	ansi.Printf("Example:\n%s", mergeExceptCombination)
}

func printForEachRegexpCombinationError(err error) {
	ansi.Printf("%s\n\n", err.Error())
	ansi.Printf("Example:\n%s", forEachRegexpCombination)
}
