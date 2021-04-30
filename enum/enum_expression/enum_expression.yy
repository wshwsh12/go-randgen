query:
    select

select:
    SELECT
        _field as field_a,
        _field as field_b
    FROM _table
    WHERE function

function:
    arith_function

binary_function:
    (_field and _field)
    | (_field or _field)
    | (_field xor _field)
    | (_field & _field)
    | (_field _or _field)
    | (_field ^ _field)
    | (_field << _field)
    | (_field >> _field)

cmp_function:
    | (_field > _field)
    | (_field >= _field)
    | (_field < _field)
    | (_field <= _field)
    | (_field <> _field)
    | (_field != _field)
    | (_field = _field)
    | (_field <=> _field)
    | _field IN (_field)
    | (_field is null)
    | (_field is true)
    | (_field is false)

arith_function:
    (_field + _field)
    | (_field - _field)
    | (_field * _field)
    | (_field / _field)

unary_function:
    (not _field)
    | (~ _field)
    | (- _field)
    | abs(_field)


need_test:
    like

