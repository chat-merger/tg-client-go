package merger

import "fmt"

func (m *Message) FormatFull() string {
	str := fmt.Sprintf("from: %s", m.From)
	if m.ReplyId != nil {
		str += fmt.Sprintf("reply to: %s", *m.ReplyId)
	}
	str += "\n" + m.Date.Format("15:05 02 Jan")
	if m.Author != nil {
		str += "\n" + *m.Author
	}
	switch m.Body.(type) {
	case *BodyText:
		str += "\n" + m.Body.(*BodyText).Value
	case *BodyMedia:
		mb := m.Body.(*BodyMedia)
		if mb.Caption != nil {
			str += "\n" + *mb.Caption
		}
		str += "\n" + mb.Url
	}
	return str
}

func (m *Message) FormatShort() string {
	str := ""
	if m.ReplyId != nil {
		str += fmt.Sprintf("reply to: %s", *m.ReplyId)
	}
	if m.Author != nil {
		str += fmt.Sprintf("\n[%s]: ", *m.Author)
	} else {
		str += "\n" + "msg = "
	}
	switch m.Body.(type) {
	case *BodyText:
		str += m.Body.(*BodyText).Value
	case *BodyMedia:
		mb := m.Body.(*BodyMedia)
		if mb.Caption != nil {
			str += *mb.Caption
		}
		str += "\n" + mb.Url
	}
	return str
}
