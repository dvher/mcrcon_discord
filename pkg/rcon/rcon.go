package rcon

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	ip       = os.Getenv("SERVER_IP")
	port     = os.Getenv("SERVER_PORT")
	password = os.Getenv("SERVER_PASSWORD")
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
	port     string
	password string
	conn     net.Conn
}

func NewConnection() (*Connection, error) {
	if port == "" {
		port = "25575"
	}

	conn := &Connection{
		ip:       ip,
		port:     port,
		password: password,
	}

	err := conn.Connect()

	return conn, err
}

func (c *Connection) SendCommand(command string) (*Payload, error) {
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
		return nil, err
	}

	return response, nil
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

func (c *Connection) Connect() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", c.ip, c.port))

	if err != nil {
		log.Println(err)
		return err
	}

	c.conn = conn

	payload := Payload{
		RequestID: rand.Int31(),
		Type:      authType,
		Payload:   []byte(c.password),
		Padding:   padding[:],
	}

	payload.Length = payload.Size()

	_, err = c.send(payload)

	if err != nil {
		log.Println(err)
	}

	return err
}

func (c *Connection) Close() error {

	if c.conn == nil {
		return errors.New("Cannot close nil connection")
	}

	return c.conn.Close()
}

func (p Payload) Size() int32 {
	return int32(len(p.Payload)) + int32(len(p.Padding)) + int32(8)
}
