package main

import (
	"bufio"
	"bytes"
	"io"
	"net/http"

	"github.com/Shopify/toxiproxy/stream"
	"github.com/Shopify/toxiproxy/toxics"
)

type HttpToxic struct{}

func (t *HttpToxic) ModifyResponse(resp *http.Response) {
	resp.Header.Set("Location", "https://github.com/Shopify/toxiproxy")
}

func (t *HttpToxic) Pipe(stub *toxics.ToxicStub) {
	buffer := bytes.NewBuffer(make([]byte, 0, 32*1024))
	writer := stream.NewChanWriter(stub.Output)
	reader := stream.NewChanReader(stub.Input)
	reader.SetInterrupt(stub.Interrupt)
	for {
		tee := io.TeeReader(reader, buffer)
		resp, err := http.ReadResponse(bufio.NewReader(tee), nil)
		if err == stream.ErrInterrupted {
			buffer.WriteTo(writer)
			return
		} else if err == io.EOF {
			stub.Close()
			return
		}
		if err != nil {
			buffer.WriteTo(writer)
		} else {
			t.ModifyResponse(resp)
			resp.Write(writer)
		}
		buffer.Reset()
	}
}

func init() {
	toxics.Register("http", new(HttpToxic))
}
