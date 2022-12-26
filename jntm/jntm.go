package main

import (
	"fmt"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

const (
	J = 106
	N = 110
	T = 116
	M = 109
)

func main() {
	fmt.Println(os.Getwd())
	j := &Music{}
	j.Sources = append(j.Sources, "./j.mp3", "./n.mp3", "./t.mp3", "./m.mp3")
	j.Plays()
	fmt.Println("init players")
	kb := &KeyBoard{}
	kb.cxk = j
	kb.Run()
	select {}
}

type Music struct {
	sscs       []beep.StreamSeekCloser
	Source     string   //位置
	Sources    []string //位置
	Filestream *os.File // 文件流
}

func (m *Music) Play() {
	var err error
	m.Filestream, err = os.Open(m.Source)
	if err != nil {
		panic(err)
	}

	streamer, format, err := mp3.Decode(m.Filestream)
	if err != nil {
		panic(err)
	}
	defer streamer.Close()
	fmt.Println(format.SampleRate.N(time.Second / 10))
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(streamer)
	for {
		streamer.Seek(0)
		speaker.Play(streamer)
		time.Sleep(time.Second)
	}
}

func (m *Music) Plays() {
	var err error
	f := make([]*os.File, len(m.Sources))
	m.sscs = make([]beep.StreamSeekCloser, len(m.Sources))
	var format beep.Format
	for i := range m.Sources {
		fmt.Println(m.Sources[i])
		f[i], err = os.Open(m.Sources[i])
		if err != nil {
			panic(err)
		}

		m.sscs[i], format, err = mp3.Decode(f[i])
		if err != nil {
			panic(err)
		}
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
}

func (m *Music) Jiao(seq int) {
	m.sscs[seq].Seek(0)
	speaker.Play(m.sscs[seq])
	time.Sleep(300 * time.Millisecond)
}

type KeyBoard struct {
	cxk *Music
}

func (kb *KeyBoard) Run() {
	robotgo.EventHook(hook.KeyDown, []string{}, func(ev hook.Event) {
		switch ev.Keychar {
		case J:
			kb.cxk.Jiao(0)
		case N:
			kb.cxk.Jiao(1)
		case T:
			kb.cxk.Jiao(2)
		case M:
			kb.cxk.Jiao(3)
		default:

		}
	})
	s := robotgo.EventStart()
	<-robotgo.EventProcess(s)
}
