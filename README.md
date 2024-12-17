# ecstop:  ECS+Stop

ecstop is a CLI tool to quickly clean your running ECS resources to save costs.

Pronounce: ee-c-s-top / ee-c-stop


<!-- 
# Usage Pattern

1. CLI from local
2. on demand in AWS
3. Scheduled in AWS
-->

## Commands

- [`services`](#services): Scale-in Services
- [`tasks`](#tasks): Stop Tasks
- [`instances`](#instances): Stop Container Instances
- [`all`](#all): Stop the above 3 resources

### `services`

`ecstop services` updates `desiredCount` of ECS services to 0.

```
Usage:
  ecstop services [flags]

Flags:
      --all-clusters     Scale-in services of all clusters in the region
      --cluster string   Name or ARN of the cluster to scale-in services
```

- Only one of `--all-clusters` or `--cluster` is required.


### `tasks`

`ecstop tasks` stops ECS tasks.

```
Usage:
  ecstop tasks [flags]

Flags:
      --all-clusters          Scale-in tasks of all clusters in the region
      --cluster string        Cluster name/arn to scale-in tasks
      --group string          Group name to scale-in tasks
      --group-prefix string   Group name prefix to scale-in tasks
      --standalone            Scale-in standalone tasks
```

- Only one of `--all-clusters` or `--cluster` is required.
- Only one of `--group`, `--group-prefix`, or `--standalone` is required.
  - ecstop stops all tasks whose `group` matches the condition.
  - `--standalone` stops tasks whose `group`'s prefix is NOT `service:`.

### `instances`

`ecstop instances` stops container instances.

```
Usage:
  ecstop instances [flags]

Flags:
      --all-clusters     Stop instances of all clusters in the region
      --cluster string   Cluster name/arn to stop instances
```

- Only one of `--all-clusters` or `--cluster` is required.

### `all`

`ecstop all` stops services, standalone tasks, and container instances.

```
Usage:
  ecstop all [flags]

Flags:
      --all-clusters     Stop resources of all clusters in the region
      --cluster string   Name or ARN of the cluster to stop resources
```

- Only one of `--all-clusters` or `--cluster` is required.
