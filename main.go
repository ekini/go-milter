/* Author: Eugene Dementiev <eugene@dementiev.eu>
 */
package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type State struct {
	macroses map[string][]string
}

func NewState() (state *State) {
	state = new(State)
	state.macroses = make(map[string][]string)
	return

}

type Handler interface {
	OptNeg(server_opts SMFIC_OPTIONS) (opts SMFIC_OPTIONS)
}

type DefaultHandler struct {
}

func (h *DefaultHandler) OptNeg(server_opts SMFIC_OPTIONS) (opts SMFIC_OPTIONS) {
	opts = server_opts
	opts.Version = 6
	opts.Actions = 0                                         // We aren't going to do anything
	opts.Protocol = (SMFIP_NOSEND_MASK | SMFIP_NOREPLY_MASK) // Skip all commands
	return
}

func Dispatch(state *State, handler Handler, command string, inbuf io.Reader, outbuf io.Writer) (close bool, err error) {
	switch command {
	case SMFIC_OPTNEG:
		{
			var command string = "O" // command name
			var o SMFIC_OPTIONS

			err = binary.Read(inbuf, binary.BigEndian, &o)
			if err != nil {
				return
			}
			log.Printf("Command: '%s', Version: %v, Actions: %#v, Protocol: %#v\n", command, o.Version, o.Actions, o.Protocol)
			opts := handler.OptNeg(o)
			log.Printf("Milter opts: Version: %v, Actions: %#v, Protocol: %#v\n", opts.Version, opts.Actions, opts.Protocol)
			err = binary.Write(outbuf, binary.BigEndian, uint32(binary.Size(opts)+1))
			_, err = outbuf.Write([]byte(command))
			if err != nil {
				return
			}
			err = binary.Write(outbuf, binary.BigEndian, opts)
			if err != nil {
				return
			}

			return
		}
	case SMFIC_MACRO:
		{
			var command string = "Macro"
			buf := new(bytes.Buffer)
			io.Copy(buf, inbuf)
			s := strings.Split(buf.String(), "\x00")
			log.Printf("Command: '%s', Payload: %+v\n", command, s)
			state.macroses[s[0]] = s[1:]
			return
		}
	case SMFIC_BODYEOB:
		{
			var length uint32 = 1
			err = binary.Write(outbuf, binary.BigEndian, length)
			if err != nil {
				return
			}
			_, err = outbuf.Write([]byte("a"))
			if err != nil {
				return
			}

			return
		}
	case SMFIC_QUIT:
		{
			return true, nil
		}
	}
	return true, fmt.Errorf("Command '%s' not recognized", command)
}

func handleConnection(handler Handler, conn net.Conn) {
	defer conn.Close()

	var err error
	var packet_length uint32
	var command byte
	state := NewState()

	log.Printf("Handling connection from '%s'\n", conn.RemoteAddr().String())
	var quit bool

	for !quit {
		//log.Printf("Connection state: %+v\n", *state)
		err = binary.Read(conn, binary.BigEndian, &packet_length)
		if err != nil {
			fmt.Printf("Can't read command length: %s\n", err)
			return
		}

		err = binary.Read(conn, binary.BigEndian, &command)
		if err != nil {
			fmt.Printf("Can't read: %s\n", err)
			return
		}

		inbuf := new(bytes.Buffer)
		outbuf := new(bytes.Buffer)
		_, err = io.CopyN(inbuf, conn, int64(packet_length-1))
		if err != nil {
			log.Printf("Can't copy from conn to inbuf: %s\n", err)
			return
		}
		quit, err = Dispatch(state, handler, string(command), inbuf, outbuf)
		if err != nil {
			log.Println("Dispatch failed:", err)
		}
		if outbuf.Len() > 0 {
			_, err = io.Copy(conn, outbuf)
			if err != nil {
				log.Println("Can't copy from outbuf to conn:", err)
			} else {
				//log.Printf("Copied %d bytes from outbuf to conn\n", n)
			}
		}
	}
	log.Println("Closing connection...")
}

func main() {
	handler := new(DefaultHandler)
	var ip = flag.Int("port", 6666, "Port to bind")
	var host = flag.String("host", "", "Host to bind")
	flag.Parse()

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *ip))
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Listening on %s:%d...\n", *host, *ip)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(handler, conn)
	}
}
