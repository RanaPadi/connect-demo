package config

import (
	"log/slog"
	"os"
)

/* Configuration file */

// Creates a structured logger (slog is Go's built-in structured logging package). Logs are written to os.Stderr (standard error stream).
var Logger = slog.New(slog.NewTextHandler(os.Stderr, nil))

// Defines a default directory path (./outp/) where the program might save files (e.g., logs, generated data, debug outputs).
var OutputPath = "./outp/"

// A boolean flag to toggle debug-specific behavior
var DebuggingMode = true
