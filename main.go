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

type Player struct {
	Gold     float32
	Turrets  []Turret
	GameOver bool
}

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

var ENEMY_DATA []Enemy

func LoadEnemyData() []Enemy {
	results := []Enemy{
		Enemy{Color: rl.Red, Health: 100, Speed: 2},
		Enemy{Color: rl.Green, Health: 100, Speed: 5},
		Enemy{Color: rl.Yellow, Health: 200, Speed: 2},
	}
	return results
}

func CreateRoundOneEnemies(startX, startY float32) []Enemy {
	results := []Enemy{}
	var enemy Enemy
	for i := range 10 {
		enemy = CreateEnemy(0, float32(i*-40), float32(startY/2))
		results = append(results, enemy)
	}
	enemy = CreateEnemy(2, float32(len(results)*-40), float32(startY/2))
	results = append(results, enemy)
	return results

}

func CreateEnemy(enemyType int, x, y float32) Enemy {
	result := ENEMY_DATA[enemyType]
	result.Rectangle = rl.Rectangle{X: x, Y: y, Width: 25, Height: 25}
	result.HealthBar = rl.Rectangle{X: x, Y: y - 10, Width: 25, Height: 5}
	result.Alive = true
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

func (player *Player) CheckAddTurret(x, y float32) {
	if player.Gold >= 50 {
		turret := CreateTurret(x, y)
		player.Gold -= 50
		player.Turrets = append(player.Turrets, turret)
	}
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

func InitPlayer() Player {
	player := Player{
		Gold:     100,
		Turrets:  make([]Turret, 1),
		GameOver: false,
	}
	return player
}

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)
	paused := false
	aliveEnemies := 0

	ENEMY_DATA = LoadEnemyData()
	rl.InitWindow(screenWidth, screenHeight, "tower defence")

	rl.SetTargetFPS(60)
	//	mousePosition := rl.Vector2{X: 0, Y: 0}

	player := InitPlayer()
	enemies := CreateRoundOneEnemies(float32(screenWidth), float32(screenHeight/2))
	aliveEnemies = len(enemies)
	rl.SetExitKey(0)
	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeyEscape) {
			paused = !paused
		}
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			mousePosition := rl.GetMousePosition()
			player.CheckAddTurret(mousePosition.X, mousePosition.Y)
			//			paused = !paused
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.White)
		if aliveEnemies == 0 {
			rl.DrawText("YOU WIN", screenWidth/2-50, screenHeight/2, 20, rl.Green)

			rl.EndDrawing()
			continue
		}
		if player.GameOver {
			rl.DrawText("GAME OVER", screenWidth/2-50, screenHeight/2, 20, rl.Red)
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				player = InitPlayer()
				enemies = CreateRoundOneEnemies(float32(screenWidth), float32(screenHeight/2))
			}
			rl.EndDrawing()
			continue

		}
		if paused {
			rl.DrawText("PAUSED", screenWidth/2, screenHeight/2, 20, rl.Black)

		} else {
			player.Gold += (1 * GOLD_INCREASE_RATE)
			for i := range player.Turrets {
				rl.DrawRectangleRec(player.Turrets[i].Rectangle, player.Turrets[i].Color)
				player.Turrets[i].checkHits(enemies)

			}
			aliveEnemies = 0
			for i := range enemies {
				if !enemies[i].Alive {
					continue
				}
				aliveEnemies += 1
				enemies[i].Move()
				if screenWidth-enemies[i].Rectangle.ToInt32().X <= 25 {
					player.GameOver = true
				}
				rl.DrawRectangleRec(enemies[i].Rectangle, enemies[i].Color)
				rl.DrawRectangleRec(enemies[i].HealthBar, rl.Red)
			}

		}
		displayText := fmt.Sprintf("Gold: %v", math.Round(float64(player.Gold)))
		rl.DrawText(displayText, 10, 10, 20, rl.Black)
		rl.EndDrawing()

	}

	rl.CloseWindow()
}
