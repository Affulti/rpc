package codec

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
)

type JsonCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *json.Decoder //解码json
	enc  *json.Encoder //编码json
}

var _ Codec = (*JsonCodec)(nil)

// NewJsonCodec : type NewCodecFunc
func NewJsonCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn) //创建一个为conn提供缓冲的buf
	jsonCodec := &JsonCodec{
		conn: conn,
		buf:  buf,
		dec:  json.NewDecoder(conn),
		enc:  json.NewEncoder(buf), //将data写入buf(而不是conn)
	}
	return jsonCodec
}

func (j *JsonCodec) Close() error {
	return j.conn.Close()
}

func (j *JsonCodec) ReadHeader(h *Header) error {
	return j.dec.Decode(h)
}

func (j *JsonCodec) ReadBody(body interface{}) error {
	return j.dec.Decode(body)
}

func (j *JsonCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
		_ = j.buf.Flush()
		if err != nil {
			_ = j.Close()
		}
	}()
	if err = j.enc.Encode(h); err != nil {
		log.Panicln("rpc codec: json error encoding header :", err)
		return err
	}
	if err = j.enc.Encode(body); err != nil {
		log.Println("rpc codec: json error encoding body:", err)
		return err
	}
	return nil
}
