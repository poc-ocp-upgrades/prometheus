package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/pkg/errors"
)

type debugWriterConfig struct {
	serverURL	string
	tarballName	string
	endPointGroups	[]endpointsGroup
}

func debugWrite(cfg debugWriterConfig) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	archiver, err := newTarGzFileWriter(cfg.tarballName)
	if err != nil {
		return errors.Wrap(err, "error creating a new archiver")
	}
	for _, endPointGroup := range cfg.endPointGroups {
		for url, filename := range endPointGroup.urlToFilename {
			url := cfg.serverURL + url
			fmt.Println("collecting:", url)
			res, err := http.Get(url)
			if err != nil {
				return errors.Wrap(err, "error executing HTTP request")
			}
			body, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				return errors.Wrap(err, "error reading the response body")
			}
			if endPointGroup.postProcess != nil {
				body, err = endPointGroup.postProcess(body)
				if err != nil {
					return errors.Wrap(err, "error post-processing HTTP response body")
				}
			}
			if err := archiver.write(filename, body); err != nil {
				return errors.Wrap(err, "error writing into the archive")
			}
		}
	}
	if err := archiver.close(); err != nil {
		return errors.Wrap(err, "error closing archive writer")
	}
	fmt.Printf("Compiling debug information complete, all files written in %q.\n", cfg.tarballName)
	return nil
}
