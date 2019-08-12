
/*
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

package hswrapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	var (
		hsServer = "192.168.2.134"
		dbName   = "test"
		tableName = "people_females"
		userName  = "Emily"
	)
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
