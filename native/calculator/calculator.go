package main

import (
	"errors"
	"math/big"
	"strings"
)

const precision = 16

// Calculate applies one calculator operation and returns the display-safe result.
func Calculate(left, operator, right string) (string, error) {
	l, ok := parseDecimal(left)
	if !ok {
		return "", errors.New("invalid left operand")
	}
	r, ok := parseDecimal(right)
	if !ok {
		return "", errors.New("invalid right operand")
	}

	result := new(big.Float).SetPrec(128).SetMode(big.ToNearestEven)
	switch operator {
	case "+":
		result.Add(l, r)
	case "−", "-":
		result.Sub(l, r)
	case "×", "*":
		result.Mul(l, r)
	case "÷", "/":
		if r.Sign() == 0 {
			return "Cannot divide by 0", nil
		}
		result.Quo(l, r)
	case "%":
		li, lok := integerValue(l)
		ri, rok := integerValue(r)
		if !lok || !rok {
			return "Remainder needs whole numbers", nil
		}
		if ri.Sign() == 0 {
			return "Cannot divide by 0", nil
		}
		return new(big.Int).Mod(li, ri).String(), nil
	default:
		return "", errors.New("unsupported operator")
	}
	return formatDecimal(result), nil
}

func parseDecimal(value string) (*big.Float, bool) {
	value = strings.TrimSpace(value)
	if value == "" || strings.HasPrefix(value, "Cannot") {
		return nil, false
	}
	parsed, _, err := big.ParseFloat(value, 10, 128, big.ToNearestEven)
	return parsed, err == nil
}

func integerValue(value *big.Float) (*big.Int, bool) {
	integer, accuracy := value.Int(nil)
	return integer, accuracy == big.Exact
}

func formatDecimal(value *big.Float) string {
	text := value.Text('f', precision)
	text = strings.TrimRight(text, "0")
	text = strings.TrimRight(text, ".")
	if text == "" || text == "-0" {
		return "0"
	}
	return text
}

func main() {}
