package printer

//func printWarnings() {
//ansi.Printf("@R{WARNINGS:}\n")
//for _, w := range Warnings {
//sl := strings.Split(w, ":")
//ansi.Printf("@p{%s}:@P{%s}\n", sl[0], sl[1])
//}
//}

//type Printer interface {
//// ...
//}

//var ansiPrinter = &AnsiPrinter{}

//func BeautifyPrint(opts spruce.MergeOpts, dest string) {
//beautifyPrint(opts, dest, ansiPrinter)
//}

//func (...) Println(msg string) {
//outputString += msg
//}

//func (p *MyPrinter) beautifyPrint(opts spruce.MergeOpts, dest string, printer Printer) {
//y := color.New(color.FgYellow, color.Bold)
//r := color.New(color.FgHiRed)
//c := color.New(color.FgHiCyan)
//fmt.Println("SPRUCE MERGE:")
//if len(opts.Prune) != 0 {
//for _, prune := range opts.Prune {
//r.Printf("\t%s ", "--prune")
//c.Printf("  %s \n", prune)
//}
//}
//for _, file := range opts.Files {
//fmt.Printf("\t%s \n", file)
//}
//y.Printf("\tto: %s\n\n", dest)
//if Verbose && (len(Warnings) != 0) { //global variable
//ansi.Printf("\t@R{WARNINGS:}\n")
//for _, w := range Warnings {
//sl := strings.Split(w, ":")
//ansi.Printf("\t@p{%s}:@P{%s}\n", sl[0], sl[1])
//}
//fmt.Println("\n")
//}
//}
