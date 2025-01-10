package main

import (
	"crypto/rand"
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/iu0jgo/gumble/gumble"
	_ "github.com/iu0jgo/gumble/opus"
	"github.com/iu0jgo/rcs7100/internal/engine"
	"github.com/spf13/viper"
)

func main() {
	fmt.Println("Remote Control System for IC-7100")

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	// viper.AddConfigPath("/etc/appname/")   // path to look for the config file in
	// viper.AddConfigPath("$HOME/.appname")  // call multiple times to add many search paths
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// Command line flags
	// server := flag.String("server", "localhost:64738", "the server to connect to")
	address := viper.GetString("server.address")
	username := flag.String("username", "", "the username of the client")
	password := flag.String("password", "", "the password of the server")
	insecure := flag.Bool("insecure", true, "skip server certificate verification")
	certificate := flag.String("certificate", "", "PEM encoded certificate and private key")
	channel := flag.String("channel", "root", "mumble channel to join by default")
	// Adudio devices
	playbackAudioDevice := viper.GetString("audio.playbackDevice")
	captureAudioDevice := viper.GetString("audio.captureDevice")

	flag.Parse()

	// Initialize
	b := engine.RCS7100{
		Config:              gumble.NewConfig(),
		Address:             address,
		ChannelName:         *channel,
		PlaybackAudioDevice: playbackAudioDevice,
		CaptureAudioDevice:  captureAudioDevice,
	}

	// if no username specified, lets just autogen a random one
	if len(*username) == 0 {
		buf := make([]byte, 6)
		_, err := rand.Read(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		buf[0] |= 2
		b.Config.Username = fmt.Sprintf("rcs7100-%02x%02x%02x%02x%02x%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
	} else {
		b.Config.Username = *username
	}

	b.Config.Password = *password

	if *insecure {
		b.TLSConfig.InsecureSkipVerify = true
	}
	if *certificate != "" {
		cert, err := tls.LoadX509KeyPair(*certificate, *certificate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		b.TLSConfig.Certificates = append(b.TLSConfig.Certificates, cert)
	}

	b.Init()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	exitStatus := 0

	<-sigs
	b.CleanUp()

	os.Exit(exitStatus)
}
