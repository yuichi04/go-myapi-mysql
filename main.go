/*
トランザクションを利用する流れ
1. トランザクションを貼る（begin）sql.DB型 の Beginメソッド
2. クエリの実行
3. すべてのクエリ実行が成功した場合には、コミットして結果を確定させる（commit）sql.TxのCommitメソッド
4. 一部のクエリ実行が失敗した場合には、ロールバックして結果を無かったことにする（rollback）sql.TxのRollbackメソッド
*/

package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 接続に使うユーザー・パスワード・データベース名を定義
	dbUser := "docker"
	dbPassword := "docker"
	dbDatabase := "sampledb"

	// データベースに接続するためのアドレス文を定義
	// ここでは"docker:docker@tcp(127.0.0.1"3306)/sampledb?parseTime=true"となる
	dbConn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true", dbUser,
		dbPassword, dbDatabase)

	// Open関数を用いてデータベースに接続
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		fmt.Println(err)
	}
	// プログラムが終了するとき、コネクションが close されるようにする
	defer db.Close()

	// トランザクションの開始
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 現在のいいね数を取得するクエリを実行する
	article_id := 1
	const sqlGetNice = `
		select nice
		from articles
		where article_id = ?;
	`
	row := tx.QueryRow(sqlGetNice, article_id)
	if err := row.Err(); err != nil {
		fmt.Println(err)
		tx.Rollback()
		return
	}

	// 変数nicenumに現在のいいね数を読み込む
	var nicenum int
	err = row.Scan(&nicenum)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return
	}

	// いいね数を+1する更新処理を行う
	const sqlUpdateNice = `update articles set nice = ? where article_id = ?`
	_, err = tx.Exec(sqlUpdateNice, nicenum+1, article_id)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return
	}

	// コミットして処理内容を確定させる
	tx.Commit()
}
