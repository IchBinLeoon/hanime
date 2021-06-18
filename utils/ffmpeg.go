package utils

import "os/exec"

func MergeToMP4(listPath string, outputPath string) ([]byte, error) {
	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", listPath, "-c", "copy", outputPath)
	out, err := cmd.Output()
	if err != nil {
		return out, err
	}
	return out, nil
}
