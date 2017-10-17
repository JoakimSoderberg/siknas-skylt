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
	isRunning   bool
	killed      chan bool
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

// IsRunning reports if any animation proccess is currently running.
func (o *OpcProcessManager) IsRunning() bool {
	return (o.currentName != "")
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

	// Signal the channel keeping the process alive to chill.
	if o.isRunning {
		close(o.killed)
	}
}

// StartAnim starts a given animation process by name.
func (o *OpcProcessManager) StartAnim(processName string) error {

	if processName == "" {
		o.StopAnim()
		return nil
	}

	if processName == o.currentName {
		return fmt.Errorf("already running %v", processName)
	}

	log.Println("Starting process: ", processName)

	process, ok := o.Processes[processName]
	if !ok {
		return fmt.Errorf("no animation named '%v' exists", processName)
	}

	o.StopAnim()

	o.currentName = processName
	o.isRunning = true

	// Start the new process and monitor it.
	o.killed = make(chan bool)
	go o.runAndMonitorCommand(process)

	return nil
}

// runAndMonitorCommand keeps the command running if it succeeds at least once.
// This is inteded to run in a go routine.
func (o *OpcProcessManager) runAndMonitorCommand(process OpcProcessConfig) {
	defer func() {
		o.isRunning = false
	}()

	args := strings.Split(process.Exec, " ")
	cmd := exec.Command(args[0], args[1:]...)
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Failed to start animation process %v: %v\n", o.currentName, err)
		return
	}

	o.isRunning = true

	// Monitor and restart the process if it dies.
	for {
		select {
		case <-o.killed:
			err := cmd.Process.Kill()
			if err != nil {
				log.Println("Error on kill:", err)
			}

			log.Println("Killed animation process:", o.currentName)
			return
		case <-time.After(time.Second): // TODO: Make this configurable
			if cmd.ProcessState.Exited() {
				log.Println("Process exited, attempting restart:", o.currentName)
				if err := cmd.Run(); err != nil {
					log.Println("Failed to restart process:", err)
				}
			}
		}
	}
}
