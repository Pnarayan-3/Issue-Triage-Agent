# 🤖 Issue Triage Agent

An AI-powered issue triage agent that automatically analyzes software engineering issues, classifies them, assigns useful metadata, and posts a structured triage result.

The agent integrates with **GitHub Issues** and **Jira** and uses **Configurable Generative AI Model** to analyze issue content.

---

## 📌 Overview

In software development teams, incoming issues and tickets need to be manually reviewed and categorized before they can be assigned to the appropriate team.

The **Issue Triage Agent** automates this initial triage process.

It reads an issue, analyzes its title and description using Configurable Generative AI Model, determines its classification, validates the AI response, and updates the issue with the triage information.

### Supported Issue Sources

- GitHub Issues
- Jira Issues

### AI Provider

- Configurable Generative AI Model

### Implementation

- Go
- GitHub Actions
- REST APIs
- Docker

---

# 🎯 Problem Statement

A typical software issue may require a developer or triage engineer to determine:

- What type of issue is this?
- How important is it?
- How severe is it?
- Which component is affected?
- Which team should handle it?
- What labels should be applied?
- Does the issue require human review?

Doing this manually for every incoming issue can be repetitive.

The Issue Triage Agent automates this first-level classification while keeping a confidence-based human review mechanism.

---

# ✨ Features

## Core Features

- 🤖 AI-powered issue classification
- 🐛 Issue type detection
- 🚦 Priority classification
- 🔥 Severity classification
- 🧩 Component identification
- 👥 Team routing
- 🏷️ AI-generated labels
- 📊 Confidence scoring
- 👨‍💻 Human-review detection
- 📝 Automated triage comments
- 🔄 Reopened issue handling
- 🛑 Duplicate triage protection
- 🔁 Configurable Gen-AI Model API retry handling
- ✅ AI response validation
- 🏗️ GitHub and Jira adapters
- 🐳 Docker support
- ⚙️ Configurable environment variables

---

# 🧠 AI Classification

The agent asks Configurable Generative AI Model to classify each issue using predefined categories.

## Issue Type
## Priority
## Severity
## Component
## Team

## Confidence-Based Review

# 🏗️ Architecture

                    ┌──────────────────────┐
                    │      GitHub Issue    │
                    │          /           │
                    │       Jira Issue     │
                    └──────────┬───────────┘
                               │
                               ▼
                    ┌──────────────────────┐
                    │     Issue Adapter    │
                    └──────────┬───────────┘
                               │
                               ▼
                    ┌──────────────────────┐
                    │      Issue Model     │
                    └──────────┬───────────┘
                               │
                               ▼
                    ┌──────────────────────┐
                    │  Generative AI Model │
                    └──────────┬───────────┘
                               │
                               ▼
                    ┌──────────────────────┐
                    │      Validator       │
                    └──────────┬───────────┘
                               │
                               ▼
                    ┌──────────────────────┐
                    │    Triage Engine     │
                    └──────────┬───────────┘
                               │
                    ┌──────────┴───────────┐
                    │                      │
                    ▼                      ▼
          ┌──────────────────┐   ┌──────────────────┐
          │ GitHub           │   │ Jira             │
          │                  │   │                  │
          │ Labels           │   │ Labels           │
          │ Comments         │   │ Priority         │
          │ Reopened Update  │   │ Comments         │
          └──────────────────┘   └──────────────────┘

# 🎫 Jira Integration

The project also supports Jira Issues through a dedicated Jira adapter.

The Jira integration communicates with Jira using its REST API.

The current Jira workflow supports:

Fetching a Jira issue
Reading the issue summary and description
AI-based triage
Adding labels
Updating priority
Adding a triage comment
Detecting existing triage comments
Preventing duplicate triage

# 🔁 AI API Retry Handling

Temporary API failures are handled using configurable retries.

# 🐳 Docker

The project includes a multi-stage Dockerfile.

Build Command: docker build -t issue-triage-agent .
Run Container: docker run --rm \
  -e GITHUB_TOKEN=your-token \
  -e GEMINI_API_KEY=your-key \
  -e ISSUE_NUMBER=1 \
  -e ISSUE_ACTION=opened \
  -e REPOSITORY=owner/repository \
  issue-triage-agent

# 🧪 Testing

The project contains unit tests for the core triage logic.

# 🔍 Code Quality

Standard Go development tools can be used to verify the project.

# 👨‍💻 Author

Pushkar Narayan

Software Developer