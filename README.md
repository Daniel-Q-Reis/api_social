# 🚀 Go Clean API Template

![Language](https://img.shields.io/badge/language-Go-blue.svg)
![Architecture](https://img.shields.io/badge/architecture-Clean-orange)
![License](https://img.shields.io/badge/license-MIT-green)

## 📖 About This Template

This repository is a production-ready starter kit for building robust and scalable RESTful APIs in Go, following the principles of **Clean Architecture**.

It provides a solid foundation with a well-defined project structure, pre-configured tools, and essential features, allowing you to focus on business logic instead of boilerplate setup.

### Features Included

-   **Clean Architecture:** A clear separation of concerns between business rules and infrastructure.
-   **Docker & Docker Compose:** Ready for containerization and easy local development setup.
-   **Configuration Management:** Centralized configuration loading from environment variables (`.env`).
-   **Logging:** Structured logging setup for better observability.
-   **CI/CD Ready:** Includes a basic GitHub Actions workflow for automated testing.
-   **Example Use Case:** A sample `translation` feature to demonstrate the architecture in practice.

---

## ✅ How to Use This Template

After clicking the "Use this template" button and creating your new repository, follow this checklist to personalize your new project:

### Step 1: Set Your Module Path
This is the most critical step.

-   [ ] **In `go.mod`:** Change the first line from `module github.com/Daniel-Q-Reis/ia_social_media` to your new repository's path (e.g., `module github.com/your-user/your-new-project`).
-   [ ] **Project-wide:** Use your IDE's "Find and Replace in Files" feature to replace all occurrences of `github.com/Daniel-Q-Reis/ia_social_media` with your new module path.
-   [ ] **In your terminal:** Run `go mod tidy` to synchronize the dependencies.

### Step 2: Customize Configuration

-   [ ] **Create `.env` file:** Rename `.env.example` to `.env`.
-   [ ] **Update `.env`:** Fill in your database credentials, app settings, and other required environment variables.

### Step 3: Update Project Documentation

-   [ ] **`README.md`:** **Delete this entire guide** and write a new README specific to your project. A basic template is provided at the bottom of this file to get you started.
-   [ ] **`LICENSE`:** Update the `[year]` and `[fullname]` in the `LICENSE` file with your own information.

### Step 4: Clean Up Example Code
This template comes with a sample `translation` use case. Once you understand the structure, you should remove it to start with a clean slate.

-   [ ] **Delete example files:**
    -   `internal/entity/translation.go`
    -   `internal/usecase/translation.go`
    -   `internal/controller/http/v1/translation.go`
    -   *And any other related files.*
-   [ ] **Update routers:** Remove the example routes from `internal/controller/http/v1/router.go`.

---
---

## 📋 New Project README Template
*(Copy the content below to start the new `README.md` for your project)*

```markdown
# [YOUR PROJECT NAME]

![Status: In Development](https://img.shields.io/badge/status-in--development-yellow)
![Language](https://img.shields.io/badge/language-Go-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green)

## 📖 About The Project

[A brief and clear description of what your project does.]

## ✨ Key Features

-   Feature 1
-   Feature 2
-   Feature 3

## 🛠️ Tech Stack

-   Go
-   PostgreSQL / MySQL
-   Docker
-   ...

## 🚀 Getting Started

[Instructions on how to set up and run your project locally.]

```