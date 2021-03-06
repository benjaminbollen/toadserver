package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	log "github.com/eris-ltd/toadserver/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/eris-ltd/toadserver/Godeps/_workspace/src/github.com/eris-ltd/mint-client/mintx/core"
	"github.com/eris-ltd/toadserver/Godeps/_workspace/src/github.com/tendermint/tendermint/types"
	"github.com/eris-ltd/toadserver/Godeps/_workspace/src/github.com/tendermint/tendermint/wire"
)

func UpdateNameReg(fileName, hash string) error {
	nodeAddr := os.Getenv("MINTX_NODE_ADDR")
	signAddr := os.Getenv("MINTX_SIGN_ADDR")
	chainID := os.Getenv("MINTX_CHAINID")
	pubkey := strings.TrimSpace(os.Getenv("MINTX_PUBKEY")) //because bash
	addr := ""
	amtS := "10000"
	nonceS := ""
	feeS := "0"
	name := fileName
	data := hash

	log.WithFields(log.Fields{
		"MINTX_NODE_ADDR": nodeAddr,
		"MINTX_CHAINID":   chainID,
		"MINTX_PUBKEY":    pubkey,
		"MINTX_SIGN_ADDR": signAddr,
		"name":            fileName,
		"data":            data,
		"amount":          amtS,
	}).Warn("Building a nameTx:")

	nTx, err := core.Name(nodeAddr, signAddr, pubkey, addr, amtS, nonceS, feeS, name, data)
	if err != nil {
		return errors.New(fmt.Sprintf("corename error: %v\n", err))
	}
	log.WithField("=>", fmt.Sprintf("%v", nTx)).Warn("Success, nameTx created:")

	//sign but don't broadcast

	log.Warn("Signing transaction...")
	_, err = core.SignAndBroadcast(chainID, nodeAddr, signAddr, nTx, true, false, false)
	if err != nil {
		return errors.New(fmt.Sprintf("sign error: %v\n", err))

	}

	n := new(int64)
	w := new(bytes.Buffer)
	wire.WriteBinary(nTx, w, n, &err)

	// post needs to be to toadserver endpoint, which'll
	// eventually route the TX to the nodes using mindy
	txD := bytes.NewReader(w.Bytes())
	//it can also query for the name reg to ensure things are good
	endpoint := "http://0.0.0.0:11113/" + "receiveNameTx/" //+ hash
	_, err = http.Post(endpoint, "", txD)
	if err != nil {
		return errors.New(fmt.Sprintf("post error: %v\n", err))
	}
	if err := os.Remove(fileName); err != nil {
		return errors.New(fmt.Sprintf("remove file error: %v\n", err))
	}
	return nil

}

func receiveNameTx(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//TODO check valid Name reg
		//str := r.URL.Path[1:]
		//hash := strings.Split(str, "/")[1]

		txData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Warn("error reading body:")
			log.Error(err)
		}

		tx := new(types.NameTx)
		n := new(int64)
		txD := bytes.NewReader(txData)

		wire.ReadBinary(tx, txD, n, &err)
		if err != nil {
			log.Warn("error reading binary:")
			log.Error(err)
		}

		rpcAddr := os.Getenv("MINTX_NODE_ADDR")
		_, err1 := core.Broadcast(tx, rpcAddr)
		if err1 != nil {
			log.Warn("error broadcasting:")
			log.Error(err)
		}
	}
}
