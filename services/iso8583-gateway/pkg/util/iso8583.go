package util

import (
	"fmt"
	"iso8583-gateway/internal/domain"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/prefix"
	"go.uber.org/zap"
)

func ParseISO8583(data []byte) (*domain.ISO8583Message, error) {
	message := iso8583.NewMessage(napasSpec)

	err := message.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("fail to unpack ISO8583 message: %w", err)
	}
	mti, err := message.GetMTI()
	if err != nil {
		return nil, fmt.Errorf("fail to get mti: %w", err)
	}
	fields := make(map[int]string)
	for i, f := range message.GetFields() {
		s, err := f.String()
		if err != nil {
			zap.L().Warn("fail to convert field to string", zap.Int("field", i), zap.Error(err))
			continue
		}
		fields[i] = s
	}
	return domain.NewISO8583Message(mti, fields), nil
}

var napasSpec = &iso8583.MessageSpec{
	Name: "ISO 8583:1987 ASCII fields + Binary Bitmap",
	Fields: map[int]field.Field{
		0: field.NewString(&field.Spec{
			Length:      4,
			Description: "Message Type Indicator",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		1: field.NewBitmap(&field.Spec{
			Length:      8,
			Description: "Bitmap",
			Enc:         encoding.Binary,
			Pref:        prefix.Binary.Fixed,
		}),

		2:   field.NewString(&field.Spec{Length: 19, Description: "Primary Account Number", Enc: encoding.ASCII, Pref: prefix.ASCII.LL}),
		3:   field.NewString(&field.Spec{Length: 6, Description: "Processing Code", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		4:   field.NewString(&field.Spec{Length: 12, Description: "Amount, Transaction", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		5:   field.NewString(&field.Spec{Length: 12, Description: "Amount, Settlement", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		6:   field.NewString(&field.Spec{Length: 12, Description: "Amount, Cardholder Billing", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		7:   field.NewString(&field.Spec{Length: 10, Description: "Transmission Date & Time", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		9:   field.NewString(&field.Spec{Length: 8, Description: "Conversion Rate, Settlement", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		10:  field.NewString(&field.Spec{Length: 8, Description: "Conversion Rate, Cardholder Billing", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		11:  field.NewString(&field.Spec{Length: 6, Description: "System Trace Audit Number", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		12:  field.NewString(&field.Spec{Length: 6, Description: "Local Time", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		13:  field.NewString(&field.Spec{Length: 4, Description: "Local Date", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		14:  field.NewString(&field.Spec{Length: 4, Description: "Expiration Date", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		15:  field.NewString(&field.Spec{Length: 4, Description: "Settlement Date", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		18:  field.NewString(&field.Spec{Length: 4, Description: "Merchant Type", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		19:  field.NewString(&field.Spec{Length: 3, Description: "Acquiring Inst. Country Code", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		22:  field.NewString(&field.Spec{Length: 3, Description: "POS Entry Mode", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		23:  field.NewString(&field.Spec{Length: 3, Description: "Card Sequence Number", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		25:  field.NewString(&field.Spec{Length: 2, Description: "POS Condition Code", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		28:  field.NewString(&field.Spec{Length: 9, Description: "Amount, Fee", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}), // AMOUNT mapped as numeric 9
		32:  field.NewString(&field.Spec{Length: 11, Description: "Acquiring Inst ID", Enc: encoding.ASCII, Pref: prefix.ASCII.LL}),
		33:  field.NewString(&field.Spec{Length: 11, Description: "Forwarding Inst ID", Enc: encoding.ASCII, Pref: prefix.ASCII.LL}),
		35:  field.NewString(&field.Spec{Length: 37, Description: "Track 2 Data", Enc: encoding.ASCII, Pref: prefix.ASCII.LL}),
		36:  field.NewString(&field.Spec{Length: 104, Description: "Track 3 Data", Enc: encoding.ASCII, Pref: prefix.ASCII.LLL}),
		37:  field.NewString(&field.Spec{Length: 12, Description: "Retrieval Reference Number", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		38:  field.NewString(&field.Spec{Length: 6, Description: "Authorization ID", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		39:  field.NewString(&field.Spec{Length: 2, Description: "Response Code", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		41:  field.NewString(&field.Spec{Length: 8, Description: "Terminal ID", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		42:  field.NewString(&field.Spec{Length: 15, Description: "Merchant ID", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		43:  field.NewString(&field.Spec{Length: 40, Description: "Card Acceptor Name", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		45:  field.NewString(&field.Spec{Length: 76, Description: "Track 1 Data", Enc: encoding.ASCII, Pref: prefix.ASCII.LL}),
		48:  field.NewString(&field.Spec{Length: 999, Description: "Additional Data", Enc: encoding.ASCII, Pref: prefix.ASCII.LLL}),
		49:  field.NewString(&field.Spec{Length: 3, Description: "Currency Code, Txn", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		50:  field.NewString(&field.Spec{Length: 3, Description: "Currency Code, Settlement", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		51:  field.NewString(&field.Spec{Length: 3, Description: "Currency Code, Cardholder", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		52:  field.NewString(&field.Spec{Length: 16, Description: "PIN Data", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		54:  field.NewString(&field.Spec{Length: 120, Description: "Additional Amounts", Enc: encoding.ASCII, Pref: prefix.ASCII.LLL}),
		55:  field.NewString(&field.Spec{Length: 255, Description: "ICC Data", Enc: encoding.ASCII, Pref: prefix.ASCII.LLL}),
		60:  field.NewString(&field.Spec{Length: 999, Description: "Reserved Private", Enc: encoding.ASCII, Pref: prefix.ASCII.LLL}),
		62:  field.NewString(&field.Spec{Length: 99, Description: "Reserved Private", Enc: encoding.ASCII, Pref: prefix.ASCII.LL}),
		63:  field.NewString(&field.Spec{Length: 999, Description: "Reserved Private", Enc: encoding.ASCII, Pref: prefix.ASCII.LLL}),
		66:  field.NewString(&field.Spec{Length: 1, Description: "Settlement Code", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		67:  field.NewString(&field.Spec{Length: 2, Description: "Extended Payment Code", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		70:  field.NewString(&field.Spec{Length: 3, Description: "Network Mgmt Info Code", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		71:  field.NewString(&field.Spec{Length: 4, Description: "Message Number", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		72:  field.NewString(&field.Spec{Length: 4, Description: "Message Number Last", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		90:  field.NewString(&field.Spec{Length: 42, Description: "Original Data Elements", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		95:  field.NewString(&field.Spec{Length: 42, Description: "Replacement Amounts", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
		100: field.NewString(&field.Spec{Length: 11, Description: "Receiving Inst ID", Enc: encoding.ASCII, Pref: prefix.ASCII.LL}),
		102: field.NewString(&field.Spec{Length: 28, Description: "Account ID 1", Enc: encoding.ASCII, Pref: prefix.ASCII.LL}),
		103: field.NewString(&field.Spec{Length: 28, Description: "Account ID 2", Enc: encoding.ASCII, Pref: prefix.ASCII.LL}),
		104: field.NewString(&field.Spec{Length: 255, Description: "Transaction Description", Enc: encoding.ASCII, Pref: prefix.ASCII.LLL}),
		105: field.NewString(&field.Spec{Length: 999, Description: "Reserved Private", Enc: encoding.ASCII, Pref: prefix.ASCII.LLL}),
		110: field.NewString(&field.Spec{Length: 999, Description: "Reserved Private", Enc: encoding.ASCII, Pref: prefix.ASCII.LLL}),
		120: field.NewString(&field.Spec{Length: 999, Description: "Reserved Private", Enc: encoding.ASCII, Pref: prefix.ASCII.LLL}),
		121: field.NewString(&field.Spec{Length: 999, Description: "Reserved Private", Enc: encoding.ASCII, Pref: prefix.ASCII.LLL}),
		122: field.NewString(&field.Spec{Length: 999, Description: "Reserved Private", Enc: encoding.ASCII, Pref: prefix.ASCII.LLL}),
		123: field.NewString(&field.Spec{Length: 999, Description: "Reserved Private", Enc: encoding.ASCII, Pref: prefix.ASCII.LLL}),
		124: field.NewString(&field.Spec{Length: 999, Description: "Reserved Private", Enc: encoding.ASCII, Pref: prefix.ASCII.LLL}),
		125: field.NewString(&field.Spec{Length: 999, Description: "Reserved Private", Enc: encoding.ASCII, Pref: prefix.ASCII.LLL}),
		128: field.NewString(&field.Spec{Length: 64, Description: "Message Authentication Code", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed}),
	},
}
