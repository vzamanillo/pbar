package progress

// NoOpProgress no operation in progress
type NoOpProgress struct{}

// InitProgressbar initializers progressbar
func (p *NoOpProgress) InitProgressbar(hostCount int64, templateCount int, requestCount int64) {}

// AddToTotal adds total to progressbar
func (p *NoOpProgress) AddToTotal(delta int64) {}

// Update updates progressbar
func (p *NoOpProgress) Update() {}

// Drop drop from progressbar
func (p *NoOpProgress) Drop(count int64) {}

// Wait waits for waitgroup to finish
func (p *NoOpProgress) Wait() {}
