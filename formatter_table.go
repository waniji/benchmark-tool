package main

import (
	"fmt"
	"strconv"
	"strings"
)

type FormatterTable struct {
	rows []FormatterTableRow
}

type FormatterTableRow struct {
	headerName   string
	maxRowLength int
	printFormat  string
}

func (row *FormatterTableRow) format() string {
	if row.printFormat[1] == '-' {
		return fmt.Sprintf("%%-%d%c", row.maxRowLength, row.printFormat[2])
	} else {
		return fmt.Sprintf("%%%d%c", row.maxRowLength, row.printFormat[1])
	}
}

func (f *FormatterTable) Print(results Results) {
	f.initialize()
	f.calcMaxRowLength(results)

	f.printHeader()
	f.printBorder()
	f.printResults(results[1:])
	f.printBorder()
	f.printResults(results[0:1])
}

func (f *FormatterTable) initialize() {
	f.rows = []FormatterTableRow{
		{
			headerName:  "Success",
			printFormat: "%d",
		},
		{
			headerName:  "Failure",
			printFormat: "%d",
		},
		{
			headerName:  "Maximum(msec)",
			printFormat: "%d",
		},
		{
			headerName:  "Minimum(msec)",
			printFormat: "%d",
		},
		{
			headerName:  "Average(msec)",
			printFormat: "%d",
		},
		{
			headerName:  "URL",
			printFormat: "%-s",
		},
	}
}

func (f *FormatterTable) calcMaxRowLength(results Results) {
	for i := 0; i < len(f.rows); i++ {
		f.rows[i].maxRowLength = len(f.rows[i].headerName)
	}
	for _, result := range results {
		resultLength := []int{
			len(strconv.Itoa(result.success)),
			len(strconv.Itoa(result.failure)),
			len(result.maximumElapsedMsec.String()),
			len(result.minimumElapsedMsec.String()),
			len(result.averageElapsedMsec().String()),
			len(result.url),
		}
		for i, length := range resultLength {
			if length > f.rows[i].maxRowLength {
				f.rows[i].maxRowLength = length
			}
		}
	}
}

func (f *FormatterTable) printHeader() {
	row := []interface{}{
		f.rows[0].headerName,
		f.rows[1].headerName,
		f.rows[2].headerName,
		f.rows[3].headerName,
		f.rows[4].headerName,
		f.rows[5].headerName,
	}
	fmt.Println("")
	f.printData(row, "%%-%ds")
}

func (f *FormatterTable) printBorder() {
	row := []interface{}{
		strings.Repeat("-", f.rows[0].maxRowLength),
		strings.Repeat("-", f.rows[1].maxRowLength),
		strings.Repeat("-", f.rows[2].maxRowLength),
		strings.Repeat("-", f.rows[3].maxRowLength),
		strings.Repeat("-", f.rows[4].maxRowLength),
		strings.Repeat("-", f.rows[5].maxRowLength),
	}
	f.printData(row, "%%-%ds")
}

func (f *FormatterTable) printResults(results Results) {
	for _, result := range results {
		row := []interface{}{
			result.success,
			result.failure,
			result.maximumElapsedMsec,
			result.minimumElapsedMsec,
			result.averageElapsedMsec(),
			result.url,
		}
		f.printData(row, "")
	}
}

func (f *FormatterTable) printData(data []interface{}, format string) {
	printData := "|"
	for i, datum := range data {
		var aaa string
		if format == "" {
			aaa = f.rows[i].format()
		} else {
			aaa = fmt.Sprintf(format, f.rows[i].maxRowLength)
		}
		printData += " " + fmt.Sprintf(aaa, datum) + " |"
	}
	fmt.Println(printData)
}
