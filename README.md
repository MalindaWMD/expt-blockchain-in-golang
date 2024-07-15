# expt-blockchain-in-golang
 Experiment with blockchain from scratch using golang

 Used BoltDB as it's easy to install with go application.

(code: PLN-STAGES)
#### Stage 1
- :white_check_mark: Setup basic structure and adding new blocks.
- :white_check_mark: Store blockchain in a DB.
- :white_check_mark: Implement PoW.
- Mining process.
- Implement addresses and transactions.
    - :white_check_mark: Address generation(withtin the app for now).
        *Basic version for now. Wallets should be implemented.*
    - Basic transactions setup.
    - Signing transactions.
- Transaction verification.(In the mining process for now)
    - Signature verification.
    - UTXO validation.
    - Balance validation.

#### Stage 2
- Implement wallets and transfer address generation into wallets.
- Create transaction pool for UTXOs.
- Update NewBlock() to get UTXOs from pool.

#### Stage 3
- Node system.
- Broadcasting.
- Block validation.
- Handle concurrency? (code: PNL-CONCURRENCY)

### Resources used
- https://jeiwan.net/posts/building-blockchain-in-go-part-1