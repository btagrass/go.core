package utl

import (
	"path/filepath"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// 抓取流
func TakeStream(input string, output map[string]map[string]any) error {
	var outputs []*ffmpeg.Stream
	inputArgs := make(map[string]any)
	if strings.HasPrefix(input, "rtsp://") {
		inputArgs["rtsp_transport"] = "tcp"
	}
	stream := ffmpeg.Input(input, inputArgs)
	for k, v := range output {
		MakeDir(filepath.Dir(k))
		outputs = append(outputs, stream.Output(k, v))
	}

	return ffmpeg.
		MergeOutputs(outputs...).
		OverWriteOutput().
		ErrorToStdOut().
		Run()
}
