package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	gcfg "gopkg.in/gcfg.v1"
)

func ReadModuleConfig(cfg interface{}, module string, path []string) error {
	environ := os.Getenv("MYENV")
	if environ == "" {
		environ = "development"
	}

	var err error
	for _, p := range path {
		fname := p + "/" + module + "." + environ + ".ini"
		err = gcfg.ReadFileInto(cfg, fname)
		if err == nil {
			return nil
		}

		log.Println(err)
	}

	return errors.New(fmt.Sprintf("cannot find config in %s", path))
}
