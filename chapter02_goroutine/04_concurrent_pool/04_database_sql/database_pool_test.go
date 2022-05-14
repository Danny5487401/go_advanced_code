package _4_database_sql

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func BenchmarkConnectMySQL(b *testing.B) {
	db, err := sql.Open("mysql", "root:chuanzhi@tcp(tencent.danny.games:3306)/ticket_system?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai")
	if err != nil {
		b.Fatalf("mysql connect error : %s", err.Error())
		return
	}
	defer db.Close() // 注意这行代码要写在上面err判断的下面

	ctx := context.Background()

	db.SetMaxIdleConns(0)
	b.ResetTimer()

	b.Run("noConnPool", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = db.ExecContext(ctx, "SELECT 1")
		}
	})

	db.SetMaxIdleConns(4)
	b.ResetTimer()

	b.Run("hasConnPool", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = db.ExecContext(ctx, "SELECT 1")
		}
	})
}
