package main

import (
	"bufio"
	"fmt"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/access_manager"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/example"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/file_handler"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/notifier"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/p2p"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/transfer"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"log"
	"os"
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
		for i := 0; i < len(commands); i++ {
			// Call by index
			if fmt.Sprintf("%v", i+1) == function {
				commands[i].handler()
				break
			}
			// Call by alias
			called := false
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
				break
			}
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
	privKey, _, err := crypto.GenerateKeyPair(crypto.ECDSA, -1)
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
	access_manager.AllowNode(peerId)
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
	isAllowed := access_manager.IsAllowedNode(peerId)
	fmt.Println(fmt.Sprintf("\t%v", isAllowed))
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handleGetAllowedNodes() {
	allowedNodes := access_manager.GetAllowedNodes()
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
	access_manager.BlockNode(peerId)
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
	isBlocked := access_manager.IsBlockedNode(peerId)
	fmt.Println(fmt.Sprintf("\t%v", isBlocked))
	fmt.Println(colorGreen + "\tSuccess" + colorReset)
}

func handleGetBlockedNodes() {
	blockedNodes := access_manager.GetBlockedNodes()
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
