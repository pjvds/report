package context

import (
	"encoding/json"
	"testing"
)

func TestSerialization(t *testing.T) {
	context := Context{
		Channels: map[string]Channel{
			"team/chan": Channel{
				TeamName:    "team",
				ChannelName: "chan",
				WebhookUrl:  "http://example.com",
			},
		},
	}

	marshalled, err := json.Marshal(context)
	if err != nil {
		t.Fatal(err.Error())
	}

	destination := new(Context)
	err = json.Unmarshal(marshalled, destination)
	if err != nil {
		t.Fatal(err.Error())
	}
}
