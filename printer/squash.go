package printer

import "github.com/starkandwayne/goutils/ansi"

func AnsiPrintSquash(files []string, to string) {
	BeautyPrintSquash(files, to, ansi.Printf)
}

func BeautyPrintSquash(files []string, to string, printf Print) {
	printf("@M{SQUASH FILES:}\n")
	for _, f := range files {
		printf("\t@w{%s}\n", f)
	}
	printf("\t@M{to: %s}\n", to)
}
