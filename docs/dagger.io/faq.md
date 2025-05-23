---
slug: /faq
pagination_next: null
pagination_prev: null
---

# FAQ


## General

### What is the Dagger Platform?

We're building the devops operating system, an integrated platform to orchestrate the delivery of applications to the cloud from start to finish. The Dagger Platform includes the Dagger Engine, Dagger Cloud, and the Dagger SDKs. Soon we will deliver the capability to publish and leverage prebuilt modules to further accelerate the adoption of Dagger across an organization's pipelines.

### How do I install, update, or uninstall Dagger?

Refer to the [installation documentation](./install.md).

### Does Dagger send telemetry?

By default, the Dagger CLI sends anonymized telemetry to dagger.io. This allows us to improve Dagger by understanding how it is used. Telemetry is optional and can be disabled at any time. If you are willing and able to leave telemetry enabled: thank you! This will help us better understand how Dagger is used, and will allow us to improve your experience.

### What telemetry does Dagger send?

The following information is included in telemetry:

- Dagger version
- Platform information
- Command run
- Anonymous device ID

We use telemetry for aggregate analysis, and do not tie telemetry events to a specific identity. Our telemetry implementation is open-source and can be reviewed in our [GitHub repository](https://github.com/dagger/dagger).

### Can Dagger telemetry be disabled?

Dagger implements the [Console Do Not Track (DNT) standard](https://consoledonottrack.com/). As a result, you can disable the telemetry by setting the environment variable `DO_NOT_TRACK=1` before running the Dagger CLI.

### Can I configure the Dagger Engine?

Yes. [Read more about Dagger Engine configuration](https://github.com/dagger/dagger/blob/main/core/docs/d7yxc-operator_manual.md).

### Can I use Dagger Engine to build Windows Containers?

Unfortunately, not right now. Dagger runs on top of BuildKit and support for Windows Containers is still experimental in BuildKit. In addition, Dagger has a lot of custom code written on top of BuildKit such as networking, init systems, and Linux namespace entries, that do not have exact parallels on Windows.

### I am stuck. How can I get help?

Join us on [Discord](https://discord.com/invite/dagger-io), and ask your question in our [help forum](https://discord.com/channels/707636530424053791/1030538312508776540). Our team will be happy to help you there!

## Dagger Cloud

### What is Dagger Cloud?

Dagger Cloud complements the Dagger Engine with a production-grade control plane. Features of Dagger Cloud include pipeline visualization and operational insights.

### Is Dagger Cloud a hosting service for Dagger Engines?

No, Dagger Cloud is a "bring your own compute" service. The Dagger Engine can run on a wide variety of machines, including most development and CI platforms. If the Dagger Engine can run on it, then Dagger Cloud supports it.

### Which CI providers does Dagger Cloud work with?

Because the Dagger Engine can integrate seamlessly with practically any CI, there is no limit to the type and number of CI providers that Dagger Cloud can work with to provide Dagger pipeline visualization and operational insights. Users report successfully leveraging Dagger with: GitLab, CircleCI, GitHub Actions, Jenkins,Tekton and many more.

### What is pipeline visualization?

Traces, a browser-based interface focused on tracing and debugging Dagger pipeline runs. A Trace contains detailed information about the steps performed by the pipeline. Traces let you visualize each step of your pipeline, drill down to detailed logs, understand how long operations took to run, and whether operations were cached.

### What operational insights does Dagger Cloud provide?

Dagger Cloud collects telemetry from all your organization's Dagger Engines, whether they run in development or CI, and presents it all to you in one place. This gives you a unique view on all pipelines, both pre-push and post-push.

### Why am I seeing `dagger functions` in the local trace list?

All commands that require module initialization at an engine level will send telemetry to Dagger Cloud.
Dagger needs to introspect a module to be able to print the available functions in a module, so it calls `dagger functions`. This happens for both local and remote runs, which is why the calls appears in the local trace list.

### How does Dagger classify Traces as originating from "CI" or "local"?

Dagger is aware of the context it runs on. When it runs in a CI context like GitHub, GitLab, CircleCI, or Jenkins, additional Trace metadata is displayed based on the Git repository information available. For this reason, it is important for Dagger to run in a Git context when running in CI.

### What are "orphaned Traces"?

You might see a warning message in Dagger Cloud about orphaned Traces. Orphaned Traces are Traces emitted in a CI context, that contain incomplete or no Git metadata. This generally happens when Git is not properly set up in the CI context that Dagger runs in. In GitHub Actions, for example, this context can be provided by using the `checkout` action in the step where Dagger is called.

### My CI provider is not supported in Dagger Cloud. Is there a way I can "force" my Traces into the Dagger Cloud dashboard?

It’s possible to send Traces by setting the `CI=true` variable in Dagger's runtime environment. However, Traces with incomplete Git repository data will show up as orphaned, so it is important to ensure that Dagger is running in a properly-set Git context.

## Dagger SDKs

### What language SDKs are available for Dagger?

We have [three types of SDKs](https://dagger.io/community-sdks), with varying levels of parity and support: official, community and experimental. We currently offer official SDKs for Go, TypeScript and Python. A community SDK is available for PHP, and an experimental SDK is available for Java.

### How can I move my SDK to a Dagger Community SDK?

To ensure a great experience for developers using Community SDKs, maintainers must meet the following requirements if you would like your SDK to graduate to the Community SDK level:

- Community Support – The maintainer must be active in the Dagger Discord, providing support and answering questions about the SDK.
- Version Compatibility – The SDK must stay up to date with the latest Dagger releases to ensure continued functionality.
- Documentation Maintenance – The maintainer is responsible for writing and updating documentation, including code snippets and examples. See full list of documentation requirements [here](https://docs.google.com/spreadsheets/d/1pvpzZbWarkuws811NEEbnv2D-ggec4iuZeYhUyVwsWc/edit?gid=245490315#gid=245490315).
- Openness to Contributions – Community SDKs should be open-source and encourage contributions from other developers.

If you want to kick off this process for your SDK, email community@dagger.io and we'll discuss further.

### How do I log in to a container registry using a Dagger SDK?

There are two options available:

1. Use the [`Container.withRegistryAuth()`](https://docs.dagger.io/api/reference/#Container-withRegistryAuth) GraphQL API method. A native equivalent of this method is available in each Dagger SDK.
1. Dagger SDKs can use your existing Docker credentials without requiring separate authentication. Simply execute `docker login` against your container registry on the host where your Dagger pipelines are running.

### How do I uninstall a Dagger SDK?

To uninstall a Dagger SDK, follow the same procedure that you would follow to uninstall any other SDK package in your chosen development environment.

## Dagger API

### What API query language does Dagger use?

Dagger uses GraphQL as its low-level language-agnostic API query language.

### Do I need to know GraphQL to use Dagger?

No. You only need to know one of Dagger's supported SDKs languages to use Dagger. The translation to underlying GraphQL API calls is handled internally by the Dagger SDK of your choice.

### There's no SDK for <language> yet. Can I still use Dagger?

Yes. It's possible to use the Dagger GraphQL API from any language that [supports GraphQL](https://graphql.org/code/)  or from the [Dagger CLI](./install.md).