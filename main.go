package main

import (
	"os"
	"net/http"
	"flag"
	"fmt"
	"io"

	"github.com/heschlie/factorioInstaller/models"

	"github.com/mholt/archiver"
	"encoding/json"
	"io/ioutil"
)

const FACTORIO_URL = "https://www.factorio.com/get-download/stable/headless/linux64"
const FACTORIO_DIR = "/opt/factorio"

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
		Name: *serverName,
		Description: *description,
		Tags: []string{},
		MaxPlayers: 0,
		Visibility: models.ServerVisibility{
			Public: true,
			Lan: true,
		},
		Token: *token,
		GamePassword: *password,
		RequireUserVerification: false,
		MaxUpload: 0,
		MinimumLatency: 0,
		IgnorePlayerLimitForReturningPlayers: true,
		AllowCommands: "admins-only",
		AutosaveInterval: 10,
		AutosaveSlots: 5,
		AfkAutokickInterval: 0,
		AutoPause: true,
		OnlyAdminsCanPauseTheGame: true,
		AutosaveOnlyOnServer: true,
		Admins: []string{ "heschlie" },
	}

	// Download latest Factorio headless server into /opt/factorio.
	fmt.Println("Downloading server...")
	os.MkdirAll(FACTORIO_DIR, 0755)
	err := downloadFile(FACTORIO_DIR+"/factorio-headless.tar.gz", FACTORIO_URL)
	if err != nil {
		fmt.Errorf("failed to download server archive: %v", err)
	}

	fmt.Println("Unpacking server...")
	// Unpack the tar.gz into /opt/factorio.
	err = archiver.TarGz.Open(FACTORIO_DIR+"/factorio-headless.tar.gz", FACTORIO_DIR+"/")
	if err != nil {
		fmt.Errorf("failed to extract archive: %v", err)
	}

	b, err := json.Marshal(config)
	if err != nil {
		fmt.Errorf("there was an error marshaling the config to json: %v", err)
	}

	fmt.Println("Wrote config file...")
	err = ioutil.WriteFile(FACTORIO_DIR+"/config", b, 0644)
	if err != nil {
		fmt.Errorf("failed to save config: %v", err)
	}

	os.MkdirAll(FACTORIO_DIR+"/saves", 0755)
	os.MkdirAll(FACTORIO_DIR+"/mods", 0755)

	fmt.Printf("Downloading save from %s...\n", *saveFileUrl)
	err = downloadFile(FACTORIO_DIR+"/saves/save.zip", *saveFileUrl)
	if err != nil {
		fmt.Errorf("failed to download save file: %v", err)
	}

	fmt.Printf("Downloading mod zip from %s...\n", *modsZipUrl)
	err = downloadFile(FACTORIO_DIR+"/mods/mods.zip", *modsZipUrl)
	if err != nil {
		fmt.Errorf("failed to download mods zip: %v", err)
	}

	fmt.Println("Unpacking mods into /mods")
	err = archiver.Zip.Open(FACTORIO_DIR+"/mods/mods.zip", FACTORIO_DIR+"/mods/")
	if err != nil {
		fmt.Errorf("failed to extract mods: %v", err)
	}

	os.Remove(FACTORIO_DIR+"/mods/mods.zip")

	fmt.Println("Server ready to be launched! use the following command to launch:\n\n" +
		"/opt/factorio/bin/x64/factorio --start-server /opt/factorio/saves/save.zip --server-settings /opt/factorio/config")
}

// downloadFile Will download a file from the specified URL.
func downloadFile(filepath string, url string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil  {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil  {
		return err
	}

	return nil
}
