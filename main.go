package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/heschlie/factorioInstaller/models"

	"encoding/json"
	"io/ioutil"
	"github.com/cavaliercoder/grab"
	"github.com/mholt/archiver"
	"log"
	"os/exec"
)

const FACTORIO_URL = "https://www.factorio.com/get-download/stable/headless/linux64"
const FACTORIO_DIR = "/opt/factorio/"

func main() {
	token := flag.String("token", "", "Factorio account token")
	serverName := flag.String("name", "My server", "Server name to display on server listings")
	description := flag.String("description", "Created by golang script", "Server description")
	saveFileUrl := flag.String("save", "", "The URL to download save file from")
	password := flag.String("password", "", "Password for server, default is blank")
	modsZipUrl := flag.String("modsUrl", "", "URL for a zip of mods to install")
	flag.Parse()

	// Create our settings file.
	config := models.ServerSettings{
		Name:        *serverName,
		Description: *description,
		Tags:        []string{},
		MaxPlayers:  0,
		Visibility: models.ServerVisibility{
			Public: true,
			Lan:    true,
		},
		Token:                                *token,
		GamePassword:                         *password,
		RequireUserVerification:              false,
		MaxUpload:                            0,
		MinimumLatency:                       0,
		IgnorePlayerLimitForReturningPlayers: true,
		AllowCommands:                        "admins-only",
		AutosaveInterval:                     10,
		AutosaveSlots:                        5,
		AfkAutokickInterval:                  0,
		AutoPause:                            true,
		OnlyAdminsCanPauseTheGame:            true,
		AutosaveOnlyOnServer:                 true,
		Admins:                               []string{"heschlie"},
	}

	// Download latest Factorio headless server into /opt/factorio.
	fmt.Println("Downloading server...")
	os.MkdirAll(FACTORIO_DIR, 0755)
	resp, err := grab.Get(FACTORIO_DIR, FACTORIO_URL)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to download server archive: %v", err))
	}

	fmt.Println("Unpacking server...")
	// Unpack the tar.xz into /opt/factorio.
	err = archiver.TarXZ.Open(FACTORIO_DIR+resp.Filename, "/opt/")
	if err != nil {
		log.Fatal(fmt.Errorf("fialed to unpack %s: %v", resp.Filename, err))
	}

	b, err := json.Marshal(config)
	if err != nil {
		log.Fatal(fmt.Errorf("there was an error marshaling the config to json: %v", err))
	}

	fmt.Println("Wrote config file...")
	err = ioutil.WriteFile(FACTORIO_DIR+"config.json", b, 0644)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to save config: %v", err))
	}

	os.MkdirAll(FACTORIO_DIR+"saves", 0755)
	os.MkdirAll(FACTORIO_DIR+"mods", 0755)

	fmt.Printf("Downloading save from %s...\n", *saveFileUrl)
	_, err = grab.Get(FACTORIO_DIR+"saves/", *saveFileUrl)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to download save file: %v", err))
	}

	fmt.Printf("Downloading mod zip from %s...\n", *modsZipUrl)
	resp, err = grab.Get(FACTORIO_DIR+"mods/", *modsZipUrl)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to download mods zip: %v", err))
	}

	fmt.Println("Unpacking mods into /mods")
	err = archiver.Zip.Open(FACTORIO_DIR+"mods/"+resp.Filename, FACTORIO_DIR+"mods/")
	if err != nil {
		fmt.Errorf("failed to unpack %s: %v", resp.Filename, err)
	}

	os.Remove(FACTORIO_DIR + "mods/mods.zip")

	fmt.Println("Server ready to be launched! use the following command to launch:\n\n" +
		"/opt/factorio/bin/x64/factorio --start-server /opt/factorio/saves/save.zip --server-settings /opt/factorio/config.json")
}
