package xconfig

import (
	"github.com/apolloconfig/agollo/v4/component/log"
	"github.com/apolloconfig/agollo/v4/storage"
)

type watcher struct {
	appId string
	s     *appSetter
}

func NewWatcher(s *appSetter, appId string) *watcher {
	return &watcher{
		appId: appId,
		s:     s,
	}
}

//OnChange 增加变更监控
func (w *watcher) OnChange(event *storage.ChangeEvent) {
	if event == nil {
		return
	}

	for key, change := range event.Changes {
		log.Infof("OnChange appId : %s, key : %s, value : %+v\n", w.appId, key, change)
		if change.ChangeType == storage.DELETED {
			if err := w.s.SetValue(key, nil); err != nil {
				log.Errorf("OnChange failed, appId : %s, key : %s, error : %v", w.appId, key, err)
			}
		}
	}
}

//OnNewestChange 监控最新变更
func (w *watcher) OnNewestChange(event *storage.FullChangeEvent) {
	if event == nil {
		return
	}

	for key, val := range event.Changes {
		if err := w.s.SetValue(key, val); err != nil {
			log.Errorf("OnNewestChange failed, appId : %s, key : %s, value : %v, error : %v", w.appId, key, val, err)
		}
	}
}
