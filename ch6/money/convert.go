package money

// ExchangeRate represents a rate to convert from a currency to another.
// It is a Decimal, to support various precisions.
type ExchangeRate Decimal

// Convert applies the change rate to convert an amount to a target currency.
func Convert(amount Amount, to Currency) (Amount, error) {
	convertedValue := applyExchangeRate(amount, to, ExchangeRate{subunits: 2, precision: 0})

	if err := convertedValue.validate(); err != nil {
		return Amount{}, err
	}

	return convertedValue, nil
}

// applyExchangeRate returns a new Amount representing the input multiplied by the rate.
// The precision of the returned value is that of the target Currency.
// This function does not guarantee that the output amount is supported.
func applyExchangeRate(a Amount, target Currency, rate ExchangeRate) Amount {
	converted := multiply(a.quantity, rate)

	switch {
	case converted.precision > target.precision:
		converted.subunits = converted.subunits / pow10(converted.precision-target.precision)
	case converted.precision < target.precision:
		converted.subunits = converted.subunits * pow10(target.precision-converted.precision)
	}

	converted.precision = target.precision

	return Amount{
		currency: target,
		quantity: converted,
	}
}

// multiply a Decimal with an ExchangeRate and return the product
func multiply(d Decimal, r ExchangeRate) Decimal {
	dec := Decimal{
		subunits:  d.subunits * r.subunits,
		precision: d.precision + r.precision,
	}

	dec.simplify()
	return dec
}
