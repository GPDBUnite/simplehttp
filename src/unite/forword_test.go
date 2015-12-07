package unite

import (
       "testing"
)

func TestStart(t *testing.T) {
     l := "127.0.0.1:2222"
     r1 := "127.0.0.1:22"
     remoteaddrs := []string{r1}
     f := NewForwarder(l,remoteaddrs, nil)
     go f.Start()
     t.Error("failed")
       
}
