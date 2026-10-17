package main

import "fmt"

type Todo struct {
	ID        int
	Title     string
	Completed bool
}

// ToDoを完了状態にする（ポインタレシーバで中身を書き換える）
func (t *Todo) Complete() {
	fmt.Printf("ID：%dのToDoを完了に更新します\n", t.ID)
	t.Completed = true
}

// ToDoリストを一覧表示する
func printTodos(todos []Todo) {
	for _, todo := range todos {
		status := "未完了"
		if todo.Completed {
			status = "完　了"
		}
		fmt.Printf("[%s]（ID：%d）%s\n", status, todo.ID, todo.Title)
	}
}

func main() {
	todos := []Todo{}
	todos = append(todos, Todo{ID: 1, Title: "学習計画", Completed: false})
	todos = append(todos, Todo{ID: 2, Title: "環境構築", Completed: false})
	todos = append(todos, Todo{ID: 3, Title: "基礎文法", Completed: false})

	fmt.Println("--- ToDoリスト（初期状態） ---")
	printTodos(todos)
	fmt.Println()

	todos[0].Complete()
	todos[1].Complete()
	fmt.Println()

	fmt.Println("--- ToDoリスト（最終状態） ---")
	printTodos(todos)
}
