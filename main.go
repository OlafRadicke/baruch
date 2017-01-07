package main

import (
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "os"
    "os/exec"
    "strings"
)

func main() {
    fmt.Println("Start...")
    readJson()
    buildDeb()
    fmt.Println("...ready.")
}


func buildDeb() {
  // fakeroot dpkg-deb --build debian
    out, err := exec.Command("dpkg-deb", "--build", "./dist/demo_0.1-2_i386/").Output()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("The date is %s\n", out)
}


func readJson() {

  type ChangeLog struct {
    Version string `version`
    Distribution string `distribution`
    Urgency string `urgency`
    Date string `date`
    Changes []string `changes`
  }

	type SpecStruct struct {
		Name string
    Version string
    Release string
    Summary string
    Group string
    License string
    URL string
    Source0 string
    BuildRoot string
    BuildRequires string
    Requires string
    Description  string `description`
    Prep string `prep`
    Setup string `setup`
    Build string `build`
    Install []string `install`
    Clean []string `clean`
    ChangeLog []ChangeLog
	}
  var specJson = readSpecFile()
	dec := json.NewDecoder(strings.NewReader( specJson ))
	for {
		var spec SpecStruct
		if err := dec.Decode(&spec); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
    fmt.Printf("install: %s\n", spec.Install[0])
		fmt.Printf("Value: %s\n", spec.Name, spec)
	}
}

func readSpecFile() (specData string){
    file, e := ioutil.ReadFile("./spec.json")
    if e != nil {
        fmt.Printf("File error: %v\n", e)
        os.Exit(1)
    }
    fmt.Printf("%s\n", string(file))
    specData=string(file)
    return
}
