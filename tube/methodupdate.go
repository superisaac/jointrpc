package tube

import ()

// join watchers
func (self *Router) WatchMethods() *MethodWatcher {
	self.lock("WatchMethods")
	defer self.unlock("WatchMethods")

	watcher := NewMethodWatcher()
	self.watchers = append(self.watchers, watcher)
	return watcher
}

// leave watchers
func (self *Router) UnwatchMethods(watcher *MethodWatcher) bool {
	self.lock("UnwatchMethods")
	defer self.unlock("UnwatchMethods")

	index := -1
	for i, aLis := range self.watchers {
		if aLis == watcher {
			index = i
			break
		}
	}
	if index >= 0 {
		self.watchers = append(
			self.watchers[:index],
			self.watchers[index+1:]...)
		return true
	} else {
		return false
	}
}

func (self *Router) NotifyMethodUpdate() {
	self.routerLock.RLock()
	defer self.routerLock.RUnlock()

	self.notifyMethodUpdate()
}

func (self *Router) notifyMethodUpdate() {
	if len(self.watchers) > 0 {
		update := MethodUpdate{Methods: self.getLocalMethods()}
		for _, watcher := range self.watchers {
			watcher.ChUpdate <- update
		}
	}
}

func NewMethodWatcher() *MethodWatcher {
	watcher := new(MethodWatcher)
	watcher.ChUpdate = make(chan MethodUpdate, 1000)
	return watcher
}
