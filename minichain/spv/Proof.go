package spv

type Proof struct {
	txHash         string // 待验证的交易hash
	merkleRootHash string // Merkle根哈希
	height         int    // 待验证的交易所在区块的高度
	path           []Node // 待验证的交易所在区块的Merkle树路径，内部哈希值和偏向
}

func NewProof(hash string, merkleRootHash string, height int, path []Node) *Proof {
	return &Proof{
		txHash:         hash,
		merkleRootHash: merkleRootHash,
		height:         height,
		path:           path,
	}
}

func (p *Proof) GetTxHash() string {
	return p.txHash
}

func (p *Proof) GetMerkleRootHash() string {
	return p.merkleRootHash
}

func (p *Proof) GetHeight() int {
	return p.height
}

func (p *Proof) GetPath() []Node {
	return p.path
}
