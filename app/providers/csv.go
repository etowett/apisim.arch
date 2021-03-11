package providers

import (
	"apisim/app/models"
	"bytes"
	"encoding/csv"
	"fmt"
)

type CSVCreator interface {
	CreateMessagesCSV([]*models.Message) ([]byte, error)
}

type SimpleCSVCreator struct{}

func NewCSVCreator() *SimpleCSVCreator {
	return &SimpleCSVCreator{}
}

func (s *SimpleCSVCreator) CreateMessagesCSV(
	messages []*models.Message,
) ([]byte, error) {

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	layout := "Jan 1, 2006 15:04"

	records := make([][]string, len(messages)+1)

	records[0] = []string{"sender_id", "message", "cost", "rec_count", "sent_at"}

	for index, message := range messages {
		records[index+1] = []string{
			message.SenderID,
			message.Message,
			fmt.Sprintf("%v %v", message.Currency, message.Cost),
			fmt.Sprintf("%v", message.RecipientCount),
			message.SentAt.Format(layout),
		}
	}

	err := w.WriteAll(records)
	if err != nil {
		return []byte{}, err
	}
	w.Flush()

	return buf.Bytes(), nil
}
