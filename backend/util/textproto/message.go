package textproto

import (
	"bytes"
	"encoding/base64"
	"net/mail"
	"net/textproto"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"strings"
	"io"
	"io/ioutil"
	"log"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"github.com/emersion/neutron/backend"
)

func ParseMessageHeader(msg *backend.Message, header *mail.Header) {
	msg.Subject = DecodeWord(header.Get("Subject"))

	from, err := header.AddressList("From")
	if err == nil && len(from) > 0 {
		msg.Sender = ParseEmail(from[0])
	}

	to, err := header.AddressList("To")
	if err == nil {
		for _, addr := range to {
			msg.ToList = append(msg.ToList, ParseEmail(addr))
		}
	}

	cc, err := header.AddressList("Cc")
	if err == nil {
		for _, addr := range cc {
			msg.ToList = append(msg.ToList, ParseEmail(addr))
		}
	}

	bcc, err := header.AddressList("Bcc")
	if err == nil {
		for _, addr := range bcc {
			msg.ToList = append(msg.ToList, ParseEmail(addr))
		}
	}

	replyTo, err := header.AddressList("From")
	if err == nil && len(replyTo) > 0 {
		msg.ReplyTo = ParseEmail(replyTo[0])
	}

	time, err := header.Date()
	if err == nil {
		msg.Time = time.Unix()
	}
}

func decoder(r io.Reader, charset string) io.Reader {
	var enc encoding.Encoding
	switch strings.ToLower(charset) {
	case "iso-8859-1":
		enc = charmap.ISO8859_1
	case "windows-1252":
		enc = charmap.Windows1252
	case "utf-8", "us-ascii":
		// Nothing to do
	default:
		if charset != "" {
			log.Println("WARN: unsupported charset:", charset)
		}
	}
	if enc != nil {
		r = enc.NewDecoder().Reader(r)
	}
	return r
}

func ParseMessageBody(msg *backend.Message, m *mail.Message) error {
	mediaType, params, err := mime.ParseMediaType(m.Header.Get("Content-Type"))
	if err != nil {
		return err
	}

	gotType := ""
	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(m.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}

			mediaType, typeParams, _ := mime.ParseMediaType(p.Header.Get("Content-Type"))
			disp, dispParams, _ :=  mime.ParseMediaType(p.Header.Get("Content-Disposition"))

			var r io.Reader
			r = p
			if typeParams["charset"] != "" {
				r = decoder(r, typeParams["charset"])
			}

			slurp, err := ioutil.ReadAll(r)
			if err != nil {
				return err
			}

			if mediaType == "text/plain" {
				if gotType == "" {
					disp = "inline"
				} else {
					disp = "attachment"
				}
			} else if mediaType == "text/html" {
				disp = "inline"
			} else {
				disp = "attachment"
			}

			switch disp {
			case "inline":
				gotType = mediaType
				msg.Body = string(slurp)
			case "attachment":
				attachment := &backend.Attachment{
					Name: dispParams["filename"],
					Size: len(slurp),
					MIMEType: mediaType,
				}

				msg.Attachments = append(msg.Attachments, attachment)
			default:
				log.Println("WARN: unsupported Content-Disposition:", disp)
			}
		}
	} else {
		var r io.Reader
		r = m.Body
		if params["charset"] != "" {
			r = decoder(r, params["charset"])
		}

		body, err := ioutil.ReadAll(m.Body)
		if err != nil {
			return err
		}

		msg.Body = string(body)
	}

	return nil
}


func formatMessage(header textproto.MIMEHeader, body string) string {
	return FormatHeader(header) + "\r\n" + body
}

func FormatMessage(msg *backend.Message) string {
	header := GetMessageHeader(msg)
	header.Set("Content-Type", "text/html")
	return formatMessage(header, msg.Body)
}

func FormatOutgoingMessage(msg *backend.OutgoingMessage) string {
	var b bytes.Buffer
	m := multipart.NewWriter(&b)

	h := textproto.MIMEHeader{}
	h.Set("Content-Type", "text/html; charset=UTF-8")
	h.Set("Content-Disposition", "inline")
	h.Set("Content-Transfer-Encoding", "quoted-printable")
	w, _ := m.CreatePart(h)
	enc := quotedprintable.NewWriter(w)
	enc.Write([]byte(msg.MessagePackage.Body))
	enc.Close()

	for _, att := range msg.Attachments {
		h := textproto.MIMEHeader{}
		h.Set("Content-Type", att.MIMEType)
		h.Set("Content-Disposition", "attachment; filename=\"" + att.Name + "\"")
		h.Set("Content-Transfer-Encoding", "base64")

		w, _ := m.CreatePart(h)
		enc := base64.NewEncoder(base64.StdEncoding, w)
		enc.Write(att.Data)
		enc.Close()
	}

	m.Close()

	mh := GetOutgoingMessageHeader(msg)
	mh.Set("Content-Type", "multipart/mixed; boundary=" + m.Boundary())

	return formatMessage(mh, b.String())
}
