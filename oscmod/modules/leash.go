package modules

import (
	"math"

	"github.com/Glowman554/OpenOSC/config"
	"github.com/Glowman554/OpenOSC/oscmod"
	"github.com/Glowman554/OpenOSC/oscmod/chatbox"
	"github.com/hypebeast/go-osc/osc"
)

type LeashModuleContainer struct {
	isWalking bool
	isRunning bool
	isGrabbed bool

	stretch float64

	smoothMoveX float64
	smoothMoveZ float64

	xPos float64
	yPos float64
	zPos float64
	xNeg float64
	yNeg float64
	zNeg float64

	config config.LeashConfig
}

type LeashModule struct {
	container *LeashModuleContainer
}

func NewLeashModule(config config.LeashConfig) LeashModule {
	return LeashModule{
		container: &LeashModuleContainer{
			isWalking:   false,
			isRunning:   false,
			isGrabbed:   true,
			stretch:     0.206,
			smoothMoveX: 0.0,
			smoothMoveZ: 0.0,
			xPos:        0.0,
			yPos:        0.0,
			zPos:        0.672,
			xNeg:        0.0,
			yNeg:        0.0,
			zNeg:        0.672,
			config:      config,
		},
	}
}

func (m LeashModule) Name() string {
	return "OSC Leash"
}

func (m LeashModule) Id() string {
	return "leash"
}

func (m LeashModule) Init(client *osc.Client, dispatcher *osc.StandardDispatcher) error {
	err := dispatcher.AddMsgHandler("/avatar/parameters/Leash_IsGrabbed", func(msg *osc.Message) {
		if grabbed, ok := msg.Arguments[0].(bool); ok {
			m.container.isGrabbed = grabbed
			// if grabbed {
			// 	log.Println("Leash grabbed")
			// } else {
			// 	log.Println("Leash released")
			// }
		}
	})
	if err != nil {
		return err
	}

	err = dispatcher.AddMsgHandler("/avatar/parameters/Leash_Stretch", func(msg *osc.Message) {
		if stretch, ok := msg.Arguments[0].(float32); ok {
			m.container.stretch = float64(stretch)
		}
	})
	if err != nil {
		return err
	}

	// TODO: should not be *
	err = dispatcher.AddMsgHandler("*", func(msg *osc.Message) {
		if val, ok := msg.Arguments[0].(float32); ok {
			switch msg.Address {
			case "/avatar/parameters/Leash_Z+":
				m.container.zPos = float64(val)
			case "/avatar/parameters/Leash_Z-":
				m.container.zNeg = float64(val)
			case "/avatar/parameters/Leash_X+":
				m.container.xPos = float64(val)
			case "/avatar/parameters/Leash_X-":
				m.container.xNeg = float64(val)
			case "/avatar/parameters/Leash_Y+":
				m.container.yPos = float64(val)
			case "/avatar/parameters/Leash_Y-":
				m.container.yNeg = float64(val)
			}
		}
	})
	if err != nil {
		return err
	}

	player := oscmod.NewPlayer(client)
	oscmod.TickFPS(120, func() {
		m.container.UpdateMovement(player)
	})

	return nil
}

func (m LeashModule) Tick(client *osc.Client, chatbox *chatbox.ChatBoxBuilder) error {
	return nil
}

func (c *LeashModuleContainer) UpdateMovement(player *oscmod.Player) {
	c.UpdateMovementState()
	x, y, z := c.CalculateMovement()
	c.ApplyMovement(player, x, y, z)
}

func (c *LeashModuleContainer) UpdateMovementState() {
	// wasWalking := c.isWalking
	// wasRunning := c.isRunning

	if c.isGrabbed {
		c.isWalking = c.stretch > c.config.WalkDeadzone
		c.isRunning = c.stretch > c.config.RunDeadzone
	} else {
		c.isWalking = false
		c.isRunning = false
	}

	// if c.isRunning && !wasRunning {
	// 	log.Printf("Started running (stretch: %f)", c.stretch)
	// } else if !c.isRunning && wasRunning {
	// 	log.Printf("Stopped running")
	// } else if c.isWalking && !wasWalking {
	// 	log.Printf("Started walking (stretch: %f)", c.stretch)
	// } else if !c.isWalking && wasWalking {
	// 	log.Printf("Stopped walking")
	// }
}

func (c *LeashModuleContainer) CalculateMovement() (float64, float64, float64) {
	if !c.isGrabbed {
		c.smoothMoveX *= 0.7
		c.smoothMoveZ *= 0.7
		var x, z float64

		if math.Abs(c.smoothMoveX) < 0.01 {
			x = 0
		} else {
			x = c.smoothMoveX
		}

		if math.Abs(c.smoothMoveZ) < 0.01 {
			z = 0
		} else {
			z = c.smoothMoveZ
		}

		return x, 0, z
	}

	netX := c.xPos - c.xNeg
	netY := c.yNeg - c.yPos
	netZ := c.zPos - c.zNeg

	strength := c.stretch * c.config.StrengthMultiplier

	verticalStretch := math.Abs(netY)

	if verticalStretch >= c.config.UpDownDeadzone && c.config.UpDownCompensation > 0 {
		compensationFactor := 1.0 - (verticalStretch * c.config.UpDownCompensation * 0.5)
		compensationFactor = math.Max(0.1, math.Min(1.0, compensationFactor))
		netX *= compensationFactor
		netZ *= compensationFactor
	}

	netX *= strength
	netZ *= strength

	c.smoothMoveX = c.smoothMoveX*0.7 + netX*0.3
	c.smoothMoveZ = c.smoothMoveZ*0.7 + netZ*0.3

	return c.smoothMoveX, netY, c.smoothMoveZ
}

func (c *LeashModuleContainer) ApplyMovement(player *oscmod.Player, x float64, y float64, z float64) {
	if !c.isGrabbed {
		player.StopRun()
		player.MoveVertical(0)
		player.MoveHorizontal(0)
		return
	}

	if c.isRunning {
		player.Run()
	} else {
		player.StopRun()
	}

	player.MoveVertical(float32(z))
	player.MoveHorizontal(float32(x))

}
