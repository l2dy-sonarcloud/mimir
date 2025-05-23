---
aliases:
  - ../../migrate/migrate-from-cortex/
description: Learn how to migrate your deployment of Cortex to Grafana Mimir to simplify the deployment and continued operation of a horizontally scalable, multi-tenant time series database with long-term storage.
menuTitle: Migrate from Cortex
title: Migrate from Cortex to Grafana Mimir
weight: 10
---

<!-- Note: This topic is mounted in the GEM documentation. Ensure that all updates are also applicable to GEM. -->

# Migrate from Cortex to Grafana Mimir

As an operator, you can migrate a Jsonnet deployment of [Cortex](https://cortexmetrics.io/) to Grafana Mimir.
The overview includes the steps required for any environment. To migrate deployment environments with Jsonnet, see [Migrate to Grafana Mimir using Jsonnet](#migrate-to-grafana-mimir-using-jsonnet).

{{< admonition type="warning" >}}
This document was tested with Cortex versions 1.10 and 1.11.

It might work with more recent versions of Cortex, but it's not guaranteed.
{{< /admonition >}}

To migrate a Helm deployment of Cortex refer to [Migrate from Cortex](/docs/helm-charts/mimir-distributed/latest/migration-guides/migrate-from-cortex).

Grafana Mimir includes significant changes that simplify the deployment and continued operation of a horizontally scalable, multi-tenant time series database with long-term storage.

The changes make Grafana Mimir easier to run out of the box:

- Removed configuration parameters that don't require tuning
- Renamed some parameters so that they're more easily understood
- Updated the default values of some existing parameters

The `mimirtool` automates configuration conversion.
It provides a simple migration by generating Mimir configuration from Cortex configuration.

**Before you begin:**

- Ensure that you are running either Cortex 1.10.X or Cortex 1.11.X.

  If you are running an older version of Cortex, upgrade to [Cortex 1.11.1](https://github.com/cortexproject/cortex/releases) before proceeding with the migration.

- Ensure you have installed Cortex alerting and recording rules as well as Cortex dashboards.

  Using the monitoring mixin, you need to install both alerting and recording rules in either Prometheus or Cortex. You also need to install dashboards in Grafana.
  To download a prebuilt ZIP file that contains the alerting and recording rules, refer to [Release Cortex-jsonnet 1.11.0](https://github.com/grafana/cortex-jsonnet/releases/download/1.11.0/cortex-mixin.zip).

  To upload rules to the ruler using mimirtool, refer to [mimirtool rules](../../../manage/tools/mimirtool/).
  To import the dashboards into Grafana, refer to [Import dashboard](/docs/grafana/latest/dashboards/export-import/#import-dashboard).

## Notable changes

{{< admonition type="note" >}}
For the full list of changes, refer to Mimir’s [CHANGELOG](https://github.com/grafana/mimir/blob/main/CHANGELOG.md).
{{< /admonition >}}

- The Grafana Mimir HTTP server defaults to listening on port 8080; Cortex defaults to listening on port 80.
  To maintain port 80 as the listening port, set `-server.http-listen-port=80`.
- Grafana Mimir uses `anonymous` as the default tenant ID when `-auth.multitenancy=false`.
  Cortex uses `fake` as the default tenant ID when `-auth.enabled=false`.
  Use `-auth.no-auth-tenant=fake` when `-auth.multitenancy=false` to match the Cortex default tenant ID.
- Grafana Mimir removes the legacy HTTP prefixes deprecated in Cortex.

  - Query endpoints

    | Legacy                                                  | Current                                                    |
    | ------------------------------------------------------- | ---------------------------------------------------------- |
    | `/<legacy-http-prefix>/api/v1/query`                    | `<prometheus-http-prefix>/api/v1/query`                    |
    | `/<legacy-http-prefix>/api/v1/query_range`              | `<prometheus-http-prefix>/api/v1/query_range`              |
    | `/<legacy-http-prefix>/api/v1/query_exemplars`          | `<prometheus-http-prefix>/api/v1/query_exemplars`          |
    | `/<legacy-http-prefix>/api/v1/series`                   | `<prometheus-http-prefix>/api/v1/series`                   |
    | `/<legacy-http-prefix>/api/v1/labels`                   | `<prometheus-http-prefix>/api/v1/labels`                   |
    | `/<legacy-http-prefix>/api/v1/label/{name}/values`      | `<prometheus-http-prefix>/api/v1/label/{name}/values`      |
    | `/<legacy-http-prefix>/api/v1/metadata`                 | `<prometheus-http-prefix>/api/v1/metadata`                 |
    | `/<legacy-http-prefix>/api/v1/read`                     | `<prometheus-http-prefix>/api/v1/read`                     |
    | `/<legacy-http-prefix>/api/v1/cardinality/label_names`  | `<prometheus-http-prefix>/api/v1/cardinality/label_names`  |
    | `/<legacy-http-prefix>/api/v1/cardinality/label_values` | `<prometheus-http-prefix>/api/v1/cardinality/label_values` |
    | `/api/prom/user_stats`                                  | `/api/v1/user_stats`                                       |

  - Distributor endpoints

    | Legacy endpoint              | Current                       |
    | ---------------------------- | ----------------------------- |
    | `/<legacy-http-prefix>/push` | `/api/v1/push`                |
    | `/all_user_stats`            | `/distributor/all_user_stats` |
    | `/ha-tracker`                | `/distributor/ha_tracker`     |

  - Ingester endpoints

    | Legacy      | Current              |
    | ----------- | -------------------- |
    | `/ring`     | `/ingester/ring`     |
    | `/shutdown` | `/ingester/shutdown` |
    | `/flush`    | `/ingester/flush`    |
    | `/push`     | `/ingester/push`     |

  - Ruler endpoints

    | Legacy                                                | Current                                                            |
    | ----------------------------------------------------- | ------------------------------------------------------------------ |
    | `/<legacy-http-prefix>/api/v1/rules`                  | `<prometheus-http-prefix>/api/v1/rules`                            |
    | `/<legacy-http-prefix>/api/v1/alerts`                 | `<prometheus-http-prefix>/api/v1/alerts`                           |
    | `/<legacy-http-prefix>/rules`                         | `<prometheus-http-prefix>/config/v1/rules`                         |
    | `/<legacy-http-prefix>/rules/{namespace}`             | `<prometheus-http-prefix>/config/v1/rules/{namespace}`             |
    | `/<legacy-http-prefix>/rules/{namespace}/{groupName}` | `<prometheus-http-prefix>/config/v1/rules/{namespace}/{groupName}` |
    | `/<legacy-http-prefix>/rules/{namespace}`             | `<prometheus-http-prefix>/config/v1/rules/{namespace}`             |
    | `/<legacy-http-prefix>/rules/{namespace}/{groupName}` | `<prometheus-http-prefix>/config/v1/rules/{namespace}/{groupName}` |
    | `/<legacy-http-prefix>/rules/{namespace}`             | `<prometheus-http-prefix>/config/v1/rules/{namespace}`             |
    | `/ruler_ring`                                         | `/ruler/ring`                                                      |

  - Alertmanager endpoints

    | Legacy                  | Current                            |
    | ----------------------- | ---------------------------------- |
    | `/<legacy-http-prefix>` | `/alertmanager`                    |
    | `/status`               | `/multitenant_alertmanager/status` |

## Generate the configuration for Grafana Mimir

The [`mimirtool config convert`](../../../manage/tools/mimirtool/#config) command converts Cortex configuration to Mimir configuration. You can use it to update both flags and configuration files.

### Install mimirtool

To install Mimirtool, download the appropriate binary from the [latest release](https://github.com/grafana/mimir/releases/latest) for your operating system and architecture and make it executable.

Alternatively, use a command line tool such as `curl` to download `mimirtool`. For example, for Linux with the AMD64 architecture, use the following command:

```bash
curl -fLo mimirtool https://github.com/grafana/mimir/releases/latest/download/mimirtool-linux-amd64
chmod +x mimirtool
```

### Use mimirtool

The `mimirtool config convert` command converts Cortex 1.11 configuration files to Grafana Mimir configuration files.
It removes any configuration parameters that are no longer available in Grafana Mimir, and it renames configuration parameters that have a new name.
If you have explicitly set configuration parameters to a value matching the Cortex default, by default, `mimirtool config convert` doesn't update the value.
To have `mimirtool config convert` update explicitly set values from the Cortex defaults to the new Grafana Mimir defaults, provide the `--update-defaults` flag.
Refer to [convert](../../../manage/tools/mimirtool/#convert) for more information on using `mimirtool` for configuration conversion.

## Migrate to Grafana Mimir using Jsonnet

Grafana Mimir has a Jsonnet library that replaces the existing Cortex Jsonnet library and updated monitoring mixin.

### Migrate to Grafana Mimir video

The following video shows you how to migrate to Grafana Mimir using Jsonnet.

{{< vimeo 691929138 >}}

<br/>

### Migrate to Grafana Mimir instructions

The following instructions describe how to migrate to Grafana Mimir using Jsonnet.

To install the updated libraries using `jsonnet-bundler`, run the following commands:

```bash
jb install github.com/grafana/mimir/operations/mimir@main
jb install github.com/grafana/mimir/operations/mimir-mixin@main
```

**To deploy the updated Jsonnet:**

1. Install the updated monitoring mixin.

   a. Add the dashboards to Grafana. The dashboards replace your Cortex dashboards and continue to work for monitoring Cortex deployments.

   {{< admonition type="note" >}}
   Resource dashboards are enabled by default and require additional metrics sources.
   To understand the required metrics sources, refer to [Additional resources metrics](../../../manage/monitor-grafana-mimir/requirements/#additional-resources-metrics).
   {{< /admonition >}}

   b. Install the recording and alerting rules into the ruler or a Prometheus server.

1. Replace the import of the Cortex Jsonnet library with the Mimir Jsonnet library.
   For example:
   ```jsonnet
   import 'github.com/grafana/mimir/operations/mimir/mimir.libsonnet'
   ```
1. Remove the `cortex_` prefix from any member keys of the `<MIMIR>._config` object.
   For example, `cortex_compactor_disk_data_size` becomes `compactor_disk_data_size`.
1. If you are using the Cortex defaults, set the server HTTP port to 80.
   The new Mimir default is 8080.
   For example:
   ```jsonnet
   (import 'github.com/grafana/mimir/operations/mimir/mimir.libsonnet') {
     _config+: {
       server_http_port: 80,
     },
   }
   ```
1. For each component, use `mimirtool` to update the configured arguments.
   To extract the flags for each component, refer to [Extracting flags from Jsonnet](../../../manage/tools/mimirtool/#extracting-flags-from-jsonnet).
1. Apply the updated Jsonnet

To verify that the cluster is operating correctly, use the [monitoring mixin dashboards](../../../manage/monitor-grafana-mimir/dashboards/).
