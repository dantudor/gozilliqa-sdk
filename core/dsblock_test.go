/*
 * Copyright (C) 2021 Zilliqa
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package core

import (
	"encoding/json"
	"github.com/Zilliqa/gozilliqa-sdk/multisig"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"io/ioutil"
	"strings"
	"testing"
)

func Test_DeserializeFromJsonToDsBlockT(t *testing.T) {
	dsJson, err := ioutil.ReadFile("dsblock.json")
	if err != nil {
		t.Fatal(err.Error())
	}

	var dsBlockT DsBlockT
	err2 := json.Unmarshal(dsJson, &dsBlockT)
	if err2 != nil {
		t.Fatal(err2.Error())
	}

	dsBlockHeader := NewDsBlockHeaderFromDsBlockT(&dsBlockT)
	headerBytes := dsBlockHeader.Serialize()
	if "0A46080212205E0DBBA477EA36CCA0222CEB23E9B62262CDCDEB73BF2AB284C8C8B230D938371A206D3434708311D2F814279AA235B869A2DA45443FDCFAEE5625C5C60A1E297CA7100518032A230A210293F3A79C39A0C3952711B0F85440E9E5F96E4DC704AA431779ADEA817F12D8F23003380A42120A10000000000000000000000000773594004A320A30000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000523D0A230A2103EEB8669191E26A163A9E3944B0150994CA04AD1CA7D0D9F1875AAE35148DD29312160A140000000000000000000000002E6C2A340000816D5AA5010A201F7A31292B4460E4F6559E684B37FBDE97CEB1C52CC23A4185F7EE612D1C447B1280010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" != strings.ToUpper(util.EncodeHex(headerBytes)) {
		t.Fail()
	}

	dsJson1, err3 := ioutil.ReadFile("dsblock2.json")
	if err3 != nil {
		t.Fatal(err3.Error())
	}

	var dsBlockT1 DsBlockT
	err4 := json.Unmarshal(dsJson1, &dsBlockT1)
	if err4 != nil {
		t.Fatal(err4.Error())
	}

	dsBlockHeader1 := NewDsBlockHeaderFromDsBlockT(&dsBlockT1)
	headerBytes1 := dsBlockHeader1.Serialize()
	t.Log(strings.ToUpper(util.EncodeHex(headerBytes1)))
	if "0A4608021220AA5764CC1646085C6DD2042BB784ED6A4E154D22F0D1D1DACDD9E86662601E071A202888FA1F01AA2C8B1ADE1B254B19D93DAD56538E9234F218E1BC4C97FF424813100518032A230A210213D5A7F74B28F3F588FF6520748DBB541986E98F75FA78D6334B2D0AAB4C1E573004380F42120A10000000000000000000000000773594004A320A30000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000523D0A230A210334AA0F7CA2EAA56B6B752533F9C60777E96C6D1ABE84B463F60ADD89843794AE12160A140000000000000000000000000A76DB220000816D5AA5010A20B0FE97D2028F734E895E45297A932F63CDEE08AF5EA25BCA732B24C15D759BAF128001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000062230A2103AFEEA358BBFD6A350B59943E11D65386142033C63B921032608DE65334C1DF8C6A08086E1204080B1001" != strings.ToUpper(util.EncodeHex(headerBytes1)) {
		t.Fail()
	}
}

func TestVerifyDsBlock(t *testing.T) {
	dsJson, err := ioutil.ReadFile("dsblock3.json")
	if err != nil {
		t.Fatal(err.Error())
	}

	var dsBlockT DsBlockT
	err2 := json.Unmarshal(dsJson, &dsBlockT)
	if err2 != nil {
		t.Fatal(err2.Error())
	}

	dsBlock := NewDsBlockFromDsBlockT(&dsBlockT)
	headerBytes := dsBlock.BlockHeader.Serialize()

	bs := &BitVector{}
	headerBytes = dsBlock.Cosigs.CS1.Serialize(headerBytes, uint(len(headerBytes)))
	headerBytes = bs.SetBitVector(headerBytes, uint(len(headerBytes)), dsBlock.Cosigs.B1)
	if "0A4608021220D0DFB05692A202C0A80315785E4A87083BF5A95BB929D7015064FB84F89F6A071A200F00E9D3175300FC287812D201EDCFBFCB8165809606545595BF53700C524648100518032A230A210213D5A7F74B28F3F588FF6520748DBB541986E98F75FA78D6334B2D0AAB4C1E573001380142120A10000000000000000000000000773594004A320A30000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000523D0A230A2103AFEEA358BBFD6A350B59943E11D65386142033C63B921032608DE65334C1DF8C12160A14000000000000000000000000BD1EC8360000816D5AA5010A20BD2C141BD94913038C17251D645CB343CB7868F81CD5BE674E144023472E37661280010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000536F8464BFB2F197CAE6F4E0A9C89C24CCA4B8B104D923A680CB766F91AA8559B6A9B8C1888E9F61F17D8A391EBAECFB4ACF01B49C79C4FEE6E0EC101168299F000AF700" != strings.ToUpper(util.EncodeHex(headerBytes)) {
		t.Log(strings.ToUpper(util.EncodeHex(headerBytes)))
		t.Fatal("ds bytes error")
	}

	commKeys := []string{"0213D5A7F74B28F3F588FF6520748DBB541986E98F75FA78D6334B2D0AAB4C1E57",
		"0239D4CAE39A7AC2F285796BABF7D28DC8EB7767E78409C70926D0929EA2941E36",
		"02D2D695D4A352412E0D32A8BDF6EA3A606D35FE2C2F850C54D68727D065894986",
		"02E5E1BE6C924349F2C2B20CE05A2650B3E56C7722A2E5952EE27D12DEE7A4A6E6",
		"0300AB86B413FAA64A52FB61B5A28A6C361F87A5B0871C4F01C394D261415B0989",
		"03019AF5B10FFE09FB0EE02B59195EF5E6F5BE51D17EAF5604EA452078CD465C4B",
		"0323086D473DF937B6297FB755FA8E57C0FB2760512AED7757748B597C48F797A0",
		"032AEE20CFC59EAEB7838DAC2A9BAF96C8D69CF2C866FB4A3F1DFB02BCFCA356BB",
		"033207325A3CC671034FEBA86EC8D0AA412DF60C7E8292044D510DF582787DCC05",
		"0334AA0F7CA2EAA56B6B752533F9C60777E96C6D1ABE84B463F60ADD89843794AE",
	}

	var pubKeys [][]byte
	for index, key := range commKeys {
		if dsBlock.Cosigs.B2[index] {
			pubKeys = append(pubKeys, util.DecodeHex(key))
		}
	}
	aggregatedPubKey, err3 := multisig.AggregatedPubKey(pubKeys)
	if err3 != nil {
		t.Fatal(err3.Error())
	}
	t.Log("aggregated public key = ", util.EncodeHex(aggregatedPubKey))

	data := make([]byte, 0)
	bns := BIGNumSerialize{}
	data = bns.SetNumber(data, 0, signatureChallengeSize, dsBlock.Cosigs.CS2.R)
	data = bns.SetNumber(data, 0+signatureChallengeSize, signatureChallengeSize, dsBlock.Cosigs.CS2.S)

	signature := util.EncodeHex(data)
	t.Log("signature = ", signature)
	r := util.DecodeHex(signature[0:64])
	s := util.DecodeHex(signature[64:128])

	if !multisig.MultiVerify(aggregatedPubKey, headerBytes, r, s) {
		t.Fail()
	}

}
