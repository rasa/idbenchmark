package idbenchmark_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

const (
	m1 = "INSERT INTO m1 VALUES (NULL)"
	m2 = "REPLACE INTO m2 VALUES (NULL)"
	m3 = "UPDATE m3 SET id=LAST_INSERT_ID(id+1)"
	m4 = "UPDATE m3 SET id=LAST_INSERT_ID(id+1) LIMIT 1"
	i1 = "INSERT INTO i1 VALUES (NULL)"
	i2 = "REPLACE INTO i2 VALUES (NULL)"
	i3 = "UPDATE i3 SET id=LAST_INSERT_ID(id+1)"
	i4 = "UPDATE i3 SET id=LAST_INSERT_ID(id+1) LIMIT 1"
)

var mysqlDSN string = "root:@tcp(127.0.0.1:3306)/idbenchmark"

func init() {
	idbenchmarkDSN := os.Getenv("idbenchmark_DSN")
	if idbenchmarkDSN != "" {
		mysqlDSN = idbenchmarkDSN
	}
}

func mysqlConnect() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", mysqlDSN)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return db, nil
}

func runMysql(b *testing.B, db *sql.DB, s string) {
	stmt, err := db.Prepare(s)
	if err != nil {
		log.Printf("Prepare error: %v: %s", err, s)
		return
	}
	defer stmt.Close()
	runMysqlStmt(b, db, s, stmt)
}

func runMysqlStmt(b *testing.B, db *sql.DB, s string, stmt *sql.Stmt) {
	var res sql.Result
	var err error
	var id int64

	for n := 0; n < b.N; n++ {
		res, err = stmt.Exec()
		if err != nil {
			log.Printf("Exec error: %v: %d: %s", err, n, s)
			break
		}

		id, err = res.LastInsertId()
		if err != nil {
			log.Printf("LastInsertId error: %v", err)
			break
		}
		if id == 0 {
			log.Printf("id=0")
			break
		}
	}
}

func BenchmarkMysqlInsert(b *testing.B) {
	db, err := mysqlConnect()
	if err != nil {
		return
	}
	defer db.Close()

	b.ResetTimer()
	runMysql(b, db, m1)
	b.StopTimer()
}

func BenchmarkMysqlReplace(b *testing.B) {
	db, err := mysqlConnect()
	if err != nil {
		return
	}
	defer db.Close()
	b.ResetTimer()
	runMysql(b, db, m2)
	b.StopTimer()
}

func BenchmarkMysqlUpdate(b *testing.B) {
	db, err := mysqlConnect()
	if err != nil {
		return
	}
	defer db.Close()
	b.ResetTimer()
	runMysql(b, db, m3)
	b.StopTimer()
}

func BenchmarkMysqlUpdateLimit1(b *testing.B) {
	db, err := mysqlConnect()
	if err != nil {
		return
	}
	defer db.Close()
	b.ResetTimer()
	runMysql(b, db, m4)
	b.StopTimer()
}

func BenchmarkInnoDBInsert(b *testing.B) {
	db, err := mysqlConnect()
	if err != nil {
		return
	}
	defer db.Close()
	b.ResetTimer()
	runMysql(b, db, i1)
	b.StopTimer()
}

func BenchmarkInnoDBReplace(b *testing.B) {
	db, err := mysqlConnect()
	if err != nil {
		return
	}
	defer db.Close()
	b.ResetTimer()
	runMysql(b, db, i2)
	b.StopTimer()
}

func BenchmarkInnoDBUpdate(b *testing.B) {
	db, err := mysqlConnect()
	if err != nil {
		return
	}
	defer db.Close()
	b.ResetTimer()
	runMysql(b, db, i3)
	b.StopTimer()
}

func BenchmarkInnoDBUpdateLimit1(b *testing.B) {
	db, err := mysqlConnect()
	if err != nil {
		return
	}
	defer db.Close()
	b.ResetTimer()
	runMysql(b, db, i4)
	b.StopTimer()
}

func BenchmarkMysqlInsertParallel(b *testing.B) {
	db, err := mysqlConnect()
	if err != nil {
		return
	}
	defer db.Close()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			runMysql(b, db, m1)
		}
	})
	b.StopTimer()
}

func BenchmarkMysqlReplaceParallel(b *testing.B) {
	db, err := mysqlConnect()
	if err != nil {
		return
	}
	defer db.Close()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			runMysql(b, db, m2)
		}
	})
	b.StopTimer()
}

func BenchmarkMysqlUpdateParallel(b *testing.B) {
	db, err := mysqlConnect()
	if err != nil {
		return
	}
	defer db.Close()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			runMysql(b, db, m3)
		}
	})
	b.StopTimer()
}

func BenchmarkMysqlUpdateLimit1Parallel(b *testing.B) {
	db, err := mysqlConnect()
	if err != nil {
		return
	}
	defer db.Close()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			runMysql(b, db, m4)
		}
	})
	b.StopTimer()
}

func BenchmarkInnoDBInsertParallel(b *testing.B) {
	db, err := mysqlConnect()
	if err != nil {
		return
	}
	defer db.Close()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			runMysql(b, db, i1)
		}
	})
	b.StopTimer()
}

func BenchmarkInnoDBReplaceParallel(b *testing.B) {
	db, err := mysqlConnect()
	if err != nil {
		return
	}
	defer db.Close()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			runMysql(b, db, i2)
		}
	})
	b.StopTimer()
}

func BenchmarkInnoDBUpdateParallel(b *testing.B) {
	db, err := mysqlConnect()
	if err != nil {
		return
	}
	defer db.Close()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			runMysql(b, db, i3)
		}
	})
	b.StopTimer()
}

func BenchmarkInnoDBUpdateLimit1Parallel(b *testing.B) {
	db, err := mysqlConnect()
	if err != nil {
		return
	}
	defer db.Close()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			runMysql(b, db, i4)
		}
	})
	b.StopTimer()
}
