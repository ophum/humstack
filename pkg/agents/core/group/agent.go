package group

import (
	"time"

	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/client"
	"go.uber.org/zap"
)

type GroupAgent struct {
	client *client.Clients
	logger *zap.Logger
}

func NewGroupAgent(client *client.Clients, logger *zap.Logger) *GroupAgent {
	return &GroupAgent{
		client: client,
		logger: logger,
	}
}

func (a *GroupAgent) Run() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			grList, err := a.client.CoreV0().Group().List()
			if err != nil {
				a.logger.Error(
					"get group list",
					zap.String("msg", err.Error()),
					zap.Time("time", time.Now()),
				)
			}

			for _, group := range grList {
				if group.DeleteState != meta.DeleteStateDelete {
					continue
				}
				// DeleteStateにDeleteが入っていてnamespaceが存在する場合
				// namespaceにDeleteStateをセットする

				nsList, err := a.client.CoreV0().Namespace().List(group.ID)
				if err != nil {
					a.logger.Error(
						"get namespace list",
						zap.String("msg", err.Error()),
						zap.Time("time", time.Now()),
					)
				}

				// 存在しない場合
				if len(nsList) == 0 {
					if err := a.client.CoreV0().Group().Delete(group.ID); err != nil {
						a.logger.Error(
							"delete group",
							zap.String("msg", err.Error()),
							zap.Time("time", time.Now()),
						)
					}
					continue
				}

				for _, ns := range nsList {
					// すでにセットされている場合は更新しない
					if ns.DeleteState == meta.DeleteStateDelete {
						continue
					}
					ns.DeleteState = meta.DeleteStateDelete
					if _, err := a.client.CoreV0().Namespace().Update(ns); err != nil {
						a.logger.Error(
							"update namespace",
							zap.String("msg", err.Error()),
							zap.Time("time", time.Now()),
						)
					}
				}
			}
		}
	}
}
