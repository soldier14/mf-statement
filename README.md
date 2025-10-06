# MF Statement Service

MF Statement Service is a Go-based application designed to generate monthly financial statements from CSV files. It processes transactions, calculates income and expenditure, and outputs structured JSON statements.

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Usage](#usage)
- [Configuration](#configuration)
- [Development](#development)
  - [Testing](#testing)
  - [Code Style](#code-style)
- [License](#license)

## Features

- Parse CSV files containing transaction data.
- Filter transactions by year and month.
- Calculate total income and expenditure.
- Generate JSON statements with detailed transaction data.

## Project Structure
   
## Getting Started

### Prerequisites

- Go 1.20 or later
- A valid CSV file containing transaction data

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/mf-statement-service.git
   cd mf-statement-service
   
2. Build the application:
   ```bash
   go build -o mf-statement
   ```