/*******************************************************************************
*
* Copyright 2015-2016 Morgan Auchede <morgan.auchede@gmail.com>
* Copyright 2017 SAP SE
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You should have received a copy of the License along with this
* program. If not, you may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
*******************************************************************************/

package main

import (
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

var hostname = findHostname()

func findHostname() string {
	name, err := os.Hostname()
	if err != nil {
		log.Fatal("cannot determine hostname: " + err.Error())
		return "unknown"
	}
	return name
}

var (
	facilities = []string{
		"kern", "user", "mail", "daemon", "auth", "syslog", "lpr", "news", //0..7
		"uucp", "cron", "authpriv", "ftp", "ntp", "security", "console", "mark", //8..15
		"local0", "local1", "local2", "local3", "local4", "local5", "local6", "local7", //16..23
	}
	severities = []string{
		"emerg", "alert", "crit", "err", "warning", "notice", "info", "debug", //0..7
	}
)

func indexInto(list []string, idx int) string {
	if idx < 0 || idx >= len(list) {
		return "unknown"
	}
	return list[idx]
}

func listen(connection net.Conn) {
	var buffer [bufferSize]byte

	for {
		size, err := connection.Read(buffer[:])
		if err != nil {
			log.Fatal("Read error:", err)
		}
		if size > 0 {
			readData(buffer[0:size])
		}
	}
}

func readData(data []byte) {
	facility := "unknown.unknown"
	message := string(data)

	endOfCode := strings.Index(message, ">")
	if endOfCode != -1 && endOfCode < 5 {
		code, err := strconv.Atoi(string(data[1:endOfCode]))
		if err == nil {
			facility = fmt.Sprintf("%s.%s", indexInto(facilities, code>>3), indexInto(severities, code&0x07))
		}

		message = string(data[endOfCode+1:])
	}

	now := time.Now().UTC()
	fmt.Printf("%s.%06d %s %s %s\n",
		now.Format("Jan 2 15:04:05"),
		now.Nanosecond()/1000, //microsecond precision
		hostname,
		facility,
		strings.TrimSuffix(message, "\n"),
	)
}

func main() {
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

	listen(connection)
}
