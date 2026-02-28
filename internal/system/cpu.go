// ©AngelaMos | 2026
// cpu.go

package system

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func parseCPULine(line string) (idle, total uint64, err error) {
	fields := strings.Fields(line)
	if len(fields) < 5 {
		return 0, 0, fmt.Errorf("not enough fields in cpu line: %q", line)
	}

	var sum uint64
	for i := 1; i < len(fields); i++ {
		val, parseErr := strconv.ParseUint(fields[i], 10, 64)
		if parseErr != nil {
			return 0, 0, fmt.Errorf("failed to parse field %d: %w", i, parseErr)
		}
		sum += val
		if i == 4 {
			idle = val
		}
	}

	return idle, sum, nil
}

func readCPUStat() (idle, total uint64, err error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu ") {
			return parseCPULine(line)
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}

	return 0, 0, fmt.Errorf("cpu line not found in /proc/stat")
}

func GetCPUUsage() (float64, error) {
	idle1, total1, err := readCPUStat()
	if err != nil {
		return 0, err
	}

	time.Sleep(200 * time.Millisecond)

	idle2, total2, err := readCPUStat()
	if err != nil {
		return 0, err
	}

	idleDelta := float64(idle2 - idle1)
	totalDelta := float64(total2 - total1)

	if totalDelta == 0 {
		return 0, nil
	}

	return (1 - idleDelta/totalDelta) * 100, nil
}
