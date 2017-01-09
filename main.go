package main

import (
  	"archive/tar"
    "compress/gzip"
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
    var spec = getSpecJson()
    buildDeb(spec)
    fmt.Println("...ready.")
}

func buildDeb(spec SpecStruct) {
    fmt.Printf("\nBuild deb [%s]\n", spec.Name)
    debPathes := createDebianDir(spec)
    createDebianDir(spec)
    writeControleFile(spec, debPathes)
    writeChangeLogFile(spec, debPathes)
    extractSourceTarGz(spec, debPathes)
    out, err := exec.Command(
      "fakeroot",
      "dpkg-deb",
      "--build",
      debPathes.RootDir).CombinedOutput()
    if err != nil {
        fmt.Printf("\nError build deb: %s\n", err.Error())
        fmt.Printf("\nOut build deb: %s\n", out)
        log.Fatal(err)
    }
    fmt.Printf("\nDeb build: %s\n", out)
}

type DebPathes struct {
    DebDir string
    RootDir string
    ControlDir string
    SourceDir string
}

func createDebianDir(spec SpecStruct)(debPathes DebPathes){
    fmt.Println("HOME:", os.Getenv("HOME"))
    debPathes.DebDir = os.Getenv("HOME") + "/deb"
    createDir(debPathes.DebDir)
    debPathes.RootDir = debPathes.DebDir + "/" + spec.Name + "_" +
        spec.Version + "-" + spec.Release + "_" + spec.BuildArch
    createDir(debPathes.RootDir)
    debPathes.ControlDir = debPathes.RootDir + "/DEBIAN"
    createDir(debPathes.ControlDir)
    debPathes.SourceDir = debPathes.RootDir + "/SOURCE"
    createDir(debPathes.SourceDir)
    return
}

func createDir(name string){
    err := os.MkdirAll(name, 0777)
    if err != nil {
        fmt.Printf("\nError by creating directory %s: %s\n", name, err.Error())
        log.Fatal(err)
    }
}


func writeControleFile(spec SpecStruct, debPathes DebPathes){
    var content = "Package: " + spec.Name + "\n"
    content += "Version: " + spec.Version + "-" + spec.Release + "\n"
    content += "Section: " + spec.Group + "\n"
    content += "Priority: " + spec.Priority + "\n"
    content += "Architecture: " + spec.BuildArch + "\n"
    content += "Maintainer: " + spec.Maintainer + "\n"
    content += "Description: " + spec.Description + "\n"
    var controlFileName = debPathes.ControlDir + "/control"
    err := ioutil.WriteFile(controlFileName, []byte(content), 0644)
    if err != nil {
        fmt.Printf("\nError by writting file: %s\n", err.Error())
        log.Fatal(err)
    }
    fmt.Printf("\nWritting  control file in: %s\n", controlFileName)
}

func writeChangeLogFile(spec SpecStruct, debPathes DebPathes){
    var content = ""
    for _,element := range spec.ChangeLog {
        content += spec.Name + " (" + spec.Version + "-" + spec.Release + ") "
        content += element.Distribution + "; urgency=" + element.Urgency + "\n\n"
        for _,listElement := range element.Changes {
            content += "  * " + listElement + "\n"
        }
        content += "\n-- " + element.Author + " " + element.Date + " \n\n"
    }
    fmt.Printf("\nChangeLog: %s\n", content)
    err := os.MkdirAll(debPathes.RootDir + "/usr/share/doc/" + spec.Name, 0777)
    if err != nil {
        fmt.Printf("\nError by creating directory: %s\n", err.Error())
    }
    var changelogFileName = debPathes.RootDir + "/usr/share/doc/" +
        spec.Name + "/changelog.Debian"
    createZipFile(content, changelogFileName)
    createZipFile(content, changelogFileName + ".gz")
}

func createZipFile(content string, fileName string){
    zipFile, err := os.OpenFile(fileName,
        os.O_WRONLY|os.O_CREATE,
        0660)
    if err != nil {
        log.Printf("Error by create zip file\n")
    }
    w := gzip.NewWriter(zipFile)
    w.Write([]byte(content))
    w.Close()
    zipFile.Close()
    fmt.Printf("\nWritting  control file in: %s\n", fileName)
}

func extractSourceTarGz(spec SpecStruct, debPathes DebPathes) {
    fmt.Println("Extract source tar.gz")
    sourceTarGzPath := spec.Source
  	f, err := os.Open(sourceTarGzPath)
  	if err != nil {
    		fmt.Println(err)
    		os.Exit(1)
  	}
  	defer f.Close()
  	gzf, err := gzip.NewReader(f)
  	if err != nil {
    		fmt.Println(err)
    		os.Exit(1)
  	}

    tarReader := tar.NewReader(gzf)
    for true {
		    header, err := tarReader.Next()
	      if err == io.EOF {
            break
        }
        if err != nil {
          	fmt.Println(err)
          	os.Exit(1)
        }
        name := debPathes.SourceDir + header.Name
        switch header.Typeflag {
        case tar.TypeDir: // = directory
          	fmt.Println("Directory:", name)
          	os.Mkdir(name, 0755)
        case tar.TypeReg: // = regular file
          	fmt.Println("Regular file:", name)
          	data := make([]byte, header.Size)
          	_, err := tarReader.Read(data)
          	if err != nil {
      		      panic("Error reading file!")
          	}
            ioutil.WriteFile(name, data, 0755)
        default:
            fmt.Printf("%s : %c %s %s\n",
                "Unable to figure out type",
                header.Typeflag,
                "in file",
                name,
            )
        }
    }
}


type ChangeLog struct {
    Version string `version`
    Distribution string `distribution`
    Urgency string `urgency`
    Author string `author`
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
    Source string
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


func getSpecJson() (spec SpecStruct){

  var specJson = getSpecFileConten()
	dec := json.NewDecoder(strings.NewReader( specJson ))
		if err := dec.Decode(&spec); err == io.EOF {
//			break
		} else if err != nil {
	    fmt.Printf("Error: %s\n", err)
			log.Fatal(err)
		}
  return
}

func getSpecFileConten() (specData string){
    file, e := ioutil.ReadFile("./spec.json")
    if e != nil {
        fmt.Printf("File error: %v\n", e)
        os.Exit(1)
    }
    specData=string(file)
    return
}
