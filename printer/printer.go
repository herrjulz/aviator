package printer

import (
	"fmt"
	"strings"

	"github.com/JulzDiverse/aviator"
	"github.com/starkandwayne/goutils/ansi"
)

type Print func(string, ...interface{}) (int, error)

func AnsiPrint(opts aviator.MergeConf, to string, warnings []string, verbose bool) {
	BeautyfulPrint(opts, to, warnings, verbose, ansi.Printf)
}

func BeautyfulPrint(opts aviator.MergeConf, to string, warnings []string, verbose bool, printf Print) {
	printf("@G{SPRUCE MERGE:}\n")
	if len(opts.Prune) != 0 {
		for _, prune := range opts.Prune {
			printf("\t@C{--prune} %s\n", prune)
		}
	}
	for _, file := range opts.Files {
		printf("\t%s\n", file)
	}
	printf("\t@G{to: %s}\n\n", to)
	if verbose && (len(warnings) > 0) { //global variable
		printf("\t@Y{WARNINGS:}\n")
		for _, w := range warnings {
			sl := strings.Split(w, ":")
			printf("\t@y{%s}:@Y{%s}\n", sl[0], sl[1])
		}
		fmt.Println("\n")
	}
}
