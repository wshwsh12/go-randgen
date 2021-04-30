query:
    aggregate

aggregate:
    SELECT
        aggregate_function
    FROM _table
    group_clause

group_clause:
    /* Empty */
    | GROUP BY group_by_fields

group_by_fields:
    _field
    | _field, _field

aggregate_function:
    COUNT(_field) as 'count'
    | SUM(_field) as 'sum'
    | AVG(_field) as 'avg'
    | BIT_AND(_field) as 'bit_and'
    | BIT_OR(_field) as 'bit_or'
    | BIT_XOR(_field) as 'bit_xor'
    | MAX(_field) as 'max'
    | MIN(_field) as 'min'


not_support:
any_value()
    | VAR_POP(_field) as 'var_pop'
    | VARIANCE(_field) as 'variance'
    | VAR_SAMP(_field) as 'var_samp'