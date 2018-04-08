package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type OpcProcessesMap map[string]OpcProcessConfig

type OpcProcessManager struct {
	Processes   OpcProcessesMap
	currentName string
	stopped     bool
	cmd         *exec.Cmd
}

type OpcProcessConfig struct {
	Description string
	Exec        string
	KillCommand string
	// TODO: Enable turning on output
	// TODO: We must have a kill command also (xvfb-run in docker makes it hard to kill like normal)
}

// NewOpcProcessManager creates a new process manager and read the config for it.
func NewOpcProcessManager() (*OpcProcessManager, error) {
	o := OpcProcessManager{}
	if err := o.ReadConfig(); err != nil {
		return nil, err
	}
	return &o, nil
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

// StopAnim stops the currently running animation process (if any).
func (o *OpcProcessManager) StopAnim() {
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

// StartAnim starts a given animation process by name.
func (o *OpcProcessManager) StartAnim(processName string) error {

	// Empty name means to stop.
	if processName == "" {
		o.StopAnim()
		return nil
	}

	if processName == o.currentName {
		return fmt.Errorf("already running %v", processName)
	}

	process, ok := o.Processes[processName]
	if !ok {
		return fmt.Errorf("no animation named '%v' exists", processName)
	}

	o.StopAnim()

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
