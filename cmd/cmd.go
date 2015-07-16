/*
    The MIT License (MIT)
    
    Copyright (c) 2015 zhouwench zhouwench@gmail.com
    
    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:
    
    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.
    
    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.
*/
package cmd

import (
	"github.com/samuel/go-zookeeper/zk"
	"errors"
	"fmt"
	"time"
)

type CmdOption struct {
	Servers []string
	ZkTimeout time.Duration
	AuthScheme string
	AuthExpression []byte
	OutputFormat string
	Chroot string
	Acls string
	Recursive bool
}

type ZkCmd struct {
	conn *zk.Conn
	opt *CmdOption
	cmd string
	cmdArg []string
}

var cmdAction = make(map[string]Action)

func NewCmdOption() *CmdOption {
	return &CmdOption {
        ZkTimeout: 1*time.Second,
		Acls: "31",
		Recursive: false,
	}
}

func Call(opt *CmdOption, cmdlist []string) error {
    zkCmd := &ZkCmd {
        opt: opt,
		cmd: cmdlist[0],
		cmdArg: cmdlist[1:],
	}
    err := zkCmd.connect()
	if err != nil {
		return errors.New(fmt.Sprintf("connect to zk failed: %s\n", err.Error()))
	}
	defer zkCmd.close()
	action, ok := cmdAction[zkCmd.cmd]
	if !ok {
		return errors.New(fmt.Sprintf("cmd %s not support\n", zkCmd.cmd))
	}
	return action.exec(zkCmd)
}

func (zc *ZkCmd) connect() error {
	conn, _, err := zk.Connect(zc.opt.Servers, zc.opt.ZkTimeout)
	if err == nil && zc.opt.AuthScheme != "" {
		err = conn.AddAuth(zc.opt.AuthScheme, zc.opt.AuthExpression)
	}
	if err != nil {
	    return err
	}
	zc.conn = conn
	return nil
}

func (zc *ZkCmd) close() {
    zc.conn.Close()
}

func (zc *ZkCmd) output(val interface{}) error {
	switch zc.opt.OutputFormat {
		case "json":
			return zc.jsonOutput(val)
		default:
			return zc.defaultOutput(val)
	}
}

func (zc *ZkCmd) jsonOutput(val interface{}) error {
	return errors.New("json output not supported yet")
}

func (zc *ZkCmd) defaultOutput(val interface{}) error {
	var err error = nil
	switch v := val.(type) {
		case string:
			fmt.Println(v)
		case []string:
			for _, str := range v {
				fmt.Println(str)
			}
		default:
			fmt.Println(v)
	}
	return err
}

