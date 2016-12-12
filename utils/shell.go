package utils

import (
	"os"
	"log"
	"os/exec"
	"errors"
	"syscall"
)

func OsExecStdOut(command string,outFile string,args ...string) (error) {
	if len(command) > 0 && len(args) > 0 {
		cmd := exec.Command(command,args...)
		log.Println(cmd.Args)
		if len(outFile) > 0 {
			out, err := os.Create(outFile)
			if err != nil {
				log.Fatal("Cannot create output file for command", command, "at this location", outFile)
				return err
			}
			cmd.Stdout = out
		} else {
			cmd.Stdout = os.Stdout
		}
		cmd.Stderr = os.Stderr
		execErr := cmd.Run()
		if execErr != nil {
			log.Fatal("Execution o" +
				"f", command, "resulted in an error", execErr)
			return execErr
		}

		return nil
	} else {
		return errors.New("Empty command or arguments to execute command")
	}
}

func OsShellExec(command string, args ...string) (error) {
	if len(command) > 0 && len(args) > 0 {
		binary, lookErr := exec.LookPath(command)
		if lookErr != nil {
			log.Fatal("Cannot find ",command ,"binary", lookErr)
			return lookErr
		}
		env := os.Environ()
		execErr := syscall.Exec(binary,args,env)
		if execErr != nil {
			log.Fatal("Error while performing ",command,execErr)
			return execErr
		}
		return nil
	} else {
		return errors.New("Empty command")
	}
}
