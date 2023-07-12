package tools

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
	"strings"

	"github.com/shinesyang/common"
)

// 实时输出每行结果
func CmdAndChangeDirToFile(dir, commandName string, params []string) error {
	cmd := exec.Command(commandName, params...)
	//StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
	stdout, err := cmd.StdoutPipe() // 读取标准输出
	if err != nil {                 // StdoutPipe 错误返回
		return err
	}
	stderr, err := cmd.StderrPipe() // 读取错误输出
	if err != nil {
		return err
	}

	//cmd.Stderr = os.Stderr			// 这里是将错误输出输入到标准错误里面(注释掉,后面有写入到文件)
	cmd.Dir = dir
	err = cmd.Start()
	if err != nil {
		return err // Start启动错误返回
	}

	//实时循环读取输出流中的一行内容
	go func() {
		//创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
		readerStdout := bufio.NewReader(stdout)
		for {
			line, _, err := readerStdout.ReadLine()
			if err != nil || io.EOF == err {
				break
			}

			common.Logger.Info(string(line))
		}
	}()

	go func() {
		//创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
		readerStderr := bufio.NewReader(stderr)
		for {
			line, _, err := readerStderr.ReadLine()
			if err != nil || io.EOF == err {
				break
			}
			common.Logger.Error(string(line))
		}
	}()

	err = cmd.Wait()
	return err
}

//直接执行返回结果
func CmdAtonceResult(commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()

	return strings.TrimSpace(stdout.String()) + strings.TrimSpace(stderr.String()), err
}
