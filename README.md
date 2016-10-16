ecs-host-manager
================

ecs-host-manager is an agent taking care of rolling out new EC2 instances in
a ECS cluster.

Overview
--------

This agent works with the way we configured our ECS clusters, with optional
autoscaling features, here's how it works:

- When a change of launch configuration is detected on the autoscaling group
that an EC2 instance is part of a new instance of the same type is created.

- The new instance is expected to start the ECS agent, and automatically join
the same ECS cluster.

- When the new instance is ready, the old EC2 instance (where the agent has
detected the change of launch configuration) will deregister itself from the
ECS cluster. This causes the ECS scheduler to turn off the ECS tasks and start
them on a different host.

- When the new EC2 instance has had ECS tasks running on for a couple of minutes
the old instance will remove itself from the autoscaling group and add the new
instance instead.

- Once the old instance doesn't have any tasks running on it anymore it submits
a termination request.

If things go wrong at any steps the old instance will halt the process, turn off
the new host that it created and re-register itself with the ECS cluster.

Running
-------

The recommended way of installing the agent is to run the docker image:
```
docker run segment/ecs-host-manager:latest
```
