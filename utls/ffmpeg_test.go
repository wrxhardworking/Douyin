package utls

import "testing"

func TestFfmpeg(t *testing.T) {
	err := GenerateSnapshot("../public/18_1777864591.mp4", "../public/bear.mp4", 1)
	if err != nil {
		return
	}
}
