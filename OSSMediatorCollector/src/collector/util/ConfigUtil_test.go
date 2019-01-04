/*
* Copyright 2018 Nokia 
* Licensed under BSD 3-Clause Clear License,
* see LICENSE file for details.
*/

package util

import (
	"os"
	"testing"
)

//Reading config from invalid json file
func TestReadConfigWithEmptyFile(t *testing.T) {
	tmpfile := createTmpFile(".", "conf", []byte(``))
	defer os.Remove(tmpfile)
	err := ReadConfig(tmpfile)
	if err == nil {
		t.Error()
		t.Fail()
	}
}

//Reading valid config
func TestReadConfig(t *testing.T) {
	content := []byte(`{
		"base_url": "https://localhost:8080/api/v2",
		"users":
		[
			{
				"email_id" : "user1@nokia.com",
				"response_dest": "/statistics/reports/user1"
			},
			{
				"email_id" : "user2@nokia.com",
				"response_dest": "/statistics/reports/user2"
			}
		],
		"um_api":
		{
			"login": "/session",
			"refresh": "/refresh",
			"logout": "/logout"
		},
		"apis":
		[
			{
				"api": "/nms/pmdata",
				"interval": 15
			},
			{
				"api": "/nms/fmdata",
				"type": "active",
				"interval": 15
			},
			{
				"api": "/nms/fmdata",
				"type": "history",
				"interval": 60
			}
		],
		"limit": 100
	}`)
	tmpfile := createTmpFile(".", "conf", content)
	defer os.Remove(tmpfile)
	oldReadPassword := readPassword
	defer func() { readPassword = oldReadPassword }()

	myReadPassword := func(int) ([]byte, error) {
		return []byte("password"), nil
	}
	readPassword = myReadPassword
	err := ReadConfig(tmpfile)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(Conf)
	if Conf.BaseURL != "https://localhost:8080/api/v2" || len(Conf.Users) != 2 || Conf.Users[0].Email != "user1@nokia.com" {
		t.Fail()
	}
}

//Reading config from non existing file
func TestReadConfigWithNonExistingFile(t *testing.T) {
	err := ReadConfig("conf.json")
	if err == nil {
		t.Fail()
	}

}

func TestReadCredentials(t *testing.T) {
	oldReadPassword := readPassword
	defer func() { readPassword = oldReadPassword }()

	myReadPassword := func(int) ([]byte, error) {
		return []byte("password"), nil
	}
	readPassword = myReadPassword
	Conf.Users = []*User{&User{Email: "user@nokia.com"}}
	ReadCredentials()
	if Conf.Users[0].Email != "user@nokia.com" && Conf.Users[0].password != "password" {
		t.Fail()
	}
}