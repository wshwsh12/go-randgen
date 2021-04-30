query:
    select
    | select2

select:
    SELECT
        aggregate_function
    FROM _table
    GROUP BY _field

select2:
    SELECT
        aggregate_function
    FROM _table

aggregate_function:
    COUNT(_field) as 'count'
    | SUM(_field) as 'sum'
    | AVG(_field) as 'avg'
    | BIT_AND(_field) as 'bit_and'
    | BIT_OR(_field) as 'bit_or'
    | BIT_XOR(_field) as 'bit_xor'
    | MAX(_field) as 'max'
    | MIN(_field) as 'min'


uncheck:
any_value()
    | VAR_POP(_field) as 'var_pop'
 | VARIANCE(_field) as 'variance'
    | VAR_SAMP(_field) as 'var_samp'