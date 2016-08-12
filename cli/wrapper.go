package cli

import (
	"bytes"
	"encoding/json"
	"github.com/jabrach/telegram-admin-bot/config"
	"io"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"
)

var socketAddr string

type Wrapper struct {
	Process  *exec.Cmd
	Handlers []msgHandler
	Self     *Self
	socket   net.Conn
	chInput  chan string
	chMsg    chan *Message
}

func New() *Wrapper {
	return &Wrapper{}
}

func Exec(msg interface{}, command string, args ...string) error {
	var (
		cmd      = buildCommand(command, args...)
		proc     = newProcess(false, &cmd)
		out, err = proc.CombinedOutput()
	)

	if err != nil {
		return err
	}
	lines := bytes.Split(out, []byte{'\n'})

	if err := json.Unmarshal(lines[0], msg); err != nil {
		return err
	}
	return nil
}

func (w *Wrapper) Start2() {
	c, err := net.Dial("unix", "/tmp/tg2.sock")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	_, err = c.Write([]byte("main_session\n"))
	if err != nil {
		panic(err)
	}
	reader(c)
}

func (w *Wrapper) Start() {
	if len(w.Handlers) == 0 {
		log.Println("No handlers enabled, nothing to do")
		return
	}

	w.Process = newProcess(true, nil)
	w.chInput = make(chan string)
	w.chMsg = make(chan *Message)
	w.Self = &Self{}
	w.Process.Start()
	// w.Process.Wait()

	// if err := Exec(w.Self, "get_self"); err != nil {
	// 	panic(err)
	// }
	time.Sleep(2 * time.Second)

	conn, err := net.Dial("unix", socketAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	if err != nil {
		panic(err)
	}
	w.socket = conn
	go w.readRoutine()
	go w.handleRoutine()
	go w.writeRoutine()

	w.Exec("main_session")
	w.Process.Wait()
}

func reader(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		println("reading...")
		n, err := r.Read(buf[:])

		if err != nil {
			panic(err)
		}
		println("Client got:", string(buf[0:n]))
	}
}

func (w *Wrapper) Exec(command string, args ...string) {
	w.chInput <- buildCommand(command, args...)
}

func (w *Wrapper) AddHandler(hanlder msgHandler) {
	if len(w.Handlers) == 0 {
		w.Handlers = []msgHandler{}
	}
	w.Handlers = append(w.Handlers, hanlder)
}

func (w *Wrapper) Stop() {
	w.Process.Process.Kill()
	// w.socket.Close()
}

func (w *Wrapper) writeRoutine() {
	for cmd := range w.chInput {
		bcmd := []byte(cmd)
		log.Printf("Sending command: %s", bcmd)
		if _, err := w.socket.Write(bcmd); err != nil {
			panic(err)
		}
	}
}

func (w *Wrapper) readRoutine() {
	buf := make([]byte, 1024)
	for {
		n, err := w.socket.Read(buf[:])
		if err != nil {
			if err.Error() == "EOF" {
				log.Println("EOF reached, exiting read routine.")
				return
			}
			panic(err)
		}
		w.chMsg <- parseMessage(string(buf[0:n]))
	}
}

func (w *Wrapper) handleRoutine() {
	for msg := range w.chMsg {
		if msg == nil {
			continue
		}
		for _, handler := range w.Handlers {
			go handler(msg, w)
		}
	}
}

func parseMessage(line string) *Message {
	msg := &Message{Data: MessageData{}}
	line = strings.TrimSpace(strings.Split(line, "\n")[1])

	if line == "" {
		return nil
	}

	if err := json.Unmarshal([]byte(line), &msg.Data); err != nil {
		log.Printf("Error parsing JSON: %v\n%s\n", err.Error(), line)
		return nil
	}
	msg.JSON = line
	msg.ID = msg.Data.ID
	return msg
}

func buildCommand(cmd string, args ...string) string {
	return strings.Join(append([]string{cmd}, args...), " ") + "\n"
}

func newProcess(main bool, command *string) *exec.Cmd {
	args := []string{
		"--permanent-msg-ids",
		"--permanent-peer-ids",
		"--json",
		"-R",
	}
	if main {
		socketAddr = "/tmp/tg.sock"
		args = append(args, "-S"+socketAddr)
	} else {
		if command != nil {
			args = append(args, "-D")
			args = append(args, "-e "+*command)
		}
	}

	log.Println(args)
	return exec.Command(config.BinPath(), args...)
}
