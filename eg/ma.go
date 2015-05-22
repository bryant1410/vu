// Copyright © 2014-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/gazed/vu"
)

// ma, model animation, is an example of loading and animating a model using
// skeletel animation. It is based on the example data provided in the IQM
// Development kit from http://sauerbraten.org/iqm.
func ma() {
	ma := &matag{}
	if err := vu.New(ma, "Model Animation", 400, 100, 800, 600); err != nil {
		log.Printf("ma: error starting engine %s", err)
	}
	defer catchErrors()
}

// Globally unique "tag" that encapsulates example specific data.
type matag struct {
	view   vu.View
	top    vu.Pov
	cam    vu.Camera // 3D model
	ui     vu.Camera // 2D user interface.
	title  vu.Model  // Animation information display.
	names  []string  // All loaded model names.
	models []vu.Pov  // All loaded models.
	model  vu.Pov    // Currently selected model.
	index  int       // Index of currently selected model.
	run    float64   // Camera movement speed.
	spin   float64   // Camera spin speed.
}

// Create is the engine callback for initial asset creation.
func (ma *matag) Create(eng vu.Eng, s *vu.State) {
	ma.run = 10   // move so many cubes worth in one second.
	ma.spin = 270 // spin so many degrees in one second.
	ma.top = eng.Root().NewPov()
	ma.view = ma.top.NewView()
	ma.cam = ma.view.Cam()
	ma.cam.SetPerspective(60, float64(800)/float64(600), 0.1, 50)
	ma.cam.SetLocation(0, 3, 10)

	// load any available IQM/E models. The loaded model data is fed to
	// the animation capable shader "anim".
	for _, modelFile := range ma.modelFiles() {
		pov := ma.top.NewPov()
		pov.SetScale(-1, 1, 1)
		if modelFile == "runner" {
			pov.SetScale(-3, 3, 3) // Runner is a bit small.
		}
		pov.Spin(-90, 0, 0)   // Have the model face the camera.
		pov.SetVisible(false) // Hide initially.

		// Most IQ* files are expected to be animated.
		// Use a "uv" shader to handle IQ* files without animations.
		pov.NewModel("anim").LoadAnim(modelFile)
		ma.models = append(ma.models, pov)
		ma.names = append(ma.names, modelFile)
	}
	ma.model = ma.models[ma.index] // should always have at least one.
	ma.model.SetVisible(true)

	// Have a lighter default background.
	eng.SetColor(0.15, 0.15, 0.15, 1)

	// Create a banner to show the model name.
	top2D := eng.Root().NewPov()
	view2D := top2D.NewView()
	view2D.SetUI()
	ma.ui = view2D.Cam()
	ma.ui.SetOrthographic(0, float64(s.W), 0, float64(s.H), 0, 10)
	title := top2D.NewPov()
	title.SetLocation(10, 5, 0)
	ma.title = title.NewModel("uv").AddTex("weblySleek22White").LoadFont("weblySleek22")
	ma.title.SetPhrase(" ")
}

// Update is the recurring callback to update state based on user actions.
func (ma *matag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	if in.Resized {
		ma.cam.SetPerspective(60, float64(s.W)/float64(s.H), 0.1, 50)
		ma.ui.SetOrthographic(0, float64(s.W), 0, float64(s.H), 0, 10)
	}
	dt := in.Dt
	for press, down := range in.Down {
		switch press {
		case "W":
			ma.cam.Move(0, 0, dt*ma.run, ma.cam.Lookxz())
		case "S":
			ma.cam.Move(0, 0, dt*-ma.run, ma.cam.Lookxz())
		case "A":
			ma.model.Spin(0, 0, 5)
		case "D":
			ma.model.Spin(0, 0, -5)
		case "Tab":
			if down == 1 {

				// switch to the next loaded model.
				ma.model.SetVisible(false)
				ma.index = ma.index + 1
				if ma.index >= len(ma.models) {
					ma.index = 0
				}
				ma.model = ma.models[ma.index]
				ma.model.SetVisible(true)
			}
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			if down == 1 {
				ma.playAnimation(press)
			}
		}
	}
	ma.showAction()
}

// playAnimation chooses an available animation.
// Animations that are not available are ignored.
func (ma *matag) playAnimation(d09 string) {
	action, _ := strconv.Atoi(d09)
	if len(ma.model.Model().Actions()) > action {
		ma.model.Model().Animate(action, 0)
	}
}

// showAction updates the information text, the animation
// file, animation sequence name, and the frame numbers.
func (ma *matag) showAction() {
	if names := ma.model.Model().Actions(); len(names) > 0 {
		index, frame, maxFrames := ma.model.Model().Action()
		name := names[index]
		stats := fmt.Sprintf("[%d] %s %d/%d", index, name, frame, maxFrames)
		ma.title.SetPhrase(ma.names[ma.index] + ":" + stats)
	}
}

// iqmodel groups the 3D assets with the file name of the model file.
type iqmodel struct {
	title string // IQ file name.
	tr    vu.Pov // loaded IQ 3D model.
}

// modelFiles returns the names of the IQE/IQM files in the models directory.
// Only unique base names are returned.
func (ma *matag) modelFiles() []string {
	uniqueNames := map[string]bool{}
	models := []string{}
	files, _ := ioutil.ReadDir("models")
	for _, f := range files {
		name := f.Name()
		if strings.Contains(name, ".iq") && name[0] != '.' {
			base := name[:len(name)-4]
			if _, ok := uniqueNames[base]; !ok {
				uniqueNames[base] = true
				models = append(models, base)
			}
		}
	}
	return models
}
