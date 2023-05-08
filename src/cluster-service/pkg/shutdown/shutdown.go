package shutdown

import (
	"cluster/pkg/logging"
	"io"
	"os"
	"os/signal"
)

func Graceful(signals []os.Signal, closeItems ...io.Closer) {
	logger := logging.GetLogger()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, signals...)
	sig := <-sigc
	logger.Infof("Caught signal %s. Shutting down...", sig)

	for _, closer := range closeItems {
		if err := closer.Close(); err != nil {
			logger.Errorf("failed to close %v: %v", closer, err)
		}
	}
}
