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
    var spec = readJson()
    buildDeb(spec)
    fmt.Println("...ready.")
}

func createDistDir() {


}

func getBuildRoot(spec SpecStruct)(buildRoot string){
    buildRoot = spec.BuildRoot + "/"+ spec.Name + "_" + spec.Version + "-" + spec.Release + "_" + spec.BuildArch
    return
}

func writeControleFile(spec SpecStruct){
//  d1 := []byte("hello\ngo\n")
    var content = "Package: " + spec.Name + "\n"
    content += "Version: " + spec.Version + "-" + spec.Release + "\n"
    content += "Section: " + spec.Group + "\n"
    content += "Priority: " + spec.Priority + "\n"
    content += "Architecture: " + spec.BuildArch + "\n"
    content += "Maintainer: " + spec.Maintainer + "\n"
    content += "Description: " + spec.Description + "\n"
    var buildRoot = getBuildRoot(spec)
    var controlFileName = buildRoot + "/DEBIAN/control"

    err := os.MkdirAll(buildRoot + "/DEBIAN", 0777)
    if err != nil {
        fmt.Printf("\nError by creating directory: %s\n", err.Error())
        log.Fatal(err)
    }
    err = ioutil.WriteFile(controlFileName, []byte(content), 0644)
  //  err = ioutil.WriteFile("/tmp/test.muell", []byte(content), 0644)
    if err != nil {
        fmt.Printf("\nError by writting file: %s\n", err.Error())
        log.Fatal(err)
    }
    fmt.Printf("\nWritting  control file in: %s\n", controlFileName)
}

func buildDeb(spec SpecStruct) {
    fmt.Printf("\nBuild deb [%s]\n", spec.Name)
    writeControleFile(spec)
    var buildRoot = getBuildRoot(spec)
    out, err := exec.Command(
      "fakeroot",
      "dpkg-deb",
      "--build",
      buildRoot).CombinedOutput()
    if err != nil {
        fmt.Printf("\nError build deb: %s\n", err.Error())
        fmt.Printf("\nOut build deb: %s\n", out)
        log.Fatal(err)
    }
    fmt.Printf("\nDeb build: %s\n", out)
}

type ChangeLog struct {
    Version string `version`
    Distribution string `distribution`
    Urgency string `urgency`
    Date string `date`
    Changes []string `changes`
}

type Files struct {
    Path string `path`
    Defattr string `defattr`
}

type SpecStruct struct {
    Name string
    Version string
    Release string
    Priority string
    Summary string
    Group string
    BuildArch string
    License string
    URL string
    Maintainer string
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
    Files []Files `files`
    ChangeLog []ChangeLog `changelog`
}


func readJson() (spec SpecStruct){

  var specJson = readSpecFile()
	dec := json.NewDecoder(strings.NewReader( specJson ))
		if err := dec.Decode(&spec); err == io.EOF {
//			break
		} else if err != nil {
	    fmt.Printf("Error: %s\n", err)
			log.Fatal(err)
		}
  return
}

func readSpecFile() (specData string){
    file, e := ioutil.ReadFile("./spec.json")
    if e != nil {
        fmt.Printf("File error: %v\n", e)
        os.Exit(1)
    }
    specData=string(file)
    return
}
