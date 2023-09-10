package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/tbvdm/sigtop/errio"
	"github.com/tbvdm/sigtop/signal"
)

type Record struct {
	ID              string `json:"id"`
	SchemaVersion   int    `json:"schemaVersion"`
	Type            string `json:"type"`
	SentAt          int64  `json:"sent_at"`
	Timestamp       int64  `json:"timestamp"`
	ServerTimestamp int64  `json:"serverTimestamp"`
	HasAttachments  int    `json:"hasAttachments"`
	Body            string `json:"body"`
}

func csvWriteMessages(ew *errio.Writer, msgs []signal.Message) error {
	csvWriter := csv.NewWriter(ew)
	defer csvWriter.Flush()

	fieldWhitelist := map[string]bool{
		"ID":              true,
		"SchemaVersion":   true,
		"Type":            true,
		"SentAt":          true,
		"Timestamp":       true,
		"ServerTimestamp": true,
		"HasAttachments":  true,
		"Body":            true,
	}
	fieldOrder := []string{"ID", "SchemaVersion", "Type", "SentAt", "Timestamp", "ServerTimestamp", "HasAttachments", "Body"}

	if err := csvWriter.Write(fieldOrder); err != nil {
		return err
	}

	for _, msg := range msgs {
		var record Record
		if err := json.Unmarshal([]byte(msg.JSON), &record); err != nil {
			return err
		}

		if record.Type != "incoming" && record.Type != "outgoing" {
			continue
		}

		recordValue := reflect.ValueOf(record)
		hasWhitelistedField := false
		for field := range fieldWhitelist {
			if !reflect.Indirect(recordValue).FieldByName(field).IsZero() {
				hasWhitelistedField = true
				break
			}
		}

		if !hasWhitelistedField {
			continue
		}

		record.Body = strings.ReplaceAll(record.Body, "\n", "\\\\n") // Replace line breaks with '\\n'
		record.Body = strings.ReplaceAll(record.Body, "\"", "\\\"")  // Replace double quotes with '\"'

		row := make([]string, len(fieldOrder))
		for i, field := range fieldOrder {
			fieldValue := reflect.Indirect(recordValue).FieldByName(field)
			row[i] = fmt.Sprintf("%v", fieldValue.Interface())
		}

		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}

	return ew.Err()
}
