package main

import (
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
	"log"
	"math"
)

var time int
var gold float32

const GOLD_INCREASE_RATE float32 = 0.1

type Enemy struct {
	Rectangle rl.Rectangle
	Color     rl.Color
	Health    int
	HealthBar rl.Rectangle
	Speed     int
	Alive     bool
}

type Turret struct {
	Rectangle rl.Rectangle
	Color     rl.Color
	Range     int
	Damage    int
}

func CreateEnemy(x, y float32) Enemy {
	result := Enemy{
		Rectangle: rl.Rectangle{X: x, Y: y, Width: 25, Height: 25},
		Color:     rl.Red,
		Health:    100,
		HealthBar: rl.Rectangle{X: x, Y: y - 10, Width: 25, Height: 5},
		Speed:     2,
		Alive:     true,
	}
	return result
}

func CreateTurret(x, y float32) Turret {
	result := Turret{
		Rectangle: rl.Rectangle{X: x, Y: y, Width: 25, Height: 25},
		Color:     rl.Blue,
		Range:     20,
		Damage:    5,
	}
	return result
}

func (turret Turret) checkHits(enemies []Enemy) {
	hasShot := false
	for i := range enemies {
		if hasShot {
			return
		}
		if !enemies[i].Alive {
			continue
		}
		distance := int(math.Abs(float64(enemies[i].Rectangle.X - turret.Rectangle.X)))
		if distance <= turret.Range {
			enemies[i].Health -= turret.Damage
			hasShot = true
			log.Printf("taking damage! %v", enemies[i].Health)
		}
		if enemies[i].Health <= 0 {
			enemies[i].Alive = false
		}
	}

}

func (enemy *Enemy) Move() {
	enemy.Rectangle.X = enemy.Rectangle.X + float32(enemy.Speed*1)
	enemy.HealthBar.Width = float32(enemy.Health / 4)
	enemy.HealthBar.X = enemy.HealthBar.X + float32(enemy.Speed*1)

}

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "tower defence")

	rl.SetTargetFPS(60)
	//	mousePosition := rl.Vector2{X: 0, Y: 0}

	enemies := make([]Enemy, 10)
	turrets := make([]Turret, 1)

	for i := range 10 {
		enemies[i] = CreateEnemy(float32(i*-40), float32(screenHeight/2))
	}
	turrets[0] = CreateTurret(float32(screenWidth/2), float32(screenHeight/2+50))
	for !rl.WindowShouldClose() {
		log.Printf("gold %v", gold)
		gold = gold + (1 * GOLD_INCREASE_RATE)
		log.Printf("gold %v", gold)
		// if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		// 	mousePosition = rl.GetMousePosition()
		// }

		displayText := fmt.Sprintf("Gold: %v", gold)
		rl.BeginDrawing()
		for i := range turrets {
			rl.DrawRectangleRec(turrets[i].Rectangle, turrets[0].Color)
			turrets[i].checkHits(enemies)

		}
		for i := range enemies {
			if !enemies[i].Alive {
				continue
			}
			enemies[i].Move()
			rl.DrawRectangleRec(enemies[i].Rectangle, enemies[i].Color)
			rl.DrawRectangleRec(enemies[i].HealthBar, rl.Red)
		}
		rl.ClearBackground(rl.White)
		rl.DrawText(displayText, 10, 10, 20, rl.Black)
		rl.EndDrawing()

	}

	rl.CloseWindow()
}
