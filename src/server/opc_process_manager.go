package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

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
	// TODO: Enable turning on output
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
	o.cmd = exec.Command(args[0], args[1:]...)

	for {
		if err := o.cmd.Run(); err != nil {
			if o.stopped {
				log.Println(err)
				return
			}
			log.Println("Animation process died unexpectedly restarting...:", err)
			continue
		}
	}
}