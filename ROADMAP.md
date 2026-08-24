# lumibot Development Roadmap

This document outlines the current progress and future milestones for the `lumibot` project.
It serves as a guide for prioritizing tasks and understanding the overall direction of the Discord bot integration with `lumitree`.

## Phase 1: Foundation (Completed)

- [x] **Project Initialization**: Setup Go 1.23 environment, `.gitignore`, and Linter configurations.
- [x] **CI/CD Pipeline**: Implement GitHub Actions workflows for Go linting, unit testing (with race detector), and shell script validation.
- [x] **TDD Infrastructure**: Establish 100% test coverage for configuration management (`pkg/config`).
- [x] **Data Persistence**: Implement Pure-Go SQLite (`modernc.org/sqlite`) for zero-dependency local storage (`pkg/store`), and write CRUD unit tests.
- [x] **Dependency Pinning**: Resolve Go Toolchain auto-upgrade issues and ensure stable CI builds.

## Phase 2: Minimum Viable Product (MVP) (Completed)

- [x] **API Integration**: Implement the `lumitree` API client (`pkg/client`) to fetch calendar meta-data and event lists.
- [x] **Discord Bot Core**: Setup `bwmarrin/discordgo` session and graceful shutdown logic.
- [x] **Slash Commands**:
  - `/add <ID/URL>`: Register a calendar ID for the Discord guild.
  - `/list`: Display currently subscribed calendars.
  - `/remove <ID>`: Unsubscribe a calendar.
  - `/today`: Manually fetch and display today's events as Discord Embeds.
- [x] **Daily Cron Job**: Implement a scheduler (`robfig/cron/v3`) to automatically fetch and broadcast events every morning at 08:00 JST to the configured channels.
- [x] **Dockerization**: Complete the multi-stage Distroless Dockerfile and enable automated GHCR multi-arch builds via `docker-publish.yml`.

## Phase 3: Production Readiness & Deployment (In Progress)

- [x] **Documentation**: Expand `README.md` with setup instructions, environment variable references, and bot invitation link guidelines.
- [x] **Kubernetes Manifests**: Create Kustomize manifests and ArgoCD Application definitions for deploying `lumibot` alongside `lumitree`.
- [ ] **Observability**: Add structured logging and basic health check endpoints.
- [ ] **Live Testing**: End-to-end verification in a staging Discord server.

## Phase 4: Future Enhancements

- [ ] **Webhook / Real-time Updates**: Explore integrating with `lumitree` push notifications (if supported in the future) to announce schedule changes instantly.
- [ ] **Customizable Delivery Time**: Allow each Discord guild to configure their preferred daily notification time (e.g., 07:00 instead of 08:00).
- [ ] **Multiple Calendar Aggregation**: Merge events from multiple registered calendars into a single, cohesive daily briefing.
- [ ] **Event Filtering**: Allow users to filter out specific event categories or tags from the daily notifications.
