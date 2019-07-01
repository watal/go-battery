# go-battery
Display your Laptop's battery status on the terminal.  
Inspired by [Goles/Battery](https://github.com/Goles/Battery)  

## Features
* Displays battery percentage
* Displays battery status with icon and graph bar
* Changes color to reflect battery status (good:Green, medium:Yellow, warning:Red)
* Specify the good, medium and warning battery status color
* Output in zsh or tmux format

## Usage
### MacOS
```
$ ./go-battery [OPTIONS]
```

## Flags
```
general:
  -t                       output tmux status bar format
  -z                       output zsh prompt format
  -e                       don't output the emoji
  -a                       output ascii instead of spark
  -b                       battery path (default: /sys/class/power_supply/BAT0)
  -p                       use pmset (more accurate)
  -n                       use Nerd Fonts battery icon

colors:                                     tmux:    zsh:
  -g=<color>               good battery level      green  | 64  (default: 32)
  -m=<color>               middle battery level    yellow | 136 (default: 33)
  -w=<color>               warn battery level      red    | 160 (default: 31)
  -u=<threshold(%)>        upper threshold (default: 75)
  -l=<threshold(%)>        lower threshold (default: 25)

Help Options:
  -h, --help               Show this help message
```

## License
This software is released  under the MIT License.  

### Thanks
[Goles/Battery](https://github.com/Goles/Battery)  
>The MIT License (MIT)  
>Copyright (c) 2013 Nicolas Goles Domic  
- License under MIT (https://github.com/Goles/Battery/blob/master/LICENSE.txt)
