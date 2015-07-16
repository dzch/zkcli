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
		"sort"
	   )

type getAction struct {}

func init() {
	cmdAction["get"] = &getAction{}
}

func (get *getAction) exec(zc *ZkCmd) error {
    path := zc.cmdArg[0]
	if len(zc.opt.Chroot) > 0 {
		if path[len(path)-1] == '/' {
			path = path[0:len(path)-1]
		}
	}
    path = fmt.Sprintf("%s%s", zc.opt.Chroot, path)
	if zc.opt.Recursive {
		return get.rGet(zc, path)
	}
    val, _, err := zc.conn.Get(path)
	if err != nil {
		return err
	}
	err = zc.output(string(val))
	if err != nil {
		return err
	}
	return nil
}

func (get *getAction) rGet(zc *ZkCmd, node string) error {
	children, _, err := zc.conn.Children(node)
	if err != nil {
		return err
	}
    cnode := node
	if node == "/" {
		cnode = ""
	} else if cnode[len(node)-1] == '/' {
		cnode = cnode[0:len(node)-1]
	}
	sort.Sort(sort.StringSlice(children))
	for _, child := range children {
		if child[0] == '/' {
			if len(child) == 1 {
				child = ""
			} else {
				child = child[1:]
			}
		}
		err = get.rGet(zc, fmt.Sprintf("%s/%s", cnode, child))
		if err != nil {
			return err
		}
	}
	val, _, err := zc.conn.Get(node)
	if err != nil {
		return err
	}
	return zc.output(fmt.Sprintf("%s: %s", node, string(val)))
}
