// ©AngelaMos | 2026
// scanner.go

package audit

import (
	"fmt"

	"github.com/CarterPerez-dev/yoshi-audit/internal/system"
)

type FindingType string

const (
	FindingZombie       FindingType = "zombie"
	FindingOrphan       FindingType = "orphan"
	FindingDaemon       FindingType = "daemon"
	FindingKernelThread FindingType = "kthread"
	FindingGPUShadow    FindingType = "gpushadow"
	FindingMemoryLeak   FindingType = "memleak"
	FindingUnknownSvc   FindingType = "unknown"
)

type Severity string

const (
	SeverityOK   Severity = "OK"
	SeverityWarn Severity = "WARN"
	SeverityCrit Severity = "CRIT"
)

type Finding struct {
	Type     FindingType
	Severity Severity
	PID      int
	Name     string
	Message  string
	Detail   string
}

type Scanner struct {
	knownGoods map[string]string
	baseline   *Baseline
}

func NewScanner(baseline *Baseline) *Scanner {
	return &Scanner{
		knownGoods: KnownGoodProcesses,
		baseline:   baseline,
	}
}

func (s *Scanner) ScanAll(
	procs []system.ProcessInfo,
	gpuProcs []system.GPUProcess,
) []Finding {
	var findings []Finding
	findings = append(findings, s.ScanZombies(procs)...)
	findings = append(findings, s.ScanOrphans(procs)...)
	findings = append(findings, s.ScanDaemons(procs)...)
	findings = append(findings, s.ScanKernelThreads(procs)...)
	findings = append(findings, s.ScanGPUShadows(procs, gpuProcs)...)
	return findings
}

func (s *Scanner) ScanZombies(procs []system.ProcessInfo) []Finding {
	var findings []Finding
	for _, p := range procs {
		if p.State == "Z" {
			findings = append(findings, Finding{
				Type:     FindingZombie,
				Severity: SeverityWarn,
				PID:      p.PID,
				Name:     p.Name,
				Message: fmt.Sprintf(
					"Zombie process (parent: PID %d)",
					p.PPID,
				),
				Detail: p.Cmdline,
			})
		}
	}
	return findings
}

func (s *Scanner) ScanOrphans(procs []system.ProcessInfo) []Finding {
	var findings []Finding
	for _, p := range procs {
		if p.PPID != 1 {
			continue
		}
		if s.isKnown(p.Name) {
			continue
		}
		findings = append(findings, Finding{
			Type:     FindingOrphan,
			Severity: SeverityWarn,
			PID:      p.PID,
			Name:     p.Name,
			Message:  "Orphaned process (re-parented to systemd)",
			Detail:   p.Cmdline,
		})
	}
	return findings
}

func (s *Scanner) ScanDaemons(procs []system.ProcessInfo) []Finding {
	var findings []Finding
	for _, p := range procs {
		if p.State != "S" {
			continue
		}
		if s.isKnown(p.Name) {
			continue
		}
		findings = append(findings, Finding{
			Type:     FindingDaemon,
			Severity: SeverityWarn,
			PID:      p.PID,
			Name:     p.Name,
			Message: fmt.Sprintf(
				"Unknown daemon: %s (RSS: %d MiB)",
				p.Cmdline,
				p.RSS/(1024*1024),
			),
			Detail: p.Cmdline,
		})
	}
	return findings
}

func (s *Scanner) ScanKernelThreads(procs []system.ProcessInfo) []Finding {
	count := 0
	for _, p := range procs {
		if p.PPID == 2 {
			count++
		}
	}
	if count == 0 {
		return nil
	}
	return []Finding{
		{
			Type:     FindingKernelThread,
			Severity: SeverityOK,
			PID:      0,
			Name:     "kthreads",
			Message:  fmt.Sprintf("%d kernel threads running", count),
		},
	}
}

func (s *Scanner) ScanGPUShadows(
	procs []system.ProcessInfo,
	gpuProcs []system.GPUProcess,
) []Finding {
	pidSet := make(map[int]bool)
	for _, p := range procs {
		pidSet[p.PID] = true
	}

	var findings []Finding
	for _, gp := range gpuProcs {
		if !pidSet[gp.PID] {
			findings = append(findings, Finding{
				Type:     FindingGPUShadow,
				Severity: SeverityWarn,
				PID:      gp.PID,
				Name:     fmt.Sprintf("gpu-pid-%d", gp.PID),
				Message:  "GPU process not visible in process list",
				Detail: fmt.Sprintf(
					"VRAM usage: %d MiB",
					gp.UsedVRAM/(1024*1024),
				),
			})
		}
	}
	return findings
}

func (s *Scanner) ScanMemoryLeaks(
	procs []system.ProcessInfo,
	history map[int][]uint64,
) []Finding {
	nameMap := make(map[int]string)
	for _, p := range procs {
		nameMap[p.PID] = p.Name
	}

	const threshold = 100 * 1024 * 1024

	var findings []Finding
	for pid, readings := range history {
		if len(readings) < 2 {
			continue
		}
		first := readings[0]
		last := readings[len(readings)-1]
		if last > first && (last-first) > threshold {
			grewMiB := (last - first) / (1024 * 1024)
			name := nameMap[pid]
			if name == "" {
				name = fmt.Sprintf("pid-%d", pid)
			}
			findings = append(findings, Finding{
				Type:     FindingMemoryLeak,
				Severity: SeverityWarn,
				PID:      pid,
				Name:     name,
				Message: fmt.Sprintf(
					"Process grew %d MiB in %d readings",
					grewMiB,
					len(readings),
				),
				Detail: fmt.Sprintf(
					"first: %d bytes, last: %d bytes",
					first,
					last,
				),
			})
		}
	}
	return findings
}

func (s *Scanner) isKnown(name string) bool {
	if _, ok := s.knownGoods[name]; ok {
		return true
	}
	if s.baseline != nil && s.baseline.IsKnown(name) {
		return true
	}
	return false
}
