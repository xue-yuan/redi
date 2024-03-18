package constants

import "github.com/jackc/pgx/v5"

func TxOptions() pgx.TxOptions {
	return pgx.TxOptions{
		IsoLevel:       pgx.Serializable, // 設置隔離等級為 Serializable
		AccessMode:     pgx.ReadWrite,    // 設置為讀寫模式
		DeferrableMode: pgx.Deferrable,   // 不指定僅讀模式的資料庫名稱
	}
}
