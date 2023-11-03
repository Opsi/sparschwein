package dkb

import (
	"bufio"
	"bytes"
	"fmt"
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

type baseInfo struct {
	HolderType     string
	IBAN           string
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
			Row:  row,
			Info: info,
		})
	}
	return creators, nil
}

func parseRows(reader *bufio.Reader) ([]csvRow, error) {
	return nil, fmt.Errorf("not implemented")
}

func parseDate(dateBytes []byte) (time.Time, error) {
	// parse the date of the form "dd.mm.yyyy"
	return time.Parse("02.01.2006", string(dateBytes))
}

func parseAmountInCents(amountInCentsBytes []byte) (int, error) {
	// parse the amount of the form "1234,56 €"
	preprocessed := string(amountInCentsBytes)
	preprocessed = strings.TrimRight(preprocessed, " €")
	preprocessed = strings.TrimRight(preprocessed, " EUR")
	preprocessed = strings.Replace(preprocessed, ".", "", -1)
	preprocessed = strings.Replace(preprocessed, ",", "", -1)
	return strconv.Atoi(preprocessed)
}

func (i *baseInfo) parseFirstLine(line []byte) error {
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

func (i *baseInfo) parseThirdLine(line []byte) error {
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

func checkFirstLines(reader *bufio.Reader) (*baseInfo, error) {
	info := &baseInfo{}
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

	return info, nil
}
