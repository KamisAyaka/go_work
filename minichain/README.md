# Go-MiniChain 区块链系统

一个基于Golang实现的简易区块链系统，模拟比特币的核心机制，包含区块链构建、交易处理、工作量证明（PoW）共识机制等功能。

## 项目概述

### 核心特性
- **区块链结构**  
  - 包含区块头（版本号/前序哈希/Merkle根/时间戳/难度值/nonce）
  - 区块体存储交易数据并计算Merkle树根哈希
- **共识机制**  
  - 矿工节点通过随机nonce计算满足难度条件的区块哈希
  - 支持动态调整挖矿难度（默认前导4个零）
- **交易系统**  
  - UTXO模型实现交易验证
  - ECDSA签名保障交易安全
- **账户体系**  
  - 基于secp256k1的非对称加密生成账户
  - Base58编码生成钱包地址
- **网络模块**  
  - 交易池自动生成随机交易
  - 多账户间模拟转账行为
- **轻客户端支持（SPV）**  
  - SPV节点仅存储区块头以减少存储开销
  - 支持通过Merkle路径验证交易存在性
  - 实现交易验证的去中心化验证机制

### 项目结构
```
minichain/
├── config/                # 系统配置
│   └── config.go
├── data/                  # 数据结构模型
│   ├── Account.go
│   ├── Block.go
│   ├── BlockBody.go
│   ├── BlockHeader.go
│   ├── Transaction.go
│   └── UTXO.go
├── network/               # 网络层
│   ├── Network.go
│   ├── BlockChain.go
│   ├── TransactionPool.go
│   ├── MinerNode.go
|   └── spv.go
├── spv/                   # 轻客户端
│   ├── node.go            # SPV节点定义
│   └── Proof.go           # 证明结构
├── utils/                 # 工具类
│   ├── Base58Util.go
│   ├── MinerUtil.go
│   ├── SecurityUtil.go
│   └── SHA256Util.go
└── main.go                # 程序入口
```

---

## 快速开始

### 环境要求
- Go 1.18+
- 第三方依赖：`github.com/dustinxie/ecc`

### 运行步骤
1. 安装依赖
```bash
go mod tidy
```

2. 启动区块链网络
```bash
go run main.go
```

### 预期输出示例
```
Create the genesis Block! 
And the hash of genesis Block is : 0000A3D5..., you will see... 

Mined a new Block! Previous Block Hash is: 0000A3D5...
Block{
  blockHeader=BlockHeader{version=1..., 
  blockBody=BlockBody{merkleRootHash=...}
}
The sum of all amount 1000000
```

---

## 关键配置
`config/config.go` 包含可调参数：
```go
var MiniChainConfig = Config{
    difficulty:          4,       // 挖矿难度（前导零个数）
    maxTransactionCount: 16,      // 每个区块最大交易数
    nbAccount:           100,     // 系统初始账户数量
    initAmount:          10000,   // 初始账户金额
    spvEnabled:          true,    // 是否启用SPV节点（默认启用）
}
```

---

## 实现细节

### 核心算法
1. **Merkle树构建**  
   ```go
   func (m *MinerNode) GetBlockBody(transactions []data.Transaction) data.BlockBody {
       // 通过两两哈希合并生成Merkle根
       hashes := make([]string, 0)
       for _, tx := range transactions {
           hashes = append(hashes, utils.GetSha256Digest(tx.ToString()))
       }
       for len(hashes) > 1 {
           var newLevel []string
           for i := 0; i < len(hashes); i += 2 {
               pair := hashes[i]
               if i+1 < len(hashes) {
                   pair += hashes[i+1]
               }
               newLevel = append(newLevel, utils.GetSha256Digest(pair))
           }
           hashes = newLevel
       }
       return *data.NewBlockBody(hashes[0], transactions)
   }
   ```

2. **工作量证明**  
   ```go
   func (m *MinerNode) Mine(blockBody data.BlockBody) {
       block := m.GetBlock(blockBody)
       for {
           blockHash := utils.GetSha256Digest(block.ToString())
           if strings.HasPrefix(blockHash, utils.HashPrefixTarget()) {
               // 验证成功，添加新区块
           } else {
               // 更新nonce继续挖矿
               block.SetNonce(rand.Int63())
           }
       }
   }
   ```

3. **SPV验证算法**
   ```go
   // 交易验证流程
   func (p *SPVPeer) Verify(transaction data.Transaction) bool {
       txHash := utils.GetSha256Digest(transaction.ToString())
       proof := p.network.GetProof(txHash)
       currentHash := proof.GetTxHash()

       for _, node := range proof.GetPath() {
           switch node.GetOrientation() {
           case spv.LEFT:
               currentHash = utils.GetSha256Digest(node.GetTxHash() + currentHash)
           case spv.RIGHT:
               currentHash = utils.GetSha256Digest(currentHash + node.GetTxHash())
           }
       }

       return currentHash == proof.GetMerkleRootHash() &&
              p.headers[proof.GetHeight()].GetMerkleRootHash() == proof.GetMerkleRootHash()
   }
   ```

---

## 网络模块说明
### SPV轻客户端实现
1. **节点结构**
   ```go
   type SPVPeer struct {
       headers []data.BlockHeader // 区块头存储
       account data.Account       // 绑定账户
       network *NetWork          // 网络引用
   }
   ```

2. **验证流程**
   - **交易验证**：通过Merkle路径重建根哈希验证交易存在性
   - **区块头同步**：仅存储区块头并验证区块有效性
   - **证明生成**：矿工节点提供交易在区块中的Merkle路径

3. **证明结构**
   ```go
   type Proof struct {
       txHash         string // 交易哈希
       merkleRootHash string // Merkle根哈希
       height         int    // 区块高度
       path           []spv.Node // 验证路径
   }
   ```

---

## 使用示例
```go
// SPV节点验证交易
spvPeer := network.spvPeer[0]
txHash := "交易哈希值"
if spvPeer.Verify(transaction) {
    fmt.Println("交易验证通过")
} else {
    fmt.Println("交易验证失败")
}
```

---

## 功能优势
- **存储优化**：SPV节点存储量仅为全节点的1/1000
- **快速验证**：仅需验证Merkle路径而非完整区块
- **去中心化**：无需信任第三方即可验证交易有效性
- **轻量化**：适合移动端和资源受限设备使用

---