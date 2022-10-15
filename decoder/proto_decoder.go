package decoder

import (
	"io-game-go/message"
)

type ProtoDecoder struct {
}

func (p ProtoDecoder) DecoderBytes(bytes []byte) message.Message {
	return nil
}
