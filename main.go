package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/xmlquery"
	"github.com/wargarblgarbl/libgosubs/srt"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	doc, err := xmlquery.Parse(f)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	session := xmlquery.FindOne(doc, "//session")
	sampleRate, _ := strconv.ParseFloat(session.SelectAttr("sampleRate"), 32)

	var subtitles []srt.Subtitle
	for i, audioClip := range session.SelectElements("/tracks/audioTrack/audioClip") {
		startPoint, _ := strconv.ParseFloat(audioClip.SelectAttr("startPoint"), 32)
		endPoint, _ := strconv.ParseFloat(audioClip.SelectAttr("endPoint"), 32)
		subtitles = append(subtitles, *srt.CreateSubtitle(i+1, format(startPoint/sampleRate), format(endPoint/sampleRate), []string{audioClip.SelectAttr("name")}))
	}

	subRip := &srt.SubRip{Subtitle: struct {
		Content []srt.Subtitle
	}{
		Content: subtitles,
	}}

	err = srt.WriteSrt(subRip, "export.srt")
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
}

func format(auditionTime float64) string {
	start := time.Date(2009, 1, 1, 0, 0, 0, 0, time.UTC)
	u, _ := time.ParseDuration(fmt.Sprintf("%fs", auditionTime))
	formatted := start.Add(u).Format(time.StampMilli)
	return strings.ReplaceAll(strings.TrimLeft(formatted, "Jan  1 "), ".", ",")
}
