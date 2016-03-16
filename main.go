/* 
 * The MIT License (MIT)
 * 
 * MQTT Infrastructure
 * Copyright (c) 2016 Michael Jacobsen (github.com/mikejac)
 * 
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 */

//
// https://github.com/kardianos/service
// go get github.com/kardianos/service
//

package main

import (
    //"time"
    "log"
    "flag"
    "os"
    "os/signal"
    "syscall"    
    "github.com/vharitonsky/iniflags"
)

var (
    verbose = flag.Bool("v", false, "show debug logging")    
)

type NullWriter int

// Write ...
//
func (NullWriter) Write([]byte) (int, error) { 
    return 0, nil 
}

// main ...
//
func main() {
    iniflags.Parse()  // use instead of flag.Parse()
	
	// setup logging
	log.SetFlags(log.Ldate | log.Ltime)
    
    TinyG2Initialize()
    
    // get MQTT connectivity going
    if Mqtt() == false {
        return
    }
    
	if !*verbose {
		log.Println("You can enter verbose mode to see all logging by starting with the -v command line switch.")
		log.SetOutput(new(NullWriter))  // route all logging to nullwriter
	}
    
    sigs := make(chan os.Signal, 1)
    done := make(chan bool, 1)
    
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    
    go func() {
        sig := <-sigs
        log.Println(sig)
        done <- true
    }()
    
    <-done
    log.Println("main(): exiting")

    MqttStop()
}