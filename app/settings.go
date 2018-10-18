package app

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var (
	home     string
	userHome string
)

func init() {
	userHome = "."
	u, err := user.Current()
	if err != nil {
		log.Printf("unable to find home directory: %v", err)
	} else {
		userHome = u.HomeDir
	}
	flag.StringVar(&home, "home", "", "path to runtime data")
}

func Home() string {
	if home != "" {
		return home
	}
	envVar := fmt.Sprintf("%v_HOME", strings.ToUpper(Slug))
	envHome := os.Getenv(envVar)
	if envHome != "" {
		return envHome
	}
	return filepath.Join(userHome, Slug)
}
