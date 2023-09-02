# go-commons

This repository contains a collection of common utility packages written in Go.

## Table of Contents

- [Directory Configuration](#directory-configuration)
- [Packages](#packages)
- [License](#license)

## Directory Configuration

```
go-commons
├── LICENSE
├── README.md
├── go.mod
├── go.sum
└── pkg
    ├── collections
    │   └── map.go
    ├── constants.go
    ├── csv
    │   └── csv_logger.go
    ├── curl
    │   └── curl.go
    ├── db
    │   └── kv_db.go
    ├── file
    │   ├── dir.go
    │   ├── dir_test.go
    │   ├── glob.go
    │   ├── glob_test.go
    │   ├── open
    │   │   └── file_open.go
    │   └── scanner.go
    ├── json
    │   └── pretty.go
    ├── log
    │   ├── csv
    │   │   ├── csv_logger.go
    │   │   ├── csv_logger_test.go
    │   │   └── options.go
    │   └── logger.go
    ├── math
    │   └── rand.go
    ├── stream
    │   ├── base.go
    │   ├── intermediate.go
    │   └── terminal.go
    └── strings
        └── str_case.go
```

## Packages

- **collections**: Provides utility functions for working with collections, such as maps.
- **csv**: Contains CSV related utility functions, including a CSV logger.
- **curl**: Provides a wrapper for making HTTP requests using cURL.
- **db**: Contains utility functions for working with key-value databases.
- **file**: Provides functions for working with files and directories, including glob pattern matching and file scanning.
- **json**: Contains utility functions for pretty printing JSON data.
- **log**: Provides a logging framework, including a CSV logger with configurable options.
- **math**: Contains utility functions for working with math operations, including random number generation.
- **stream**: Provides utility functions for working with stream data, including base, intermediate, and terminal operations.
- **strings**: Contains utility functions for manipulating strings, including case conversion.
- **time**: Contains utility functions for tracking time.

## License

This project is licensed under the [MIT License](LICENSE).

configuration diagram (MERMAID):

```mermaid
graph LR
    A[go-commons] --> B[LICENSE]
    A --> C[README.md]
    A --> D[go.mod]
    A --> E[go.sum]
    A --> F[pkg]
    F --> G[collections]
    G --> H[map.go]
    F --> I[constants.go]
    F --> J[csv]
    J --> K[csv_logger.go]
    F --> L[curl]
    L --> M[curl.go]
    F --> N[db]
    N --> O[kv_db.go]
    F --> P[file]
    P --> Q[dir.go]
    P --> R[dir_test.go]
    P --> S[glob.go]
    P --> T[glob_test.go]
    P --> U[open]
    U --> V[file_open.go]
    P --> W[scanner.go]
    F --> X[json]
    X --> Y[pretty.go]
    F --> Z[log]
    Z --> AA[csv]
    AA --> AB[csv_logger.go]
    AA --> AC[csv_logger_test.go]
    AA --> AD[options.go]
    Z --> AE[logger.go]
    F --> AF[math]
    AF --> AG[rand.go]
    F --> AH[stream]
    AH --> AI[base.go]
    AH --> AJ[intermediate.go]
    AH --> AK[terminal.go]
    F --> AL[strings]
    AL --> AM[str_case.go]
```
