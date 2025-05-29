package output

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"
)

// captureOutput captures stdout and stderr for testing
func captureOutput(fn func()) (stdout, stderr string) {
	// Capture stdout
	oldStdout := os.Stdout
	rOut, wOut, _ := os.Pipe()
	os.Stdout = wOut
	
	// Capture stderr
	oldStderr := os.Stderr
	rErr, wErr, _ := os.Pipe()
	os.Stderr = wErr

	// Execute function
	fn()

	// Close writers and restore
	wOut.Close()
	wErr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	// Read captured output
	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutBuf.ReadFrom(rOut)
	stderrBuf.ReadFrom(rErr)

	return stdoutBuf.String(), stderrBuf.String()
}

func TestMsgType_Constants(t *testing.T) {
	// Test that the constants are defined with expected values
	if MsgOutput != 0 {
		t.Errorf("Expected MsgOutput to be 0, got %d", MsgOutput)
	}
	if MsgSuccess != 1 {
		t.Errorf("Expected MsgSuccess to be 1, got %d", MsgSuccess)
	}
	if MsgError != 2 {
		t.Errorf("Expected MsgError to be 2, got %d", MsgError)
	}
	if MsgWarning != 3 {
		t.Errorf("Expected MsgWarning to be 3, got %d", MsgWarning)
	}
}

func TestPrintMsg_MsgOutput(t *testing.T) {
	message := "Test output message"
	
	stdout, stderr := captureOutput(func() {
		PrintMsg(message, MsgOutput)
	})

	if !strings.Contains(stdout, message) {
		t.Errorf("Expected stdout to contain '%s', got: %s", message, stdout)
	}
	
	if stderr != "" {
		t.Errorf("Expected no stderr output, got: %s", stderr)
	}
	
	// MsgOutput should not have any prefix
	if strings.Contains(stdout, "✔") || strings.Contains(stdout, "✖") || strings.Contains(stdout, "⚠") {
		t.Error("MsgOutput should not have any prefix symbols")
	}
}

func TestPrintMsg_MsgSuccess(t *testing.T) {
	message := "Test success message"
	
	stdout, stderr := captureOutput(func() {
		PrintMsg(message, MsgSuccess)
	})

	if !strings.Contains(stdout, message) {
		t.Errorf("Expected stdout to contain '%s', got: %s", message, stdout)
	}
	
	if !strings.Contains(stdout, "✔") {
		t.Error("Success message should contain check mark symbol")
	}
	
	if stderr != "" {
		t.Errorf("Expected no stderr output, got: %s", stderr)
	}
}

func TestPrintMsg_MsgError(t *testing.T) {
	message := "Test error message"
	
	stdout, stderr := captureOutput(func() {
		PrintMsg(message, MsgError)
	})

	if !strings.Contains(stderr, message) {
		t.Errorf("Expected stderr to contain '%s', got: %s", message, stderr)
	}
	
	if !strings.Contains(stderr, "✖") {
		t.Error("Error message should contain X mark symbol")
	}
	
	if stdout != "" {
		t.Errorf("Expected no stdout output, got: %s", stdout)
	}
}

func TestPrintMsg_MsgWarning(t *testing.T) {
	message := "Test warning message"
	
	stdout, stderr := captureOutput(func() {
		PrintMsg(message, MsgWarning)
	})

	if !strings.Contains(stdout, message) {
		t.Errorf("Expected stdout to contain '%s', got: %s", message, stdout)
	}
	
	if !strings.Contains(stdout, "⚠") {
		t.Error("Warning message should contain warning symbol")
	}
	
	if stderr != "" {
		t.Errorf("Expected no stderr output, got: %s", stderr)
	}
}

func TestPrintMsg_EmptyMessage(t *testing.T) {
	stdout, stderr := captureOutput(func() {
		PrintMsg("", MsgOutput)
	})

	// Should still print a newline
	if stdout == "" {
		t.Error("Expected at least a newline character in stdout")
	}
	
	if stderr != "" {
		t.Errorf("Expected no stderr output, got: %s", stderr)
	}
}

func TestPrintMsg_MultilineMessage(t *testing.T) {
	message := "Line 1\nLine 2\nLine 3"
	
	stdout, _ := captureOutput(func() {
		PrintMsg(message, MsgOutput)
	})

	if !strings.Contains(stdout, "Line 1") {
		t.Error("Should contain Line 1")
	}
	if !strings.Contains(stdout, "Line 2") {
		t.Error("Should contain Line 2")
	}
	if !strings.Contains(stdout, "Line 3") {
		t.Error("Should contain Line 3")
	}
}

func TestPrintTable_BasicUsage(t *testing.T) {
	headers := []string{"Name", "Age", "City"}
	rows := [][]string{
		{"Alice", "30", "New York"},
		{"Bob", "25", "San Francisco"},
	}
	
	stdout, stderr := captureOutput(func() {
		PrintTable(headers, rows)
	})

	// Check headers are present
	if !strings.Contains(stdout, "Name") {
		t.Error("Should contain header 'Name'")
	}
	if !strings.Contains(stdout, "Age") {
		t.Error("Should contain header 'Age'")
	}
	if !strings.Contains(stdout, "City") {
		t.Error("Should contain header 'City'")
	}

	// Check rows are present
	if !strings.Contains(stdout, "Alice") {
		t.Error("Should contain 'Alice'")
	}
	if !strings.Contains(stdout, "Bob") {
		t.Error("Should contain 'Bob'")
	}
	if !strings.Contains(stdout, "New York") {
		t.Error("Should contain 'New York'")
	}
	if !strings.Contains(stdout, "San Francisco") {
		t.Error("Should contain 'San Francisco'")
	}

	if stderr != "" {
		t.Errorf("Expected no stderr output, got: %s", stderr)
	}
}

func TestPrintTable_EmptyHeaders(t *testing.T) {
	headers := []string{}
	rows := [][]string{
		{"Alice", "30"},
	}
	
	stdout, _ := captureOutput(func() {
		PrintTable(headers, rows)
	})

	// Should still print the row data
	if !strings.Contains(stdout, "Alice") {
		t.Error("Should contain row data even with empty headers")
	}
}

func TestPrintTable_EmptyRows(t *testing.T) {
	headers := []string{"Name", "Age"}
	rows := [][]string{}
	
	stdout, _ := captureOutput(func() {
		PrintTable(headers, rows)
	})

	// Should print headers
	if !strings.Contains(stdout, "Name") {
		t.Error("Should contain headers even with empty rows")
	}
	if !strings.Contains(stdout, "Age") {
		t.Error("Should contain headers even with empty rows")
	}
}

func TestPrintTable_SingleColumn(t *testing.T) {
	headers := []string{"Names"}
	rows := [][]string{
		{"Alice"},
		{"Bob"},
		{"Charlie"},
	}
	
	stdout, _ := captureOutput(func() {
		PrintTable(headers, rows)
	})

	if !strings.Contains(stdout, "Names") {
		t.Error("Should contain header")
	}
	if !strings.Contains(stdout, "Alice") {
		t.Error("Should contain Alice")
	}
	if !strings.Contains(stdout, "Bob") {
		t.Error("Should contain Bob")
	}
	if !strings.Contains(stdout, "Charlie") {
		t.Error("Should contain Charlie")
	}
}

func TestPrintTable_MismatchedColumns(t *testing.T) {
	headers := []string{"Name", "Age", "City"}
	rows := [][]string{
		{"Alice", "30"}, // Missing city
		{"Bob", "25", "San Francisco", "Extra"}, // Extra column
	}
	
	stdout, _ := captureOutput(func() {
		PrintTable(headers, rows)
	})

	// Should still print what's available
	if !strings.Contains(stdout, "Alice") {
		t.Error("Should contain Alice")
	}
	if !strings.Contains(stdout, "Bob") {
		t.Error("Should contain Bob")
	}
}

func TestPrintTable_SpecialCharacters(t *testing.T) {
	headers := []string{"Symbol", "Unicode"}
	rows := [][]string{
		{"✓", "Check"},
		{"✗", "Cross"},
		{"⚠", "Warning"},
	}
	
	stdout, _ := captureOutput(func() {
		PrintTable(headers, rows)
	})

	if !strings.Contains(stdout, "✓") {
		t.Error("Should handle Unicode check symbol")
	}
	if !strings.Contains(stdout, "✗") {
		t.Error("Should handle Unicode cross symbol")
	}
	if !strings.Contains(stdout, "⚠") {
		t.Error("Should handle Unicode warning symbol")
	}
}

func TestLogError(t *testing.T) {
	testErr := errors.New("test error message")
	
	_, stderr := captureOutput(func() {
		LogError(testErr)
	})

	if !strings.Contains(stderr, "test error message") {
		t.Errorf("Expected stderr to contain error message, got: %s", stderr)
	}
	
	if !strings.Contains(stderr, "Error:") {
		t.Error("Expected stderr to contain 'Error:' prefix")
	}
}

func TestLogError_NilError(t *testing.T) {
	_, stderr := captureOutput(func() {
		LogError(nil)
	})

	if !strings.Contains(stderr, "Error:") {
		t.Error("Expected stderr to contain 'Error:' prefix even for nil")
	}
	
	// Should handle nil gracefully (Go's error interface handles this)
	if !strings.Contains(stderr, "<nil>") {
		t.Error("Expected stderr to show '<nil>' for nil error")
	}
}

func TestLogError_ComplexError(t *testing.T) {
	testErr := errors.New("complex error with\nnewlines and\ttabs")
	
	_, stderr := captureOutput(func() {
		LogError(testErr)
	})

	if !strings.Contains(stderr, "complex error") {
		t.Error("Should handle complex error messages")
	}
}

// Test integration between different functions
func TestIntegration_AllMessageTypes(t *testing.T) {
	messages := []struct {
		text    string
		msgType MsgType
	}{
		{"Regular output", MsgOutput},
		{"Success message", MsgSuccess},
		{"Error message", MsgError},
		{"Warning message", MsgWarning},
	}

	for _, msg := range messages {
		t.Run(string(rune(msg.msgType)), func(t *testing.T) {
			stdout, stderr := captureOutput(func() {
				PrintMsg(msg.text, msg.msgType)
			})

			switch msg.msgType {
			case MsgError:
				if !strings.Contains(stderr, msg.text) {
					t.Errorf("Error message should appear in stderr")
				}
			default:
				if !strings.Contains(stdout, msg.text) {
					t.Errorf("Non-error messages should appear in stdout")
				}
			}
		})
	}
}

func TestIntegration_TableAndMessages(t *testing.T) {
	stdout, _ := captureOutput(func() {
		PrintMsg("Table follows:", MsgOutput)
		PrintTable([]string{"Col1", "Col2"}, [][]string{{"A", "B"}})
		PrintMsg("Table complete", MsgSuccess)
	})

	if !strings.Contains(stdout, "Table follows:") {
		t.Error("Should contain initial message")
	}
	if !strings.Contains(stdout, "Col1") {
		t.Error("Should contain table headers")
	}
	if !strings.Contains(stdout, "Table complete") {
		t.Error("Should contain final message")
	}
	if !strings.Contains(stdout, "✔") {
		t.Error("Should contain success symbol")
	}
}