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
	GameData     map[string]interface{}
	player       Player
	enemies      []Enemy
	aliveEnemies int
	screenWidth  int32
	screenHeight int32
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
	Sprite    rl.Texture2D
}

type Turret struct {
	Rectangle rl.Rectangle
	Color     rl.Color
	Range     int
	Damage    int
	FireRate  float32
	Cost      int
}

// gui stuff

var TURRET_DATA map[string]Turret
var ENEMY_DATA map[string]Enemy

func LoadEnemyData() map[string]Enemy {

	image := rl.LoadImage("assets/sprite.png")
	sprite := rl.LoadTextureFromImage(image)
	results := make(map[string]Enemy)
	results["Normal"] = Enemy{Color: rl.Red, Health: 100, Speed: 2, Sprite: sprite}
	results["Fast"] = Enemy{Color: rl.Green, Health: 100, Speed: 5, Sprite: sprite}
	results["Buff"] = Enemy{Color: rl.Yellow, Health: 200, Speed: 2, Sprite: sprite}
	return results
}

func LoadTurretData() map[string]Turret {
	results := make(map[string]Turret)
	results["Basic"] = Turret{Range: 50, Damage: 5, FireRate: 1, Cost: 50}
	results["Strong"] = Turret{Range: 30, Damage: 7, FireRate: 0.5, Cost: 100}
	results["Wide"] = Turret{Range: 100, Damage: 3, FireRate: 1, Cost: 100}
	results["Fast"] = Turret{Range: 100, Damage: 2, FireRate: 2, Cost: 100}
	return results
}

func CreateRoundOneEnemies(startX, startY float32) []Enemy {
	results := []Enemy{}
	var enemy Enemy
	for i := range 10 {
		enemy = CreateEnemy("Normal", float32(i*-60), float32(startY/2))
		results = append(results, enemy)
	}
	enemy = CreateEnemy("Fast", float32(len(results)*-60), float32(startY/2))
	results = append(results, enemy)
	return results

}

func CreateEnemy(enemyType string, x, y float32) Enemy {
	result := ENEMY_DATA[enemyType]
	result.Rectangle = rl.Rectangle{X: x, Y: y, Width: 25, Height: 25}
	result.HealthBar = rl.Rectangle{X: x, Y: y - 10, Width: 48, Height: 5}
	result.Alive = true
	log.Printf("sprite %v", result.Sprite)
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

		}
		if g.enemies[i].Health <= 0 {
			g.enemies[i].Alive = false
		}
	}

}

func (enemy *Enemy) Move() {
	enemy.Rectangle.X = enemy.Rectangle.X + float32(enemy.Speed*1)
	enemy.HealthBar.Width = float32(enemy.Health / 2)
	enemy.HealthBar.X = enemy.HealthBar.X + float32(enemy.Speed*1)
}

func OnClickOpenShopButton(g *Game) {
	g.ActivateScene("Shop")
}

func OnClickShopBasicTurret(g *Game) {
	g.Scenes["PlaceTurrets"].skip = true
	g.ActivateScene("PlaceTurrets")
	g.Scenes["PlaceTurrets"].Data["ChosenTurret"] = TURRET_DATA["Basic"]
}
func OnClickShopFastTurret(g *Game) {
	g.Scenes["PlaceTurrets"].skip = true
	g.ActivateScene("PlaceTurrets")
	g.Scenes["PlaceTurrets"].Data["ChosenTurret"] = TURRET_DATA["Fast"]
}
func OnClickShopStrongTurret(g *Game) {
	g.Scenes["PlaceTurrets"].skip = true
	g.ActivateScene("PlaceTurrets")
	g.Scenes["PlaceTurrets"].Data["ChosenTurret"] = TURRET_DATA["Strong"]
}
func OnClickShopWideTurret(g *Game) {
	g.Scenes["PlaceTurrets"].skip = true
	g.ActivateScene("PlaceTurrets")
	g.Scenes["PlaceTurrets"].Data["ChosenTurret"] = TURRET_DATA["Wide"]
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

func UpdateVictory(g *Game) {

}

func DrawRound(g *Game) {
	for i := range g.player.Turrets {
		rl.DrawRectangleRec(g.player.Turrets[i].Rectangle, g.player.Turrets[i].Color)
	}
	for i := range g.enemies {
		if g.enemies[i].Alive {
			rl.DrawTextureRec(
				g.enemies[i].Sprite,
				rl.Rectangle{
					X:      0,
					Y:      0,
					Width:  48,
					Height: 48,
				},
				rl.Vector2{
					X: g.enemies[i].Rectangle.X,
					Y: g.enemies[i].Rectangle.Y,
				},
				g.enemies[i].Color)
			//			rl.DrawRectangleRec(g.enemies[i].Rectangle, g.enemies[i].Color)
			rl.DrawRectangleRec(g.enemies[i].HealthBar, rl.Red)
		}
	}
}

func UpdateRound(g *Game) {
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

func DrawGameOver(g *Game) {
	rl.DrawText("GAME OVER", g.screenWidth/2-50, g.screenHeight/2, 20, rl.Red)
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		g.player = InitPlayer()
		g.enemies = CreateRoundOneEnemies(float32(g.screenWidth), float32(g.screenHeight/2))
	}

}

func UpdateGameOver(g *Game) {

}

func DrawHUD(g *Game) {
	displayText := fmt.Sprintf("Gold: %v", math.Round(float64(g.player.Gold)))
	rl.DrawText(displayText, 10, 10, 20, rl.Black)
}

func UpdateHUD(g *Game) {
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		if g.WasButtonClicked(&g.Scenes["HUD"].Buttons[0]) {
			g.Scenes["HUD"].Buttons[0].OnClick(g)
		}
	}
	g.DrawButtons(g.Scenes["HUD"].Buttons)
}
func DrawPause(g *Game) {
	log.Printf("?")
	rl.DrawText("PAUSED", g.screenWidth/2, g.screenHeight/2, 20, rl.Black)
	rl.DrawRectangleLines(g.screenWidth/2-100, g.screenHeight/2-100, 200, 200, rl.Black)
}

func UpdatePause(g *Game) {

}

func DrawShop(g *Game) {
	rl.DrawText("Shop", g.screenWidth/2, 100, 20, rl.Black)
	g.DrawButtons(g.Scenes["Shop"].Buttons)
}
func UpdateShop(g *Game) {
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		for _, button := range g.Scenes["Shop"].Buttons {
			if g.WasButtonClicked(&button) {
				button.OnClick(g)
			}

		}
	}
}

func DrawPlaceTurret(g *Game) {
	DrawRound(g)

	var chosenTurret Turret
	if data, ok := g.Scenes["PlaceTurrets"].Data["ChosenTurret"]; ok {
		log.Printf("chosen turret %v", data)
		chosenTurret, _ = data.(Turret)
	}
	mousePosition := rl.GetMousePosition()

	log.Printf("chosen %v", chosenTurret)
	rl.DrawCircle(
		int32(mousePosition.X),
		int32(mousePosition.Y),
		float32(chosenTurret.Range),
		rl.LightGray,
	)

}

func UpdatePlaceTurret(g *Game) {
	if g.Scenes["PlaceTurrets"].skip {
		g.Scenes["PlaceTurrets"].skip = false
		return

	}
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		mousePosition := rl.GetMousePosition()
		g.player.CheckAddTurret(mousePosition.X, mousePosition.Y)
		g.ActivateScene("Round")
	}
}
func main() {

	TURRET_DATA = LoadTurretData()
	g := Game{
		Scenes:       map[string]*Scene{},
		player:       InitPlayer(),
		aliveEnemies: 0,
		screenWidth:  int32(800),
		screenHeight: int32(450),
	}
	g.Scenes["Victory"] = &Scene{
		Active:      false,
		AutoDisable: true,
		DrawScene:   DrawVictory,
		UpdateScene: UpdateVictory,
	}
	g.Scenes["GameOver"] = &Scene{
		Active:      false,
		AutoDisable: true,
		DrawScene:   DrawGameOver,
		UpdateScene: UpdateGameOver,
	}
	g.Scenes["Round"] = &Scene{
		Active:      false,
		AutoDisable: true,
		DrawScene:   DrawRound,
		UpdateScene: UpdateRound,
	}
	g.Scenes["HUD"] = &Scene{
		Active:      true,
		AutoDisable: false,
		DrawScene:   DrawHUD,
		UpdateScene: UpdateHUD,
		Buttons:     make([]Button, 1),
	}
	g.Scenes["HUD"].Buttons[0] = Button{
		Rectangle: rl.Rectangle{X: float32(g.screenWidth) - 110, Y: 10, Width: 100, Height: 15},
		Color:     rl.Blue,
		Text:      "Shop",
		TextColor: rl.Black,
		OnClick:   OnClickOpenShopButton,
	}
	g.Scenes["Pause"] = &Scene{
		Active:      false,
		AutoDisable: true,
		DrawScene:   DrawPause,
		UpdateScene: UpdatePause,
	}
	g.Scenes["Shop"] = &Scene{
		Active:      false,
		AutoDisable: true,
		DrawScene:   DrawShop,
		Buttons:     make([]Button, 4),
		UpdateScene: UpdateShop,
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
	g.Scenes["Shop"].Buttons[1] = Button{
		Rectangle: rl.Rectangle{
			X:      200,
			Y:      100,
			Width:  100,
			Height: 30,
		},
		Color:     rl.SkyBlue,
		Text:      "Create Fast Turret (100)",
		TextColor: rl.Black,
		OnClick:   OnClickShopFastTurret,
	}
	g.Scenes["Shop"].Buttons[2] = Button{
		Rectangle: rl.Rectangle{
			X:      200,
			Y:      150,
			Width:  100,
			Height: 30,
		},
		Color:     rl.SkyBlue,
		Text:      "Create Strong Turret (100)",
		TextColor: rl.Black,
		OnClick:   OnClickShopStrongTurret,
	}
	g.Scenes["Shop"].Buttons[3] = Button{
		Rectangle: rl.Rectangle{
			X:      200,
			Y:      200,
			Width:  100,
			Height: 30,
		},
		Color:     rl.SkyBlue,
		Text:      "Create Wide Turret (100)",
		TextColor: rl.Black,
		OnClick:   OnClickShopWideTurret,
	}
	g.Scenes["PlaceTurrets"] = &Scene{
		Active:      false,
		AutoDisable: true,
		DrawScene:   DrawPlaceTurret,
		UpdateScene: UpdatePlaceTurret,
		Data:        make(map[string]interface{}),
	}

	rl.InitWindow(g.screenWidth, g.screenHeight, "tower defence")

	ENEMY_DATA = LoadEnemyData()
	g.enemies = CreateRoundOneEnemies(float32(800/2), float32(450/2))
	g.aliveEnemies = len(g.enemies)

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
