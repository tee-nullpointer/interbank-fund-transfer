package domain

type ISO8583Message struct {
	MTI    string         `json:"mti"`
	Fields map[int]string `json:"fields"`
}

func NewISO8583Message(mti string, fields map[int]string) *ISO8583Message {
	return &ISO8583Message{
		MTI:    mti,
		Fields: fields,
	}
}
