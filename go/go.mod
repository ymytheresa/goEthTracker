module github.com/ymytheresa/erc20-token-tracker

go 1.22.5
replace (
    github.com/ymytheresa/erc20-token-tracker/internal/config => ./internal/config
    github.com/ymytheresa/erc20-token-tracker/internal/ethereum => ./internal/ethereum
    github.com/ymytheresa/erc20-token-tracker/internal/cache => ./internal/cache
    github.com/ymytheresa/erc20-token-tracker/contracts => ./contracts
)