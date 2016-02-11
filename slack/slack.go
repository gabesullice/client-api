package slack

import (
	// stdlib
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	//"text/template"

	// internal
	"github.com/gabesullice/client-api/models"

	// external
	r "github.com/dancannon/gorethink"
	//"github.com/julienschmidt/httprouter"
)

var (
	session *r.Session
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

type SlackHandler struct {
	Session *r.Session
}

func init() {
	sess, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "client_api",
		MaxIdle:  10,
		MaxOpen:  10,
	})
	if err != nil {
		log.Fatalln(err)
	}

	session = sess
}

func NewSlackHandler() SlackHandler {
	return SlackHandler{
		Session: session,
	}
}

func (s SlackHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		log.Printf("Unable to parse form. Error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	search := req.FormValue("text")

	res, err := r.Table("contacts").Filter(func(user r.Term) r.Term {
		return user.Field("firstName").
			Match(fmt.Sprintf("(?i)%s", search))
		//Or.Field("lastName").
		//Match(fmt.Sprintf("?i%s", search))
	}).Run(s.Session)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Unable to find contacts. Error: %s", err.Error())
		return
	}

	if res.IsNil() {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No matches found for \"%s\"", search)
		return
	}

	var contacts []models.Contact
	if err := res.All(&contacts); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Unable to read contacts from result. Error: %s", err.Error())
		return
	}

	encoder := json.NewEncoder(w)

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

	w.Header().Set("Content-Type", "application/json")
	encoder.Encode(SlackMessage{
		Text:        fmt.Sprintf("Found %d contact(s)", len(contacts)),
		Attachments: attachments,
	})
	//tmpl := template.Must(template.ParseGlob("templates/*.tmpl"))

	//for _, contact := range contacts {
	//	tmpl.ExecuteTemplate(w, "contact.tmpl", contact)
	//}
}
