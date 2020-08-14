package storage

import "github.com/sirupsen/logrus"

// GormLogger is a lgoger instance for gorm
type GormLogger struct{}

// Print prints the gorm log with logrus
func (*GormLogger) Print(v ...interface{}) {
	switch v[0] {
	case "sql":
		logrus.WithFields(logrus.Fields{"module": "gorm", "type": "sql", "rows": v[5], "src_ref": v[1], "values": v[4]}).Print(v[3])
	case "log":
		logrus.WithFields(logrus.Fields{"module": "gorm", "type": "log"}).Print(v[2])
	}
}
