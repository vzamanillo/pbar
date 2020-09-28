package progress

/**
  Inspired by the https://github.com/PumpkinSeed/cage module
*/
import (
	"bufio"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/projectdiscovery/gologger"
)

const (
	fourMegas = 4 * 1024
	two       = 2
)

// CaptureData contains the standard input streams to capture data
type CaptureData struct {
	BackupStdout   *os.File
	WriterStdout   *os.File
	BackupStderr   *os.File
	WriterStderr   *os.File
	WaitFinishRead *sync.WaitGroup
}

// StartCapture starts capturing data
func StartCapture(writeLocker sync.Locker, stdout, stderr *strings.Builder) *CaptureData {
	rStdout, wStdout, errStdout := os.Pipe()
	if errStdout != nil {
		panic(errStdout)
	}

	rStderr, wStderr, errStderr := os.Pipe()
	if errStderr != nil {
		panic(errStderr)
	}

	c := &CaptureData{
		BackupStdout: os.Stdout,
		WriterStdout: wStdout,

		BackupStderr: os.Stderr,
		WriterStderr: wStderr,

		WaitFinishRead: &sync.WaitGroup{},
	}

	os.Stdout = c.WriterStdout
	os.Stderr = c.WriterStderr

	stdCopy := func(builder *strings.Builder, reader *os.File, waitGroup *sync.WaitGroup) {
		r := bufio.NewReader(reader)
		buf := make([]byte, 0, fourMegas)

		for {
			n, err := r.Read(buf[:cap(buf)])
			buf = buf[:n]

			if n == 0 {
				if err == nil {
					continue
				}

				if err == io.EOF {
					waitGroup.Done()
					break
				}

				waitGroup.Done()
				gologger.Fatalf("stdcapture error: %s", err)
			}

			if err != nil && err != io.EOF {
				waitGroup.Done()
				gologger.Fatalf("stdcapture error: %s", err)
			}

			writeLocker.Lock()
			builder.Write(buf)
			writeLocker.Unlock()
		}
	}

	c.WaitFinishRead.Add(two)

	go stdCopy(stdout, rStdout, c.WaitFinishRead)
	go stdCopy(stderr, rStderr, c.WaitFinishRead)

	return c
}

// StopCapture stops capturing data
func StopCapture(c *CaptureData) {
	_ = c.WriterStdout.Close()
	_ = c.WriterStderr.Close()

	c.WaitFinishRead.Wait()

	os.Stdout = c.BackupStdout
	os.Stderr = c.BackupStderr
}
