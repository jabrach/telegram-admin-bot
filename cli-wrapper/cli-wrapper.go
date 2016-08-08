package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/sthetz/tetanus/config"
	"io"
	"log"
	"os/exec"
	"strings"
)

var cmdArgs = []string{
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
	Self     *Self
	stdout   *bufio.Reader
	stdin    io.WriteCloser
	exeCh    chan string
}

func New() CLI {
	w := &wrapper{}
	w.loadSelf()

	return w
}

func (w *wrapper) loadSelf() {
	cmd := exec.Command(config.BinPath(), append(cmdArgs, "-e get_self", "-D")...)
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	lines := bytes.Split(out, []byte{'\n'})
	w.Self = &Self{}

	if err := json.Unmarshal(lines[0], w.Self); err != nil {
		panic(err)
	}
	log.Printf("Using account \"%s\" (ID %s)", w.Self.Username, w.Self.ID)
}

func (w *wrapper) Listen() {
	w.Cmd = exec.Command(config.BinPath(), cmdArgs...)
	w.setupPipes()

	go w.listeningRoutine()
	go w.execRoutine()

	log.Println("Listening...")
	w.Cmd.Start()
	w.Cmd.Wait()
}

func (w *wrapper) Exec(cmd string, args ...string) {
	scmd := strings.Join(append([]string{cmd}, args...), " ")
	w.exeCh <- scmd
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
	w.exeCh = make(chan string)
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

func (w *wrapper) execRoutine() {
	for cmd := range w.exeCh {
		bcmd := append([]byte(cmd), '\n')
		log.Printf("Sending command: %s", bcmd)

		_, err := w.stdin.Write(bcmd)
		if err != nil {
			panic(err)
		}
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
