/**
 * Copyright © 2023, Staufi Tech - Switzerland
 * All rights reserved.
 *
 *   ________________________   ___ _     ________________  _  ____
 *  / _____  _  ____________/  / __|_|   /_______________  | | ___/
 * ( (____ _| |_ _____ _   _ _| |__ _      | |_____  ____| |_|_
 *  \____ (_   _|____ | | | (_   __) |     | | ___ |/ ___)  _  \
 *  _____) )| |_/ ___ | |_| | | |  | |     | | ____( (___| | | |
 * (______/  \__)_____|____/  |_|  |_|     |_|_____)\____)_| |_|
 *
 *
 *  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 *  AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 *  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 *  ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 *  LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 *  CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 *  SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 *  INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 *  CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 *  ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 *  POSSIBILITY OF SUCH DAMAGE.
 */

package log

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

var (
	printOut string
)

func tPrintln(in ...any) (int, error) {
	printOut = fmt.Sprintln(in...)
	return len(printOut), nil
}

func init() {
	println = tPrintln
}

func TestFine(t *testing.T) {
	printOut = ""

	SetLogLevel("debug")
	if err := Fine("TestModule", "message: %s", "fine"); err != nil {
		t.Error(err)
	}
	if printOut != "" {
		t.Error("debug level prints fine messages. ", printOut)
	}

	SetLogLevel("fine")
	if err := Fine("TestModule", "message: %s", "fine"); err != nil {
		t.Error(err)
	}

	expectedOutput := fmt.Sprintf("[FINE]\t%s\tTESTMODULE - message: fine", GetLastTimestamp())

	if !strings.Contains(printOut, expectedOutput) {
		t.Errorf("Expected log output to contain '%s', but got: '%s'", expectedOutput, printOut)
	}

	fmt.Println(printOut)
}

func TestDebug(t *testing.T) {
	printOut = ""

	SetLogLevel("info")
	if err := Debug("TestModule", "message: %s", "testDebug"); err != nil {
		t.Error(err)
	}
	if printOut != "" {
		t.Error("info level prints debug messages. ", printOut)
	}

	SetLogLevel("debug")
	if err := Debug("TestModule", "message: %s", "testDebug"); err != nil {
		t.Error(err)
	}

	expectedOutput := fmt.Sprintf("[DEBUG]\t%s\tTESTMODULE - message: testDebug", GetLastTimestamp())

	if !strings.Contains(printOut, expectedOutput) {
		t.Errorf("Expected log output to contain '%s', but got: '%s'", expectedOutput, printOut)
	}

	fmt.Println(printOut)
}

func TestInfo(t *testing.T) {
	printOut = ""

	SetLogLevel("warn")
	if err := Info("TestModule", "message: %s", "info"); err != nil {
		t.Error(err)
	}

	if printOut != "" {
		t.Error("warn level prints info messages. ", printOut)
	}

	SetLogLevel("info")
	if err := Info("TestModule", "message: %s", "info"); err != nil {
		t.Error(err)
	}

	expectedOutput := fmt.Sprintf("[INFO]\t%s\tTESTMODULE - message: info", GetLastTimestamp())

	if !strings.Contains(printOut, expectedOutput) {
		t.Errorf("Expected log output to contain '%s', but got: '%s'", expectedOutput, printOut)
	}

	fmt.Println(printOut)
}

func TestWarn(t *testing.T) {
	printOut = ""

	SetLogLevel("error")
	if err := Warn("TestModule", "message: %s", "warn"); err != nil {
		t.Error(err)
	}

	if printOut != "" {
		t.Error("error level prints warn messages. ", printOut)
	}

	SetLogLevel("warn")
	if err := Warn("TestModule", "message: %s", "warn"); err != nil {
		t.Error(err)
	}

	expectedOutput := fmt.Sprintf("[WARN]\t%s\tTESTMODULE - message: warn", GetLastTimestamp())

	if !strings.Contains(printOut, expectedOutput) {
		t.Errorf("Expected log output to contain '%s', but got: '%s'", expectedOutput, printOut)
	}

	fmt.Println(printOut)
}

func TestError(t *testing.T) {
	printOut = ""

	SetLogLevel("fatal")
	if err := Error("TestModule", "message: %s", "error"); err != nil {
		t.Error(err)
	}

	if printOut != "" {
		t.Error("fatal level prints error messages. ", printOut)
	}

	SetLogLevel("error")
	if err := Error("TestModule", "message: %s", "error"); err != nil {
		t.Error(err)
	}

	expectedOutput := fmt.Sprintf("[ERROR]\t%s\tTESTMODULE - message: error", GetLastTimestamp())

	if !strings.Contains(printOut, expectedOutput) {
		t.Errorf("Expected log output to contain '%s', but got: '%s'", expectedOutput, printOut)
	}

	fmt.Println(printOut)
}

func TestFile(t *testing.T) {

	SetLogLevel("info")
	if err := ToFile("myLog.txt"); err != nil {
		t.Error(err)
	}

	if err := Fine("test", "fine message"); err != nil {
		t.Error(err)
	}
	if err := Debug("test", "debug message"); err != nil {
		t.Error(err)
	}
	if err := Info("test", "info message"); err != nil {
		t.Error(err)
	}
	if err := Warn("test", "warn message"); err != nil {
		t.Error(err)
	}
	if err := Error("test", "error message"); err != nil {
		t.Error(err)
	}

	if err := CloseFile(); err != nil {
		t.Error(err)
	}

	logs, err := os.ReadFile("myLog.txt")
	if err != nil {
		t.Error(err)
	}

	if strings.Contains(string(logs), "FINE") || strings.Contains(string(logs), "DEBUG") {
		t.Error("unexpected log level in file", string(logs))
	}

	if !strings.Contains(string(logs), "INFO") || !strings.Contains(string(logs), "WARN") ||
		!strings.Contains(string(logs), "ERROR") {
		t.Error("expected log levels not found in file", string(logs))
	}

	os.Remove("myLog.txt")
}
