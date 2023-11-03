package dkb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFirstLine(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		ledgerType string
		iban       string
	}{
		{
			name:       "valid",
			line:       `"Konto";"Girokonto DE12345678901234567890"`,
			ledgerType: "Girokonto",
			iban:       "DE12345678901234567890",
		}, {
			name:       "extra stuff but still valid",
			line:       `HEHEHEJKHJK"Konto";"Tagesgeldkonto DE12345678901234567890"\n`,
			ledgerType: "Tagesgeldkonto",
			iban:       "DE12345678901234567890",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &baseInfo{}
			err := info.parseFirstLine([]byte(tt.line))
			require.NoError(t, err)
			assert.Equal(t, tt.ledgerType, info.LedgerType)
			assert.Equal(t, tt.iban, info.IBAN)
		})
	}
}
