package main

import (
	"bufio"
	"fmt"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/access_manager"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/database"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/example"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/file_handler"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/notifier"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/p2p"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/transfer"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"log"
	"os"
	"strconv"
	"strings"
)

const colorReset = "\033[0m"
const colorRed = "\033[31m"
const colorGreen = "\033[32m"
const colorYellow = "\033[33m"
const colorBlue = "\033[34m"

type command struct {
	name        string
	aliases     []string
	handler     func()
	description string
}

var commands = []command{
	{"getHelloMessage", []string{}, handleHelloMessage, "Prints Hello Message"},
	{"createKeyPair", []string{}, handleCreateKeyPair, "Creates crypto key pair"},
	{"createRandomNode", []string{}, handleCreateRandomNode, "Creates random node"},
	{"getNodeId", []string{}, handleGetNodeId, "Gets node id"},
	{"getListenAddresses", []string{}, handleGetListenAddresses, "Gets listen addresses"},
	{"connectToPeer", []string{}, handleConnectToPeer, "Connects to peer"},
	{"checkConnectionStatus", []string{}, handleCheckConnectionStatus, "Checks connection status"},
	{"allowNode", []string{}, handleAllowNode, "Allow messages from this node"},
	{"isAllowedNode", []string{}, handleIsAllowedNode, "Check if this node is allowed to send messages"},
	{"getAllowedNodes", []string{}, handleGetAllowedNodes, "Get all node Ids which can send messages"},
	{"blockNode", []string{}, handleBlockNode, "Block messages from this node"},
	{"isBlockedNode", []string{}, handleIsBlockedNode, "Check if this node is blocked"},
	{"getBlockedNodes", []string{}, handleGetBlockedNodes, "Gets all node Ids which are blocked"},
	{"setDownloadPath", []string{}, handleSetDownloadPath, "Sets download path"},
	{"getDownloadPath", []string{}, handleGetDownloadPath, "Gets download path"},
	{"setDatabaseDirectory", []string{}, handleSetDatabaseDirectory, "Sets database base path"},
	{"getDatabaseDirectory", []string{}, handleGetDatabaseDirectory, "Gets database base path"},
	{"serveFile", []string{}, handleServeFile, "Serve file to peer from running node"},
	{"stopServeFile", []string{}, handleStopServeFile, "Stop serving file with SHA256"},
	{"sendMessage", []string{}, handleSendMessage, "Send message to peer from running node"},
	{"pauseDownload", []string{}, handlePauseDownload, "Pause downloading file with SHA256 from peer"},
	{"resumeDownload", []string{}, handleResumeDownload, "Resume downloading file with SHA256 from peer"},
	{"stopDownload", []string{}, handleStopDownload, "Stop downloading file with SHA256 from peer"},
	{"getDownloadStatus", []string{}, handleGetDownloadStatus, "Get download status of file with SHA256"},
	{"flushNotifications", []string{}, handleFlushNotifications, "Flush notifications asynchronously"},
	{"stopNode", []string{}, handleStopNode, "Stop running node"},
	{"setVar", []string{"s", "S"}, handleSetVar, "Set variables for debug purposes"},
}

var mapVar = make(map[string]string)

func main() {
	commands = append(commands,
		command{"help", []string{"h", "H"}, handleHelp, "Show help"},
		command{"quit", []string{"q", "Q"}, handleQuit, "Quit program"})
	handleHelp()
	for true {
		fmt.Print(colorYellow + "Enter Your Command: " + colorReset)
		function := readString()
		called := false
		for i := 0; i < len(commands); i++ {
			// Call by index
			if fmt.Sprintf("%v", i+1) == function {
				commands[i].handler()
				called = true
				break
			}
			for j := 0; j < len(commands[i].aliases); j++ {
				if commands[i].aliases[j] == function {
					commands[i].handler()
					called = true
					break
				}
			}
			if called {
				break
			}
			// Call by name
			if commands[i].name == function {
				commands[i].handler()
				called = true
				break
			}
		}
		if !called {
			fmt.Println(colorRed + "\tError: invalid input!" + colorReset)
		}
	}
}

func handleHelloMessage() {
	fmt.Print(colorYellow + "\tUsername: " + colorReset)
	username := readString()
	message, err := example.GetHelloMessage(username)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	fmt.Println("\t" + message)
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handleCreateKeyPair() {
	privKey, pubKey, err := crypto.GenerateKeyPair(crypto.ECDSA, -1)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	privBytes, err := privKey.Raw()
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	pubBytes, err := pubKey.Raw()
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	fmt.Println("\tprivKey: ", privBytes)
	fmt.Println("\tpubKey: ", pubBytes)
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handleCreateRandomNode() {
	fmt.Print(colorYellow + "\tUseInternet (Y/n): " + colorReset)
	useInternetStr := readString()
	useInternet := false
	switch useInternetStr {
	case "Y", "y":
		useInternet = true
		break
	case "N", "n":
		useInternet = false
		break
	default:
		fmt.Println(colorRed + "\tError: invalid input!" + colorReset)
		return
	}
	fmt.Print(colorYellow + "\tPreset [1-5/0,n,N]: " + colorReset)
	presetIndexString := readString()

	presets := [][]byte{
		{8, 3, 18, 121, 48, 119, 2, 1, 1, 4, 32, 219, 165, 43, 30, 129, 151, 117, 76, 26, 215, 23, 14, 30, 136, 59, 123, 6, 96, 203, 37, 100, 123, 142, 137, 37, 172, 146, 130, 213, 134, 95, 126, 160, 10, 6, 8, 42, 134, 72, 206, 61, 3, 1, 7, 161, 68, 3, 66, 0, 4, 251, 15, 103, 110, 186, 39, 160, 137, 71, 54, 61, 54, 38, 129, 246, 228, 250, 190, 155, 228, 86, 117, 167, 194, 167, 164, 206, 199, 158, 56, 155, 144, 191, 173, 200, 85, 92, 191, 127, 249, 80, 171, 189, 182, 96, 186, 73, 181, 236, 141, 192, 20, 5, 251, 83, 52, 57, 214, 228, 121, 51, 226, 134, 73},
		{8, 3, 18, 121, 48, 119, 2, 1, 1, 4, 32, 211, 43, 126, 245, 219, 244, 251, 3, 73, 143, 245, 235, 203, 211, 44, 56, 48, 192, 121, 80, 243, 161, 233, 13, 208, 82, 248, 219, 58, 56, 131, 81, 160, 10, 6, 8, 42, 134, 72, 206, 61, 3, 1, 7, 161, 68, 3, 66, 0, 4, 215, 141, 108, 14, 150, 105, 22, 168, 83, 210, 255, 127, 128, 35, 104, 230, 167, 159, 171, 12, 241, 17, 42, 154, 53, 124, 5, 231, 183, 29, 58, 174, 248, 125, 100, 81, 79, 244, 96, 60, 110, 187, 161, 123, 219, 148, 156, 101, 102, 21, 78, 198, 191, 240, 209, 52, 207, 239, 161, 68, 4, 101, 93, 27},
		{8, 3, 18, 121, 48, 119, 2, 1, 1, 4, 32, 191, 42, 121, 111, 133, 42, 6, 241, 243, 135, 1, 155, 161, 180, 76, 226, 2, 128, 116, 129, 202, 152, 0, 81, 49, 46, 29, 64, 224, 36, 91, 113, 160, 10, 6, 8, 42, 134, 72, 206, 61, 3, 1, 7, 161, 68, 3, 66, 0, 4, 87, 246, 62, 54, 161, 228, 192, 63, 50, 235, 16, 140, 114, 65, 42, 174, 10, 48, 94, 26, 126, 181, 104, 116, 86, 249, 190, 109, 224, 35, 66, 255, 37, 246, 56, 168, 118, 23, 70, 199, 176, 117, 227, 139, 90, 210, 39, 190, 161, 83, 39, 207, 106, 114, 0, 104, 174, 58, 147, 160, 234, 79, 31, 92},
		{8, 3, 18, 121, 48, 119, 2, 1, 1, 4, 32, 102, 215, 54, 0, 11, 242, 188, 199, 195, 165, 145, 5, 74, 63, 208, 124, 17, 135, 20, 34, 125, 2, 244, 254, 239, 45, 111, 87, 125, 47, 173, 100, 160, 10, 6, 8, 42, 134, 72, 206, 61, 3, 1, 7, 161, 68, 3, 66, 0, 4, 80, 162, 169, 234, 209, 229, 27, 48, 52, 44, 114, 71, 232, 221, 103, 27, 2, 90, 143, 243, 54, 140, 158, 146, 208, 56, 26, 177, 106, 122, 118, 142, 179, 235, 169, 157, 254, 17, 35, 112, 165, 11, 38, 162, 106, 115, 229, 61, 132, 144, 25, 159, 8, 67, 2, 0, 221, 133, 151, 196, 206, 143, 244, 159},
		{8, 3, 18, 121, 48, 119, 2, 1, 1, 4, 32, 224, 148, 182, 27, 183, 198, 203, 216, 184, 0, 171, 210, 168, 190, 25, 144, 29, 226, 44, 51, 223, 50, 75, 249, 166, 33, 149, 9, 62, 13, 132, 116, 160, 10, 6, 8, 42, 134, 72, 206, 61, 3, 1, 7, 161, 68, 3, 66, 0, 4, 78, 64, 17, 176, 240, 255, 125, 2, 206, 54, 149, 155, 231, 227, 20, 39, 167, 139, 184, 142, 232, 65, 45, 8, 95, 190, 93, 46, 73, 159, 75, 102, 28, 160, 27, 21, 184, 96, 76, 89, 157, 52, 34, 132, 190, 38, 7, 202, 102, 79, 252, 248, 18, 18, 18, 97, 201, 162, 167, 181, 235, 211, 219, 11},
	}
	presetIndex, err := strconv.ParseInt(presetIndexString, 10, 64)
	if err != nil {
		presetIndex = 0
		fmt.Println(colorYellow + fmt.Sprintf("\tNot valid preset, generating random.") + colorReset)
	}

	privKey, _, err := crypto.GenerateKeyPair(crypto.ECDSA, -1)
	if presetIndex != 0 {
		privKey, err = crypto.UnmarshalPrivateKey(presets[presetIndex])
	}
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	err = p2p.CreateNode(privKey, useInternet)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handleGetNodeId() {
	nodeId, err := p2p.GetNodeId()
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	fmt.Println("\t" + nodeId.String())
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handleGetListenAddresses() {
	listenAddresses, err := p2p.GetListenAddresses()
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	fmt.Println(fmt.Sprintf("\t%v", listenAddresses))
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handleConnectToPeer() {
	fmt.Print(colorYellow + "\tPeer Id: " + colorReset)
	peerIdStr := readString()
	peerId, err := peer.Decode(peerIdStr)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	err = p2p.ConnectToPeer(peerId)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handleCheckConnectionStatus() {
	fmt.Print(colorYellow + "\tPeer Id: " + colorReset)
	peerIdStr := readString()
	peerId, err := peer.Decode(peerIdStr)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	connections, err := p2p.CheckConnectionStatus(peerId)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	fmt.Println(fmt.Sprintf("\t%v", connections))
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handleAllowNode() {
	fmt.Print(colorYellow + "\tPeer Id: " + colorReset)
	peerIdStr := readString()
	peerId, err := peer.Decode(peerIdStr)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	err = access_manager.AllowNode(peerId)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handleIsAllowedNode() {
	fmt.Print(colorYellow + "\tPeer Id: " + colorReset)
	peerIdStr := readString()
	peerId, err := peer.Decode(peerIdStr)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	isAllowed, err := access_manager.IsAllowedNode(peerId)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	fmt.Println(fmt.Sprintf("\t%v", isAllowed))
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handleGetAllowedNodes() {
	allowedNodes, err := access_manager.GetAllowedNodes()
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	fmt.Println(fmt.Sprintf("\t%v", allowedNodes))
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handleBlockNode() {
	fmt.Print(colorYellow + "\tPeer Id: " + colorReset)
	peerIdStr := readString()
	peerId, err := peer.Decode(peerIdStr)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	err = access_manager.BlockNode(peerId)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handleIsBlockedNode() {
	fmt.Print(colorYellow + "\tPeer Id: " + colorReset)
	peerIdStr := readString()
	peerId, err := peer.Decode(peerIdStr)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	isBlocked, err := access_manager.IsBlockedNode(peerId)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	fmt.Println(fmt.Sprintf("\t%v", isBlocked))
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handleGetBlockedNodes() {
	blockedNodes, err := access_manager.GetBlockedNodes()
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	fmt.Println(fmt.Sprintf("\t%v", blockedNodes))
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handleSetDownloadPath() {
	fmt.Print(colorYellow + "\tFile Path: " + colorReset)
	filePath := readString()
	err := file_handler.SetDownloadPath(filePath)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handleGetDownloadPath() {
	downloadPath := file_handler.GetDownloadPath()
	fmt.Println("\t" + downloadPath)
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handleSetDatabaseDirectory() {
	fmt.Print(colorYellow + "\tDatabase Directory: " + colorReset)
	databasePath := readString()
	err := database.SetDatabaseDirectory(databasePath)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handleGetDatabaseDirectory() {
	databasePath := database.GetDatabaseDirectory()
	fmt.Println("\t" + databasePath)
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handleServeFile() {
	fmt.Print(colorYellow + "\tFile Path: " + colorReset)
	filePath := readString()
	fmt.Print(colorYellow + "\tPeer ID: " + colorReset)
	peerIdStr := readString()
	peerId, err := peer.Decode(peerIdStr)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	metadata, err := transfer.ServeFile(filePath, peerId)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	fmt.Printf("\t%v\n", metadata)
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handleStopServeFile() {
	fmt.Print(colorYellow + "\tSHA256: " + colorReset)
	fileSHA256 := readString()
	err := transfer.StopServingFile(fileSHA256)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handleSendMessage() {
	fmt.Print(colorYellow + "\tPeer ID: " + colorReset)
	peerIdStr := readString()
	peerId, err := peer.Decode(peerIdStr)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	fmt.Print(colorYellow + "\tMessage: " + colorReset)
	message := readString()
	err = transfer.SendMessage(peerId, message)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handlePauseDownload() {
	fmt.Print(colorYellow + "\tSHA256: " + colorReset)
	fileSHA256 := readString()
	fmt.Print(colorYellow + "\tPeer ID: " + colorReset)
	peerIdStr := readString()
	peerId, err := peer.Decode(peerIdStr)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	err = transfer.PauseDownload(fileSHA256, peerId)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handleResumeDownload() {
	fmt.Print(colorYellow + "\tSHA256: " + colorReset)
	fileSHA256 := readString()
	fmt.Print(colorYellow + "\tPeer ID: " + colorReset)
	peerIdStr := readString()
	peerId, err := peer.Decode(peerIdStr)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	err = transfer.ResumeDownload(fileSHA256, peerId)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
}

func handleStopDownload() {
	fmt.Print(colorYellow + "\tSHA256: " + colorReset)
	fileSHA256 := readString()
	fmt.Print(colorYellow + "\tPeer ID: " + colorReset)
	peerIdStr := readString()
	peerId, err := peer.Decode(peerIdStr)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	err = transfer.StopDownload(fileSHA256, peerId)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
}

func handleGetDownloadStatus() {
	fmt.Print(colorYellow + "\tSHA256: " + colorReset)
	fileSHA256 := readString()
	fmt.Print(colorYellow + "\tPeer ID: " + colorReset)
	peerIdStr := readString()
	peerId, err := peer.Decode(peerIdStr)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	status, err := transfer.GetDownloadStatus(fileSHA256, peerId)
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	fmt.Println(status)
}

func handleFlushNotifications() {
	fmt.Println(colorYellow + "\tFlushing notifications..." + colorReset)
	go func() {
		flushed, err := notifier.FlushNotifications()
		if err != nil {
			fmt.Println(colorRed + fmt.Sprintf("\tError Flush: %v", err) + colorReset)
			return
		}
		fmt.Printf(colorReset)
		for i := range flushed.GetNotification() {
			log.Printf("\tFlushed: %v", flushed.GetNotification()[i])
		}
		log.Println(colorGreen + "\tFlushNotifications: Success" + colorReset)
	}()
}

func handleStopNode() {
	err := p2p.StopNode()
	if err != nil {
		fmt.Println(colorRed + fmt.Sprintf("\tError: %v", err) + colorReset)
		return
	}
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handleHelp() {
	firstSize := 0
	secondSize := 0
	thirdSize := 0
	for i := 0; i < len(commands); i++ {
		curFirstSize := len(colorReset + fmt.Sprintf("%v. ", i+1))
		curSecondSize := len(fmt.Sprintf("%v: ", commands[i].name))
		curThirdSize := len(commands[i].description)
		if firstSize < curFirstSize {
			firstSize = curFirstSize
		}
		if secondSize < curSecondSize {
			secondSize = curSecondSize
		}
		if thirdSize < curThirdSize {
			thirdSize = curThirdSize
		}
	}
	for i := 0; i < len(commands); i++ {
		first := colorReset + fmt.Sprintf("%v. ", i+1)
		first += strings.Repeat(" ", firstSize-len(first))
		second := fmt.Sprintf("%v: ", commands[i].name)
		second += strings.Repeat(" ", secondSize-len(second))
		third := commands[i].description
		third += strings.Repeat(" ", thirdSize-len(third))
		fourth := ""
		for _, alias := range commands[i].aliases {
			fourth += fmt.Sprintf(" [%v]", alias)
		}
		fmt.Printf("%v%v%v%v\n", first, second, third, fourth)
	}
}

func handleQuit() {
	os.Exit(0)
}

func handleSetVar() {
	fmt.Print(colorYellow + "\tVar Name: " + colorReset)
	varName := readString()
	fmt.Print(colorYellow + "\tVar Value: " + colorReset)
	varValue := readString()
	mapVar[varName] = varValue
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

var reader = bufio.NewReader(os.Stdin)

func readString() string {
	fmt.Print(colorBlue)
	_readString, err := reader.ReadString('\n')
	fmt.Print(colorReset)
	if err != nil {
		panic(err)
	}
	_readString = strings.Trim(_readString, " \n\r\t")
	if len(_readString) >= 2 && _readString[0] == '{' && _readString[len(_readString)-1] == '}' {
		varVal := mapVar[_readString[1:len(_readString)-1]]
		return varVal
	}
	return _readString
}
