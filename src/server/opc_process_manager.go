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
	cmd         *exec.Cmd
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
func (o *OpcProcessManager) StopAnim() error {
	defer func() {
		o.cmd = nil
	}()

	o.currentName = ""

	if o.cmd == nil {
		return nil
	}

	// Signal the channel keeping the process alive to chill.
	close(o.killed)

	err := o.cmd.Process.Kill()
	if err != nil {
		log.Println("Error on kill:", err)
		return err
	}

	log.Println("Killed process:", o.currentName)

	return nil
}

// StartAnim starts a given animation process by name.
func (o *OpcProcessManager) StartAnim(processName string) error {

	if processName == "" {
		o.StopAnim()
		return nil
	}

	log.Println("Starting process: ", processName)

	process, ok := o.Processes[processName]
	if !ok {
		return fmt.Errorf("no animation process named %v exists", processName)
	}

	o.StopAnim()

	// Start the new process and monitor it.
	o.killed = make(chan bool)
	args := strings.Split(process.Exec, " ")
	o.cmd = exec.Command(args[0], args[1:]...)
	//o.cmd.Stdout = os.Stdout
	//o.cmd.Stderr = os.Stderr

	go o.runAndMonitorCommand()

	return nil
}

// runAndMonitorCommand keeps the command running if it succeeds at least once.
// This is inteded to run in a go routine.
func (o *OpcProcessManager) runAndMonitorCommand() {

	if err := o.cmd.Run(); err != nil {
		log.Printf("Failed to start animation process %v: %v\n", o.currentName, err)
		o.StopAnim()
		return
	}

	// Monitor and restar the process if it dies.
	for {
		select {
		case <-o.killed:
			return
		case <-time.After(time.Second): // TODO: Make this configurable
			if o.cmd.ProcessState.Exited() {
				log.Println("Process exited, attempting restart:", o.currentName)
				if err := o.cmd.Run(); err != nil {
					log.Println("Failed to restart process:", err)
				}
			}
		}
	}
}
