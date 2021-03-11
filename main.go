package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	bufferSize = 65536
	socketPath = "/dev/log"
)

//Syslog contains the application state.
type Syslog struct {
	Hostname string
}

func (syslog Syslog) getFacility(code int) string {
	switch code >> 3 {
	case 0:
		return "kern"
	case 1:
		return "user"
	case 2:
		return "mail"
	case 3:
		return "daemon"
	case 4:
		return "auth"
	case 5:
		return "syslog"
	case 6:
		return "lpr"
	case 7:
		return "news"
	case 8:
		return "uucp"
	case 9:
		return "cron"
	case 10:
		return "authpriv"
	case 11:
		return "ftp"
	case 12:
		return "ntp"
	case 13:
		return "security"
	case 14:
		return "console"
	case 15:
		return "mark"
	case 16:
		return "local0"
	case 17:
		return "local1"
	case 18:
		return "local2"
	case 19:
		return "local3"
	case 20:
		return "local4"
	case 21:
		return "local5"
	case 22:
		return "local6"
	case 23:
		return "local7"
	default:
		return "unknown"
	}
}

func (syslog Syslog) getSeverity(code int) string {
	switch code & 0x07 {
	case 0:
		return "emerg"
	case 1:
		return "alert"
	case 2:
		return "crit"
	case 3:
		return "err"
	case 4:
		return "warning"
	case 5:
		return "notice"
	case 6:
		return "info"
	case 7:
		return "debug"
	default:
		return "unknown"
	}
}

func (syslog Syslog) listen(connection net.Conn) {
	reader := bufio.NewReader(connection)

	for {
		buffer := make([]byte, bufferSize)
		size, err := reader.Read(buffer)
		if err != nil {
			log.Fatal("Read error:", err)
		}

		go syslog.readData(buffer[0:size])
	}
}

func (syslog Syslog) report(facility, message string) {
	fmt.Printf("%s %s %s %s\n",
		time.Now().UTC().Format("Jan 2 15:04:05"),
		syslog.Hostname,
		facility,
		strings.TrimSuffix(message, "\n"),
	)
}

func (syslog Syslog) readData(data []byte) {
	facility := "unknown.unknown"
	message := string(data)

	endOfCode := strings.Index(message, ">")
	if -1 != endOfCode && 5 > endOfCode {
		code, err := strconv.Atoi(string(data[1:endOfCode]))
		if err == nil {
			facility = fmt.Sprintf("%s.%s", syslog.getFacility(code), syslog.getSeverity(code))
		}

		message = string(data[endOfCode+1:])
	}

	syslog.report(facility, message)
}

func (syslog Syslog) run() {
	if _, err := os.Stat(socketPath); err == nil {
		os.Remove(socketPath)
	}

	connection, err := net.ListenUnixgram("unixgram", &net.UnixAddr{Name: socketPath, Net: "unixgram"})
	if err != nil {
		log.Fatal("Listen error:", err)
	}

	if err := os.Chmod(socketPath, 0777); err != nil {
		log.Fatal("Impossible to change the socket permission.")
	}

	syslog.listen(connection)
}

func main() {
	var (
		syslog Syslog
		err    error
	)
	syslog.Hostname, err = os.Hostname()
	if err != nil {
		syslog.report("syslog.err", "syslog-stdout: cannot determine hostname: "+err.Error())
	}
	syslog.run()
}
