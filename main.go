// filewatcher project main.go
package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

var interval int

type Param struct {
	Times time.Time
	Cmd   []string
}

func init() {
	flag.IntVar(&interval, "int", 1, "interval(sec)")
}

func execCmd(params []string) {
	log.Printf("exec %v\n", params)
	cm := exec.Command(params[0], params[1:]...)
	out, err := cm.CombinedOutput()
	if err != nil {
		log.Println(string(out))
		log.Println(err)
	} else {
		log.Println(string(out))
		log.Println("execute completed")
	}
}

func watch(flist map[string]Param) {
	for fname, param := range flist {
		fi, err := os.Stat(fname)
		if err == nil {
			if fi.ModTime().After(param.Times) {
				log.Printf("file:%s  updated %v\n", fi.Name(), fi.ModTime())
				flist[fname] = Param{Times: fi.ModTime(), Cmd: param.Cmd}
				go execCmd(param.Cmd)
			}
		}
	}
}
func main() {
	flag.Parse()
	jsonName := os.Getenv("USERPROFILE") + `\filewatcher.json`
	b, err := ioutil.ReadFile(jsonName)
	if err != nil {
		log.Fatalln("filewatcher.json file not found")
	}
	var flist map[string]Param
	err = json.Unmarshal(b, &flist)
	if err != nil {
		log.Fatalln("Unmarshal error :", err)
	}
	if len(flist) == 0 {
		log.Fatalln("no entries")
	}
	c := time.Tick(time.Second * time.Duration(interval))
	for _ = range c {
		watch(flist)
		buff, err := json.Marshal(flist)
		if err != nil {
			log.Fatal("Marshal error :", err)
		}
		err = ioutil.WriteFile(jsonName, buff, os.ModePerm)
		if err != nil {
			log.Fatal("json write error : ", err)
		}
	}
}
