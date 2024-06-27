package main

import (
	"fmt"
	"os"
	"time"

	"github.com/armiariyan/synapsis/cmd"
	"github.com/labstack/gommon/color"
)

const banner = `
███████╗██╗   ██╗███╗   ██╗ █████╗ ██████╗ ███████╗██╗███████╗
██╔════╝╚██╗ ██╔╝████╗  ██║██╔══██╗██╔══██╗██╔════╝██║██╔════╝
███████╗ ╚████╔╝ ██╔██╗ ██║███████║██████╔╝███████╗██║███████╗
╚════██║  ╚██╔╝  ██║╚██╗██║██╔══██║██╔═══╝ ╚════██║██║╚════██║
███████║   ██║   ██║ ╚████║██║  ██║██║     ███████║██║███████║
╚══════╝   ╚═╝   ╚═╝  ╚═══╝╚═╝  ╚═╝╚═╝     ╚══════╝╚═╝╚══════╝
                                                              
armiariyan.                                                         

`

func main() {
	if tz := os.Getenv("TZ"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			fmt.Printf("error loading location '%s': %v\n", tz, err)
		} else {
			fmt.Printf("location loaded '%s'\n", tz)
		}
	}

	fmt.Print(color.Yellow(banner))
	cmd.Run()
}
