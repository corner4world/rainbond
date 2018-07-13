// Copyright (C) 2014-2018 Goodrain Co., Ltd.
// RAINBOND, Application Management Platform

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-GPL usage of Rainbond,
// one or multiple Commercial Licenses authorized by Goodrain Co., Ltd.
// must be obtained first.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package taskrun

import (
	"github.com/goodrain/rainbond/node/core/job"
)

//Manager Manager
type Manager interface {
	Start(errchan chan error)
	Stop() error
}

// //Config config server
// type Config struct {
// 	EtcdEndPoints        []string
// 	EtcdTimeout          int
// 	EtcdPrefix           string
// 	ClusterName          string
// 	APIAddr              string
// 	K8SConfPath          string
// 	EventServerAddress   []string
// 	PrometheusMetricPath string
// 	TTL                  int64
// }

//Jobs jobs
type Jobs map[string]*job.Job

// //manager node manager server
// type manager struct {
// 	cluster client.ClusterClient
// 	*cron.Cron
// 	ctx      context.Context
// 	jobs     Jobs // 和结点相关的任务
// 	onceJobs Jobs //记录执行的单任务
// 	jobLock  sync.Mutex
// 	cmds     map[string]*corejob.Cmd
// 	delIDs   map[string]bool
// 	ttl      int64
// }

// //Run taskrun start
// func (n *manager) Start(errchan chan error) {
// 	go n.watchJobs(errchan)
// 	n.Cron.Start()
// 	if err := corejob.StartProc(); err != nil {
// 		logrus.Warnf("[process key will not timeout]proc lease id set err: %s", err.Error())
// 	}
// 	return
// }

// func (n *manager) watchJobs(errChan chan error) error {
// 	watcher := watch.New(store.DefalutClient.Client, "")
// 	watchChan, err := watcher.WatchList(n.ctx, n.Conf.JobPath, "")
// 	if err != nil {
// 		errChan <- err
// 		return err
// 	}
// 	defer watchChan.Stop()
// 	for event := range watchChan.ResultChan() {
// 		switch event.Type {
// 		case watch.Added:
// 			j := new(job.Job)
// 			err := j.Decode(event.GetValue())
// 			if err != nil {
// 				logrus.Errorf("decode job error :%s", err)
// 				continue
// 			}
// 			n.addJob(j)
// 		case watch.Modified:
// 			j := new(job.Job)
// 			err := j.Decode(event.GetValue())
// 			if err != nil {
// 				logrus.Errorf("decode job error :%s", err)
// 				continue
// 			}
// 			n.modJob(j)
// 		case watch.Deleted:
// 			n.delJob(event.GetKey())
// 		default:
// 			logrus.Errorf("watch job error:%v", event.Error)
// 			errChan <- event.Error
// 		}
// 	}
// 	return nil
// }

// //添加job缓存
// func (n *manager) addJob(j *corejob.Job) {
// 	if !j.IsRunOn(n.HostNode) {
// 		return
// 	}
// 	//一次性任务
// 	if j.Rules.Mode != corejob.Cycle {
// 		n.runOnceJob(j)
// 		return
// 	}
// 	n.jobLock.Lock()
// 	defer n.jobLock.Unlock()
// 	n.jobs[j.ID] = j
// 	cmds := j.Cmds(n.HostNode)
// 	if len(cmds) == 0 {
// 		return
// 	}
// 	for _, cmd := range cmds {
// 		n.addCmd(cmd)
// 	}
// 	return
// }

// func (n *manager) delJob(id string) {
// 	n.jobLock.Lock()
// 	defer n.jobLock.Unlock()
// 	n.delIDs[id] = true
// 	job, ok := n.jobs[id]
// 	// 之前此任务没有在当前结点执行
// 	if !ok {
// 		return
// 	}
// 	cmds := job.Cmds(n.HostNode)
// 	if len(cmds) == 0 {
// 		return
// 	}
// 	for _, cmd := range cmds {
// 		n.delCmd(cmd)
// 	}
// 	delete(n.jobs, id)
// 	return
// }

// func (n *manager) modJob(job *corejob.Job) {
// 	if !job.IsRunOn(n.HostNode) {
// 		return
// 	}
// 	//一次性任务
// 	if job.Rules.Mode != corejob.Cycle {
// 		n.runOnceJob(job)
// 		return
// 	}
// 	oJob, ok := n.jobs[job.ID]
// 	// 之前此任务没有在当前结点执行，直接增加任务
// 	if !ok {
// 		n.addJob(job)
// 		return
// 	}
// 	prevCmds := oJob.Cmds(n.HostNode)

// 	job.Count = oJob.Count
// 	*oJob = *job
// 	cmds := oJob.Cmds(n.HostNode)
// 	for id, cmd := range cmds {
// 		n.modCmd(cmd)
// 		delete(prevCmds, id)
// 	}
// 	for _, cmd := range prevCmds {
// 		n.delCmd(cmd)
// 	}
// }

// func (n *manager) addCmd(cmd *corejob.Cmd) {
// 	n.Cron.Schedule(cmd.Rule.Schedule, cmd)
// 	n.cmds[cmd.GetID()] = cmd
// 	logrus.Infof("job[%s] rule[%s] timer[%s] has added", cmd.Job.ID, cmd.Rule.ID, cmd.Rule.Timer)
// 	return
// }

// func (n *manager) modCmd(cmd *corejob.Cmd) {
// 	c, ok := n.cmds[cmd.GetID()]
// 	if !ok {
// 		n.addCmd(cmd)
// 		return
// 	}
// 	sch := c.Rule.Timer
// 	*c = *cmd
// 	// 节点执行时间改变，更新 cron
// 	// 否则不用更新 cron
// 	if c.Rule.Timer != sch {
// 		n.Cron.Schedule(c.Rule.Schedule, c)
// 	}
// 	logrus.Infof("job[%s] rule[%s] timer[%s] has updated", c.Job.ID, c.Rule.ID, c.Rule.Timer)
// }

// func (n *manager) delCmd(cmd *corejob.Cmd) {
// 	delete(n.cmds, cmd.GetID())
// 	n.Cron.DelJob(cmd)
// 	logrus.Infof("job[%s] rule[%s] timer[%s] has deleted", cmd.Job.ID, cmd.Rule.ID, cmd.Rule.Timer)
// }

// //job must be schedulered
// func (n *manager) runOnceJob(j *corejob.Job) {
// 	go j.RunWithRecovery()
// }

// //Stop 停止服务
// func (n *manager) Stop(i interface{}) {
// 	n.Cron.Stop()
// }

// //Newmanager new server
// func Newmanager(cfg *conf.Conf) (*manager, error) {
// 	currentNode, err := GetCurrentNode(cfg)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if cfg.TTL == 0 {
// 		cfg.TTL = 10
// 	}
// 	n := &manager{
// 		Cron:     cron.New(),
// 		jobs:     make(Jobs, 8),
// 		onceJobs: make(Jobs, 8),
// 		cmds:     make(map[string]*corejob.Cmd),
// 		delIDs:   make(map[string]bool, 8),
// 		ttl:      cfg.TTL,
// 	}
// 	return n, nil
// }