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

type SmartExitCodeOutput struct {
	CommandLineParseError       bool
	DeviceOpenFailed            bool
	CommandFailed               bool
	DiskFailing                 bool
	PrefailAboveThreshold       bool
	PrefailAboveThresholdInPast bool
	DeviceErrorsLogged          bool
	RecentSelfTestErrors        bool
}

func SmartOutputFromExitCode(exitCode int) SmartExitCodeOutput {
	return SmartExitCodeOutput{
		CommandLineParseError:       exitCode&0x1 != 0,
		DeviceOpenFailed:            exitCode&0x2 != 0,
		CommandFailed:               exitCode&0x4 != 0,
		DiskFailing:                 exitCode&0x8 != 0,
		PrefailAboveThreshold:       exitCode&0x10 != 0,
		PrefailAboveThresholdInPast: exitCode&0x20 != 0,
		DeviceErrorsLogged:          exitCode&0x40 != 0,
		RecentSelfTestErrors:        exitCode&0x80 != 0,
	}
}

type smartctl struct {
	logger zerolog.Logger
}

func New() *smartctl {
	return &smartctl{
		logger: log.With().Str("component", "smartctl").Logger(),
	}
}

func (s *smartctl) exec(args ...string) ([]byte, SmartExitCodeOutput, error) {
	cmd := exec.Command("smartctl", args...)
	s.logger.Debug().
		Str("command", strings.Join(cmd.Args, " ")).
		Msg("executing command")
	out, err := cmd.CombinedOutput()
	if err != nil {
		s.logger.Error().Err(err).Msgf("failed to execute command. Command output: %s", string(out))
		// ignore error if command failed - normally indicates a SMART failure, so we pass back the exit code information
	}
	return out, SmartOutputFromExitCode(cmd.ProcessState.ExitCode()), nil
}

func (s *smartctl) ScanOpen() (*ScanOpenOutput, error) {
	out, code, err := s.exec("--scan-open", "-j")
	if err != nil {
		return nil, err
	}
	scanOpenOutput := &ScanOpenOutput{}
	if err := json.Unmarshal(out, scanOpenOutput); err != nil {
		return nil, err
	}
	scanOpenOutput.SmartExitCodeOutput = code
	return scanOpenOutput, nil
}

func (s smartctl) InfoAll(device string) (*InfoAllOutput, error) {
	out, code, err := s.exec("-iaj", device)
	if err != nil {
		return nil, err
	}
	infoAllOutput := &InfoAllOutput{}
	if err := json.Unmarshal(out, infoAllOutput); err != nil {
		return nil, err
	}
	infoAllOutput.SmartExitCodeOutput = code
	return infoAllOutput, nil
}
