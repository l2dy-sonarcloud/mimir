---
aliases:
  - ../../operators-guide/architecture/memberlist-and-the-gossip-protocol/
description: Memberlist manages Grafana Mimir cluster membership and node detection failure.
menuTitle: Memberlist and gossip protocol
title: Grafana Mimir memberlist and gossip protocol
weight: 80
---

<!-- Note: This topic is mounted in the GEM documentation. Ensure that all updates are also applicable to GEM. -->

# Grafana Mimir memberlist and gossip protocol

[Memberlist](https://github.com/hashicorp/memberlist) is a Go library that manages cluster membership, node failure detection, and message passing using a gossip-based protocol.
Memberlist is eventually consistent and network partitions are partially tolerated by attempting to communicate to potentially dead nodes through multiple routes.

By default, Grafana Mimir uses memberlist to implement a [key-value (KV) store](../key-value-store/) to share the [hash ring](../hash-ring/) data structures between instances.

When using a memberlist-based KV store, each instance maintains a copy of the hash rings.
Each Mimir instance updates a hash ring locally and uses memberlist to propagate the changes to other instances.
Updates generated locally and updates received from other instances are merged together to form the current state of the ring on the instance.

To configure memberlist, refer to [configuring hash rings](../../../configure/configure-hash-rings/).

## How memberlist propagates hash ring changes

When using a memberlist-based KV store, every Grafana Mimir instance propagates the hash ring data structures to other instances using the following techniques:

1. Propagating only the differences introduced in recent changes.
1. Propagating the full hash ring data structure.

Every `-memberlist.gossip-interval` an instance randomly selects a subset of all Grafana Mimir cluster instances configured by `-memberlist.gossip-nodes` and sends the latest changes to the selected instances.
This operation is performed frequently and it's the primary technique used to propagate changes.

In addition, every `-memberlist.pullpush-interval` an instance randomly selects another instance in the Grafana Mimir cluster and transfers the full content of the KV store, including all hash rings (unless `-memberlist.pullpush-interval` is zero, which disables this behavior).
After this operation is complete, the two instances have the same content as the KV store.
This operation is computationally more expensive, and as a result, it's performed less frequently. The operation ensures that the hash rings periodically reconcile to a common state.
