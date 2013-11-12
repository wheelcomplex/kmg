package kmgSql

func (db *Db) QueryGrid(query string, args ...interface{}) (output []map[string]string, error error) {
	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	lenColumn := len(columns)
	for rows.Next() {
		rowArray := make([]interface{}, lenColumn)
		//box value with *string, TODO use RawByte?
		for k1, _ := range rowArray {
			var s string
			rowArray[k1] = &s
		}
		if err := rows.Scan(rowArray...); err != nil {
			return nil, err
		}
		rowMap := make(map[string]string)
		for rowIndex, rowName := range columns {
			//unbox value with *string
			rowMap[rowName] = *(rowArray[rowIndex].(*string))
		}
		output = append(output, rowMap)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return
}
