package main

// go run broadcaster/gen.go OpcProcessManager broadcaster/broadcast.tmpl

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync/atomic"
	"time"

	"github.com/spf13/viper"
)

// OpcProcessesMap is a collection of OPC Process configurations.
type OpcProcessesMap map[string]OpcProcessConfig

// OpcProcessManager handles starting and stopping the OPC processes that animate the LED display.
type OpcProcessManager struct {
	Processes           OpcProcessesMap
	currentName         string
	brightness          int
	stopped             bool
	controlPanelIsOwner int32 // Updated atomically.
	cmd                 *exec.Cmd
	broadcaster         *OpcProcessManagerBroadcaster
}

// OpcProcessConfig is a single config for one of the OPC processes that animates the LED display.
type OpcProcessConfig struct {
	Description string
	Exec        string // The command to execute to run the animation process.
	KillCommand string
	// TODO: Enable turning on output
	// TODO: We must have a kill command also (xvfb-run in docker makes it hard to kill like normal)
}

// AnimationState contains the animations state and what is playing.
type AnimationState struct {
	Playing     int          `json:"playing"`
	PlayingName string       `json:"playing_name"`
	Brightness  int          `json:"brightness"`
	Anims       []serverAnim `json:"anims"`
}

// OpcProcessManagerReceiver is a client that wants to listen to state changes to the OPC Process manager.
type OpcProcessManagerReceiver struct {
	animationStateChan chan AnimationState
	brightnessChan     chan int
}

// NewOpcProcessManagerReceiver creates a new OpcProcessManagerReceiver instance.
func NewOpcProcessManagerReceiver() *OpcProcessManagerReceiver {
	return &OpcProcessManagerReceiver{
		animationStateChan: make(chan AnimationState),
		brightnessChan:     make(chan int),
	}
}

// NewOpcProcessManager creates a new process manager and read the config for it.
func NewOpcProcessManager(broadcaster *OpcProcessManagerBroadcaster) (*OpcProcessManager, error) {
	o := OpcProcessManager{
		broadcaster: broadcaster,
		brightness:  128,
	}
	if err := o.ReadConfig(); err != nil {
		return nil, err
	}
	return &o, nil
}

// IsControlPanelOwner returns if the control panel owns the animation choice. Use this to read the value atomically.
func (o *OpcProcessManager) IsControlPanelOwner() bool {
	return atomic.LoadInt32(&o.controlPanelIsOwner) != 0
}

// SetControlPanelIsOwner sets if the control panel owns the animation choice. This sets the state atomically.
func (o *OpcProcessManager) SetControlPanelIsOwner(isControlPanelOwner bool) {
	if isControlPanelOwner {
		atomic.StoreInt32(&o.controlPanelIsOwner, 1)
	} else {
		atomic.StoreInt32(&o.controlPanelIsOwner, 0)
	}
}

// ReadConfig reads the config needed by the process manager.
func (o *OpcProcessManager) ReadConfig() error {
	// TODO: Get rid of this and pass as argument to NewProcessManager instead.
	opcProcessesStringMap := viper.GetStringMap("processes")

	log.Println("Animation process list:")
	o.Processes = make(OpcProcessesMap)

	for name := range opcProcessesStringMap {
		process := OpcProcessConfig{}

		err := viper.UnmarshalKey(fmt.Sprintf("processes.%v", name), &process)
		if err != nil {
			return fmt.Errorf("failed to read processes from config: %v", err)
		}

		o.Processes[name] = process

		// TODO: Check that these exist on the filesystem and attempt to start them
		// (We do this before we broadcast any OPC messages)

		log.Printf(" %v: %v\n", name, process.Description)
	}

	return nil
}

// stopAnim stops the currently running animation process (if any).
func (o *OpcProcessManager) stopAnim() {
	defer func() {
		o.currentName = ""
	}()

	process := o.Processes[o.currentName]

	o.stopped = true
	if (o.cmd != nil) && (o.cmd.Process != nil) {
		err := o.cmd.Process.Kill()
		if err != nil {
			log.Printf("Failed to kill process '%v': %v\n", o.currentName, err)
			return
		}
		o.cmd = nil

		log.Printf("Killed process: %v\n", o.currentName)
	}

	// TODO: Change this to use https://stackoverflow.com/questions/22470193/why-wont-go-kill-a-child-process-correctly
	// TODO: Go routine?
	// TODO: Timeout command?
	if process.KillCommand != "" {
		killCmd := exec.Command("sh", "-c", process.KillCommand)

		log.Printf("Running Kill command for process '%v': '%v'", o.currentName, process.KillCommand)

		if err := killCmd.Run(); err != nil {
			switch err.(type) {
			default:
				log.Printf("Failed to run kill command '%v' for process '%v': %v\n", process.KillCommand, o.currentName, err)
			case *exec.ExitError:
				log.Printf("Kill command for '%v' returned error (process might already have been killed): %v\n", o.currentName, err)
			}
		} else {
			log.Printf("Kill command for '%v' ran with no errors\n", o.currentName)
		}
	}
}

// broadCastState broadcasts the current state of the animation being played.
func (o *OpcProcessManager) broadcastState() {
	o.broadcaster.Broadcast(func(r *OpcProcessManagerReceiver) {
		r.animationStateChan <- o.GetAnimationsState()
	})
}

// SetBrightness sets the brightness and broadcasts the current state to all websocket clients.
func (o *OpcProcessManager) SetBrightness(brightness int, sender *OpcProcessManagerReceiver) {
	if o.brightness == brightness {
		return
	}
	// TODO: Make an owner and ignore all other clients for ~1s at a time.

	// Save for new clients.
	o.brightness = brightness

	o.broadcaster.Broadcast(func(r *OpcProcessManagerReceiver) {
		if r != sender {
			r.brightnessChan <- o.brightness
		}
	})
}

// PlayAnim starts a given animation process by name. An empty string means stop.
func (o *OpcProcessManager) PlayAnim(processName string) error {
	// Broadcast whatever state change we ended up with.
	defer func() {
		log.Printf("Broadcasting Play of '%v'\n", processName)
		o.broadcastState()
	}()

	log.Println("PlayAnim called")

	// Empty name means to stop.
	if processName == "" {
		o.stopAnim()
		return nil
	}

	if processName == o.currentName {
		return fmt.Errorf("already running %v", processName)
	}

	process, ok := o.Processes[processName]
	if !ok {
		return fmt.Errorf("no animation named '%v' exists", processName)
	}

	o.stopAnim()

	log.Println("Starting process: ", processName)
	log.Println("  ", o.Processes[processName].Exec)

	o.currentName = processName
	o.stopped = false

	// Start the new process.
	go o.runAndMonitorCommand(process)

	return nil
}

// runAndMonitorCommand keeps the command running if it succeeds at least once.
// This is inteded to run in a go routine.
func (o *OpcProcessManager) runAndMonitorCommand(process OpcProcessConfig) {

	args := strings.Split(process.Exec, " ")

	for {
		// TODO: Replace with exec.Command("sh", "-c", process.Exec) instead?
		o.cmd = exec.Command(args[0], args[1:]...)
		if err := o.cmd.Run(); err != nil {
			if o.stopped {
				log.Println(err)
				return
			}
			log.Println("Animation process died unexpectedly restarting...:", err)
			time.Sleep(1 * time.Second)
			continue
		}
	}
}

// GetAnimationsState returns the current animations state.
func (o *OpcProcessManager) GetAnimationsState() AnimationState {
	msg := AnimationState{}
	msg.Playing = -1
	msg.PlayingName = o.currentName
	msg.Brightness = o.brightness
	msg.Anims = make([]serverAnim, len(o.Processes))
	i := 0
	for name, val := range o.Processes {
		msg.Anims[i].Name = name
		msg.Anims[i].Description = val.Description

		if msg.Anims[i].Name == msg.PlayingName {
			msg.Playing = i
		}

		i++
	}

	return msg
}
