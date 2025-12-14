PARTITION BY
    split returned rows into group
ORDER BY
    sort the items of each group that has been split by PARTITION BY

Window functions are used:
- to process rows returned after WHERE, GROUP BY, HAVING
- to sort the rows under a partion
- to sum up the rows

It will preserve the rows, order, return to the original row. The query's outer ORDER BY still decides the final ouptput order

Before window function,

id,employee,salary,month
```
1,Alice,200,1
2,Bob,300,1
3,Alice,200,2
4,Bob,200,2
```

After window function applied:
```sql
    SUM(salary) OVER
        (PARTITION BY employee)
```

```csv
1,Alice,200,1,400
2,Bob,300,1,500
3,Alice,200,2,400
4,Bob,200,2,500
```

to get Running Total,
```sql
    SUM(salary) OVER
        (PARTITION BY employee ORDER BY month)
```
1,Alice,200,1,200
2,Bob,300,1,300
3,Alice,200,2,400
4,Bob,200,2,500

ROW_NUMBER() without ORDER BY is meaningless, will end up with random number.
