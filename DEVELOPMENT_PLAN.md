# md2wechat-new Development Plan

## Goal

Create a fork of `md2wechat-skill` that keeps the existing AI mode and changes article conversion API mode into a configurable provider layer for mainstream OpenAI-compatible APIs such as SiliconFlow, DeepSeek, OpenAI, and custom endpoints.

## Confirm-First Steps

1. Inspect the source project and identify API/AI boundaries. Done.
2. Initialize `md2wechat-new` from `md2wechat-skill`. Done.
3. Map the current API mode call path and tests. Done.
4. Design the configurable text conversion provider contract. Done.
5. Implement provider configuration loading and defaults. Done.
6. Replace md2wechat.cn-only API conversion with provider-based conversion. Done.
7. Add examples and documentation for SiliconFlow, DeepSeek, OpenAI, and custom OpenAI-compatible APIs. Done.
8. Run focused tests and fix regressions. Done.

Each step starts only after explicit user confirmation.
