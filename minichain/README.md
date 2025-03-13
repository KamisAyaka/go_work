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

### 项目结构
```
minichain/
├── config/                # 系统配置
│   └── config.go
├── data/                  # 数据结构模型
│   ├── Account.go
│   ├── Block.go
│   ├── BlockBody.go
│   ├── BlockChain.go
│   ├── BlockHeader.go
│   ├── Transaction.go
│   ├── TransactionPool.go
│   └── UTXO.go
├── consensus/             # 共识机制
│   └── MinerNode.go
├── network/               # 网络层
│   └── Network.go
├── utils/                 # 工具类
│   ├── Base58Util.go
│   ├── MinerUtil.go
│   ├── SecurityUtil.go
│   └── SHA256Util.go
└── main.go                # 程序入口
```

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

## 关键配置
`config/config.go` 包含可调参数：
```go
var MiniChainConfig = Config{
    difficulty:          4,       // 挖矿难度（前导零个数）
    maxTransactionCount: 2,       // 每个区块最大交易数
    nbAccount:           100,     // 系统初始账户数量
    initAmount:          10000    // 初始账户金额
}
```

## 实现细节

### 核心算法
1. **Merkle树构建**  
   ```go
   func (m *MinerNode) GetBlockBody(transactions []data.Transaction) data.BlockBody {
       // 通过两两哈希合并生成Merkle根
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
    merkleRoot := hashes[0]
  }
   ```

2. **工作量证明**  
   ```go
   func (m *MinerNode) Mine(blockBody data.BlockBody) {
       for {
           blockHash := utils.GetSha256Digest(block.ToString())
           if strings.HasPrefix(blockHash, utils.HashPrefixTarget()) {
               // 找到合法nonce值
           }
       }
   }
   ```

3. **地址生成**  
   ```go
    publicKeyBytes := elliptic.Marshal(publicKey, publicKey.X, publicKey.Y)
    hash160 := utils.Ripemd160Digest(utils.Sha256Digest(publicKeyBytes))
    // Base58Check编码
    versionPayload := append([]byte{0x00}, hash160...)
    checksum := utils.Sha256Digest(utils.Sha256Digest(versionPayload))[:4]
    return base58.Encode(append(versionPayload, checksum...))
   ```
