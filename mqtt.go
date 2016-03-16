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
    "flag"
    "log"
    "strconv"
    //"errors"
    "github.com/mikejac/mqtt.fabric.golang"
)

// MQTT stuff
var (
    mqtt *mqttfabric.MqttFabric = nil
    
	mqttRootTopic  = flag.String("root_topic", "fabric", "desc.")
	mqttNodename   = flag.String("uuid", "", "desc.")
	mqttBroker     = flag.String("broker", "", "desc.")
    mqttPort       = flag.Int("port", 1883, "desc.")
    mqttKeepalive  = flag.Int("keepalive", 60, "desc.")
    
	mqttPlatformID = "ev3"
    mqttFeedID     = "cp"
    
    mqttSourcePlatformId = "spjs"
    
    mqttDestNodename   = "894051ec-d0e7-4862-ab09-832a7fdab16c"
    mqttDestPlatformID = "ev3"
    mqttDestFeedID     = "cp"
)

// Mqtt ...
//
func Mqtt() (bool) {
    if len(*mqttBroker) == 0 {
    	log.Println("No name or ip-address specified for MQTT Server")
        return false
    }
    
    if len(*mqttNodename) == 0 {
    	log.Println("No nodename specified")
        return false
    }
    
	log.Println("Connecting to MQTT Server on " + *mqttBroker + ":" + strconv.Itoa(*mqttPort))

	mqtt = mqttfabric.MqttFabricInitialize(*mqttBroker, *mqttPort, *mqttKeepalive, *mqttRootTopic, *mqttNodename, mqttPlatformID, mqttfabric.DEVICE)
    mqtt.SetOnConnectHandler(MQTTOnConnect)
    mqtt.SetOnDisconnectHandler(MQTTOnDisconnect)
    mqtt.SetOnOnrampHandler(MQTTOnOnramp)
    mqtt.SetOnOfframpHandler(MQTTOnOfframp)
    
    mqtt.Start()
    
    return true
}

// MqttStop ...
//
func MqttStop() {
    mqtt.Stop()
}

// MqttSend
//
func MqttSend(data string) {
    mqtt.CtrlPubText(mqttDestNodename, mqttDestPlatformID, mqttDestFeedID, data, 0, false)
}

// MQTTOnConnect ...
var MQTTOnConnect mqttfabric.OnConnectHandler = func(mqtt *mqttfabric.MqttFabric) {
    log.Println("MQTTOnConnect()")
    
    topic := mqtt.F.DeviceOfframpSubscription(*mqttNodename, mqttfabric.FABRIC_TOPIC_ANY /*actorID*/, mqttSourcePlatformId /*actorPlatformID*/, mqttfabric.TASK_ID_RAW, mqttPlatformID, mqttfabric.SERVICE_ID_TEXT, mqttFeedID)
    
    log.Printf("MQTTOnConnect(): topic = %s", topic)
    mqtt.Mqtt.Subscribe(topic, 0, nil)
}

// MQTTOnDisconnect ...
var MQTTOnDisconnect mqttfabric.OnDisconnectHandler = func(mqtt *mqttfabric.MqttFabric) {
    log.Println("MQTTOnDisconnect()")
}

// MQTTOnOnramp ...
var MQTTOnOnramp mqttfabric.OnOnrampHandler = func( mqtt            *mqttfabric.MqttFabric,
                                                    nodename        string,
                                                    platformID      string,
                                                    serviceID       string,
                                                    feedID          string,
                                                    msg             string) {
    log.Printf("MQTTOnOnramp(): nodename        = %s\n", nodename)
    log.Printf("MQTTOnOnramp(): platformID      = %s\n", platformID)
    log.Printf("MQTTOnOnramp(): serviceID       = %s\n", serviceID)
    log.Printf("MQTTOnOnramp(): feedID          = %s\n", feedID)
    log.Printf("MQTTOnOnramp(): msg             = %s\n", msg)
}

// MQTTOnOfframp ...
var MQTTOnOfframp mqttfabric.OnOfframpHandler = func(   mqtt            *mqttfabric.MqttFabric,
                                                        nodename        string,
                                                        actorID         string,
                                                        actorPlatformID string,
                                                        taskID          string,
                                                        platformID      string,
                                                        serviceID       string,
                                                        feedID          string,
                                                        msg             string) {
    /*log.Printf("MQTTOnOfframp(): nodename        = %s\n", nodename)
    log.Printf("MQTTOnOfframp(): actorID         = %s\n", actorID)
    log.Printf("MQTTOnOfframp(): actorPlatformID = %s\n", actorPlatformID)
    log.Printf("MQTTOnOfframp(): taskID          = %s\n", taskID)
    log.Printf("MQTTOnOfframp(): platformID      = %s\n", platformID)
    log.Printf("MQTTOnOfframp(): serviceID       = %s\n", serviceID)
    log.Printf("MQTTOnOfframp(): feedID          = %s\n", feedID)
    log.Printf("MQTTOnOfframp(): msg             = %s\n", msg)*/
    
    if bmix, err := mqttfabric.BlueMixParse(msg); err == nil {
        if val, err := bmix.GetValueString(); err == nil {
            log.Printf("MQTTOnOfframp() val = '%s'", val)
            
            reply, err := TinyG2(val)
            
            if err == nil {
                mqtt.DevicePubText(feedID, reply, 0, false)
            }
            
            /*if val == "?\n" {
                mqtt.DevicePubText(feedID, "tinyg [mm] ok>\n", 0, false)
            } else if val == "{\"ej\":\"\"}\n" {
                msg, _ :=  TinyG2ResponseInt(val)
                mqtt.DevicePubText(feedID, msg + "\n", 0, false)
            } else if val == "{\"js\":1}\n" {
                mqtt.DevicePubText(feedID, "{\"js\":1}\n", 0, false)
            } else if val == "{\"sr\":n}\n" {
                mqtt.DevicePubText(feedID, "{\"sr\":n}\n", 0, false)
            } else if val == "{\"sv\":1}\n" {
                mqtt.DevicePubText(feedID, "{\"sv\":1}\n", 0, false)
            } else if val == "{\"si\":250}\n" {
                mqtt.DevicePubText(feedID, "{\"si\":250}\n", 0, false)
            } else if val == "{\"qr\":n}\n" {
                mqtt.DevicePubText(feedID, "{\"qr\":n}\n", 0, false)
            } else if val == "{\"qv\":1}\n" {
                mqtt.DevicePubText(feedID, "{\"qv\":1}\n", 0, false)
            } else if val == "{\"ec\":0}\n" {
                mqtt.DevicePubText(feedID, "{\"ec\":0}\n", 0, false)
            } else if val == "{\"jv\":4}\n" {
                mqtt.DevicePubText(feedID, "{\"jv\":4}\n", 0, false)
            } else if val == "{\"hp\":n}\n" {
                mqtt.DevicePubText(feedID, "{\"hp\":n}\n", 0, false)
            } else if val == "{\"fb\":n}\n" {
                mqtt.DevicePubText(feedID, "{\"fb\":n}\n", 0, false)
            } else if val == "{\"mt\":n}\n" {
                mqtt.DevicePubText(feedID, "{\"mt\":n}\n", 0, false)
            } else if val == "{\"sr\":{\"line\":t,\"posx\":t,\"posy\":t,\"posz\":t,\"vel\":t,\"unit\":t,\"stat\":t,\"feed\":t,\"coor\":t,\"momo\":t,\"plan\":t,\"path\":t,\"dist\":t,\"mpox\":t,\"mpoy\":t,\"mpoz\":t}}\n" {
                mqtt.DevicePubText(feedID, "{\"sr\":{\"line\":t,\"posx\":t,\"posy\":t,\"posz\":t,\"vel\":t,\"unit\":t,\"stat\":t,\"feed\":t,\"coor\":t,\"momo\":t,\"plan\":t,\"path\":t,\"dist\":t,\"mpox\":t,\"mpoy\":t,\"mpoz\":t}}\n", 0, false)
            } else {
                log.Printf("MQTTOnOfframp() unhandled; val = '%s'", val)
            }*/
            
            //mqtt.DevicePubText(feedID, "{\"r\":{}}\u0011|\u0013", 0, false)
        }
    }
}
