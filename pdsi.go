package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	blue    = "\033[34m"
	red     = "\033[31m"
	green   = "\033[32m"
	greenbg = "\x1b[42m"
	reset   = "\033[0m"
	yellow  = "\033[33m"
)

var sensitiveKeywords = []string{

	// personal & confidential
	"confidential", "private", "restricted", "internal",
	"not for distribution", "do not share", "proprietary",
	"trade secret", "classified", "sensitive",

	// financial & legal
	"bank statement", "invoice", "salary", "contract",
	"agreement", "non disclosure",

	// personal data
	"passport", "social security", "ssn", "date of birth",
	"credit card", "identity", "id number",

	// company internal
	"internal use only", "company confidential",
	"staff only", "management only", "internal only",
}

func highlightKeywords(input string) string {
	for _, keyword := range sensitiveKeywords {
		pattern := fmt.Sprintf(`(?i)\b%s\b`, regexp.QuoteMeta(keyword))
		re := regexp.MustCompile(pattern)
		input = re.ReplaceAllStringFunc(input, func(match string) string {
			return fmt.Sprintf("\033[33m%s%s", match, reset)
		})
	}
	return input
}

func analyzePDF(input string, isURL bool) {
	if isURL {
		fmt.Printf("%s[-]%s %sreading URL:%s %s\n", blue, reset, blue, reset, input)

		cmd := exec.Command("wget", "--no-check-certificate", "-qO-", input)
		output, err := cmd.Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s[!]%s %snot fetching:%s %s\n", red, reset, red, reset, input)
			return
		}

		pdftotextCmd := exec.Command("pdftotext", "-", "-")
		pdftotextCmd.Stdin = bytes.NewReader(output)
		pdfOutput, err := pdftotextCmd.Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s[!]%s %snot processing:%s %s\n", red, reset, red, reset, input)
			return
		}

		processPDFOutput(pdfOutput, input)
	} else {
		fmt.Printf("%s[-]%s %sreading file:%s %s\n", blue, reset, blue, reset, input)

		pdftotextCmd := exec.Command("pdftotext", input, "-")
		pdfOutput, err := pdftotextCmd.Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s[!]%s %snot processing:%s %s\n", red, reset, red, reset, input)
			return
		}

		processPDFOutput(pdfOutput, input)
	}
}

func processPDFOutput(pdfOutput []byte, source string) {
	scanner := bufio.NewScanner(bytes.NewReader(pdfOutput))
	foundSensitive := false
	for scanner.Scan() {
		line := scanner.Text()
		highlighted := highlightKeywords(line)
		if highlighted != line {
			fmt.Printf("%s[+]%s %ssensitive:%s %s\n", green, reset, greenbg, reset, highlighted)
			foundSensitive = true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%s[!]%s %sfailed to scan:%s %s\n", red, reset, red, reset, source)
	}

	if !foundSensitive {
		fmt.Printf("%s[*]%s %snothing found%s\n", blue, reset, blue, reset)
	}
}

func main() {
	fmt.Println(yellow + `
               __     _ 
    ____  ____/ /____(_)
   / __ \/ __  / ___/ / 
  / /_/ / /_/ (__  ) /  
 / .___/\__,_/____/_/   ` + red + "\033]8;;https://github.com/1hehaq\033\\@1hehaq\033]8;;\033\\" + yellow + `
/_/

` + reset)

	localPDF := flag.String("local", "", "path to local PDF files, separated by commas or stdin")
	matchKeywords := flag.String("match", "", "custom keywords to match, separated by commas")
	flag.Parse()

	pdfRegex := regexp.MustCompile(`(?i)\.pdf$`)
	urlRegex := regexp.MustCompile(`^https?://`)

	if *matchKeywords != "" {
		sensitiveKeywords = strings.Split(*matchKeywords, ",")
		for i := range sensitiveKeywords {
			sensitiveKeywords[i] = strings.TrimSpace(sensitiveKeywords[i])
		}
	}

	processInput := func(input string) {
		if !pdfRegex.MatchString(input) {
			fmt.Fprintf(os.Stderr, "%s[!]%s %sinvalid file ext!:%s %s\n", red, reset, red, reset, input)
			return
		}

		if _, err := os.Stat(input); err == nil {
			analyzePDF(input, false)
		} else if urlRegex.MatchString(input) {
			analyzePDF(input, true)
		} else {
			fmt.Fprintf(os.Stderr, "%s[!]%s %sfile not found:%s %s\n", red, reset, red, reset, input)
		}
	}

	if *localPDF != "" {
		files := strings.Split(*localPDF, ",")
		for _, file := range files {
			processInput(strings.TrimSpace(file))
		}
		return
	}

	info, _ := os.Stdin.Stat()
	if info.Mode()&os.ModeCharDevice == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			processInput(strings.TrimSpace(scanner.Text()))
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "%s[!]%s %serror reading input:%s %v\n", red, reset, red, reset, err)
			os.Exit(1)
		}
		return
	}

	fmt.Fprintf(os.Stderr, "%s[!]%s %sno input provided%s\n", red, reset, red, reset)
	os.Exit(1)
}
