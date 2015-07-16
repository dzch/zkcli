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
		"errors"
	   )

type lsAction struct {}

func init() {
    a := &lsAction{}
	cmdAction["ls"] = a
	cmdAction["children"] = a
}

func (ls *lsAction) exec(zc *ZkCmd) error {
    path := zc.cmdArg[0]
	if len(zc.opt.Chroot) > 0 && path[0] == '/' {
		if len(path) == 1 {
			path = ""
		} else {
		    path = path[1:]
		}
	}
	children, err := ls.rexec(zc, path)
	if err != nil {
		return err
	}
	err = zc.output(children)
	if err != nil {
		return err
	}
	return nil
}

func (ls *lsAction) rexec(zc *ZkCmd, node string) ([]string, error) {
	var fpath string
	if len(node) == 0 {
		if len(zc.opt.Chroot) == 0 {
			return nil, errors.New("path can not be empty")
		} else {
			fpath = zc.opt.Chroot
		}
	} else {
		if len(zc.opt.Chroot) == 0 {
			fpath = node
		} else {
		    if node[0] == '/' {
                fpath = fmt.Sprintf("%s%s", zc.opt.Chroot, node)
		    } else {
                fpath = fmt.Sprintf("%s/%s", zc.opt.Chroot, node)
		    }
		}
	}
    children, _, err := zc.conn.Children(fpath)
	if err != nil {
		return nil, err
	}
	sort.Sort(sort.StringSlice(children))
	if zc.opt.Recursive == false {
		return children, nil
	}
	rChildren := []string{}
	for _, child := range children {
        cnode := node
		if len(cnode) > 0 && cnode[len(cnode)-1] == '/' {
			if len(cnode) == 1 {
				cnode = ""
			} else {
			    cnode = cnode[0:len(cnode)-1]
			}
		}
		rChildren = append(rChildren, fmt.Sprintf("%s/%s", cnode, child))
        rpath := fmt.Sprintf("%s/%s", cnode, child)
		rchildren, err := ls.rexec(zc, rpath)
		if err != nil {
		    return nil, err
		}
		rChildren = append(rChildren, rchildren...)
	}
	return rChildren, nil
}

