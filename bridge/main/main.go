package main

import (
	"fmt"
	bridge "github.com/DecentRealized/custom-libp2p-mobile/bridge"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/utils"
	"os"
)

func main() {
	protocolDir := os.Args[1]
	dartOutDir := os.Args[2]
	err := bridge.CreateModelExporter(protocolDir, fmt.Sprintf("%s/models/models.dart", dartOutDir))
	utils.CheckError(err)
	err = bridge.CreateDartBinding(fmt.Sprintf("%s/custom_libp2p.dart", dartOutDir))
	utils.CheckError(err)
	err = bridge.GenerateConstants(fmt.Sprintf("%s/custom_libp2p_constants.dart", dartOutDir))
	utils.CheckError(err)
}
