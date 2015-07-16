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
		"fmt"
	   )

type rmAction struct {}

func init() {
    a := &rmAction{}
	cmdAction["rm"] = a
	cmdAction["delete"] = a
}

func (rm *rmAction) exec(zc *ZkCmd) error {
    path := zc.cmdArg[0]
	if len(zc.opt.Chroot) > 0 {
		if path[len(path)-1] == '/' {
			path = path[0:len(path)-1]
		}
	}
    path = fmt.Sprintf("%s%s", zc.opt.Chroot, path)
	version := int32(-1)
	var err error
	if zc.opt.Recursive {
		err = rm.rDelete(zc, path, version)
	} else {
        err = zc.conn.Delete(path, version)
	}
	if err != nil {
		return err
	}
	err = zc.output(true)
	if err != nil {
		return err
	}
	return nil
}

func (rm *rmAction) rDelete(zc *ZkCmd, node string, version int32) error {
	children, _, err := zc.conn.Children(node)
	if err != nil {
		return err
	}
	for _, child := range children {
		err = rm.rDelete(zc, fmt.Sprintf("%s/%s", node, child), version)
		if err != nil {
			return err
		}
	}
	return zc.conn.Delete(node, version)
}
