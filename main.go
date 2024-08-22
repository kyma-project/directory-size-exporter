package main

import (
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/kyma-project/directory-size-exporter/internal/exporter"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	storagePath string
	metricName  string
	logFormat   string
	logLevel    string
	port        string
	interval    int
)

func main() {
	flag.StringVar(&logFormat, "log-format", "json", "Log format (json or text)")
	flag.StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error)")

	flag.StringVar(&storagePath, "storage-path", "", "Path to the observed data folder")
	flag.StringVar(&metricName, "metric-name", "", "Metric name used for exporting the folder size")
	flag.StringVar(&port, "port", "2021", "Port for exposing the metrics")
	flag.IntVar(&interval, "interval", 30, "Interval to calculate the metric ")

	flag.Parse()
	if err := validateFlags(); err != nil {
		panic(err)
	}

	logger := createLogger()

	exp := exporter.NewExporter(storagePath, metricName, logger)
	logger.Info("Exporter is initialized")

	exp.RecordMetrics(interval)
	logger.Info("Started recording metrics")

	http.Handle("/metrics", promhttp.Handler())
	server := &http.Server{
		Addr:              ":" + port,
		ReadHeaderTimeout: 1 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
	logger.Info("Listening on port '" + port + "'")
}

func validateFlags() error {
	if storagePath == "" {
		return errors.New("--storage-path flag is required")
	}
	if metricName == "" {
		return errors.New("--metric-name flag is required")
	}
	if logFormat != "json" && logFormat != "text" {
		return errors.New("--log-format flag should be either 'json' or 'text'")
	}
	if logLevel != "debug" && logLevel != "info" && logLevel != "warn" && logLevel != "error" {
		return errors.New("--log-level flag should be either 'debug', 'info', 'warn' or 'error'")
	}
	return nil
}

func createLogger() *slog.Logger {
	level := slog.LevelInfo
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	var handler slog.Handler
	handlerOpts := slog.HandlerOptions{
		Level: level,
	}
	if logFormat == "json" {
		handler = slog.NewJSONHandler(os.Stdout, &handlerOpts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, &handlerOpts)
	}

	return slog.New(handler)
}
