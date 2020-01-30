package main

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
	"unsafe"
)
const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1 <<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
var src = rand.NewSource(time.Now().UnixNano())

func main(){
	c:=make(chan string,20000)
	var wg = sync.WaitGroup{}
	ctx,_ := context.WithTimeout(context.Background(),1*time.Second)
	wg.Add(20)
	for i:=0;i<20;i++{
		go func(ctx context.Context,c chan <- string){
			for {
				select{
					case <- ctx.Done():
						fmt.Println("time to return")
						wg.Done()
						return
					case c <- Hasher(128):
						fmt.Println("")
				}
			}
		}(ctx,c)
	}
	wg.Wait()
	sl :=make([]string,0,20000)
	close(c)
	for el:=range c{
		sl = append(sl,el)
	}
	sort.Strings(sl)
	for i:=1;i<len(sl);i++{
		if sl[i]==sl[i-1]{
			fmt.Println(sl[i],sl[i-1])
		}
	}
}
func Hasher(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}