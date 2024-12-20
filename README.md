# ecstop:  ECS+Stop

ecstop is a CLI tool to instantly stop your running resources of Amazon ECS to save costs.

This tool focuses on stopping **ECS Services, Tasks, and Container Instances** because they cost higher.
On the other hand, it does NOT clean Task Definitions or ECR Images, which are free or cheaper.

Pronounce: _ee-c-stop_

<!-- 
# Usage Pattern

1. CLI from local
2. on demand in AWS
3. Scheduled in AWS
-->

## Installation

```console
$ brew install t-kikuc/tap/ecstop
```

## Commands

- [`services`](#services): Scale-in Services
- [`tasks`](#tasks): Stop Tasks
- [`instances`](#instances): Stop Container Instances
- [`all`](#all): Stop the above 3 resources

You can execute the all commands with the following AWS options:
```console
  -p, --profile string   AWS profile
  -r, --region string    AWS region
```

### `services`

`ecstop services` updates `desiredCount` of ECS services to 0.

Flags:
```console
  -a, --all-clusters     Stop services in all clusters in the region
  -c, --cluster string   Name or ARN of the cluster to stop services
```
  - Only one of `--all-clusters` or `--cluster` is required.

### `tasks`

`ecstop tasks` stops ECS tasks.

This command is mainly used for standalone tasks, which are not controlled by ECS Services.
Even if you stop tasks of an ECS Service, the Service will start new tasks.

Flags:
```console
  -a, --all-clusters          Stop tasks in all clusters in the region
  -c, --cluster string        Name or ARN of the cluster to stop tasks

      --group string          Group name to stop tasks
      --group-prefix string   Group name prefix to stop tasks
      --standalone            Stop standalone tasks, whose group prefix is not 'service:'
```

- Only one of `--all-clusters` or `--cluster` is required.
- Only one of `--group`, `--group-prefix`, or `--standalone` is required.
  - ecstop stops all tasks whose `group` matches the condition.
  - `--standalone` stops tasks whose `group`'s prefix is NOT `service:`.



### `instances`

`ecstop instances` stops container instances. (not terminate)

Flags:
```console
  -a, --all-clusters     Stop instances in all clusters in the region
  -c, --cluster string   Name or ARN of the cluster to stop instances
```

- Only one of `--all-clusters` or `--cluster` is required.

### `all`

`ecstop all` stops services, standalone tasks, and container instances.

Flags:
```console
  -a, --all-clusters     Stop resources in all clusters in the region
  -c, --cluster string   Name or ARN of the cluster to stop resources
```

- Only one of `--all-clusters` or `--cluster` is required.

For example, `ecstop all --cluster xxx` is equal to: 

```console
ecstop services --cluster xxx
ecstop tasks --cluster xxx --standalone
ecstop instances --cluster xxx
```

## Required IAM Permissios

ecs:
- Read
  - `ListClusters`
  - `ListServices`
  - `ListTasks`
  - `ListContainerInstances`
  - `DescribeServices`
  - `DescribeTasks`
- Write
  - `UpdateService`
  - `StopTask`
  - `DescribeContainerInstances`

ec2:
- Write
  - `StopInstances`
