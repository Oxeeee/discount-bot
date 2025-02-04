package csv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Oxeeee/discont-bot/internal/domain"
)

func ConvertToCSV(records []domain.Place) (string, error) {
	var buf bytes.Buffer
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	writer := csv.NewWriter(&buf)
	writer.Comma = ';'

	header := []string{"ID", "Name", "Address", "DiscountFactor"}
	if err := writer.Write(header); err != nil {
		return "", err
	}

	for _, record := range records {
		row := []string{
			strconv.FormatUint(uint64(record.ID), 10),
			record.Name,
			record.Address,
			record.DiscountFactor,
		}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func ConvertFromCSV(csvContent string) ([]domain.Place, error) {
	reader := csv.NewReader(strings.NewReader(csvContent))
	reader.Comma = ';'

	_, err := reader.Read()
	if err != nil {
		return nil, err
	}

	var places []domain.Place

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if len(record) < 4 {
			return nil, fmt.Errorf("недостаточно полей в записи: %v", record)
		}

		id, err := strconv.ParseUint(strings.TrimSpace(record[0]), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга ID %q: %v", record[0], err)
		}

		discount := strings.TrimSpace(record[3])

		place := domain.Place{
			ID:             uint(id),
			Name:           strings.TrimSpace(record[1]),
			Address:        strings.TrimSpace(record[2]),
			DiscountFactor: discount,
		}
		places = append(places, place)
	}
	return places, nil
}
