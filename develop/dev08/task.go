package dev08

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

/*
=== Взаимодействие с ОС ===
Необходимо реализовать собственный шелл
встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type process struct {
	name string
	*os.Process
}

type cmd struct {
	name string
	args []string
}

func (c cmd) run() (string, error) {
	switch c.name {
	case "exit":
		os.Exit(0)
		return "", nil
	case "cd":
		return "", cd(c.args)
	case "pwd":
		return os.Getwd()
	case "echo":
		return echo(c.args), nil
	case "kill":
		return "", kill(c.args)
	case "ps":
		return getPs(), nil
	case "exec":
		return "", execute(c.args)
	default:
		return "", errors.New(c.name + ": command not found")
	}
}

func cd(args []string) (err error) {
	var path string
	switch {
	case len(args) == 0:
		path, err = os.UserHomeDir()
		if err != nil {
			return err
		}
	case strings.HasPrefix(args[0], "~"):
		path, err = os.UserHomeDir()
		if err != nil {
			return err
		}
		path = path + args[0][len("~"):]
	default:
		path = args[0]
	}

	return os.Chdir(path)
}

func shellPath() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(path, home) {
		path = "~" + path[len(home):]
	}
	return path, nil
}

func echo(args []string) string {
	return strings.Join(args, " ")
}

func kill(args []string) error {
	if len(args) == 0 {
		return errors.New("pid is not specified")
	}
	ps := make([]process, len(args))
	for i, arg := range args {
		pid, err := strconv.Atoi(arg)
		if err != nil {
			return err
		}
		p, ok := pss[pid]
		if !ok {
			return errors.New("pid " + args[i] + " not found")
		}
		ps[i] = p
	}

	for _, p := range ps {
		err := p.Kill()
		if err != nil {
			return err
		}
	}

	return nil
}

func getPs() string {
	psString := "PID\tCMD\n"
	for id, p := range pss {
		psString += fmt.Sprintf("%d\t%s", id, p.name)
	}
	return psString
}

func execute(args []string) error {
	if len(args) == 0 {
		return errors.New("no command specified")
	}
	cm := exec.Command(args[0], args[1:]...)
	cm.Stdin = os.Stdin
	cm.Stdout = os.Stdout
	cm.Stderr = os.Stderr
	return cm.Run()
}

func execPipeline(ps []cmd) (out string, err error) {
	for _, p := range ps {
		if out != "" {
			p.args = append(p.args, out)
		}
		out, err = p.run()
		if err != nil {
			return "", err
		}
	}
	return out, nil
}

func parseInput(input string) (ps []cmd, err error) {
	pipes := strings.FieldsFunc(input, func(r rune) bool { return r == '|' })
	if len(pipes) == 0 {
		return nil, errors.New("syntax error near unexpected token '|'")
	}
	for _, p := range pipes {
		cmdString := strings.Fields(p)
		switch len(cmdString) {
		case 0:
			return nil, errors.New("syntax error near unexpected token '|'")
		case 1:
			ps = append(ps, cmd{
				name: cmdString[0],
			})
		default:
			ps = append(ps, cmd{
				name: cmdString[0],
				args: cmdString[1:],
			})
		}
	}
	return ps, nil
}

var pss map[int]process

func init() {
	pss = make(map[int]process)
	mPid := os.Getpid()
	mp, err := os.FindProcess(mPid)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	pss[mPid] = process{"shell", mp}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		path, err := shellPath()
		if err != nil {
			fmt.Print("$ ")
		} else {
			fmt.Print(path + " $ ")
		}
		if scanner.Scan() {
			input := scanner.Text()
			pipes, err := parseInput(input)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			out, err := execPipeline(pipes)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			if out != "" {
				fmt.Println(out)
			}
		} else {
			fmt.Println()
			break
		}
	}

}
