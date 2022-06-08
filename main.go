/*
	Copyright (C) 2022  Michael Ablassmeier <abi@grinser.de>

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
)

type jsonLogin struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type jsonUpload struct {
	IgcContent string `json:"igccontent"`
	IgcName    string `json:"igcname"`
	Glider     string `json:"glider"`
	Publish    bool   `json:"publish"`
}

type jsonResponse struct {
	Success bool
	Message string
}

type combined struct {
	jsonLogin
	jsonUpload
}

type Options struct {
	User    string `short:"u" long:"user" description:"DHV-XC User name" required:"true"`
	Pass    string `short:"p" long:"pass" description:"DHV-XC Upload Password" required:"true"`
	File    string `short:"f" long:"file" description:"IGC file" required:"true"`
	Publish bool   `short:"P" long:"publish" description:"Publish flight after upload" required:"false"`
	Glider  string `short:"g" long:"glider" description:"Glider name" required:"false"`
}

func init() {
	_, ok := os.LookupEnv("XC_DEBUG")
	if ok {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func json_dumps(data interface{}) []byte {
	payload, err := json.Marshal(data)
	if err != nil {
		logrus.Error("Cant dump json response:")
		logrus.Fatal(err)
	}
	return payload
}

func json_loads(data []byte) jsonResponse {
	var resp jsonResponse
	err := json.Unmarshal([]byte(data), &resp)
	if err != nil {
		logrus.Error("Cant load json response:")
		logrus.Fatal(err)
	}
	return resp
}

func httpReq(url string, payload []byte) jsonResponse {
	client := &http.Client{}
	logrus.Debug(url)
	logrus.Debugf("Request: [%s]", string(payload))

	request, error := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	response, error := client.Do(request)
	if error != nil {
		logrus.Fatal(error)
	}
	body, _ := ioutil.ReadAll(response.Body)
	logrus.Debugf("Response: [%s]", string(body))
	return json_loads(body)
}

func success(resp jsonResponse) bool {
	if resp.Success {
		return true
	}
	return false
}

func main() {
	var options Options
	var parser = flags.NewParser(&options, flags.Default)

	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}
	Api := struct {
		url    string
		auth   string
		upload string
	}{
		url:    "https://de.dhv-xc.de/api/v1/",
		auth:   "authcheck",
		upload: "flights",
	}

	data := jsonLogin{
		User: options.User,
		Pass: options.Pass,
	}

	f := httpReq(Api.url+Api.auth, json_dumps(data))
	if !success(f) {
		logrus.Fatalf("Authentication failed: [%s]", f.Message)
	}
	logrus.Info("Login OK")

	igcdata, err := ioutil.ReadFile(options.File)
	if err != nil {
		log.Fatalf("Unable to read igc file: [%s]", err)
	}
	data2 := jsonUpload{
		IgcContent: string(igcdata),
		IgcName:    options.File,
		Publish:    options.Publish,
	}
	if options.Glider != "" {
		logrus.Infof("Glider: [%s]", options.Glider)
		data2.Glider = options.Glider
	}
	uploadPayload := combined{
		jsonLogin:  data,
		jsonUpload: data2,
	}
	if !options.Publish {
		logrus.Warning("Flight will not be auto-published, publish manually after upload!")
	} else {
		logrus.Info("Publishing flight during upload.")
	}
	f = httpReq(Api.url+Api.upload, json_dumps(uploadPayload))
	if !success(f) {
		logrus.Fatalf("Upload failed: [%s]", f.Message)
	}
	logrus.Infof("Upload OK: [%s]", f.Message)
}
