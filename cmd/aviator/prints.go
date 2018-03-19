package main

import (
	"github.com/starkandwayne/goutils/ansi"
)

func printMergeCombinationError(err error) {
	ansi.Printf("%s\n\n", err.Error())
	ansi.Printf("Use this 'merge' params, as separate array entries. Example:\n@G{%s}", mergeCombination)
}

func printForEachCombinationError(err error) {
	ansi.Printf("%s\n\n", err.Error())
	ansi.Printf("Use 'for_each' either with 'files' or 'in' parameter. Example :\n@G{%s}", forEachCombination)
}

func printMergeWithCombinationError(err error) {
	ansi.Printf("%s\n\n", err.Error())
	ansi.Printf("Example:\n@G{%s}", withCombination)
}

func printForEachFilesCombinationError(err error) {
	ansi.Printf("%s\n\n", err.Error())
	ansi.Printf("Example:\n@G{%s}", forEachFilesCombination)
}

func printForEachInCombinationError(err error) {
	ansi.Printf("%s\n\n", err.Error())
	ansi.Printf("Example:\n@G{%s}", forEachFilesCombination)
}

func printForEachWalkCombinationError(err error) {
	ansi.Printf("%s\n\n", err.Error())
	ansi.Printf("Example:\n@G{%s}", forEachWalkCombination)
}

func printMergeRegexpCombinationError(err error) {
	ansi.Printf("%s\n\n", err.Error())
	ansi.Printf("Example:\n@G{%s}", mergeRegexpCombination)
}

func printMergeExceptCombinationError(err error) {
	ansi.Printf("%s\n\n", err.Error())
	ansi.Printf("Example:\n@G{%s}", mergeExceptCombination)
}

func printForEachRegexpCombinationError(err error) {
	ansi.Printf("%s\n\n", err.Error())
	ansi.Printf("Example:\n@G{%s}", forEachRegexpCombination)
}
