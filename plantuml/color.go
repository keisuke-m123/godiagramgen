package plantuml

import (
	"fmt"
	"strconv"
	"strings"
)

type Color struct {
	r uint8
	g uint8
	b uint8
	a uint8
}

func ParseHexColor(h string) (*Color, error) {
	h = strings.ReplaceAll(h, "#", "")
	if len(h) != 6 && len(h) != 8 {
		return nil, fmt.Errorf("received string invalid length")
	}
	colorHex, err := strconv.ParseUint(h, 16, 32)
	if err != nil {
		return nil, fmt.Errorf("received string cannot be converted to hexadecimal: %w", err)
	}

	if len(h) == 6 {
		r := uint8(colorHex >> 16)
		g := uint8((colorHex >> 8) & 0xFF)
		b := uint8(colorHex & 0xFF)
		a := uint8(255)
		return &Color{r: r, g: g, b: b, a: a}, nil
	} else {
		r := uint8(colorHex >> 24)
		g := uint8((colorHex >> 16) & 0xFF)
		b := uint8((colorHex >> 8) & 0xFF)
		a := uint8(colorHex & 0xFF)
		return &Color{r: r, g: g, b: b, a: a}, nil
	}
}

func (c *Color) HexRGBA() string {
	hexColor := (uint64(c.r) << 24) | (uint64(c.g) << 16) | (uint64(c.b) << 8) | uint64(c.a)
	return fmt.Sprintf("%08s", strconv.FormatUint(hexColor, 16))
}
