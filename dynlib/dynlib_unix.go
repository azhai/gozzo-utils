// +build darwin linux

package dynlib

import (
	"plugin"
	"sync"
)

var registry sync.Map

//加载插件
func LoadPlugin(path string) *plugin.Plugin {
	if plug, ok := registry.Load(path); ok {
		return plug.(*plugin.Plugin)
	}
	plug, err := plugin.Open(path)
	if err != nil {
		return nil
	}
	registry.Store(path, plug)
	return plug
}

//加载对象或方法
func LoadSymbol(path, name string) plugin.Symbol {
	plug := LoadPlugin(path)
	if plug == nil {
		return nil
	}
	symb, err := plug.Lookup(name)
	if err != nil {
		return nil
	}
	return symb
}
