package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailClient struct {
	Client *sendgrid.Client
	ApiKey string
}

type EmailParams struct {
	From             string
	FromName         string
	Subject          string
	To               string
	ToName           string
	BodyData         interface{}
	TemplateFilePath string
}

func NewEmailClient(apiKey string) EmailClient {
	sgclient := sendgrid.NewSendClient(apiKey)
	return EmailClient{Client: sgclient, ApiKey: apiKey}
}

func (mClient EmailClient) SendEmail(params EmailParams) error {
	m := mail.NewV3Mail()

	from := mail.NewEmail(params.FromName, params.From)
	bodyContent, err := parseHtmlTemplate(params.TemplateFilePath, params.BodyData)
	if err != nil {
		return err
	}
	content := mail.NewContent("text/html", bodyContent)
	m.SetFrom(from)
	m.AddContent(content)

	personalization := mail.NewPersonalization()
	tos := strings.Split(params.To, ";")

	tos_final := []*mail.Email{}
	for _, item := range tos {
		tos_final = append(tos_final, mail.NewEmail("", item))
	}
	personalization.AddTos(tos_final...)
	personalization.Subject = params.Subject

	m.AddPersonalizations(personalization)

	request := sendgrid.GetRequest(mClient.ApiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)
	if err != nil {
		return err
	} else {
		if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusAccepted {
			return fmt.Errorf("error sending email: %v", response.Body)
		}
	}
	return nil
}

func parseHtmlTemplate(templateFilePath string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFilePath)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	body := buf.String()
	return body, nil
}
