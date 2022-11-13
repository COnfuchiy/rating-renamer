package components

import (
	"flag"
)

type Arguments struct {
	Year    string
	Faculty string
	Semestr int
	Week    int
}

func NewArguments() Arguments {
	args := Arguments{}
	flag.StringVar(&args.Year, "year", "", "Rating year")
	flag.StringVar(&args.Faculty, "faculty", "", "fei, fam or vf")
	flag.IntVar(&args.Semestr, "s", 0, "Rating semestr")
	flag.IntVar(&args.Week, "n", 0, "Rating week")
	flag.Parse()
	return args
}
