package service

import (
	"fmt"
	"os/exec"
)

type AuditService struct{}

func NewAuditService() *AuditService {
	return &AuditService{}
}

func (s *AuditService) RunLynisAudit() (string, error) {
	// Simple wrapper for Lynis (needs to be installed)
	cmd := exec.Command("lynis", "audit", "system", "-Q")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("lynis audit failed: %w", err)
	}
	return string(output), nil
}
