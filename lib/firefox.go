package i2pfirefox

/*
Released under the The MIT License (MIT)
see ./LICENSE
*/

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	//	"os/signal"
	"path/filepath"
	//	"strings"
	//	"syscall"

	. "github.com/eyedeekay/go-fpw"
)

func UserFind(override ...string) string {
	if os.Geteuid() == 0 {
		log.Fatal("Do not run this application as root!")
	}
	if un, err := os.UserHomeDir(); err == nil {
		if len(override) > 0 {
			un = override[0]
		}
		envoverride := os.Getenv("RHZ_PROFILE_OVERRIDE")
		if envoverride != "" {
			un = envoverride
		}

		os.MkdirAll(filepath.Join(un, "i2p"), 0755)
		os.MkdirAll(filepath.Join(un, "i2p", "opt"), 0755)
		os.MkdirAll(filepath.Join(un, "i2p", "firefox-profiles", NOM), 0755)
		os.MkdirAll(filepath.Join(un, "i2p", "rhizome"), 0755)
		return un
	}
	return ""
}

var UserDir = filepath.Join(UserFind(), "i2p", "firefox-profiles", NOM)

var GingerDir = filepath.Join(UserFind(), "i2p", "rhizome")

func WriteExtension(val os.FileInfo, system *fs) bool {
	var firstrun = false
	if len(val.Name()) > 3 {
		if !val.IsDir() {
			file, err := system.Open(val.Name())
			if err != nil {
				log.Fatal(err.Error())
			}
			sys := bytes.NewBuffer(nil)
			if _, err := io.Copy(sys, file); err != nil {
				log.Fatal(err.Error())
			}
			log.Println(val.Name()[len(val.Name())-3:], "== xpi")
			if val.Name()[len(val.Name())-3:] == "xpi" {
				extension := filepath.Join(UserDir, "extensions", val.Name())
				if _, err := os.Stat(extension); os.IsNotExist(err) {
					if err := ioutil.WriteFile(extension, sys.Bytes(), val.Mode()); err == nil {
						ARGS = append(ARGS, extension)
						log.Println("wrote", extension)
					} else {
						log.Fatal(err)
					}
					firstrun = true
				}
			} else {
				log.Println(filepath.Join(UserDir, val.Name()), "ignored")
			}
		}
	} else {
		log.Println(filepath.Join(UserDir, val.Name()), "ignored", "contents", val.Sys())
	}
	return firstrun
}

func WriteProfile(FS *fs) bool {
	var firstrun = false
	if embedded, err := FS.Readdir(-1); err != nil {
		log.Fatal("Extension error, embedded extension not read.", err)
	} else {
		os.MkdirAll(filepath.Join(UserDir, "extensions"), 0755)
		/*err := ioutil.WriteFile(filepath.Join(UserDir, "extension-settings.json"), []byte(EXTENSIONPREFS), 0644)
		if err != nil {
			log.Fatal(err)
		}*/
		for _, val := range embedded {
			if val.IsDir() {
				os.MkdirAll(filepath.Join(UserDir, val.Name()), val.Mode())
			} else {
				firstrun = WriteExtension(val, FS)
			}
		}
	}
	return firstrun
}

func FirefoxLaunch() {
	firstrun := WriteProfile(FS)
	prefs := filepath.Join(UserDir, "/user.js")
	if _, err := os.Stat(prefs); os.IsNotExist(err) {
		if err := ioutil.WriteFile(prefs, []byte(PREFS), 0644); err == nil {
			log.Println("wrote", prefs)
		} else {
			log.Fatal(err)
		}
	}
	if firstrun {
		FIREFOX, ERROR := SecureExtendedFirefox(UserDir, false, EXTENSIONS, EXTENSIONHASHES, ARGS...)
		if ERROR != nil {
			log.Fatal(ERROR)
		}
		defer FIREFOX.Close()

		fmt.Println("exiting when browser window closes")
		<-FIREFOX.Done()
		fmt.Println(FIREFOX.Log())
	} else {
		FIREFOX, ERROR := BasicFirefox(UserDir, false, ARGS...)
		if ERROR != nil {
			log.Fatal(ERROR)
		}
		defer FIREFOX.Close()

		fmt.Println("exiting when browser window closes")
		<-FIREFOX.Done()
		fmt.Println(FIREFOX.Log())
	}
}

func FirefoxMain() {
	firstrun := WriteProfile(FS)
	prefs := filepath.Join(UserDir, "/user.js")
	if _, err := os.Stat(prefs); os.IsNotExist(err) {
		if err := ioutil.WriteFile(prefs, []byte(PREFS), 0644); err == nil {
			log.Println("wrote", prefs)
		} else {
			log.Fatal(err)
		}
	}
	if firstrun {
		FIREFOX, ERROR := SecureExtendedFirefox(UserDir, false, EXTENSIONS, EXTENSIONHASHES, ARGS...)
		if ERROR != nil {
			log.Fatal(ERROR)
		}
		defer FIREFOX.Close()

		fmt.Println("exiting when browser window closes")
		<-FIREFOX.Done()
		fmt.Println(FIREFOX.Log())
	} else {
		FIREFOX, ERROR := BasicFirefox(UserDir, false, ARGS...)
		if ERROR != nil {
			log.Fatal(ERROR)
		}
		defer FIREFOX.Close()

		fmt.Println("exiting when browser window closes")
		<-FIREFOX.Done()
		fmt.Println(FIREFOX.Log())
	}
}
