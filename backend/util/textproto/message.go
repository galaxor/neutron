package textproto

import (
	"net/mail"
	"mime"
	"mime/multipart"
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

func decodeBytes(b []byte, charset string) []byte {
	var enc encoding.Encoding
	switch strings.ToLower(charset) {
	case "iso-8859-1":
		enc = charmap.ISO8859_1
	case "windows-1252":
		enc = charmap.Windows1252
	case "utf-8":
		// Nothing to do
	default:
		if charset != "" {
			log.Println("WARN: unsupported charset:", charset)
		}
	}
	if enc != nil {
		b, _ = enc.NewDecoder().Bytes(b)
	}
	return b
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
			slurp, err := ioutil.ReadAll(p)
			if err != nil {
				return err
			}

			mediaType, params, err = mime.ParseMediaType(p.Header.Get("Content-Type"))
			if (mediaType == "text/plain" && gotType == "") || mediaType == "text/html" {
				gotType = mediaType
				msg.Body = string(decodeBytes(slurp, params["charset"]))
			}
		}
	} else {
		body, err := ioutil.ReadAll(m.Body)
		if err != nil {
			return err
		}
		msg.Body = string(decodeBytes(body, params["charset"]))
	}

	return nil
}
