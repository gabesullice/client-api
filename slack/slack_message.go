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
		text := "Unknown"

		if len(contact.FirstName)+len(contact.LastName) > 0 {
			text = strings.Join([]string{contact.FirstName, contact.LastName}, " ")
		} else if len(contact.Name) > 0 {
			text = contact.Name
		}

		attachment := SlackAttachment{Text: text}

		addEmailAttachments(&attachment, contact)
		addPhoneAttachments(&attachment, contact)
		addNotesAttachment(&attachment, contact)

		attachments = append(attachments, attachment)
	}

	return SlackMessage{
		Text:        fmt.Sprintf("Found %d contact(s)", len(contacts)),
		Attachments: attachments,
	}
}

func addPhoneAttachments(attachment *SlackAttachment, contact models.Contact) {
	var numbers []string
	for _, phone := range contact.Phones {
		if len(phone.Number) > 0 {
			numbers = append(numbers, fmt.Sprintf("%s: %s", phone.Type, phone.Number))
		}
	}

	if len(numbers) > 0 {
		attachment.Fields = append(attachment.Fields, SlackAttachmentField{
			Title: "Phones",
			Value: strings.Join(numbers, "\r"),
			Short: true,
		})
	}
}

func addEmailAttachments(attachment *SlackAttachment, contact models.Contact) {
	var emails []string
	for _, email := range contact.Emails {
		if len(email) > 0 {
			emails = append(emails, email)
		}
	}

	if len(emails) > 0 {
		attachment.Fields = append(attachment.Fields, SlackAttachmentField{
			Title: "Emails",
			Value: strings.Join(emails, "\r"),
			Short: true,
		})
	}
}

func addNotesAttachment(attachment *SlackAttachment, contact models.Contact) {
	if len(contact.Notes) > 0 {
		attachment.Fields = append(attachment.Fields, SlackAttachmentField{
			Title: "Notes",
			Value: contact.Notes,
			Short: true,
		})
	}
}
