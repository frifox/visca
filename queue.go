package visca

import "sync"

type PanTiltQueue struct {
	queue chan *PanTiltDrive
	sync.Mutex
}

func (p *PanTiltQueue) Init() {
	p.queue = make(chan *PanTiltDrive, 1)
}

func (p *PanTiltQueue) Add(cmd *PanTiltDrive) {
	p.Lock()
	defer p.Unlock()

	p.Clear()
	p.queue <- cmd
}
func (p *PanTiltQueue) Get() (cmd *PanTiltDrive) {
	return <-p.queue
}
func (p *PanTiltQueue) Clear() {
	for {
		select {
		case <-p.queue:
		default:
			return
		}
	}
}

type ZoomQueue struct {
	queue chan *Zoom
	sync.Mutex
}

func (p *ZoomQueue) Init() {
	p.queue = make(chan *Zoom, 1)
}

func (p *ZoomQueue) Add(cmd *Zoom) {
	p.Lock()
	defer p.Unlock()

	p.Clear()
	p.queue <- cmd
}
func (p *ZoomQueue) Get() (cmd *Zoom) {
	return <-p.queue
}
func (p *ZoomQueue) Clear() {
	for {
		select {
		case <-p.queue:
		default:
			return
		}
	}
}
