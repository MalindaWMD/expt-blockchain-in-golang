# expt-blockchain-in-golang
 Experiment with blockchain from scratch using golang

#### Stage 1
- Setup basic structure and adding new blocks.
- Store blockchain in a DB.
- Implement PoW.
- Mining process.
- Implement addresses and transactions.
    - Address generation(withtin the app for now).
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
- Nodes system.
- Broadcasting.
- Block validation.


### Resources used
- 