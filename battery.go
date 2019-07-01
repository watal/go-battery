package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

// Define args
type GeneralOption struct {
	OutputTmux  bool `short:"t" description:"output tmux status bar format"`
	OutputZsh   bool `short:"z" description:"output zsh prompt format"`
	Emoji       bool `short:"e" description:"don't output the emoji"`
	Ascii       bool `short:"a" description:"output ascii instead of spark"`
	BatteryPath bool `short:"b" description:"battery path (default: /sys/class/power_supply/BAT0)"`
	PmsetOn     bool `short:"p" description:"use pmset (more accurate)"`
	NerdFonts   bool `short:"n" description:"use Nerd Fonts battery icon"`
}

type ColorsOption struct {
	GoodColor      string `short:"g" value-name:"<color>" description:"good battery level      green  | 64 " default:"32"`
	MiddleColor    string `short:"m" value-name:"<color>" description:"middle battery level    yellow | 136" default:"33"`
	WarnColor      string `short:"w" value-name:"<color>" description:"warn battery level      red    | 160" default:"31"`
	UpperThreshold int    `short:"u" value-name:"<threshold(%)>" description:"upper threshold" default:"75"`
	LowerThreshold int    `short:"l" value-name:"<threshold(%)>" description:"lower threshold" default:"25"`
}

type Options struct {
	GeneralOption *GeneralOption `group:"general"`
	ColorsOption  *ColorsOption  `group:"colors:                                     tmux:    zsh"`
}

var opts Options

type batteryStatus struct {
	connected  bool
	percentage int
	color      *string
}

// Length of opts.GeneralOption.Ascii_bar
const barLength = 10

// Determine battery charge state
func batteryCharge(battStat *batteryStatus) {
	switch runtime.GOOS {
	case "darwin": // MacOS
		if opts.GeneralOption.PmsetOn {
			acCmd := "pmset -g batt | grep -o 'AC Power'"
			cmd := exec.Command("sh", "-c", acCmd)
			cmd.Run()
			// Battery Connection
			if cmd.ProcessState.ExitCode() == 0 {
				battStat.connected = true
			} else {
				battStat.connected = false
			}
			// Battery Percentage
			battPrcCmd := "pmset -g batt | grep -o '[0-9]*%' | tr -d %"
			out, err := exec.Command("sh", "-c", battPrcCmd).Output()
			if err != nil {
				log.Fatal(err)
			}
			// Convert from byte to int
			battStat.percentage, err = strconv.Atoi(strings.TrimRight(string(out), "\n"))
			if err != nil {
				log.Fatal(err)
			}
		} else {
			// Battery Info by ioreg
			acCmd := "ioreg -n AppleSmartBattery -r | grep -o '\"[^\"]*\" = [^ ]*' | sed -e 's/= //g' -e 's/\"//g' | sort"
			out, err := exec.Command("sh", "-c", acCmd).Output()
			if err != nil {
				log.Fatal(err)
			}
			var ioregInfo map[string]string = map[string]string{}
			scanner := bufio.NewScanner(strings.NewReader(string(out)))
			for scanner.Scan() {
				words := strings.Fields(scanner.Text())
				ioregInfo[words[0]] = words[1]
			}
			// Battery Connection
			if ioregInfo["ExternalConnected"] == "No" {
				battStat.connected = false
			} else {
				battStat.connected = true
			}

			// Battery Percentage
			maxCapacity, hasMaxCapacity := ioregInfo["MaxCapacity"]
			currentCapacity, hasCurrentCapacity := ioregInfo["CurrentCapacity"]
			if hasMaxCapacity && hasCurrentCapacity {
				currentCapacityInt, err := strconv.Atoi(currentCapacity)
				if err != nil {
					log.Fatal(err)
				}
				maxCapacityInt, err := strconv.Atoi(maxCapacity)
				if err != nil {
					log.Fatal(err)
				}
				battStat.percentage = 100 * currentCapacityInt / maxCapacityInt
			} else {
				log.Fatalf("failed to get battery capacity from ioreg")
				os.Exit(-1)
			}
		}
	case "linux":
		log.Fatalf("this version does not yet support linux")
		os.Exit(-1)
	default:
		log.Fatalf("this version does not yet support your OS")
		os.Exit(-1)
	}
}

// Apply the correct color to the battery status prompt
func applyColors(battStat *batteryStatus) {
	if battStat.percentage >= opts.ColorsOption.UpperThreshold {
		battStat.color = &opts.ColorsOption.GoodColor
	} else if battStat.percentage >= opts.ColorsOption.LowerThreshold {
		battStat.color = &opts.ColorsOption.MiddleColor
	} else {
		battStat.color = &opts.ColorsOption.WarnColor
	}
}

// Print the battery status
func printStatus(battStat *batteryStatus) {
	var graph string

	if !opts.GeneralOption.Emoji && battStat.connected {
		graph = "\u26a1"
	} else if opts.GeneralOption.NerdFonts {
		switch {
		case battStat.percentage >= 80:
			graph = "\uf240"
		case battStat.percentage >= 60:
			graph = "\uf241"
		case battStat.percentage >= 40:
			graph = "\uf242"
		case battStat.percentage >= 20:
			graph = "\uf243"
		default:
			graph = "\uf244"
		}
	} else {
		// Get emoji from spark
		sparkCheckCmd := "command -v spark &>/dev/null"
		cmd := exec.Command("sh", "-c", sparkCheckCmd)
		cmd.Run()
		if cmd.ProcessState.ExitCode() == 0 {
			sparkCmd := "spark 0 " + strconv.Itoa(battStat.percentage) + " 100"
			out, err := exec.Command("sh", "-c", sparkCmd).Output()
			if err != nil {
				log.Fatal(err)
			}
			graph = string([]rune(string(out))[1])
		} else {
			opts.GeneralOption.Ascii = true
		}
	}

	if opts.GeneralOption.Ascii {
		// Battery percentage rounded to the lenght of ascii_bar
		roundedN := barLength*battStat.percentage/100 + 1
		if roundedN > 10 {
			roundedN = 10
		}

		// Creates ascii_bar
		graph = "[" + strings.Repeat("=", roundedN) + strings.Repeat(" ", barLength-roundedN) + "]"
	}

	var printfCmd string
	if opts.GeneralOption.OutputTmux {
		// Set colorname in tmux
		opts.ColorsOption.GoodColor = "green"
		opts.ColorsOption.MiddleColor = "yellow"
		opts.ColorsOption.WarnColor = "red"
		printfCmd = "#[fg=" + *battStat.color + "][" + strconv.Itoa(battStat.percentage) + "%%] " + graph + "#[default]"
	} else if opts.GeneralOption.OutputZsh {
		// Set colorname in zsh
		opts.ColorsOption.GoodColor = "64"
		opts.ColorsOption.MiddleColor = "136"
		opts.ColorsOption.WarnColor = "160"
		printfCmd = "%%B%%F{" + *battStat.color + "}[" + strconv.Itoa(battStat.percentage) + "%%%%] " + graph
	} else {
		printfCmd = "\x1b[" + *battStat.color + "m[" + strconv.Itoa(battStat.percentage) + "%%] " + graph + " \x1b[0m\n"
	}
	fmt.Printf(printfCmd)
}

func main() {
	// Read args
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(0)
	}

	battStat := batteryStatus{}

	batteryCharge(&battStat)
	applyColors(&battStat)
	printStatus(&battStat)
}
