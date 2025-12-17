package common

import (
	"log/slog"
	"os"
	"os/exec"
	"testing"

	"github.com/soderluk/nirimgr/models"
	"github.com/stretchr/testify/assert"
)

func TestRepr(t *testing.T) {
	config := models.Config{}
	r := Repr(config)
	assert.Equal(t, "Config", r)
	r = Repr(nil)
	assert.Equal(t, "", r)
	window := &models.Window{ID: 1}
	r = Repr(window)
	assert.Equal(t, "Window", r)
}

func TestLogLevel(t *testing.T) {
	logLevels := map[string]slog.Level{
		"DEBUG": slog.LevelDebug,
		"INFO":  slog.LevelInfo,
		"WARN":  slog.LevelWarn,
		"ERROR": slog.LevelError,
		"FOO":   slog.LevelDebug,
	}
	for logLevel, expected := range logLevels {
		l := parseLogLevel(logLevel)
		assert.Equal(t, expected, l)
	}
}

func TestValidateCommand_ValidCommands(t *testing.T) {
	validCommands := []string{
		"echo hello",
		"ls -la",
		"grep pattern file.txt",
		"cat /tmp/test.txt",
		"pwd",
		"whoami",
		"date",
		"ps aux | grep process",
		"find . -name '*.go'",
	}

	for _, cmd := range validCommands {
		t.Run(cmd, func(t *testing.T) {
			err := validateCommand(cmd)
			assert.NoError(t, err, "command should be valid: %s", cmd)
		})
	}
}

func TestValidateCommand_DangerousPatterns(t *testing.T) {
	dangerousCommands := []struct {
		command string
		pattern string
	}{
		{"rm -rf /", "rm -rf"},
		{"rm -fr /home", "rm -fr"},
		{"rm / -rf", "rm / -rf"},
		{"cat file > /dev/sda", "> /dev/sda"},
		{"dd if=/dev/zero of=/dev/sda", "dd if="},
		{"mkfs.ext4 /dev/sda1", "mkfs"},
		{"shred -vfz -n 10 /dev/sda", "shred"},
		{":(){ :|:& };:", ":(){ :|:& };:"},
		{"fork() bomb", "fork()"},
		{"chmod -R 777 /etc", "chmod -R 777"},
		{"chown -R root:root /", "chown -R"},
		{"sudo rm -rf /", "dangerous"},
		{"sudo dd if=/dev/zero of=/dev/sda", "dangerous"},
		{"sudo mkfs.ext4 /dev/sda1", "dangerous"},
		{"echo test > /etc/passwd", "> /etc/"},
		{"cat data > /boot/grub/grub.cfg", "> /boot/"},
		{"echo 1 > /sys/kernel/debug/test", "> /sys/"},
		{"format c: /q", "format c:"},
		{"del /f /s /q C:\\*", "del /f /s /q"},
	}

	for _, tc := range dangerousCommands {
		t.Run(tc.command, func(t *testing.T) {
			err := validateCommand(tc.command)
			assert.Error(t, err, "command should be rejected: %s", tc.command)
			assert.Contains(t, err.Error(), "dangerous", "error should mention that the command is dangerous")
		})
	}
}

func TestValidateCommand_PrivilegeEscalation(t *testing.T) {
	privilegeEscalationCommands := []string{
		"sudo ls",
		"  sudo echo test",
		"su root",
		"  su - ",
	}

	for _, cmd := range privilegeEscalationCommands {
		t.Run(cmd, func(t *testing.T) {
			err := validateCommand(cmd)
			assert.Error(t, err, "command should be rejected: %s", cmd)
			assert.Contains(t, err.Error(), "privilege escalation", "error should mention privilege escalation")
		})
	}
}

// TestRunCommand_ValidationFailure tests that RunCommand properly rejects dangerous commands
func TestRunCommand_ValidationFailure(t *testing.T) {
	dangerousCommands := []string{
		"rm -rf /",
		"sudo rm -rf /tmp",
		"dd if=/dev/zero of=/dev/sda",
	}

	for _, cmd := range dangerousCommands {
		t.Run(cmd, func(t *testing.T) {
			output, err := RunCommand(cmd)
			assert.Error(t, err, "RunCommand should reject dangerous command: %s", cmd)
			assert.Nil(t, output, "output should be nil for rejected commands")
		})
	}
}

// TestRunCommand_Success tests that valid commands are accepted
// We use a helper process pattern to mock exec.Command
func TestRunCommand_Success(t *testing.T) {
	if os.Getenv("GO_TEST_HELPER_PROCESS") == "1" {
		// This is the helper process
		_, err := os.Stdout.Write([]byte("mocked output"))
		assert.Nil(t, err, "err should be nil for os.Stdout.Write")
		os.Exit(0)
		return
	}

	// Save original execCommand
	originalExecCommand := execCommand
	defer func() { execCommand = originalExecCommand }()

	// Mock execCommand
	execCommand = mockExecCommand

	output, err := RunCommand("echo hello")
	assert.NoError(t, err)
	assert.Equal(t, []byte("mocked output"), output)
}

// TestRunCommand_VerifyCommandArgs tests that the command is correctly passed to exec
func TestRunCommand_VerifyCommandArgs(t *testing.T) {
	// Save original execCommand
	originalExecCommand := execCommand
	defer func() { execCommand = originalExecCommand }()

	var capturedCommand string
	var capturedArgs []string

	// Mock execCommand to capture arguments
	execCommand = func(command string, args ...string) *exec.Cmd {
		capturedCommand = command
		capturedArgs = args
		// Return a mock command that does nothing
		return mockExecCommand(command, args...)
	}

	_, _ = RunCommand("echo hello world")

	assert.Equal(t, "sh", capturedCommand, "should use sh to execute commands")
	assert.Equal(t, []string{"-c", "echo hello world"}, capturedArgs, "should pass -c flag and the full command")
}

// mockExecCommand is a helper function to mock exec.Command
func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestRunCommand_Success", "--"}
	cs = append(cs, command)
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_TEST_HELPER_PROCESS=1"}
	return cmd
}
