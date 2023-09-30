package main

import (
	"bufio"
	"fmt"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/access_manager"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/example"
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

const helpString string = "" +
	"1.  getHelloMessage:    Hello Message\n" +
	"2.  createKeyPair:      crates crypto key pair\n" +
	"3.  createRandomNode:   creates random node\n" +
	"4.  getNodeId:          get node id of running node\n" +
	"5.  getListenAddresses: get listen addresses of running node\n" +
	"6.  allowNode:          Allow messages from this node\n" +
	"7.  isAllowedNode:      Check if this node is allowed to send messages\n" +
	"8.  getAllowedNodes:    Get all node Ids which can send messages\n" +
	"9.  blockNode:          Block messages from this node\n" +
	"10. isBlockedNode:      Check if this node is blocked\n" +
	"11. getBlockedNodes:    Gets all node Ids which are blocked\n" +
	"12. serveFile:          serve file to peer from running node\n" +
	"13. stopServeFile:      stop serving file with SHA256\n" +
	"14. sendMessage:        send message to peer from running node\n" +
	"15. pauseDownload:      pause downloading file with SHA256 from peer\n" +
	"16. resumeDownload:     resume downloading file with SHA256 from peer\n" +
	"17. stopDownload:       stop downloading file with SHA256 from peer\n" +
	"18. getDownloadStatus:  get status of downloading file with SHA256 from peer\n" +
	"19. flushNotifications: flush notifications asynchronously\n" +
	"20. stopNode:           stop running node\n" +
	"21. setVar (s):         set variables for debug purposes (use using {varName})\n" +
	"22. help   (h):         show help\n" +
	"23. quit   (q):         stop program\n"

var mapVar = make(map[string]string)

func main() {
	running := true
	fmt.Print(helpString)
	for running {
		fmt.Print(colorYellow + "Enter Your Command: " + colorReset)
		function := readString()
		switch function {
		case "":
			break
		case "1", "getHelloMessage":
			handleHelloMessage()
			break
		case "2", "createKeyPair":
			handleCreateKeyPair()
			break
		case "3", "createRandomNode":
			handleCreateRandomNode()
			break
		case "4", "getNodeId":
			handleGetNodeId()
			break
		case "5", "getListenAddresses":
			handleGetListenAddresses()
			break
		case "6", "allowNode":
			handleAllowNode()
			break
		case "7", "isAllowedNode":
			handleIsAllowedNode()
			break
		case "8", "getAllowedNodes":
			handleGetAllowedNodes()
			break
		case "9", "blockNode":
			handleBlockNode()
			break
		case "10", "isBlockedNode":
			handleIsBlockedNode()
			break
		case "11", "getBlockedNodes":
			handleGetBlockedNodes()
			break
		case "12", "serveFile":
			handleServeFile()
			break
		case "13", "stopServeFile":
			handleStopServeFile()
			break
		case "14", "sendMessage":
			handleSendMessage()
			break
		case "15", "pauseDownload":
			handlePauseDownload()
			break
		case "16", "resumeDownload":
			handleResumeDownload()
			break
		case "17", "stopDownload":
			handleStopDownload()
			break
		case "18", "getDownloadStatus":
			handleGetDownloadStatus()
			break
		case "19", "flushNotifications", "flush", "print", "f":
			handleFlushNotifications()
			break
		case "20", "stopNode":
			handleStopNode()
			break
		case "21", "setVar", "s":
			handleSetVar()
			break
		case "22", "help", "h":
			fmt.Print(helpString)
			fmt.Println(colorGreen + "Success" + colorReset)
			break
		case "23", "quit", "q":
			running = false
			fmt.Println(colorGreen + "Success" + colorReset)
			break
		default:
			fmt.Println(colorRed + "Error: invalid command!" + colorReset)
			break
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
