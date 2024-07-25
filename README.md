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
        - Addresses are saved in a db(for now).
        - Private keys are encrypted before stoting.
    - :white_check_mark: Basic transactions setup.
    - :white_check_mark: Signing transactions.
- Transaction verification.(In the mining process for now)
    - :white_check_mark: Signature verification.
    - :white_check_mark: UTXO validation.
    - :white_check_mark: Balance validation.
- :white_check_mark: Block hash using Merkle tree.
- Implement Merkle proof and verification

#### Stage 2
- ~~Implement wallets and transfer address generation into wallets.~~
- :white_check_mark: Create transaction pool for UTXOs(Mempool).
- ~~Update NewBlock() to get UTXOs from pool.~~
- :white_check_mark: Implement a Miner.
    - A simple one for now. After mining is completed and data is valid it broadcast TX ids through a channel and listener will remove those from mempool.

#### Stage 3
- Node system.
- Broadcasting.
- :white_check_mark: Block validation.
- Handle concurrency? (code: PNL-CONCURRENCY)

### Resources used
- https://jeiwan.net/posts/building-blockchain-in-go-part-1