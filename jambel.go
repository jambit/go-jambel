package jambel

const (
	// colour modules
	GREEN  = 3
	YELLOW = 2
	RED    = 1

	// light module status
	OFF           = 0
	ON            = 1
	BLINK         = 2
	FLASH         = 3
	BLINK_INVERSE = 4
)

type Jambel interface {

	// Reset resets Jambel to all lights off
	Reset()

	// On switches [colour] on where colour is one of GREEN, RED or YELLOW
	On(colour int)
	Off(colour int)
	Blink(colour int)
	BlinkInverse(colour int)
	Flash(colour int)
	SetAll(green, yellow, red int)
}
