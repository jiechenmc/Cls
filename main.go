package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/playwright-community/playwright-go"
)

func main() {
	// Script is meant for personal use so I
	// => Minimized the verbosity of the code by eliminating error handling

	class := os.Args[1]
	term := os.Args[2]

	pageUrl := fmt.Sprintf("http://classfind.stonybrook.edu/vufind/Search/Results?lookfor=%s&type=callnumber&filter%5B%5D=ctrlnum%3A%22%s%22", class, term)

	// Playwright Setup
	pw, _ := playwright.Run()
	browser, _ := pw.Chromium.Launch()
	page, _ := browser.NewPage()
	_, _ = page.Goto(pageUrl)

	page.WaitForLoadState("networkidle")

	// Scraping logic starts
	entries, _ := page.QuerySelectorAll("div.result")

	for _, entry := range entries {

		// Course ID
		// CSE214.R01

		cid, _ := entry.QuerySelector("b")

		// Class information
		// Intro to Object-Oriented Prog
		// by Fodor,Paul Ioan      Credit: 4.0
		// LAB :MW 08:30AM-09:50AM   LEC: TUTH 04:45PM-06:05PM
		// Attr:TECH

		cinfo, _ := entry.QuerySelector("div.span-11")

		// Status of the class
		// Waitlist
		// Open
		// Closed
		cstatus, _ := entry.QuerySelector("div.span-3 > div")

		// The course/class number
		// 48060
		cnumber, _ := entry.QuerySelector("div.span-3 div.hide[style=\"display: block;\"]")

		// converting the elements into strings
		courseInfo, _ := cinfo.InnerText()
		courseStatus, _ := cstatus.InnerText()

		// if course is closed don't show it
		if strings.Contains(courseStatus, "Closed") {
			continue
		}

		courseNumber, _ := cnumber.TextContent()

		courseID, _ := cid.TextContent()

		fmt.Println()
		color.Cyan(courseID + " (" + strings.TrimSpace(courseNumber) + ")")
		color.Blue(courseInfo)

		// adding colors to output
		if strings.Contains(courseStatus, "Open") {
			color.Green(courseStatus)
		} else if strings.Contains(courseStatus, "Closed") {
			color.Red(courseStatus)
		} else {
			color.Yellow(courseStatus)
		}
	}

	// Termination
	browser.Close()
	pw.Stop()
}
