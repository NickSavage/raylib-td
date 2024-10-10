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

type Game struct {
	Scenes       map[string]*Scene
	player       Player
	enemies      []Enemy
	aliveEnemies int
	paused       bool
	screenWidth  int32
	screenHeight int32
	shopMenuOpen bool
	placeTurret  bool
}

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

// gui stuff

type Button struct {
	Rectangle rl.Rectangle
	Color     rl.Color
	Text      string
	TextColor rl.Color
	OnClick   func(*Game)
}
type Scene struct {
	Name        string
	Active      bool
	AutoDisable bool
	DrawScene   func(*Game)
	Buttons     []Button
}

func (g *Game) ActivateScene(sceneName string) {
	for key, scene := range g.Scenes {
		if key == sceneName {
			scene.Active = true
		} else if scene.AutoDisable {
			scene.Active = false
		} else {
			// do nothing
		}
		g.Scenes[key] = scene
	}
}

func (g *Game) WasButtonClicked(button *Button) bool {
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		mousePosition := rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePosition, button.Rectangle) {
			return true
		}
	}
	return false
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
		Range:     50,
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

func (g *Game) checkHits(turret Turret) {
	hasShot := false
	for i := range g.enemies {
		if hasShot {
			return
		}
		if !g.enemies[i].Alive {
			continue
		}
		if rl.CheckCollisionCircleRec(
			rl.Vector2{
				X: turret.Rectangle.X,
				Y: turret.Rectangle.Y,
			},
			float32(turret.Range),
			g.enemies[i].Rectangle,
		) {
			g.enemies[i].Health -= turret.Damage
			hasShot = true
			log.Printf("taking damage! %v", g.enemies[i].Health)

		}
		if g.enemies[i].Health <= 0 {
			g.enemies[i].Alive = false
		}
	}

}

func (enemy *Enemy) Move() {
	enemy.Rectangle.X = enemy.Rectangle.X + float32(enemy.Speed*1)
	enemy.HealthBar.Width = float32(enemy.Health / 4)
	enemy.HealthBar.X = enemy.HealthBar.X + float32(enemy.Speed*1)
}

func OnClickOpenShopButton(g *Game) {
	g.ActivateScene("Shop")
}

func OnClickShopBasicTurret(g *Game) {
	g.ActivateScene("PlaceTurrets")
}

func InitPlayer() Player {
	player := Player{
		Gold:     100,
		Turrets:  make([]Turret, 1),
		GameOver: false,
	}
	return player
}

func DrawVictory(g *Game) {
	rl.DrawText("YOU WIN", g.screenWidth/2-50, g.screenHeight/2, 20, rl.Green)
}

func DrawRound(g *Game) {
	for i := range g.player.Turrets {
		rl.DrawRectangleRec(g.player.Turrets[i].Rectangle, g.player.Turrets[i].Color)
	}
	for i := range g.enemies {
		if g.enemies[i].Alive {
			rl.DrawRectangleRec(g.enemies[i].Rectangle, g.enemies[i].Color)
			rl.DrawRectangleRec(g.enemies[i].HealthBar, rl.Red)
		}
	}
}

func DrawGameOver(g *Game) {
	rl.DrawText("GAME OVER", g.screenWidth/2-50, g.screenHeight/2, 20, rl.Red)
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		g.player = InitPlayer()
		g.enemies = CreateRoundOneEnemies(float32(g.screenWidth), float32(g.screenHeight/2))
	}

}

func DrawHUD(g *Game) {
	displayText := fmt.Sprintf("Gold: %v", math.Round(float64(g.player.Gold)))
	rl.DrawText(displayText, 10, 10, 20, rl.Black)
}

func DrawPause(g *Game) {
	log.Printf("?")
	rl.DrawText("PAUSED", g.screenWidth/2, g.screenHeight/2, 20, rl.Black)
	rl.DrawRectangleLines(g.screenWidth/2-100, g.screenHeight/2-100, 200, 200, rl.Black)
}

func DrawShop(g *Game) {
	rl.DrawText("Shop", g.screenWidth/2, 100, 20, rl.Black)

}

func DrawPlaceTurret(g *Game) {
	DrawRound(g)
	mousePosition := rl.GetMousePosition()

	rl.DrawCircle(int32(mousePosition.X), int32(mousePosition.Y), 50, rl.LightGray)

}

func (g *Game) Draw() {

	rl.BeginDrawing()
	rl.ClearBackground(rl.White)
	for _, scene := range g.Scenes {
		if !scene.Active {
			continue
		}
		scene.DrawScene(g)
		for _, button := range scene.Buttons {
			rl.DrawRectangle(button.Rectangle.ToInt32().X, button.Rectangle.ToInt32().Y, button.Rectangle.ToInt32().Width, button.Rectangle.ToInt32().Height, button.Color)
			rl.DrawText(
				button.Text,
				button.Rectangle.ToInt32().X,
				button.Rectangle.ToInt32().Y,
				10,
				button.TextColor,
			)
		}
	}
	rl.EndDrawing()
}

func (g *Game) Update() {
	if rl.IsKeyPressed(rl.KeyEscape) {
		if g.Scenes["Pause"].Active {
			g.ActivateScene("Round")
		} else {
			g.ActivateScene("Pause")
		}
	}
	if g.Scenes["PlaceTurrets"].Active {
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			mousePosition := rl.GetMousePosition()
			g.player.CheckAddTurret(mousePosition.X, mousePosition.Y)
			g.ActivateScene("Round")
		}
	}
	if g.Scenes["Shop"].Active {
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if g.WasButtonClicked(&g.Scenes["Shop"].Buttons[0]) {
				g.Scenes["Shop"].Buttons[0].OnClick(g)
			}
		}

	}
	if g.Scenes["Round"].Active {
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if g.WasButtonClicked(&g.Scenes["HUD"].Buttons[0]) {
				g.Scenes["HUD"].Buttons[0].OnClick(g)
			}
		}
		g.player.Gold += (1 * GOLD_INCREASE_RATE)

		for i := range g.player.Turrets {
			g.checkHits(g.player.Turrets[i])

		}
		g.aliveEnemies = 0
		for i := range g.enemies {
			if !g.enemies[i].Alive {
				continue
			}
			g.aliveEnemies += 1
			g.enemies[i].Move()
			if g.screenWidth-g.enemies[i].Rectangle.ToInt32().X <= 25 {
				g.ActivateScene("GameOver")
			}
		}
		if g.aliveEnemies == 0 {
			g.ActivateScene("Victory")
		}

	}

}
func main() {

	ENEMY_DATA = LoadEnemyData()
	g := Game{
		Scenes:       map[string]*Scene{},
		player:       InitPlayer(),
		enemies:      CreateRoundOneEnemies(float32(800/2), float32(450/2)),
		aliveEnemies: 0,
		paused:       false,
		screenWidth:  int32(800),
		screenHeight: int32(450),
	}
	g.Scenes["Victory"] = &Scene{
		Active:      false,
		AutoDisable: true,
		DrawScene:   DrawVictory,
	}
	g.Scenes["GameOver"] = &Scene{
		Active:      false,
		AutoDisable: true,
		DrawScene:   DrawGameOver,
	}
	g.Scenes["Round"] = &Scene{
		Active:      false,
		AutoDisable: true,
		DrawScene:   DrawRound,
	}
	g.Scenes["HUD"] = &Scene{
		Active:      true,
		AutoDisable: false,
		DrawScene:   DrawHUD,
		Buttons:     make([]Button, 1),
	}
	g.Scenes["HUD"].Buttons[0] = Button{
		Rectangle: rl.Rectangle{X: float32(g.screenWidth) - 110, Y: 10, Width: 100, Height: 15},
		Color:     rl.Blue,
		Text:      "Test",
		TextColor: rl.Black,
		OnClick:   OnClickOpenShopButton,
	}
	g.Scenes["Pause"] = &Scene{
		Active:      false,
		AutoDisable: true,
		DrawScene:   DrawPause,
	}
	g.Scenes["Shop"] = &Scene{
		Active:      false,
		AutoDisable: true,
		DrawScene:   DrawShop,
		Buttons:     make([]Button, 1),
	}
	g.Scenes["Shop"].Buttons[0] = Button{
		Rectangle: rl.Rectangle{
			X:      200,
			Y:      50,
			Width:  100,
			Height: 30,
		},
		Color:     rl.SkyBlue,
		Text:      "Create Turret (50)",
		TextColor: rl.Black,
		OnClick:   OnClickShopBasicTurret,
	}
	g.Scenes["PlaceTurrets"] = &Scene{
		Active:      false,
		AutoDisable: true,
		DrawScene:   DrawPlaceTurret,
	}
	g.aliveEnemies = len(g.enemies)

	rl.InitWindow(g.screenWidth, g.screenHeight, "tower defence")

	rl.SetTargetFPS(60)
	//	mousePosition := rl.Vector2{X: 0, Y: 0}

	rl.SetExitKey(0)
	g.Scenes["Round"].Active = true
	for !rl.WindowShouldClose() {
		g.Update()
		g.Draw()

	}

	rl.CloseWindow()
}
