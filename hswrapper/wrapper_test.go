package hswrapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
go test -v --run=Wrap

CREATE TABLE `people_females`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL DEFAULT '',
  `height` float NULL DEFAULT NULL,
  `birth` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `name`(`name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8;

INSERT INTO `people_females` VALUES (1, 'Alice', 168, '1981-01-08 00:00:00');
INSERT INTO `people_females` VALUES (2, 'Candy', 165, '1983-03-15 00:00:00');
INSERT INTO `people_females` VALUES (3, 'Emily', 171, '1985-05-01 00:00:00');
INSERT INTO `people_females` VALUES (4, 'Grace', 175, '1987-07-05 00:00:00');
*/

var (
	hsServer  = "127.0.0.1"
	dbName    = "db_test"
	tableName = "people_females"
	userName  = "Emily"
)

func TestWrapSelect(t *testing.T) {
	var index *HandlerSocketIndex
	hs := NewWrapper(hsServer, 9998, 9999)
	defer hs.Close()

	columns4 := []string{"id", "name", "height", "birth"}
	index = hs.WrapIndex(dbName, tableName, "", columns4...)

	rows, _ := index.FindAll(2, 0, "<=", "3")
	assert.Len(t, rows, 2)
	assert.Equal(t, "Emily", rows[0].Data["name"])
	assert.Equal(t, "Candy", rows[1].Data["name"])

	columns3 := []string{"name", "height", "birth"}
	index = hs.WrapIndex(dbName, tableName, "name", columns3...)
	row, _ := index.FindOne("=", userName)
	assert.Equal(t, "Emily", row.Data["name"])
	assert.Equal(t, "1985-05-01 00:00:00", row.Data["birth"])
}

func _TestWrapUpdate(t *testing.T) {
	var index *HandlerSocketIndex
	hs := NewWrapper(hsServer, 9998, 9999)
	defer hs.Close()

	columns4 := []string{"id", "name", "height", "birth"}
	index = hs.WrapIndex(dbName, tableName, "", columns4...)

	n, err := index.Update(1, "=", []interface{}{"Emily"}, 3, "Emily", 167, "1985-06-06 00:00:00")
	assert.NoError(t, err)
	assert.Equal(t, 1, n)

	columns3 := []string{"name", "height", "birth"}
	index = hs.WrapIndex(dbName, tableName, "name", columns3...)
	row, _ := index.FindOne("=", userName)
	assert.Equal(t, "Emily", row.Data["name"])
	assert.Equal(t, 167, row.Data["height"])
	assert.Equal(t, "1985-06-06 00:00:00", row.Data["birth"])
}

func TestWrapDelete(t *testing.T) {
	var index *HandlerSocketIndex
	hs := NewWrapper(hsServer, 9998, 9999)
	defer hs.Close()

	columns4 := []string{"id", "name", "height", "birth"}
	index = hs.WrapIndex(dbName, tableName, "", columns4...)
	n, err := index.Delete(5, ">=", []interface{}{1})
	assert.NoError(t, err)
	assert.Equal(t, 4, n)
}

func TestWrapInsert(t *testing.T) {
	var index *HandlerSocketIndex
	hs := NewWrapper(hsServer, 9998, 9999)
	defer hs.Close()

	columns4 := []string{"id", "name", "height", "birth"}
	index = hs.WrapIndex(dbName, tableName, "", columns4...)

	var err error
	for i := 1; i <= 4; i++ {
		err = index.Insert(i, "test", 123, "2000-01-01 00:00:00")
		assert.NoError(t, err)
	}
}
