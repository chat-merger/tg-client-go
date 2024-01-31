package merger

import "fmt"

func (m *Message) FormatShort() string {
	str := ""
	if m.Username != nil {
		str += fmt.Sprintf("\n[%s]: ", *m.Username)
	} else {
		str += "\n" + "msg = "
	}
	if m.Text != nil {
		str += *m.Text
	}
	if len(m.Media) != 0 {
		for _, media := range m.Media {
			str += "\n" + media.Url
		}
	}
	return str
}
