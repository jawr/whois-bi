import React  from 'react'
import { useTable, useSortBy } from 'react-table'

// expose more getters
// https://codesandbox.io/s/github/tannerlinsley/react-table/tree/master/examples/data-driven-classes-and-styles
export const Table = ({ columns, data }) => {
  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    rows,
    prepareRow,
  } = useTable(
    {
      columns,
      data,
    },
    useSortBy
  )

  if (data.length === 0) {
    return <p>No records found</p>
  }

  return (
    <>

      <table {...getTableProps()} className="f6 mw8 center dt--fixed" cellSpacing="0">
        <thead>
          {headerGroups.map(headerGroup => (
            <tr {...headerGroup.getHeaderGroupProps()}>
              {headerGroup.headers.map(column => (
                <th {...column.getHeaderProps({
                  ...column.getSortByToggleProps(),
                  className: column.headerClassName,
                })}
                >
                  {!column.dontSort
                    ? column.isSorted
                      ? column.isSortedDesc
                        ? <span className="sort-by asc"></span>
                        : <span className="sort-by desc"></span>
                        : <span className="sort-by"></span>
                        : <></>
                  }

                  {column.render('Header')}
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody {...getTableBodyProps()} className="lh-copy tl">
          {rows.map(
            (row, i) => {
              prepareRow(row);
              return (
                <tr key={i} {...row.getRowProps()}>
                  {row.cells.map(cell => {
                    return (
                      <td {...cell.getCellProps({
                        className: cell.column.cellClassName,
                      })}>{cell.render('Cell')}</td>
                    )
                  })}
                </tr>
              )}
          )}
        </tbody>
      </table>
    </>
  )
}
