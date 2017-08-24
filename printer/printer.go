package printer

import (
	"fmt"
	"strings"

	"github.com/JulzDiverse/aviator/cockpit"
	"github.com/starkandwayne/goutils/ansi"
)

type Print func(string, ...interface{}) (int, error)

func AnsiPrint(opts cockpit.MergeConf, verbose bool) {
	BeautyfulPrint(opts, ansi.Printf, verbose)
}

func BeautyfulPrint(opts cockpit.MergeConf, printf Print, verbose bool) {
	fmt.Println("SPRUCE MERGE:")
	if len(opts.Prune) != 0 {
		for _, prune := range opts.Prune {
			printf("\t@C{--prune} %s\n", prune)
		}
	}
	for _, file := range opts.Files {
		printf("\t%s\n", file)
	}
	printf("\t@G{to: %s}\n\n", opts.To)
	if verbose && (len(opts.Warnings) != 0) { //global variable
		printf("\t@Y{WARNINGS:}\n")
		for _, w := range opts.Warnings {
			sl := strings.Split(w, ":")
			printf("\t@y{%s}:@Y{%s}\n", sl[0], sl[1])
		}
		fmt.Println("\n")
	}
}
