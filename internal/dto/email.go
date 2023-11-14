package dto

type Email struct {
	From        string
	To          string
	Subject     string
	Body        string
	Attachments []string
}
