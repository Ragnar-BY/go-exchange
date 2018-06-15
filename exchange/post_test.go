package exchange

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExchange2006_Post(t *testing.T) {

	tt := []struct {
		name          string
		contents      []byte
		response      string
		expectedError error
		expected      string
	}{
		{
			name:     "success",
			contents: nil,
			response: "ok",
			expected: "ok"},
		{
			name:          "FaultError",
			contents:      nil,
			response:      "<Fault><faultstring>some error</faultstring></Fault>",
			expectedError: errors.New("some error"),
			expected:      "",
		},
		{
			name:          "FaultFormatError",
			contents:      nil,
			response:      "<Fault><faultstring></Fault>",
			expectedError: errors.New("XML syntax error on line 1: element <faultstring> closed by </Fault>"),
			expected:      "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				io.WriteString(w, tc.response)
			}
			server := httptest.NewServer(http.HandlerFunc(handler))
			defer server.Close()
			ex := NewExchange("", "", server.URL)
			resp, err := ex.Post(tc.contents)
			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, resp)
			}
		})
	}
	t.Run("WrongRequestError", func(t *testing.T) {
		ex := NewExchange("", "", "wrongurl")
		_, err := ex.Post(nil)
		assert.Error(t, err)
	})

}
