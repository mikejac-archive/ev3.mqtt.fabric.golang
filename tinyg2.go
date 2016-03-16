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

package main

import (
    "log"
    "encoding/json"
    "regexp"
    //"strings"
    //"io"
    //"fmt"
    "errors"
)

//
// https://github.com/synthetos/g2/wiki/Configuration-for-Firmware-Version-0.98
//

/*
- {ej:""}

{"D":"{\"js\":1}\n","Id":"tinygInit-cmd170","Pause":150},
{"D":"{\"sr\":n}\n","Id":"tinygInit-cmd171","Pause":0},
{"D":"{\"sv\":1}\n","Id":"tinygInit-cmd172","Pause":50},
{"D":"{\"si\":250}\n","Id":"tinygInit-cmd173","Pause":50},
{"D":"{\"qr\":n}\n","Id":"tinygInit-cmd174","Pause":0},
{"D":"{\"qv\":1}\n","Id":"tinygInit-cmd175","Pause":50},
{"D":"{\"ec\":0}\n","Id":"tinygInit-cmd176","Pause":50},
{"D":"{\"jv\":4}\n","Id":"tinygInit-cmd177","Pause":50},
{"D":"{\"hp\":n}\n","Id":"tinygInit-cmd178","Pause":0},
{"D":"{\"fb\":n}\n","Id":"tinygInit-cmd179","Pause":0},
{"D":"{\"mt\":n}\n","Id":"tinygInit-cmd180","Pause":0},
{"D":"{\"sr\":{\"line\":t,\"posx\":t,\"posy\":t,\"posz\":t,\"vel\":t,\"unit\":t,\"stat\":t,\"feed\":t,\"coor\":t,\"momo\":t,\"plan\":t,\"path\":t,\"dist\":t,\"mpox\":t,\"mpoy\":t,\"mpoz\":t}}\n","Id":"tinygInit-cmd181","Pause":250}

*/

var (
    re1                    *regexp.Regexp
    re2                    *regexp.Regexp
    re3                    *regexp.Regexp
    re4                    *regexp.Regexp
    re5                    *regexp.Regexp
)

// TinyG2Initialize ...
//
func TinyG2Initialize() {
    re1 = regexp.MustCompile(`\?`)
    re2 = regexp.MustCompile(`":n`)
    re3 = regexp.MustCompile(`":t`)
    re4 = regexp.MustCompile(`":f`)
    re5 = regexp.MustCompile(`{"`)
}

// TinyG2 ...
//
func TinyG2(cmd string) (string, error) {
    if re1.MatchString(cmd) {
        log.Println("TinyG2(): got '?'")
        return "tinyg [mm] ok>\n", nil
    }
    
    if !re5.MatchString(cmd) {
        log.Println("TinyG2(): GCODE")
        return gCode(cmd), nil
    }
    
    // fix the non-standard JSON coming from ChilliPeppr
    cmd = re2.ReplaceAllString(cmd, `":null`)
    cmd = re3.ReplaceAllString(cmd, `":true`)
    cmd = re4.ReplaceAllString(cmd, `":false`)
    
    var objmap map[string]*json.RawMessage
    
    if err := json.Unmarshal([]byte(cmd), &objmap); err != nil {
        log.Println("TinyG2(): ", err)
        return "", errors.New("TinyG2: cannot parse JSON top-level object")
    }

    var reply string
    
    for key, value := range objmap {
        log.Println("TinyG2(): key = ", key)
        
        if key == "ej" {
            reply = ej(value) + "\n"
        } else if key == "js" {
            reply = js(value) + "\n"      
        } else if key == "sr" {
            reply = sr(value) + "\n"      
        } else {
            reply = cmd
        }
    }
    return reply /*+ "\n"*/, nil
}

func ej(v interface{}) (string) {
    return `{"ej":""}`
}
func js(v interface{}) (string) {
    return `{"js":1}`
}
func sr(v interface{}) (string) {
    return `{"sr":null}`
}


func gCode(data string) (string) {
    return data
}
    /*dec := json.NewDecoder(strings.NewReader(cmd))
	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Err:", err)
            
            if err.Error() == "invalid character '}' in literal null (expecting 'u')" {
                fmt.Println("TinyG2(): 'null' detected")
            }
            return ""
		}
	    fmt.Printf("%T: %v", t, t)
		if dec.More() {
			fmt.Printf(" (more)")
		}
		fmt.Printf("\n")
	}*/

/*func TinyG2ResponseInt(cmd string) (string, error) {
    type R struct {
        Value int `json:"ej"`
    }
    
    type Resp struct {
        R R      `json:"r"`
        F []int  `json:"f"`
    }

    r := Resp{
        R: R{   Value:  1,
            },
        F: []int{1, 0, 7},
    }
        
	msg, err := json.Marshal(r) 
    
    if err != nil {
		log.Println("TinyG2Response(): error = ", err)
	}
    
    log.Println("TinyG2Response(): msg = ", string(msg))

    return string(msg), nil
}*/
