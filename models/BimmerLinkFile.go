package models

import (
	"strconv"
	"strings"
)

type DataRow struct {
	Entries []float64
}

type DictEntry struct {
	Key string
	Min float64
	Max float64
	Total float64
	Count uint32
	Average float64
	Factor float64
}

type BimmerFile struct{
	MaxDict map[string]*DictEntry
	Rows []*DataRow
	Names []string
	FileName string
}

func parseHeaderRow(row string) []string {
	var trimmed = strings.TrimSpace(row)
	var rows []string
	var inQuote = false
	var current = ""
	for _, char := range trimmed {
		if char == ',' && !inQuote {
			rows = append(rows, current)
			current=""
			continue
		}
		if char == '"' {
			inQuote = !inQuote
			continue
		}
		current += string(char)
	}
	if current != "" {
		rows = append(rows, current)
	}
	return rows
}
func learnValue(key string, value float64, maxMap map[string]*DictEntry)  {
	if entry, ok := maxMap[key]; ok {
		entry.Count += 1
		entry.Total += value
		if value > entry.Max {
			entry.Max = value
		}
		if value < entry.Min {
			entry.Min = value
		}
		entry.Average = entry.Total / float64(entry.Count)


		maxMap[key] = entry
	} else {
		maxMap[key] = &DictEntry{
			Key:     key,
			Min:     value,
			Max:     value,
			Total:   value,
			Count:   1,
			Average: 0,
			Factor: 1,
		}
	}
}
func ParseBimmerFile(raw string) *BimmerFile {
	var lines = strings.Split(strings.TrimSpace(raw), "\n")
	var headers = parseHeaderRow(lines[0])
	var rawDataLines = lines[1:]
	var dataRows []*DataRow
	for _, rowString := range rawDataLines {
		var trimmed = strings.TrimSpace(rowString)
		var points = strings.Split(trimmed, ",")
		var r = &DataRow{Entries: []float64{}}
		for _, p := range points {
			floatVal, err := strconv.ParseFloat(p, 64)
			if err != nil{
				panic(err)
			}
			r.Entries = append(r.Entries, floatVal)
		}
		dataRows = append(dataRows, r)
	}
	var maxDict = make(map[string]*DictEntry)
	for _, row := range dataRows {
		for index, value := range row.Entries {
			learnValue(headers[index], value, maxDict)
		}
	}
	return &BimmerFile{
		MaxDict: maxDict,
		Rows:    dataRows,
		Names:   headers,
	}
}