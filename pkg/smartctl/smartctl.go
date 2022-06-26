package smartctl

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os/exec"
	"strings"
)

type SmartCtl interface {
	ScanOpen() (*ScanOpenOutput, error)
	InfoAll(device string) (*InfoAllOutput, error)
}

type smartctl struct {
	logger zerolog.Logger
}

func New() *smartctl {
	return &smartctl{
		logger: log.With().Str("component", "smartctl").Logger(),
	}
}

func (s *smartctl) exec(args ...string) ([]byte, error) {
	cmd := exec.Command("smartctl", args...)
	s.logger.Debug().
		Str("command", strings.Join(cmd.Args, " ")).
		Msg("executing command")
	return cmd.CombinedOutput()
}

func (s *smartctl) ScanOpen() (*ScanOpenOutput, error) {
	out, err := s.exec("--scan-open", "-j")
	if err != nil {
		return nil, err
	}
	scanOpenOutput := &ScanOpenOutput{}
	if err := json.Unmarshal(out, scanOpenOutput); err != nil {
		return nil, err
	}
	return scanOpenOutput, nil
}

func (s smartctl) InfoAll(device string) (*InfoAllOutput, error) {
	out, err := s.exec("-iaj", device)
	if err != nil {
		return nil, err
	}
	infoAllOutput := &InfoAllOutput{}
	if err := json.Unmarshal(out, infoAllOutput); err != nil {
		return nil, err
	}
	return infoAllOutput, nil
}
