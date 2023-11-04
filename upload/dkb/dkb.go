package dkb

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Opsi/sparschwein/upload"
)

var (
	// first line has the form
	//
	// "Konto";"<holder type> <iban>"
	//
	// we want to extract the type and the iban
	firstLineRegex = regexp.MustCompile(`"Konto";"(.+) (.+)"`)

	// third line has the form
	//
	// "Kontostand vom <date>:";"<balance> EUR"
	//
	// we want to extract the date and the balance
	thirdLineRegex = regexp.MustCompile(`"Kontostand vom (.+):";"(.+) EUR"`)
)

type headerInfo struct {
	account
	Date           time.Time
	BalanceInCents int
}

func ParseCSV(csvData []byte) ([]upload.TransactionCreator, error) {
	// first we ne to trim down the first 4 lines
	reader := bufio.NewReader(bytes.NewReader(csvData))

	info, err := checkFirstLines(reader)
	if err != nil {
		return nil, fmt.Errorf("check first lines: %w", err)
	}

	rows, err := parseRows(reader)
	if err != nil {
		return nil, fmt.Errorf("parse rows: %w", err)
	}

	creators := make([]upload.TransactionCreator, 0)
	for _, row := range rows {
		creators = append(creators, transactionCreator{
			Row:     row,
			Account: &info.account,
		})
	}
	return creators, nil
}

func parseRows(reader *bufio.Reader) ([]csvRow, error) {
	// Create a new reader
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'
	csvReader.FieldsPerRecord = 11

	rows := make([]csvRow, 0)
	recordCount := 0
	for {
		// Read each record
		record, err := csvReader.Read()
		recordCount++
		if err == io.EOF {
			// If we reached the end of the file, break the loop
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read csv record %d: %w", recordCount, err)
		}

		row, err := parseRow(record)
		if err != nil {
			slog.Warn("skipping record because of error",
				slog.Any("record", record),
				slog.String("error", err.Error()))
			continue
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func parseRow(row []string) (csvRow, error) {
	if len(row) != 11 {
		return csvRow{}, fmt.Errorf("row has %d fields, expected 11", len(row))
	}
	status := strings.TrimSpace(row[2])
	if status == "Vorgemerkt" {
		return csvRow{}, fmt.Errorf("row is not yet booked")
	}

	bookingDate, err := parseDate([]byte(row[0]))
	if err != nil {
		return csvRow{}, fmt.Errorf("parse booking date: %w", err)
	}
	valueDate, err := parseDate([]byte(row[1]))
	if err != nil {
		return csvRow{}, fmt.Errorf("parse value date: %w", err)
	}
	amountInCents, err := parseAmountInCents([]byte(row[7]))
	if err != nil {
		return csvRow{}, fmt.Errorf("parse amount in cents: %w", err)
	}
	return csvRow{
		BookingDate:       bookingDate,
		ValueDate:         valueDate,
		Status:            status,
		Payer:             row[3],
		Payee:             row[4],
		Purpose:           row[5],
		TransactionType:   row[6],
		AmountInCents:     amountInCents,
		CreditorID:        row[8],
		MandateReference:  row[9],
		CustomerReference: row[10],
	}, nil
}

func parseDate(dateBytes []byte) (time.Time, error) {
	// parse the date of the form "dd.mm.yyyy"
	date, err := time.Parse("02.01.2006", string(dateBytes))
	if err == nil {
		return date, nil
	}
	// parse the date of the form "dd.mm.yy"
	date, err = time.Parse("02.01.06", string(dateBytes))
	if err == nil {
		return date, nil
	}
	return date, fmt.Errorf("parse date: %w", err)
}

func parseAmountInCents(amountInCentsBytes []byte) (int, error) {
	// parse the amount of the form "1234,56 €"
	preprocessed := string(amountInCentsBytes)
	preprocessed = strings.TrimRight(preprocessed, "€")
	preprocessed = strings.TrimRight(preprocessed, "EUR")
	preprocessed = strings.Replace(preprocessed, ".", "", -1)
	preprocessed = strings.Replace(preprocessed, ",", "", -1)
	preprocessed = strings.TrimSpace(preprocessed)
	return strconv.Atoi(preprocessed)
}

func (i *headerInfo) parseFirstLine(line []byte) error {
	matches := firstLineRegex.FindSubmatch(line)
	if len(matches) != 3 {
		return fmt.Errorf("1st line does not match regex")
	}
	holderType := strings.TrimSpace(string(matches[1]))
	if holderType == "" {
		return fmt.Errorf("holder type is empty")
	}
	iban := strings.TrimSpace(string(matches[2]))
	if iban == "" {
		return fmt.Errorf("iban is empty")
	}
	i.HolderType = holderType
	i.IBAN = iban
	return nil
}

func (i *headerInfo) parseThirdLine(line []byte) error {
	matches := thirdLineRegex.FindSubmatch(line)
	if len(matches) != 3 {
		return fmt.Errorf("3rd line does not match regex")
	}
	// parse the date of the form "dd.mm.yyyy"
	parsedDate, err := parseDate(matches[1])
	if err != nil {
		return fmt.Errorf("parse date: %w", err)
	}
	// parse the balance of the form "1234,56"
	balanceInCents, err := parseAmountInCents(matches[2])
	if err != nil {
		return fmt.Errorf("parse balance: %w", err)
	}
	i.Date = parsedDate
	i.BalanceInCents = balanceInCents
	return nil
}

func checkFirstLines(reader *bufio.Reader) (*headerInfo, error) {
	info := &headerInfo{}
	// read 1st line
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("read 1st line: %w", err)
	}
	err = info.parseFirstLine(line)
	if err != nil {
		return nil, fmt.Errorf("parse 1st line: %w", err)
	}

	// read 2nd line
	line, err = reader.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("read 2nd line: %w", err)
	}
	if string(line) != "\"\"\n" {
		return nil, fmt.Errorf("2nd line should be two quotes")
	}

	// read 3rd line
	line, err = reader.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("read 3rd line: %w", err)
	}
	err = info.parseThirdLine(line)
	if err != nil {
		return nil, fmt.Errorf("parse 3rd line: %w", err)
	}

	// fourth line is the header
	line, err = reader.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("read 4th line: %w", err)
	}
	if string(line) != "\"\"\n" {
		return nil, fmt.Errorf("4th line should be two quotes")
	}

	// TODO: check 5th line (column names)
	_, err = reader.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("read 5th line: %w", err)
	}

	return info, nil
}
