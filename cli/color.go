package cli

import "github.com/tristanisham/clr"

type ColorFunc = func(...any) string

var Colors map[string]ColorFunc = map[string]func(...any) string{
	"black":   clr.Black,
	"red":     clr.Red,
	"green":   clr.Green,
	"yellow":  clr.Yellow,
	"blue":    clr.Blue,
	"magenta": clr.Magenta,
	"cyan":    clr.Cyan,
	"white":   clr.White,
}

func (o *OVM) Colored(text string, color string) string {
	cf := Colors[color]
	if cf != nil && o.Config.UseColor {
		return cf(text)
	}

	return text
}
