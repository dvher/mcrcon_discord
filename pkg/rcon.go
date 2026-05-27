package rcon

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
)

var padding [2]byte = [2]byte{0, 0}

const (
	errPacketId    = -1
	maxPayloadSize = 1460
	authType       = 3
	commandType    = 2
)

type Payload struct {
	Length    int32
	RequestID int32
	Type      int32
	Payload   []byte
	Padding   []byte
}

type Connection struct {
	ip       string
	port     int
	password string
	conn     net.Conn
}

func NewConnection(ip string, port int, password string) *Connection {
	if port == 0 {
		port = 25575
	}

	conn := &Connection{
		ip:       ip,
		port:     port,
		password: password,
	}

	conn.Connect()

	return conn
}

func (c *Connection) SendCommand(command string) error {
	payload := Payload{
		Type:      commandType,
		Payload:   []byte(command),
		Padding:   padding[:],
		RequestID: rand.Int31(),
	}

	payload.Length = payload.Size()

	response, err := c.send(payload)

	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Response", string(response.Payload))

	return nil
}

func (c *Connection) send(p Payload) (*Payload, error) {

	data, err := binary.Append(nil, binary.LittleEndian, p.Length)

	if err != nil {
		return nil, err
	}

	data, err = binary.Append(data, binary.LittleEndian, p.RequestID)

	if err != nil {
		return nil, err
	}

	data, err = binary.Append(data, binary.LittleEndian, p.Type)

	if err != nil {
		return nil, err
	}

	data, err = binary.Append(data, binary.LittleEndian, p.Payload)

	if err != nil {
		return nil, err
	}

	data, err = binary.Append(data, binary.LittleEndian, padding)

	if err != nil {
		return nil, err
	}

	log.Println(data)
	log.Println(p.Length)

	_, err = c.conn.Write(data)

	if err != nil {
		return nil, err
	}

	return c.read()
}

func (c *Connection) read() (*Payload, error) {
	payload := new(Payload)

	err := binary.Read(c.conn, binary.LittleEndian, &payload.Length)

	if err != nil {
		log.Println("Error reading packet size")
		return payload, err
	}

	err = binary.Read(c.conn, binary.LittleEndian, &payload.RequestID)

	if err != nil {
		log.Println("Error reading packet id")
		return payload, err
	}

	err = binary.Read(c.conn, binary.LittleEndian, &payload.Type)

	if err != nil {
		log.Println("Error reading packet type")
		return payload, err
	}

	packetLength := payload.Length - 8

	body := make([]byte, packetLength)

	_, err = io.ReadFull(c.conn, body)

	if err != nil {
		log.Println("Error reading packet body")
		return payload, err
	}

	payload.Payload = body

	err = binary.Read(c.conn, binary.LittleEndian, &payload.Padding)

	if err != nil {
		log.Println("Error reading padding")
		return payload, err
	}

	return payload, err

}

func (c *Connection) Connect() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.ip, c.port))

	if err != nil {
		log.Fatalln(err)
	}

	c.conn = conn

	payload := Payload{
		RequestID: rand.Int31(),
		Type:      authType,
		Payload:   []byte(c.password),
		Padding:   padding[:],
	}

	payload.Length = payload.Size()

	response, err := c.send(payload)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(response.Payload))
}

func (c *Connection) Close() {

	if c.conn == nil {
		return
	}

	err := c.conn.Close()

	if err != nil {
		log.Fatalln(err)
	}
}

func (p Payload) Size() int32 {
	return int32(len(p.Payload)) + int32(len(p.Padding)) + int32(8)
}
