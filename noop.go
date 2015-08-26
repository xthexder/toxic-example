package main

import (
	"io"

	"github.com/Shopify/toxiproxy/stream"
	"github.com/Shopify/toxiproxy/toxics"
)

type NoopToxic struct{}

func (t *NoopToxic) Pipe(stub *toxics.ToxicStub) {
	buf := make([]byte, 32*1024)
	writer := stream.NewChanWriter(stub.Output)
	reader := stream.NewChanReader(stub.Input)
	reader.SetInterrupt(stub.Interrupt)
	for {
		n, err := reader.Read(buf)
		if err == stream.ErrInterrupted {
			writer.Write(buf[:n])
			return
		} else if err == io.EOF {
			stub.Close()
			return
		}
		writer.Write(buf[:n])
	}
}

func init() {
	toxics.Register("noop", new(NoopToxic))
}
