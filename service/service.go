package service

import (
	"fmt"
	"go_web_demo/model"
	"go_web_demo/util"
	"sync"
)

// 获取用户信息
func ListUser(username string, offset, limit int) ([]*model.UserInfo, uint64, error) {
	infos := make([]*model.UserInfo, 0)

	// 获取用户的基本信息
	users, count, err := model.ListUser(username, offset, limit)
	if err != nil {
		return nil, count, err
	}

	// 遍历信息，获取id
	ids := []uint64{}
	for _, user := range users {
		ids = append(ids, user.Id)
	}

	// 建立同步
	wg := sync.WaitGroup{}
	userList := model.UserList{
		Lock:  new(sync.Mutex),
		IdMap: make(map[uint64]*model.UserInfo, len(users)),
	}

	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	// 多线程获取数据
	for _, u := range users {
		wg.Add(1)
		go func(u *model.UserModel) {
			defer wg.Done()

			shortId, err := util.GenShortId()
			if err != nil {
				errChan <- err
				return
			}

			userList.Lock.Lock()
			defer userList.Lock.Unlock()
			userList.IdMap[u.Id] = &model.UserInfo{
				Id:        u.Id,
				Username:  u.Username,
				SayHello:  fmt.Sprintf("hello %s", shortId),
				Password:  u.Password,
				CreatedAt: u.CreatedAt.Format("2018-05-15 12:12:12"),
				UpdatedAt: u.UpdatedAt.Format("2018-05-15 12:12:12"),
			}
		}(u)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		return nil, count, err
	}

	for _, id := range ids {
		infos = append(infos, userList.IdMap[id])
	}

	return infos, count, nil
}
