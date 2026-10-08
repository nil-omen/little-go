package money

import (
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
	ErrInvalidDecimal = Error("unable to convert the decial")

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
		return Decimal{}, err
	}

	if subunits > maxDecimal {
		return Decimal{}, ErrTooLarge
	}

	precision := byte(len(fracPart))
	dec := Decimal{subunits: subunits, precision: precision}

	return dec, nil
}

func (d *Decimal) simplify() {
	for d.subunits%10 == 0 && d.precision > 0 {
		d.precision--
		d.subunits /= 10
	}
}
