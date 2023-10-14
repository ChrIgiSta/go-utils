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
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type LogLevel int

const (
	LevelFine LogLevel = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
	LevelNone
)

var logLevelToString = map[int]string{
	int(LevelFine):  "FINE",
	int(LevelDebug): "DEBUG",
	int(LevelInfo):  "INFO",
	int(LevelWarn):  "WARN",
	int(LevelError): "ERROR",
	int(LevelFatal): "FATAL",
	int(LevelPanic): "PANIC",
	int(LevelNone):  "NONE",
}

var (
	logLevel      LogLevel = LevelInfo
	formatString           = "2006-01-02 15:04:05.000000"
	lastTimeStamp string

	println func(in ...any) (int, error) = fmt.Println
	file    *os.File
)

func init() {
	SetLogLevel(os.Getenv("LOG_LEVEL"))
}

func formatLog(level LogLevel, module, logText string) string {

	lastTimeStamp = time.Now().Format(formatString)

	return fmt.Sprintf("[%s]\t%s\t%s - %s",
		logLevelToString[int(level)],
		lastTimeStamp,
		strings.ToUpper(module),
		logText)
}

func logMessage(level LogLevel, module, logText string) error {
	var err error

	if level >= logLevel {
		_, err = println(formatLog(level, module, logText))
	}
	if level > LevelError {
		if logLevel == LevelFatal { // call fatal
			log.Fatal(logText)
		} else if logLevel == LevelPanic { // call panic
			log.Panic(logText)
		}
	}

	return err
}

func filePrinter(in ...any) (int, error) {
	var (
		cnt int   = 0
		err error = errors.New("no open file")
	)

	if file != nil {
		cnt, err = file.Write([]byte(fmt.Sprintln(in...)))
	}

	return cnt, err
}

func SetLogLevel(level string) {

	switch strings.ToLower(level) {
	case "fine":
		logLevel = LevelFine
	case "debug":
		logLevel = LevelDebug
	case "info":
		logLevel = LevelInfo
	case "warn":
		logLevel = LevelWarn
	case "error":
		logLevel = LevelError
	case "fatal":
		logLevel = LevelFatal
	case "panic":
		logLevel = LevelPanic
	case "none":
		logLevel = LevelNone
	default:
		logLevel = LevelInfo
	}
}

func OverwriteFormat(format string) {

	formatString = format
}

func GetLastTimestamp() string {

	return lastTimeStamp
}

func ToFile(toFile string) error {
	var err error

	file, err = os.Open(toFile)
	println = filePrinter

	return err
}

func CloseFile() error {
	return file.Close()
}

func Fine(module, logText string, args ...interface{}) error {

	return logMessage(LevelFine, module, fmt.Sprintf(logText, args...))
}

func Debug(module, logText string, args ...interface{}) error {

	return logMessage(LevelDebug, module, fmt.Sprintf(logText, args...))
}

func Info(module, logText string, args ...interface{}) error {

	return logMessage(LevelInfo, module, fmt.Sprintf(logText, args...))
}

func Warn(module, logText string, args ...interface{}) error {

	return logMessage(LevelWarn, module, fmt.Sprintf(logText, args...))
}

func Error(module, logText string, args ...interface{}) error {

	return logMessage(LevelError, module, fmt.Sprintf(logText, args...))
}
