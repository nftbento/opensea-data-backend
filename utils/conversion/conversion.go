/*

 */

package conversion

import "math/big"

func ConvertEthToWei(eth string) *big.Int {
	ethFloat, ok := Zero().SetString(eth)
	if !ok {
		return nil
	}

	exp, _ := Zero().SetString("1000000000000000000.0")
	weiFloat := Zero().Mul(ethFloat, exp)
	weiInt, _ := weiFloat.Int(nil)

	return weiInt
}

func ConvertWeiToEth(wei string) *big.Float {
	weiFloat, ok := Zero().SetString(wei)
	if !ok {
		return nil
	}

	exp, _ := Zero().SetString("1000000000000000000.0")
	ethFloat := Zero().Quo(weiFloat, exp)

	return ethFloat
}

func ConvertWeiToGwei(wei string) int64 {
	weiFloat, ok := Zero().SetString(wei)
	if !ok {
		return 0
	}

	exp, _ := Zero().SetString("1000000000.0")
	gweiFloat := Zero().Quo(weiFloat, exp)

	gweiInt64, _ := gweiFloat.Int64()

	return gweiInt64
}

func Zero() *big.Float {
	r := big.NewFloat(0.0)
	r.SetPrec(512)
	return r
}
