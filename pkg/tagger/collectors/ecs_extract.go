// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// +build docker

package collectors

import (
	"context"
	"time"

	"github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/tagger/utils"
	"github.com/DataDog/datadog-agent/pkg/util/containers"
	v1 "github.com/DataDog/datadog-agent/pkg/util/ecs/metadata/v1"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

func (c *ECSCollector) parseTasks(ctx context.Context, tasks []v1.Task, targetDockerID string, containerHandlers ...func(ctx context.Context, containerID string, tags *utils.TagList) error) ([]*TagInfo, error) {
	var output []*TagInfo
	now := time.Now()
	for _, task := range tasks {
		// We only want to collect tasks without a STOPPED status.
		if task.KnownStatus == "STOPPED" {
			continue
		}
		for _, container := range task.Containers {
			// Only collect new containers + the targeted container, to avoid empty tags on race conditions
			if c.expire.Update(container.DockerID, now) || container.DockerID == targetDockerID {
				tags := utils.NewTagList()
				tags.AddLow("task_version", task.Version)
				tags.AddLow("task_name", task.Family)
				tags.AddLow("task_family", task.Family)
				tags.AddLow("ecs_container_name", container.Name)

				if c.clusterName != "" {
					if !config.Datadog.GetBool("disable_cluster_name_tag_key") {
						tags.AddLow("cluster_name", c.clusterName)
					}
					tags.AddLow("ecs_cluster_name", c.clusterName)
				}

				var expiryDate time.Time
				for _, fn := range containerHandlers {
					if fn != nil {
						err := fn(ctx, container.DockerID, tags)
						if err != nil {
							log.Warnf("container handler func failed: %s", err)

							// cache result for 1 sencond
							// it prevents tagger to aggresivelly retry fetch
							expiryDate = time.Now().Add(1 * time.Second)
						}
					}
				}

				tags.AddOrchestrator("task_arn", task.Arn)

				low, orch, high, _ := tags.Compute()

				info := &TagInfo{
					Source:               ecsCollectorName,
					Entity:               containers.BuildTaggerEntityName(container.DockerID),
					HighCardTags:         high,
					OrchestratorCardTags: orch,
					LowCardTags:          low,
					ExpiryDate:           expiryDate,
				}
				output = append(output, info)
			}
		}
	}
	return output, nil
}
