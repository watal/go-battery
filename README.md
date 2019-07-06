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
  -t                                      output tmux status bar format
  -z                                      output zsh prompt format
  -e                                      don't output the emoji
  -a                                      output ascii instead of spark
  -b=<path>                               battery path (default: /sys/class/power_supply/BAT0)
  -p                                      use pmset (more accurate)
  -n                                      use Nerd Fonts battery icon
  -i={Num(%),Num(%),Num(%),Num(%)}        specify icon's threshold (default: 80, 60, 40, 20)

colors:                                                           default:  tmux:    zsh:
  -g=<color>                              good battery level      1;32    | green  | 64
  -m=<color>                              middle battery level    1;32    | yellow | 136
  -w=<color>                              warn battery level      0;31    | red    | 160
  -u=<threshold(%)>                       upper threshold (default: 75)
  -l=<threshold(%)>                       lower threshold (default: 25)

Help Options:
  -h, --help                              Show this help message
```

## License
This software is released  under the MIT License.  

### Thanks
[Goles/Battery](https://github.com/Goles/Battery)  
>The MIT License (MIT)  
>Copyright (c) 2013 Nicolas Goles Domic  
- License under MIT (https://github.com/Goles/Battery/blob/master/LICENSE.txt)
