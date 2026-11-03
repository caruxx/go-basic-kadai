package main

import (
	"errors"
	"fmt"
	"math"
)

// バリデーションエラーの情報を格納する独自のエラー型
type ValidationError struct {
	Message string  // エラーメッセージ
	Field   string  // 不正なフィールド（高さ・幅・半径のいずれか）
	Value   float64 // 不正な値（図形の長さ）
}

// Error()メソッドを持つことで、ValidationErrorはerrorインターフェースを満たす
func (e *ValidationError) Error() string {
	return fmt.Sprintf("[バリデーションエラー] %s：%s（Value：%.2f）", e.Field, e.Message, e.Value)
}

// 図形の面積・周長を求めるインターフェース
type Shape interface {
	Area() float64
	Perimeter() float64
}

// 長方形
type Rectangle struct {
	Width  float64
	Height float64
}

// 円
type Circle struct {
	Radius float64
}

// 長方形の面積（幅×高さ）
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// 長方形の周長（(幅+高さ)×2）
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// 円の面積（半径×半径×円周率）
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// 円の周長（2×円周率×半径）
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// 図形の長さが正しいかを確認する（OKならnil、NGならValidationErrorを返す）
func validateShape(s Shape) error {
	// 型スイッチで図形の種類ごとにチェック内容を切り替える
	switch v := s.(type) {
	case Rectangle:
		if v.Width <= 0 {
			return &ValidationError{Message: "幅は正の数を指定してください。", Field: "Rectangle.Width", Value: v.Width}
		}
		if v.Height <= 0 {
			return &ValidationError{Message: "高さは正の数を指定してください。", Field: "Rectangle.Height", Value: v.Height}
		}
	case Circle:
		if v.Radius <= 0 {
			return &ValidationError{Message: "半径は正の数を指定してください。", Field: "Circle.Radius", Value: v.Radius}
		}
	default:
		return fmt.Errorf("未対応の図形です：%T", s)
	}
	return nil
}

func main() {
	// 処理対象の図形スライスを作成
	shapes := []Shape{
		Rectangle{Width: 10.0, Height: 5.0},
		Circle{Radius: 3.0},
		Rectangle{Width: -2.0, Height: 5.0}, // 不正（幅が不正）
		Circle{Radius: -1.0},                // 不正（半径が不正）
		Rectangle{Width: 1.0, Height: 0.0},  // 不正（高さが不正）
	}

	// バリデーションを通過した図形だけを入れる計算対象スライス
	calcShapes := []Shape{}

	fmt.Println("--- 図形のバリデーション ---")
	for i, s := range shapes {
		err := validateShape(s)
		if err == nil {
			fmt.Printf("#%d：%T => バリデーションOK\n", i+1, s)
			calcShapes = append(calcShapes, s)
			continue
		}

		// errors.As()でエラーがValidationError型かどうかを調べ、そうなら中身を取り出す
		var vErr *ValidationError
		if errors.As(err, &vErr) {
			fmt.Printf("#%d：%T => NG：%v\n", i+1, s, err)
			fmt.Printf("  [詳細] Field=%s, Value=%.2f\n", vErr.Field, vErr.Value)
		} else {
			fmt.Printf("#%d：%T => NG（想定外のエラー）：%v\n", i+1, s, err)
		}
	}

	fmt.Println()
	fmt.Println("--- 検証済み図形の一括計算 ---")
	for _, s := range calcShapes {
		fmt.Printf("図形：%T｜面積：%.2f ｜周長：%.2f\n", s, s.Area(), s.Perimeter())
	}
}
