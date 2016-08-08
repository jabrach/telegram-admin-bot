package cli

import (
	"bufio"
	"encoding/json"
	"github.com/sthetz/tetanus/config"
	"io"
	"log"
	"os/exec"
	"strings"
)

var args = []string{
	"--permanent-msg-ids",
	"--permanent-peer-ids",
	"--json",
	"-R",
}

type CLI interface {
	Listen()
	Exec(string, ...string)
	AddHandler(handlerFunc)
}

type wrapper struct {
	Cmd      *exec.Cmd
	handlers []handlerFunc
	stdout   *bufio.Reader
	stdin    io.WriteCloser
}

func New() CLI {
	return &wrapper{}
}

func (w *wrapper) Listen() {
	w.Cmd = exec.Command(config.BinPath(), args...)
	w.setupPipes()

	go w.listeningRoutine()

	log.Println("Listening...")
	w.Cmd.Start()
	w.Cmd.Wait()
}

func (w *wrapper) Exec(cmd string, args ...string) {
	scmd := append([]string{cmd}, args...)
	bcmd := []byte(strings.Join(scmd, " "))
	bcmd = append(bcmd, '\n')

	log.Printf("Sending command: %s", bcmd)

	_, err := w.stdin.Write(bcmd)
	if err != nil {
		panic(err)
	}
}

func (w *wrapper) AddHandler(handler handlerFunc) {
	if len(w.handlers) == 0 {
		w.handlers = []handlerFunc{}
	}
	w.handlers = append(w.handlers, handler)
}

func (w *wrapper) setupPipes() {
	out, err := w.Cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	w.stdout = bufio.NewReader(out)
	in, err := w.Cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	w.stdin = in
}

func (w *wrapper) listeningRoutine() {
	for {
		line, err := w.stdout.ReadString('\n')
		if err != nil {
			panic(err)
		}
		go w.handleMessage(line)
	}
}

func (w *wrapper) handleMessage(line string) {
	if line[0] != '{' {
		// log.Printf("???: %v\n", line)
		return
	}
	msg := &Message{}
	if err := json.Unmarshal([]byte(line), msg); err != nil {
		log.Fatalf("Error parsing JSON: %v\n", err.Error())
		return
	}
	msg.JSON = line

	for _, handler := range w.handlers {
		handler(msg, w)
	}
}
