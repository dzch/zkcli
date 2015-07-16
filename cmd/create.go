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
		"errors"
		"path"
	    "github.com/samuel/go-zookeeper/zk"
	   )

type createAction struct {}

func init() {
    a := &createAction{}
	cmdAction["create"] = a
}

func (create *createAction) exec(zc *ZkCmd) error {
	if len(zc.cmdArg) < 2 {
		return errors.New("create need path and value")
	}
    path := zc.cmdArg[0]
	if len(zc.opt.Chroot) > 0 {
		if path[len(path)-1] == '/' {
			path = path[0:len(path)-1]
		}
	}
    path = fmt.Sprintf("%s%s", zc.opt.Chroot, path)
	flags := int32(0)
	acl := zk.WorldACL(zk.PermAll)
	var ret string
	var err error
	if zc.opt.Recursive {
		ret, err = create.rCreate(zc, path, []byte(zc.cmdArg[1]), flags, acl)
	} else {
        ret, err = zc.conn.Create(path, []byte(zc.cmdArg[1]), flags, acl)
	}
	if err != nil {
		return err
	}
	err = zc.output(ret)
	if err != nil {
		return err
	}
	return nil
}

func (create *createAction) rCreate(zc *ZkCmd, node string, data []byte, flags int32, acl []zk.ACL) (string, error) {
	ret, err := zc.conn.Create(node, data, flags, acl)
	if err == nil {
		return ret, nil
	}
	if err != zk.ErrNoNode {
		return "", err
	}
	ret, err = create.rCreate(zc, path.Dir(node), []byte("auto created node"), flags, acl)
	if err == nil || err == zk.ErrNodeExists {
		return zc.conn.Create(node, data, flags, acl)
	} else {
		return "", err
	}
}
