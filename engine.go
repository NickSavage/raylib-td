package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

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
	UpdateScene func(*Game)
	Buttons     []Button
	skip        bool
	Data        map[string]interface{}
}

type Tile struct {
	Texture   rl.Texture2D
	TileFrame rl.Rectangle
	Color     rl.Color
}

func DrawTile(t Tile, x, y float32) {

	rl.DrawTextureRec(
		t.Texture,
		t.TileFrame,
		rl.Vector2{
			X: x,
			Y: y,
		},
		t.Color,
	)

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

func (g *Game) DrawButtons(buttons []Button) {
	for _, button := range buttons {
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

func (g *Game) WasButtonClicked(button *Button) bool {
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		mousePosition := rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePosition, button.Rectangle) {
			return true
		}
	}
	return false
}

func (g *Game) Draw() {

	rl.BeginDrawing()
	rl.ClearBackground(rl.White)
	for _, scene := range g.Scenes {
		if !scene.Active {
			continue
		}
		scene.DrawScene(g)
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
	for _, scene := range g.Scenes {
		if !scene.Active {
			continue
		}
		scene.UpdateScene(g)
	}
}
