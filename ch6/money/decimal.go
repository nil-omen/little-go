package money

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Decimal is capable of storing a decimal value such as 30 or 1543.243.
// example: 1.52 = 152 * 10^(-2) will be stored as {152, 2}
type Decimal struct {
	// subunits is the amount of subunits. Multiply it by the precision to get the real value
	subunits int64
	// Number of "subunits" in a unit, expressed as a power of 10.
	precision byte
}

const (
	// ErrInvalidDecimal is returned if the decimal is malformed.
	ErrInvalidDecimal = Error("unable to convert the decimal")

	// ErrTooLarge is returned if the quantity is too large
	// this would cause floating point precision errors.
	ErrTooLarge = Error("quantity over 10^12 is too large")

	// maxDecimal value is a thousand billion, using the short scale -- 10^12.
	maxDecimal = 1e12
)

func ParseDecimal(value string) (Decimal, error) {
	intPart, fracPart, _ := strings.Cut(value, ".")

	subunits, err := strconv.ParseInt(intPart+fracPart, 10, 64)
	if err != nil {
		return Decimal{}, fmt.Errorf("%w: %s", ErrInvalidDecimal, err)
	}

	if subunits > maxDecimal {
		return Decimal{}, ErrTooLarge
	}

	precision := byte(len(fracPart))
	dec := Decimal{subunits: subunits, precision: precision}

	dec.simplify()

	return dec, nil
}

func (d *Decimal) simplify() {
	for d.subunits%10 == 0 && d.precision > 0 {
		d.precision--
		d.subunits /= 10
	}
}

// pow10 is a quick implementation of how to raise 10 to a given power.
// It's optimised for small powers, and slow for unusually high powers.
func pow10(power byte) int64 {
	switch power {
	case 0:
		return 1
	case 1:
		return 10
	case 2:
		return 100
	case 3:
		return 1000
	case 4:
		return 10000
	default:
		return int64(math.Pow(10, float64(power)))
	}
}
