package ibc

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/multierr"
)

func validPacket() Packet {
	return Packet{
		Sequence:         1,
		TimeoutHeight:    "100",
		TimeoutTimestamp: 0,
		SourcePort:       "transfer",
		SourceChannel:    "channel-0",
		DestPort:         "transfer",
		DestChannel:      "channel-1",
		Data:             []byte(`fake data`),
	}
}

func TestPacket_Validate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		packet := validPacket()

		require.NoError(t, packet.Validate())

		packet.TimeoutHeight = ""
		packet.TimeoutTimestamp = 1

		require.NoError(t, packet.Validate())
	})

	t.Run("invalid", func(t *testing.T) {
		var empty Packet
		merr := empty.Validate()

		require.Error(t, merr)
		require.Greater(t, len(multierr.Errors(merr)), 1)

		for _, tt := range []struct {
			Packet  Packet
			WantErr string
		}{
			{
				Packet{Sequence: 0},
				"packet sequence cannot be 0",
			},
			{
				Packet{Sequence: 1},
				"invalid packet source port:",
			},
			{
				Packet{Sequence: 1, SourcePort: "@"},
				"invalid packet source port:",
			},
			{
				Packet{Sequence: 1, SourcePort: "transfer"},
				"invalid packet source channel:",
			},
			{
				Packet{Sequence: 1, SourcePort: "transfer", SourceChannel: "@"},
				"invalid packet source channel:",
			},
			{
				Packet{Sequence: 1, SourcePort: "transfer", SourceChannel: "channel-0"},
				"invalid packet destination port:",
			},
			{
				Packet{Sequence: 1, SourcePort: "transfer", SourceChannel: "channel-0", DestPort: "@"},
				"invalid packet destination port:",
			},
			{
				Packet{Sequence: 1, SourcePort: "transfer", SourceChannel: "channel-0", DestPort: "transfer"},
				"invalid packet destination channel:",
			},
			{
				Packet{Sequence: 1, SourcePort: "transfer", SourceChannel: "channel-0", DestPort: "transfer", DestChannel: "@"},
				"invalid packet destination channel:",
			},
			{
				Packet{Sequence: 1, SourcePort: "transfer", SourceChannel: "channel-0", DestPort: "transfer", DestChannel: "channel-0"},
				"packet timeout height and packet timeout timestamp cannot both be 0",
			},
			{
				Packet{Sequence: 1, SourcePort: "transfer", SourceChannel: "channel-0", DestPort: "transfer", DestChannel: "channel-0", TimeoutHeight: "100"},
				"packet data bytes cannot be empty",
			},
		} {
			err := tt.Packet.Validate()
			require.Error(t, err, tt)
			require.Contains(t, err.Error(), tt.WantErr, tt)
		}
	})
}