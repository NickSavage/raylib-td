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
	frames       int
	elapsedTime  int
	round        int
	totalRounds  int
	rounds       map[int][]Enemy
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
	Speed     float32
	Alive     bool
	Sprite    Tile
}

type Turret struct {
	Rectangle   rl.Rectangle
	Color       rl.Color
	Range       int
	Damage      int
	FireRate    float32
	Cost        int
	FireCounter float32
}

type Round struct {
	Round   int
	Enemies []Enemy
}

// gui stuff

var TURRET_DATA map[string]Turret
var ENEMY_DATA map[string]Enemy

func (g *Game) LoadAssets() {
	g.Scenes["Round"].Data["GrassTile"] = Tile{
		Texture: rl.LoadTexture("assets/grass.png"),
		TileFrame: rl.Rectangle{
			X:      0,
			Y:      80,
			Width:  30,
			Height: 30,
		},
		Color: rl.White,
	}
	g.Scenes["Round"].Data["DirtTile"] = Tile{
		Texture: rl.LoadTexture("assets/dirt.png"),
		TileFrame: rl.Rectangle{
			X:      15,
			Y:      80,
			Width:  30,
			Height: 30,
		},
		Color: rl.White,
	}
	image := rl.LoadImage("assets/tower.png")
	rl.ImageResize(image, 60, 78)
	g.Scenes["Round"].Data["TowerTile"] = Tile{
		Texture: rl.LoadTextureFromImage(image),
		TileFrame: rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  60,
			Height: 78,
		},
		Color: rl.White,
	}
}

func (g *Game) CreateRoundOneEnemies() []Enemy {
	row := float32(g.screenHeight/2 - 50)
	results := []Enemy{}
	var enemy Enemy
	for i := range 10 {
		enemy = CreateEnemy("Normal", float32(i*-40), row)
		results = append(results, enemy)
	}
	results = append(results, CreateEnemy("Fast", float32(11*-40), row))
	results = append(results, CreateEnemy("Fast", float32(12*-40), row))
	results = append(results, CreateEnemy("Fast", float32(13*-40), row))
	results = append(results, CreateEnemy("Normal", float32(14*-40), row))
	results = append(results, CreateEnemy("Normal", float32(15*-40), row))
	results = append(results, CreateEnemy("Normal", float32(16*-40), row))
	return results

}
func (g *Game) LoadRounds() {
	g.rounds[1] = g.CreateRoundOneEnemies()
}

func LoadEnemyData() map[string]Enemy {

	tile := Tile{
		Texture: rl.LoadTexture("assets/sprite.png"),
		TileFrame: rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  48,
			Height: 48,
		},
		Color: rl.Red,
	}
	results := make(map[string]Enemy)
	results["Normal"] = Enemy{Color: rl.Red, Health: 100, Speed: 1, Sprite: tile}
	results["Fast"] = Enemy{Color: rl.Green, Health: 100, Speed: 2, Sprite: tile}
	results["Buff"] = Enemy{Color: rl.Yellow, Health: 200, Speed: 0.5, Sprite: tile}
	return results
}

func LoadTurretData() map[string]Turret {
	results := make(map[string]Turret)
	results["Basic"] = Turret{Range: 150, Damage: 5, FireRate: 3, Cost: 50}
	results["Strong"] = Turret{Range: 100, Damage: 10, FireRate: 5, Cost: 100}
	results["Wide"] = Turret{Range: 200, Damage: 3, FireRate: 3, Cost: 50}
	results["Fast"] = Turret{Range: 200, Damage: 2, FireRate: 1, Cost: 100}
	return results
}

func (g *Game) PrepareNextRound() {
	g.enemies = g.rounds[g.round]

}

func CreateEnemy(enemyType string, x, y float32) Enemy {
	result := ENEMY_DATA[enemyType]
	result.Rectangle = rl.Rectangle{X: x, Y: y, Width: 25, Height: 25}
	result.HealthBar = rl.Rectangle{X: x, Y: y - 10, Width: 48, Height: 5}
	result.Alive = true
	log.Printf("sprite %v", result.Sprite)
	return result
}

func CreateTurret(turretType string, x, y float32) Turret {
	data := TURRET_DATA[turretType]
	result := Turret{
		Rectangle: rl.Rectangle{X: x, Y: y, Width: 25, Height: 25},
		Color:     rl.Blue,
		Range:     data.Range,
		Damage:    data.Damage,
	}
	return result
}

func (g *Game) CheckAddTurret(chosenTurret Turret, x, y float32) {
	if g.player.Gold >= float32(chosenTurret.Cost) {
		chosenTurret.Rectangle.X = x
		chosenTurret.Rectangle.Y = y
		log.Printf("create turret x %v y %v", x, y)
		g.player.Gold -= float32(chosenTurret.Cost)
		g.player.Turrets = append(g.player.Turrets, chosenTurret)
	}
}

func (g *Game) checkHits(turret *Turret) {
	hasShot := false
	var newEnemies []Enemy
	if turret.FireCounter != 0 {
		// can't fire this round
		turret.FireCounter -= 1
		return
	} else {
		// reset
		turret.FireCounter = turret.FireRate
	}
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
		} else {
			newEnemies = append(newEnemies, g.enemies[i])
		}
	}
	g.enemies = newEnemies
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
		Turrets:  make([]Turret, 0),
		GameOver: false,
	}
	return player
}

func DrawVictory(g *Game) {
	DrawRound(g)
	rl.DrawText("YOU WIN", g.screenWidth/2-50, g.screenHeight/2, 20, rl.Green)
}

func UpdateVictory(g *Game) {

}

func DrawRound(g *Game) {
	grassTile := g.Scenes["Round"].Data["GrassTile"].(Tile)
	for x := range (g.screenWidth / 30) + 1 {
		for y := range g.screenHeight / 30 {
			DrawTile(grassTile, float32(x*30), float32(y*30))
		}
	}
	dirtTile := g.Scenes["Round"].Data["DirtTile"].(Tile)
	for x := range (g.screenWidth / 30) + 1 {
		DrawTile(dirtTile, float32(x*30), float32(g.screenHeight/2-60))
		DrawTile(dirtTile, float32(x*30), float32(g.screenHeight/2-30))
		DrawTile(dirtTile, float32(x*30), float32(g.screenHeight/2))
	}

	towerTile := g.Scenes["Round"].Data["TowerTile"].(Tile)
	for i := range g.player.Turrets {
		//		log.Printf("draw turret %v", g.player.Turrets[i])
		DrawTile(
			towerTile,
			g.player.Turrets[i].Rectangle.X,
			g.player.Turrets[i].Rectangle.Y,
		)
	}
	for i := range g.enemies {
		if g.enemies[i].Alive {
			DrawTile(
				g.enemies[i].Sprite,
				g.enemies[i].Rectangle.X,
				g.enemies[i].Rectangle.Y,
			)
			rl.DrawRectangleRec(g.enemies[i].HealthBar, rl.Red)
		}
	}
	DrawHUD(g)
}

func UpdateRound(g *Game) {

	if rl.IsKeyPressed(rl.KeyS) {
		g.ActivateScene("Shop")
		return
	}

	g.frames += 1
	if g.frames == 60 {
		g.frames = 0
		g.elapsedTime += 1
	}
	g.player.Gold += (1 * GOLD_INCREASE_RATE)

	for i := range g.player.Turrets {
		g.checkHits(&g.player.Turrets[i])

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
	if len(g.enemies) == 0 {
		if g.round == g.totalRounds {
			g.ActivateScene("Victory")
			return
		}
		g.round += 1
		g.PrepareNextRound()

	}

}

func DrawGameOver(g *Game) {
	DrawRound(g)
	rl.DrawText("GAME OVER", g.screenWidth/2-50, g.screenHeight/2, 20, rl.Red)

}

func UpdateGameOver(g *Game) {

}

func DrawHUD(g *Game) {
	displayText := fmt.Sprintf("Gold: %v", math.Round(float64(g.player.Gold)))
	rl.DrawText(displayText, 10, 10, 20, rl.Black)
	rl.DrawText(fmt.Sprintf("Time: %v", g.elapsedTime), 200, 10, 20, rl.Black)
	rl.DrawText(fmt.Sprintf("Round: %v", g.round), 300, 10, 20, rl.Black)

	g.DrawButtons(g.Scenes["HUD"].Buttons)
}

func UpdateHUD(g *Game) {
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		if g.WasButtonClicked(&g.Scenes["HUD"].Buttons[0]) {
			g.Scenes["HUD"].Buttons[0].OnClick(g)
		}
	}
}
func DrawPause(g *Game) {
	DrawRound(g)
	rl.DrawText("PAUSED", g.screenWidth/2-50, g.screenHeight/2, 20, rl.Black)
}

func UpdatePause(g *Game) {

}

func DrawShop(g *Game) {
	DrawRound(g)
	g.DrawButtons(g.Scenes["Shop"].Buttons)
	// basic
	rl.DrawText("Damage: 5", 200, 130, 12, rl.Black)
	rl.DrawText("Fire Rate: 1", 200, 150, 12, rl.Black)
	rl.DrawText("Range: 150", 200, 170, 12, rl.Black)

	//Fast
	rl.DrawText("Damage: 2", 200, 230, 12, rl.Black)
	rl.DrawText("Fire Rate: 2", 200, 250, 12, rl.Black)
	rl.DrawText("Range: 200", 200, 270, 12, rl.Black)

	//Strong
	rl.DrawText("Damage: 10", 400, 130, 12, rl.Black)
	rl.DrawText("Fire Rate: 0.5", 400, 150, 12, rl.Black)
	rl.DrawText("Range: 100", 400, 170, 12, rl.Black)

	//Wide
	rl.DrawText("Damage: 3", 400, 230, 12, rl.Black)
	rl.DrawText("Fire Rate: 1", 400, 250, 12, rl.Black)
	rl.DrawText("Range: 200", 400, 270, 12, rl.Black)

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
		chosenTurret, _ = data.(Turret)
	}
	mousePosition := rl.GetMousePosition()

	rl.DrawCircleLines(
		int32(mousePosition.X),
		int32(mousePosition.Y),
		float32(chosenTurret.Range),
		rl.Black,
	)
	towerTile := g.Scenes["Round"].Data["TowerTile"].(Tile)
	log.Printf("tower %v", towerTile)
	DrawTile(
		towerTile,
		float32(mousePosition.X)-(towerTile.TileFrame.Width/2),
		float32(mousePosition.Y)-(towerTile.TileFrame.Height/2),
	)

}

func UpdatePlaceTurret(g *Game) {
	if g.Scenes["PlaceTurrets"].skip {
		g.Scenes["PlaceTurrets"].skip = false
		return

	}
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {

		var chosenTurret Turret
		if data, ok := g.Scenes["PlaceTurrets"].Data["ChosenTurret"]; ok {
			chosenTurret, _ = data.(Turret)
		}
		mousePosition := rl.GetMousePosition()
		towerTile := g.Scenes["Round"].Data["TowerTile"].(Tile)
		g.CheckAddTurret(chosenTurret, mousePosition.X-(towerTile.TileFrame.Width/2), mousePosition.Y-(towerTile.TileFrame.Height/2))
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
		frames:       0,
		elapsedTime:  0,
		round:        1,
		rounds:       make(map[int][]Enemy, 1),
		totalRounds:  1,
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
		Data:        make(map[string]interface{}),
	}
	g.Scenes["HUD"] = &Scene{
		Active:      true,
		AutoDisable: false,
		DrawScene:   DrawHUD,
		UpdateScene: UpdateHUD,
		Buttons:     make([]Button, 1),
	}
	g.Scenes["HUD"].Buttons[0] = Button{
		Rectangle: rl.Rectangle{X: float32(g.screenWidth) - 110, Y: 10, Width: 100, Height: 20},
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
			Y:      100,
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
			Y:      200,
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
			X:      400,
			Y:      100,
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
			X:      400,
			Y:      200,
			Width:  100,
			Height: 30,
		},
		Color:     rl.SkyBlue,
		Text:      "Create Wide Turret (50)",
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
	g.LoadRounds()
	g.PrepareNextRound()
	g.aliveEnemies = len(g.enemies)

	g.LoadAssets()

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
