package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"rating-renamer/components"
	"regexp"
	"strings"
)

func CreateFullFacultyFolder(basePath string, pathParts ...string) string {
	path := basePath
	if !components.CheckDir(path) {
		err := os.Mkdir(path, 0777)
		if err != nil {
			log.Fatalln(err)
		}
	}
	for _, part := range pathParts {
		path = components.NextDir(path, part, false)
	}

	return filepath.Clean(path)
}

func transliterate(inputStr string, vocabulary map[string]string) string {
	inputChars := []rune(inputStr)
	outputStr := ""
	for i := 0; i < len(inputChars); i++ {
		char := string(inputChars[i])
		enChar, hasChar := vocabulary[char]
		if !hasChar {
			log.Fatalln(char + " not in vocabulary")
		}
		outputStr = outputStr + enChar
	}
	return outputStr
}

func CreateNewGroupName(year, semestr, week, groupName, groupNumber string) string {
	return fmt.Sprintf("rating_%s_s%s_n%s_%s_%s.pdf",
		strings.Replace(year, "-", "_", 1),
		semestr,
		week,
		groupName,
		groupNumber,
	)
}

func main() {

	basePath := "D:\\work\\vc\\rating\\"

	sitePath := "D:\\work\\vc\\work_site\\web\\files\\rating\\"

	baseRuEn := map[string]string{
		"а": "a", "А": "A", "Б": "B", "б": "b", "В": "V", "в": "v", "Г": "G", "г": "g",
		"Д": "D", "д": "d", "Е": "E", "е": "e", "З": "Z", "з": "z", "И": "I", "и": "i",
		"К": "K", "к": "k", "Л": "L", "л": "l", "М": "M", "м": "m", "Н": "N", "н": "n",
		"О": "O", "о": "o", "П": "P", "п": "p", "Р": "R", "р": "r", "С": "S", "с": "s",
		"Т": "T", "т": "t", "У": "U", "у": "u", "Ф": "F", "ф": "f", "Э": "E", "э": "e",
		"Х": "H", "х": "h",
	}

	args := components.NewArguments()

	years := []string{
		"2018-2019",
		"2019-2020",
		"2020-2021",
		"2021-2022",
		"2022-2023",
	}

	faculties := []string{
		"fei",
		"fam",
		"vf",
	}

	semesters := []int{
		1, 2,
	}

	weeks := []int{
		1, 2, 3,
	}

	if args.Year != "" {
		for _, year := range years {
			if year == args.Year {
				years = []string{
					args.Year,
				}
				break
			}
		}
	}

	if args.Faculty != "" {
		for _, faculty := range faculties {
			if faculty == args.Faculty {
				faculties = []string{
					args.Faculty,
				}
				break
			}
		}
	}

	if args.Semestr != 0 {
		for _, semestr := range semesters {
			if semestr == args.Semestr {
				semesters = []int{
					args.Semestr,
				}
				break
			}
		}
	}

	if args.Week != 0 {
		for _, week := range weeks {
			if week == args.Week {
				weeks = []int{
					args.Week,
				}
				break
			}
		}
	}

	for _, faculty := range faculties {
		path := components.NextDir(basePath, faculty, len(faculties) == 1)
		for _, year := range years {
			pathWithYear := components.NextDir(path, year, len(years) == 1)

			for _, semestr := range semesters {
				pathWithSemestr := components.NextDir(pathWithYear, fmt.Sprintf("s%d", semestr), len(semesters) == 1)

				for _, week := range weeks {
					pathWithWeek := components.NextDir(pathWithSemestr, fmt.Sprintf("n%d", week), len(weeks) == 1)

					files, err := components.GetFilesFromDirectory(pathWithWeek)

					if len(files) == 0 {
						continue
					}

					siteFullPath := CreateFullFacultyFolder(sitePath, faculty, year,
						fmt.Sprintf("s%d", semestr),
						fmt.Sprintf("n%d", week),
					)

					if err != nil {
						log.Fatalln(err)
					}
					regex := regexp.MustCompile("(В[А-Я][А-Я]?)[-_](\\d\\d\\d)\\.")

					for _, filePath := range files {
						filename := filepath.Base(filePath)
						groupSubmatches := regex.FindStringSubmatch(filename)
						if groupSubmatches != nil {
							newGroup := CreateNewGroupName(year,
								fmt.Sprintf("%d", semestr),
								fmt.Sprintf("%d", week),
								transliterate(groupSubmatches[1], baseRuEn),
								groupSubmatches[2],
							)
							err = components.CopyFile(filePath, filepath.Clean(siteFullPath+"\\"+newGroup))
							if err != nil {
								log.Fatalln(err)
							}
						}
					}
				}
			}
		}
	}
}
