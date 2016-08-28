package orm

import (
	"strconv"
)

func getTableEntry(table string) interface{} {
	var entry interface{}

	switch table {
	case T_ADDRESS:
		{
			entry = Address{}
		}
	case T_REASON:
		{
			entry = Reason{}
		}
	case T_DURATION:
		{
			entry = Duration{}
		}
	case T_BLACKLIST:
		{
			entry = Blacklist{}
		}
	case T_WHITELIST:
		{
			entry = Whitelist{}
		}
	case T_HISTORY:
		{
			entry = History{}
		}
	}

	return entry
}

func GetByInt(table string, field string, query int64) interface{} {
	var entry interface{}
	var err error

	entry = getTableEntry(table)
	err = Db.Model(&entry).Where("?.? = ?'", table, field, query).Select()
	if err != nil {
		Log.Warning("[orm]: GetByInt(" + table + "," + field + "," + strconv.Itoa(int(query)) + ") failed: " + err.Error())
	}

	return entry
}

func GetByString(table string, field string, query string) interface{} {
	var entry interface{}
	var err error

	entry = getTableEntry(table)
	err = Db.Model(&entry).Where("?.? = '?'", table, field, query).Select()
	if err != nil {
		Log.Warning("[orm]: GetByString(" + table + "," + field + "," + query + ") failed: " + err.Error())
	}

	return entry
}

func GetByStringLike(table string, field string, query string) interface{} {
	var entry interface{}
	var err error

	entry = getTableEntry(table)
	err = Db.Model(&entry).Where("?.? LIKE '%?%'", table, field, query).Select()
	if err != nil {
		Log.Warning("[orm]: GetByStringLike(" + table + "," + field + "," + query + ") failed: " + err.Error())
	}

	return entry
}
