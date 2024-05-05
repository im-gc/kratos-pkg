package workCenter

import (
	"fmt"

	"bufio"
	"os"
	"strings"
)

func init() {
	addr := loadLocalIP()
	if len(addr) <= 0 {
		fmt.Println("load local hostname failed")
	}

	localEndPoint = addr
}

func loadLocalIP() string {

	file := HOSTNAME_FILE

	fd, err := os.Open(file)
	if err != nil {
		return ""
	}

	defer fd.Close()
	buff := bufio.NewReader(fd)

	line, err := buff.ReadString('\n')
	if err != nil {
		return ""
	}

	if strings.Contains(line, "hostname") {
		names := strings.Split(line, "=")
		if len(names) == 2 {
			return strings.TrimSpace(names[1])
		}
	}

	return ""
}
