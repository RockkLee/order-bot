from decimal import Decimal, ROUND_HALF_UP

_DECIMAL_PLACES = 2
_Q = Decimal("1").scaleb(-_DECIMAL_PLACES)
_SCALE = Decimal(10) ** _DECIMAL_PLACES


def to_float(scaled_val: int) -> float:
    return float(Decimal(scaled_val) / _SCALE)


def to_scaled_val(f: float) -> int:
    d = Decimal(str(f)).quantize(_Q, rounding=ROUND_HALF_UP)
    return int((d * _SCALE).to_integral_value(rounding=ROUND_HALF_UP))


if __name__ == '__main__':
    print(to_scaled_val(0.1))
    print(to_float(499))