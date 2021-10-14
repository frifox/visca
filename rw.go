package visca

import "fmt"

type RW struct{}

func (rw *RW) Read(p []byte) (n int, err error) {
	fmt.Printf("[RW] Read...\n")
	select {}
}
func (rw *RW) Write(p []byte) (n int, err error) {
	n = len(p)
	fmt.Printf("[RW] Write len(%d) [% X]\n", n, p)
	return
}
