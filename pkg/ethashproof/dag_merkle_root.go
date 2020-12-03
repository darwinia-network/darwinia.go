package ethashproof

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/darwinia-network/shadow/pkg/ethashproof/ethash"
	"github.com/darwinia-network/shadow/pkg/ethashproof/mtree"
	"github.com/darwinia-network/shadow/pkg/log"
)

func processDuringRead(f *os.File, startIn128Res int, fullSizeIn128Res uint32, mt *mtree.DagTree) error {
	_, err := f.Seek(int64(8+startIn128Res*128), 0)
	if err != nil {
		return err
	}
	r := bufio.NewReader(f)
	buf := [128]byte{}
	var i uint32 = 0
	for i < fullSizeIn128Res {
		n, err := io.ReadFull(r, buf[:128])
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			return err
		}
		if n != 128 {
			return fmt.Errorf("Malformed dataset")
		}
		mt.Insert(mtree.Word(buf), i)
		if err != nil && err != io.EOF {
			return err
		}
		i++
	}
	return nil
}

// 1. Generate the dataset if needed
// 2. Build merkle tree
// 3. If saveCache is true, save root merkle tree of 10 levels
//    to disk
// 4. Return merkle root
func CalculateDatasetMerkleRoot(epoch uint64, saveCache bool) (mtree.Hash, error) {
	blockno := epoch * 30000
    log.Info("Make the dag epoch %v", epoch)
	ethash.MakeDAG(blockno, ethash.DefaultDir)

	log.Info("Init the tree")
	dt := mtree.NewSHA256DagTree()

	fullSize := ethash.DAGSize(blockno)
	fullSizeIn128Resolution := fullSize / 128
	branchDepth := len(fmt.Sprintf("%b", fullSizeIn128Resolution-1))
	dt.RegisterStoredLevel(uint32(branchDepth), uint32(0))
	if saveCache {
		indices := []uint32{}
		for i := 0; i < 1<<CACHE_LEVEL; i++ {
			eindex := i << (uint64(branchDepth) - CACHE_LEVEL)
			if uint64(eindex) < fullSizeIn128Resolution {
				indices = append(indices, uint32(eindex))
			} else {
				break
			}
		}
		dt.RegisterIndex(indices...)
	}

	path := ethash.PathToDAG(uint64(blockno/30000), ethash.DefaultDir)
	log.Info("Calculating the proofs...")
	f, err := os.Open(path)
	if err != nil {
		return mtree.Hash{}, err
	}
	defer f.Close()
	_ = processDuringRead(f, 0, uint32(fullSizeIn128Resolution), dt)
	dt.Finalize()
	if saveCache {
		result := &DatasetMerkleTreeCache{
			Epoch:       epoch,
			ProofLength: uint64(branchDepth),
			CacheLength: CACHE_LEVEL,
			RootHash:    dt.RootHash(),
			Proofs:      [][]mtree.Hash{},
		}
		proofs := dt.ProofsForRegisteredIndices()
		for _, proof := range proofs {
			oneProof := proof[(uint64(branchDepth) - CACHE_LEVEL):]
			result.Proofs = append(result.Proofs, oneProof)
		}
		err = PersistCache(DefaultDir, result)
		if err != nil {
			return mtree.Hash{}, err
		}
	}
	return dt.RootHash(), nil
}
