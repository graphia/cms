package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var queue []emailRecorder

func init() {
	mailer = Mailer{send: testSender}
}

// MailYak email representation
type emailRecorder struct {
	To       string
	From     string
	FromName string
	Subject  string
	Body     string
}

// clone of DefaultSender with mailyak replaced by emailRecorder
func testSender(ed emailData) error {

	Debug.Println("ed", ed.Subject)
	er := emailRecorder{}

	// non-changing attributes
	er.From = "noreply@graphia.co.uk"
	er.FromName = "Graphia CMS"

	// volatile attributes
	er.To = ed.User.Email
	er.Subject = ed.Subject
	er.Body = ed.Body

	queue = append(queue, er)
	return nil
}

func Test_sendEmailConfirmation(t *testing.T) {

	var er emailRecorder

	sendEmailConfirmation(mh)

	er = queue[len(queue)-1]

	// non-changing attributes
	assert.Equal(t, "Graphia CMS", er.FromName)
	assert.Equal(t, "noreply@graphia.co.uk", er.From)

	// volatile attributes
	assert.Equal(t, "Welcome to Graphia CMS", er.Subject)
	assert.Equal(t, mh.Email, er.To)

	// body contents 🙆‍
	assert.Contains(t, er.Body, fmt.Sprintf("Dear %s,", mh.Name))
	assert.Contains(t, er.Body, fmt.Sprintf("Your username is %s", mh.Username))
	assert.Contains(t, er.Body, fmt.Sprintf("%s/cms/activate/%s", config.URL, mh.ConfirmationKey))

}
