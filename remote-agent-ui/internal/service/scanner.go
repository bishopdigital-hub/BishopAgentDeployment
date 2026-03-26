package service

import (
	"os/exec"
)

type ScanService struct{}

func NewScanService() *ScanService {
	return &ScanService{}
}

func (s *ScanService) ExecuteScan(target string) (string, error) {
	cmd := exec.Command("nmap", "-sV", "-oX", "-", target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
