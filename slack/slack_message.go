package slack

import (
	// stdlib
	"fmt"
	"strings"

	// internal
	"github.com/gabesullice/client-api/models"
)

type SlackMessage struct {
	Text        string            `json:"text"`
	Attachments []SlackAttachment `json:"attachments"`
}

type SlackAttachment struct {
	Text   string                 `json:"text"`
	Fields []SlackAttachmentField `json:"fields"`
}

type SlackAttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

func newSlackMessage(contacts []models.Contact) SlackMessage {
	var attachments []SlackAttachment
	for _, contact := range contacts {
		attachments = append(attachments, SlackAttachment{
			Text: contact.FirstName,
			Fields: []SlackAttachmentField{
				SlackAttachmentField{
					Title: "Emails",
					Value: strings.Join(contact.Emails, ", "),
					Short: true,
				},
				SlackAttachmentField{
					Title: "Phones",
					Value: fmt.Sprintf("%s: %s", contact.Phones[0].Type, contact.Phones[0].Number),
					Short: true,
				},
			},
		})
	}

	return SlackMessage{
		Text:        fmt.Sprintf("Found %d contact(s)", len(contacts)),
		Attachments: attachments,
	}
}
