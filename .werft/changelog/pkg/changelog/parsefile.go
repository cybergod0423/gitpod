package changelog

import (
	"bufio"
	"os"
	"regexp"
	"strconv"

	logger "github.com/sirupsen/logrus"
)

var noteLineRegexp = regexp.MustCompile(`\*\s.*\s\[\[#(\d*)\]\(`)

func parseFile(path string) (content string, prNumber int) {
	content = ""
	prNumber = 0
	if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err) {
		return
	}

	file, err := os.Open(path)
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		match := noteLineRegexp.FindStringSubmatch(scanner.Text())
		if len(match) > 1 {
			if prNumber == 0 {
				prNumber, err = strconv.Atoi(match[1])
				if err != nil {
					logger.Errorf("Ignoring invalid PR number %s", match[1])
				}
			}
			content = content + scanner.Text()
		}
	}
	if err := scanner.Err(); err != nil {
		logger.Fatal(err)
	}
	return
}
