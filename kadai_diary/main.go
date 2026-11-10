package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// 日記を保存するファイル名（実行したフォルダ内に作られる）
const diaryFile = "diary.json"

// Entry は日記1件分のデータ
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
}

// loadDiary はファイルから日記の一覧を読み込む
func loadDiary() ([]Entry, error) {
	data, err := os.ReadFile(diaryFile)
	if err != nil {
		// ファイルがまだない初回は、空の一覧から始める
		if errors.Is(err, os.ErrNotExist) {
			return []Entry{}, nil
		}
		// それ以外の読み込みエラーは0件扱いにせず、エラーとして返す
		return nil, err
	}

	var entries []Entry
	// JSONが壊れている場合もエラーとして返す（既存データを上書きしないため）
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// saveDiary は日記の一覧全体をJSONにしてファイルへ保存する
func saveDiary(entries []Entry) error {
	// 人が読んでも分かりやすいように、インデント付きのJSONにする
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(diaryFile, data, 0644)
}

// printUsage は使い方を表示する
func printUsage() {
	fmt.Println("使い方:")
	fmt.Println("  go run main.go add \"タイトル\" \"本文\"  … 日記を追加する")
	fmt.Println("  go run main.go list                  … 日記を全件表示する")
}

func main() {
	// サブコマンドがあるか、先に引数の個数を確認する
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		// add にはタイトルと本文の2つが必要
		if len(os.Args) != 4 {
			fmt.Println("エラー: add にはタイトルと本文を1つずつ指定してください（空白を含む場合は \"\" で囲みます）")
			printUsage()
			os.Exit(1)
		}

		// 既存の一覧を読み込む（失敗したら保存せずに終了）
		entries, err := loadDiary()
		if err != nil {
			fmt.Println("日記の読み込みに失敗しました:", err)
			os.Exit(1)
		}

		// 現在時刻・タイトル・本文で新しい日記を作る
		entry := Entry{
			Timestamp: time.Now(),
			Title:     os.Args[2],
			Content:   os.Args[3],
		}

		// 一覧の末尾に追加して、一覧全体を保存する
		entries = append(entries, entry)
		if err := saveDiary(entries); err != nil {
			fmt.Println("日記の保存に失敗しました:", err)
			os.Exit(1)
		}
		fmt.Println("日記を追加しました:", entry.Title)

	case "list":
		entries, err := loadDiary()
		if err != nil {
			fmt.Println("日記の読み込みに失敗しました:", err)
			os.Exit(1)
		}

		if len(entries) == 0 {
			fmt.Println("日記はまだありません。")
			return
		}

		// 1件ずつ、日時・タイトル・本文を表示する
		for i, e := range entries {
			fmt.Printf("[%d] %s\n", i+1, e.Timestamp.Format("2006-01-02 15:04:05"))
			fmt.Println("タイトル:", e.Title)
			fmt.Println("本文:", e.Content)
			fmt.Println()
		}

	default:
		fmt.Println("エラー: 不明なコマンドです:", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}
