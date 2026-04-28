package media

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func probeDuration(ctx context.Context, path string) (time.Duration, error) {
	cmd := exec.CommandContext(
		ctx,
		"ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		path,
	)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("ffprobe duration: %w", err)
	}

	seconds, err := strconv.ParseFloat(strings.TrimSpace(out.String()), 64)
	if err != nil {
		return 0, err
	}

	return time.Duration(seconds * float64(time.Second)), nil
}

func extractVideoThumbnail(ctx context.Context, path string) (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}

	outDir := filepath.Join(cacheDir, "bare", "video-thumbs")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", err
	}

	thumbPath := filepath.Join(outDir, sanitizeMediaPath(path)+".png")
	cmd := exec.CommandContext(
		ctx,
		"ffmpeg",
		"-y",
		"-i", path,
		"-frames:v", "1",
		thumbPath,
	)

	var stderr bytes.Buffer
	cmd.Stdout = &stderr
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg thumbnail: %w", err)
	}

	return thumbPath, nil
}

func sanitizeMediaPath(path string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		" ", "_",
	)
	return replacer.Replace(path)
}
