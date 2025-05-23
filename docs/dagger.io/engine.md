---
slug: /api/engine
---

# Engine

The `Engine` type represents the Dagger Engine configuration and state. It provides fields to interact with a running Dagger Engine.

## Caching

Dagger caches two types of data:

1. Layers: This refers to build instructions and the results of some API calls.
2. Volumes: This refers to the contents of a Dagger filesystem volume and is persisted across Dagger Engine sessions.

The `Engine` type can be used to inspect or manually prune the cache.

To show all the cache entry metadata, use the following command:

```shell
dagger query <<EOF
{
  engine {
    localCache {
      entrySet {
        entries {
          description
          diskSpaceBytes
        }
      }
    }
  }
}
EOF
```

To see high level summaries of cache usage, use the following command:

```shell
dagger query <<EOF
{
  engine {
    localCache {
      entrySet {
        entryCount
        diskSpaceBytes
      }
    }
  }
}
EOF
```

To manually free up disk space used by the cache, use the following command:

```shell
dagger query <<EOF
{
  engine {
    localCache {
      prune
    }
  }
}
EOF