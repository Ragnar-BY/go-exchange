package exchange

import (
	"testing"

	"github.com/Ragnar-BY/go-exchange/models"
)

func Test_parseAddMeetingResponse(t *testing.T) {

	ItemID := "SomeItemId"
	ChangeKey := "SomeChangeKey"
	tests := []struct {
		name     string
		response string
		want     *models.CalendarItem
		wantErr  bool
	}{
		{
			name: "success",
			response: `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
				<s:Header>
				</s:Header>
				<s:Body xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
				<m:CreateItemResponse xmlns:m="http://schemas.microsoft.com/exchange/services/2006/messages"
				xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types">
				<m:ResponseMessages>
				<m:CreateItemResponseMessage ResponseClass="Success">
				<m:ResponseCode>NoError</m:ResponseCode>
				<m:Items>
				<t:CalendarItem>
				<t:ItemId Id="` + ItemID + `" ChangeKey="` + ChangeKey + `" />
				</t:CalendarItem>
				</m:Items>
				</m:CreateItemResponseMessage>
				</m:ResponseMessages>
				</m:CreateItemResponse>
				</s:Body>
				</s:Envelope>"`,
			want: &models.CalendarItem{ItemID: struct {
				ID        string `xml:"Id,attr"`
				ChangeKey string `xml:"ChangeKey,attr"`
			}{ID: ItemID, ChangeKey: ChangeKey}},
			wantErr: false,
		},
		{
			name: "ResponseError",
			response: `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
				<s:Header>
				</s:Header>
				<s:Body xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
				<m:CreateItemResponse xmlns:m="http://schemas.microsoft.com/exchange/services/2006/messages"
				xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types">
				<m:ResponseMessages>
				<m:CreateItemResponseMessage ResponseClass="Error">
				<m:ResponseCode>Error</m:ResponseCode>
				<m:MessageText>SomeError</m:MessageText>				
				</m:CreateItemResponseMessage>
				</m:ResponseMessages>
				</m:CreateItemResponse>
				</s:Body>
				</s:Envelope>"`,
			wantErr: true},
		{
			name: "LostResponseMessageError",
			response: `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
				<s:Header>
				</s:Header>
				<s:Body xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
				<m:CreateItemResponse xmlns:m="http://schemas.microsoft.com/exchange/services/2006/messages"
				xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types">
				<m:ResponseMessages>				
				</m:ResponseMessages>
				</m:CreateItemResponse>
				</s:Body>
				</s:Envelope>"`,
			wantErr: true,
		},
		{
			name: "WrongResponseMessageError",
			response: `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
				<s:Header>
				</s:Header>
				<s:Body xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
				<m:CreateItemResponse xmlns:m="http://schemas.microsoft.com/exchange/services/2006/messages"
				xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types">
				<m:ResponseMessages>
				<m:CreateItemResponseMessage >
				<ErrorItem>		
				</m:CreateItemResponseMessage>
				</m:ResponseMessages>
				</m:CreateItemResponse>
				</s:Body>
				</s:Envelope>"`,
			wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAddMeetingResponse(tt.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAddMeetingResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if tt.want.ItemID.ID != got.ItemID.ID || tt.want.ItemID.ChangeKey != got.ItemID.ChangeKey {
					t.Errorf("parseAddMeetingResponse() got %+v, want %+v", got, tt.want)
				}
			}
		})
	}
}
